package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
)

// APIServer represents an HTTP server for handling API requests.
type APIServer struct {
	listenAddr string
	store      Storage
}

// NewAPIServer creates a new instance of APIServer.
func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Run starts the API server and listens for incoming requests.
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandlerFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleAccountAndID), s.store))
	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.handleTransfer))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodPost:
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return err
		}

		acc, err := s.store.GetAccountByNumber(req.Number)
		if err != nil {
			return err // TODO: handle this resp as JSON
		}

		//fmt.Printf("%+v\n", acc)

		if !acc.ValidatePassword(req.Password) {
			return fmt.Errorf("not authorized")
		}

		token, err := createJWT(acc)
		if err != nil {
			return err
		}

		resp := LoginResponse{
			Number: acc.Number,
			Token:  token,
		}

		return WriteJSON(w, http.StatusOK, resp)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleAccount handles requests related to account management.
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccounts(w, r)
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleGetAccountByID handles requests to retrieve an account by ID.
func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

// handleCreateAccount handles requests to create a new account.
func (s *APIServer) handleAccountAndID(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	case http.MethodPut:
		return s.handleUpdateAccount(w, r)
	default:
		return fmt.Errorf("unsupported method: %s", r.Method)
	}
}

// handleDeleteAccount handles requests to delete an account.
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Println("JWT token: ", tokenString)

	// Retrieve the account from the database to get the most up-to-date information
	// about the account, including any database-generated fields or default values
	createdAccount, err := s.store.GetAccountByID(account.ID)
	if err != nil {
		return err
	}

	// Return the newly created account in the response
	return WriteJSON(w, http.StatusOK, createdAccount)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	var account Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		return err
	}

	account.ID = id

	if err := s.store.UpdateAccount(&account); err != nil {
		return err
	}

	// Retrieve the account from the database to get the most up-to-date information
	// about the account, including any database-generated fields or default values
	updatedAccount, err := s.store.GetAccountByID(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, updatedAccount)
}

// handleTransfer handles requests to transfer funds between accounts.
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]string{"deleted": id})
}

// handleTransfer handles requests to transfer funds between accounts.
// TODO: finish handler
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferRequest := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferRequest); err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(r.Body)

	return WriteJSON(w, http.StatusOK, transferRequest)
}

// WriteJSON writes the given data as JSON to the HTTP response with the provided status code.
// It sets the "Content-Type" header to "application/json; charset=utf-8".
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

// createJWT generates a JSON Web Token (JWT) containing the specified account information.
// It returns the signed JWT token as a string and any error encountered during token generation.
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	return token.SignedString([]byte(secret))
}

func permissionDeniedError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "permission denied"})
}

// withJWTAuth adds JWT authentication to the provided HTTP handler.
// It validates the included JWT and authorizes the request.
// If the JWT is invalid or the request is unauthorized, it responds with a permission denied error.
// Returns an HTTP handler that wraps the original handler.
func withJWTAuth(fn http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDeniedError(w)
			return
		}

		if !token.Valid {
			permissionDeniedError(w)
			return
		}

		userID, err := getID(r)
		if err != nil {
			permissionDeniedError(w)
			return
		}
		account, err := s.GetAccountByID(userID)
		if err != nil {
			permissionDeniedError(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int(claims["accountNumber"].(float64)) {
			permissionDeniedError(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "invalid token"})
			return
		}

		fn(w, r)
	}
}

// validateJWT validates the given JWT token string. It verifies the signature
// and checks if the token is well-formed and valid.
func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

// apiFunc is a type representing a function that handles HTTP requests and returns an error.
type apiFunc func(http.ResponseWriter, *http.Request) error

// ApiError represents an error response in JSON format.
type ApiError struct {
	Error string `json:"error"`
}

// makeHTTPHandlerFunc creates an HTTP handler function from the given apiFunc.
// It calls the provided function f to handle HTTP requests, and if an error occurs, it writes
// the error response as JSON with status code http.StatusBadRequest.
func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// getID extracts the ID parameter from the URL path of the HTTP request r.
// It returns the extracted ID and an error if the ID is invalid or not found in the request.
func getID(r *http.Request) (string, error) {
	id := mux.Vars(r)["id"]

	_, err := uuid.Parse(id)
	if err != nil {
		return id, fmt.Errorf("invalid account id %s: %v", id, err)
	}
	return id, nil
}
