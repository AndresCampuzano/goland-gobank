<p align="center">
  <img src="https://cdn-icons-png.flaticon.com/512/6295/6295417.png" width="100" />
</p>
<h1 align="center">GOLAND-GOBANK</h1>
<p align="center">
    <em>Should be called Golang instead of Goland, I know...</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/AndresCampuzano/goland-gobank?style=flat&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/AndresCampuzano/goland-gobank?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/AndresCampuzano/goland-gobank?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/AndresCampuzano/goland-gobank?style=flat&color=0080ff" alt="repo-language-count">
<p>
<p align="center">
		<em>Developed with the software and tools below.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
</p>
<hr>

##  Quick Links

> - [ Overview](#-overview)
> - [ Features](#-features)
> - [ Technologies Used](#-technologies-used)
> - [ Repository Structure](#-repository-structure)
> - [ Getting Started](#-getting-started)

---

##  Overview

Bank API is a RESTful API server for managing user accounts and transactions. It provides endpoints for user authentication, account management, and fund transfers.

---

##  Features

- User authentication using JSON Web Tokens (JWT)
- Create, read, update, and delete (CRUD) operations on user accounts
- Routes validation
  - `/account/{account-id}` can be just executed by the same user with a valid token
- Transfer funds between accounts (not finished)


---

##  Technologies Used

- Go (Golang) programming language
- PostgreSQL database for storing account information
- Gorilla Mux for routing HTTP requests
- JWT for user authentication

---

##  Repository Structure

```sh
└── goland-gobank/
    ├── Makefile
    ├── api.go
    ├── auth.go
    ├── go.mod
    ├── go.sum
    ├── main.go
    ├── storage.go
    ├── types.go
    ├── types_test.go
    └── README.md
```

---

##  Getting Started

###  Installation

1. Clone the goland-gobank repository:

```sh
git clone https://github.com/AndresCampuzano/goland-gobank
```

2. Change to the project directory:

```sh
cd goland-gobank
```

3. Run the project:

```sh
make run
```

---
