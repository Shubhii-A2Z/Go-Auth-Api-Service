# Go Authentication API Service - Architecture & Generic Template Guide

A robust, production-ready RESTful API boilerplate written in Go. This repository demonstrates **Clean Layered Architecture**, **JWT Authentication & Authorization**, **Bcrypt Password Hashing**, **Chi Router Configuration**, **Goose Database Migrations**, **Rate Limiting**, and **Reverse Proxying**.

Use this repository as a generic template for bootstrapping new Go backend services.

---

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Project Structure](#project-structure)
3. [Prerequisites & Environment Setup](#prerequisites--environment-setup)
4. [How to Run the Project](#how-to-run-the-project)
5. [Database & Goose Migrations](#database--goose-migrations)
6. [HTTP Routing & Chi Router](#http-routing--chi-router)
7. [Authentication & Security Logic](#authentication--security-logic)
8. [Middleware Architecture](#middleware-architecture)
9. [Data Access & Repository Pattern](#data-access--repository-pattern)
10. [Dependency Injection & App Assembly](#dependency-injection--app-assembly)
11. [API Endpoints Summary](#api-endpoints-summary)
12. [Node.js & OOP Developer's Guide to Go & System Design Concepts](#nodejs--oop-developers-guide-to-go--system-design-concepts)
13. [External Dependencies & Libraries Breakdown](#external-dependencies--libraries-breakdown)
    - [12.1 Classes & Object-Oriented Programming](#121-classes--object-oriented-programming)
    - [12.2 Interfaces & Polymorphism](#122-interfaces--polymorphism)
    - [12.3 Enums in Go](#123-enums-in-go)
    - [12.4 Visibility: Public, Private, Protected](#124-visibility-public-private-protected)
    - [12.5 Static Members & Methods](#125-static-members--methods)
    - [12.6 Final, Readonly, and Immutability](#126-final-readonly-and-immutability)
    - [12.7 How Middlewares Work in Go vs. Node.js (Express)](#127-how-middlewares-work-in-go-vs-nodejs-express)
    - [12.8 Dependency Inversion Principle (DIP) & Clean Architecture](#128-dependency-inversion-principle-dip--clean-architecture)
    - [12.9 Starter Commands: Node.js vs. Go Comparison](#129-starter-commands-nodejs-vs-go-comparison)

---

## Architecture Overview

This project strictly adheres to **Separation of Concerns** using a 4-tier clean architecture:

```
    ┌─────────────────────────────────────────┐
    │           HTTP / Router Layer           │  (routers, middlewares)
    └────────────────────┬────────────────────┘
                         │
    ┌────────────────────▼────────────────────┐
    │            Controller Layer             │  (controllers, dtos)
    └────────────────────┬────────────────────┘
                         │
    ┌────────────────────▼────────────────────┐
    │              Service Layer              │  (services - Business Logic)
    └────────────────────┬────────────────────┘
                         │
    ┌────────────────────▼────────────────────┐
    │       Repository / Storage Layer        │  (db/repositories, models)
    └────────────────────┬────────────────────┘
                         │
    ┌────────────────────▼────────────────────┐
    │            Database Driver              │  (config/db - MySQL)
    └─────────────────────────────────────────┘
```

---

## Project Structure

```
.
├── app/                  # Application bootstrapping and server lifecycle
│   └── application.go    # Dependency injection wiring & http.Server init
├── config/               # Configuration loading
│   ├── db/               # MySQL driver setup and connection pinging
│   │   └── db.go
│   └── env/              # Environment variable loading (.env wrapper)
│       └── env.go
├── controllers/          # HTTP Request handlers
│   ├── ping.go           # Health check handler
│   └── user.go           # User registration, login, and retrieval handlers
├── db/                   # Database scripts and data access
│   ├── migrations/       # Goose SQL migration files
│   └── repositories/     # Raw SQL queries & interface implementations
│       └── users.go
├── dtos/                 # Data Transfer Objects with validation tags
│   └── auth.go
├── middlewares/          # HTTP Middlewares
│   ├── jwt.auth.go       # Bearer token JWT authentication guard
│   ├── rate.limiter.go   # Token-bucket request rate limiting
│   └── validator.go      # Request body decoding & Context payload injection
├── models/               # Domain entities mapping database tables
│   └── user.go
├── routers/              # Chi router setup and route mapping
│   ├── router.go         # Global router, middleware mounting, reverse proxy
│   └── userRouter.go     # Sub-router mapping user endpoints
├── services/             # Core business logic
│   └── userService.go    # Password hashing, JWT signing, user validation
├── utils/                # Helper utilities
│   ├── auth.go           # Bcrypt hashing & password matching
│   ├── json.go           # Standardized HTTP JSON responses
│   └── proxy.go          # Reverse proxy request forwarder
├── docker-compose.yml    # Docker configuration for local database
├── Makefile              # Automation tasks (migrations, build)
├── go.mod                # Module dependencies manifest
└── main.go               # Application entry point
```

---

## Prerequisites & Environment Setup

### Prerequisites
* **Go** (version 1.25.0 or later)
* **Docker & Docker Compose**
* **Goose CLI** (`go install github.com/pressly/goose/v3/cmd/goose@latest`)
* **Make** tool

### Environment Variables (.env)
Create a `.env` file in the root directory:

```env
PORT=:8080
DB_USER=root
DB_PASS=root
DB_NET=tcp
DB_ADDR=127.0.0.1:3306
DB_NAME=auth_dev
DB_URL=root:root@tcp(127.0.0.1:3306)/auth_dev?parseTime=true
MIGRATIONS_FOLDER=db/migrations
JWT_SECRET=super_secret_jwt_key
```

---

## How to Run the Project

### Step 0: Download Dependencies
Download all project dependencies specified in `go.mod`:
```bash
go mod download
```

### Step 1: Start Database Container
```bash
docker-compose up -d
```

### Step 2: Run Database Migrations
Apply pending migrations to set up the MySQL database schema:
```bash
make migrate-up
```

### Step 3: Run the Application
```bash
go run main.go
```
The server starts on `http://localhost:8080`.

---

## Database & Goose Migrations

### Concept
Database migrations provide version control for your database schema. **Goose** is used to create and apply SQL migration files (`.sql`) with `-- +goose Up` and `-- +goose Down` annotations.

### Code Implementation

#### 1. Database Connection (`config/db/db.go`)
Constructs a MySQL Data Source Name (DSN) from environment variables, opens a connection pool via `database/sql`, and verifies connectivity with `Ping()`:

```go
package config

import (
    env "AuthInGo/config/env"
    "database/sql"
    "fmt"
    "github.com/go-sql-driver/mysql"
)

func SetupDB() (*sql.DB, error) {
    cfg := mysql.NewConfig()
    cfg.User = env.GetString("DB_USER", "root")
    cfg.Passwd = env.GetString("DB_PASS", "root")
    cfg.Net = env.GetString("DB_NET", "tcp")
    cfg.Addr = env.GetString("DB_ADDR", "127.0.0.1:3306")
    cfg.DBName = env.GetString("DB_NAME", "auth_dev")
    cfg.ParseTime = true

    // Open connection
    db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        return nil, err
    }

    // Verify connection alive
    if pingErr := db.Ping(); pingErr != nil {
        return nil, pingErr
    }

    return db, nil
}
```

#### 2. Goose Makefile Commands (`Makefile`)
The `Makefile` automatically loads `.env` variables and exposes convenient targets for managing schema changes:

```makefile
include .env
export

# Create a new migration file: make migrate-create name="create_users_table"
migrate-create:
	goose -dir $(MIGRATIONS_FOLDER) create $(name) sql

# Apply all pending migrations: make migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_FOLDER) mysql "$(DB_URL)" up
```

#### 3. Migration File Example (`db/migrations/20260227174441_create_user_table.sql`)
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

---

## HTTP Routing & Chi Router

### Concept
`github.com/go-chi/chi/v5` is a lightweight, idiomatic HTTP router for Go. It supports modular sub-routers, inline middleware chaining via `.With()`, and path parameter extraction.

### Code Implementation

#### 1. Main Router Wireup (`routers/router.go`)
```go
package router

import (
    "AuthInGo/controllers"
    "AuthInGo/middlewares"
    "AuthInGo/utils"
    "github.com/go-chi/chi/v5"
)

type Router interface {
    Register(r chi.Router)
}

func SetUpRouter(UserRouter Router) *chi.Mux {
    chiRouter := chi.NewRouter()

    // Global Rate Limiter Middleware
    chiRouter.Use(middlewares.RateLimitMiddleware)

    // Public Endpoint
    chiRouter.Get("/ping", controllers.PingHandler)

    // Reverse Proxy Endpoint
    chiRouter.HandleFunc("/fakestoreservice/*", utils.ProxyToService("https://fakestoreapi.com", "/fakestoreservice"))

    // Mount Sub-Router
    UserRouter.Register(chiRouter)

    return chiRouter
}
```

#### 2. Endpoint Mapping (`routers/userRouter.go`)
```go
package router

import (
    "AuthInGo/controllers"
    "AuthInGo/middlewares"
    "github.com/go-chi/chi/v5"
)

type UserRouter struct {
    UserController controllers.UserController
}

func NewUserRouter(userController controllers.UserController) Router {
    return &UserRouter{
        UserController: userController,
    }
}

func (ur *UserRouter) Register(r chi.Router) {
    // Signup (with Request Body Validation middleware)
    r.With(middlewares.UserCreateRequestValidator).Post("/signup", ur.UserController.RegisterUser)

    // Login (with Request Body Validation middleware)
    r.With(middlewares.UserLoginRequestValidator).Post("/login", ur.UserController.LoginUser)

    // Protected Route (Guarded with JWT Authentication middleware)
    r.With(middlewares.JWTAuthMiddleware).Get("/profile/{id}", ur.UserController.GetUserById)

    // Public List
    r.Get("/profiles", ur.UserController.GetUsers)
}
```

---

## Authentication & Security Logic

### Concept
* **Password Hashing**: Plaintext passwords are never stored. Passwords are salted and hashed using `golang.org/x/crypto/bcrypt`.
* **JWT Token Generation**: Upon successful login, a signed JSON Web Token (JWT) is generated containing claims (`email`, `id`) signed with HS256 (`github.com/golang-jwt/jwt/v5`).

### Code Implementation

#### 1. Bcrypt Password Utilities (`utils/auth.go`)
```go
package utils

import (
    "golang.org/x/crypto/bcrypt"
)

func HashPassword(plainPassword string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

func CheckPasswordHash(plainPassword string, hashPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
    return err == nil
}
```

#### 2. User Service - Registration & Login (`services/userService.go`)
```go
package services

import (
    env "AuthInGo/config/env"
    db "AuthInGo/db/repositories"
    "AuthInGo/dtos"
    "AuthInGo/models"
    "AuthInGo/utils"
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

type UserServiceImpl struct {
    userRepository db.UserRepository
}

// User Registration: Hash password and save to DB
func (u *UserServiceImpl) CreateUser(payload *dtos.CreateUserRequestDTO) error {
    hashPass, err := utils.HashPassword(payload.Password)
    if err != nil {
        return err
    }
    return u.userRepository.Create(payload.Username, payload.Email, hashPass)
}

// User Login: Verify password and issue JWT token
func (u *UserServiceImpl) LoginUser(payload *dtos.LoginUserRequestDTO) (string, error) {
    user, err := u.userRepository.GetByEmail(payload.Email)
    if err != nil || user == nil {
        return "", fmt.Errorf("User not found with given email")
    }

    // Verify Password Hash
    if !utils.CheckPasswordHash(payload.Password, user.Password) {
        return "", fmt.Errorf("Incorrect Password")
    }

    // Build JWT Claims
    jwtPayload := jwt.MapClaims{
        "email": user.Email,
        "id":    user.Id,
    }

    // Sign Token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtPayload)
    tokenString, err := token.SignedString([]byte(env.GetString("JWT_SECRET", "SECRET")))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
```

---

## Middleware Architecture

### Concept
Middlewares sit between the router and HTTP handlers to inspect, validate, rate-limit, or restrict incoming HTTP requests.

### Code Implementation

#### 1. JWT Auth Middleware Guard (`middlewares/jwt.auth.go`)
Validates the `Authorization: Bearer <token>` header on protected endpoints:

```go
package middlewares

import (
    env "AuthInGo/config/env"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims := jwt.MapClaims{}

        // Verify token signature against secret key
        _, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
            return []byte(env.GetString("JWT_SECRET", "SECRET")), nil
        })
        if err != nil {
            http.Error(w, "Invalid Token", http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

#### 2. Request Validator & Context Payload Middleware (`middlewares/validator.go`)
Decodes the request body once in the middleware and passes the typed DTO through Go's `context.WithValue`:

```go
package middlewares

import (
    "AuthInGo/dtos"
    "AuthInGo/utils"
    "context"
    "net/http"
)

func UserCreateRequestValidator(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var payload dtos.CreateUserRequestDTO

        if err := utils.ReadJsonBody(r, &payload); err != nil {
            utils.WriteJsonErrorResponse(w, http.StatusBadRequest, "Invalid request body")
            return
        }

        // Store validated DTO in context
        ctx := context.WithValue(r.Context(), "payload", payload)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

#### 3. Token-Bucket Rate Limiter Middleware (`middlewares/rate.limiter.go`)
Uses `golang.org/x/time/rate` to enforce request limits (e.g., max 5 requests per 30 seconds):

```go
package middlewares

import (
    "net/http"
    "time"
    "golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Every(30*time.Second), 5)

func RateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

---

## Data Access & Repository Pattern

### Concept
The repository layer encapsulates raw SQL operations using standard `database/sql`, keeping database logic completely isolated from business logic.

### Code Implementation (`db/repositories/users.go`)
```go
package db

import (
    "AuthInGo/models"
    "database/sql"
)

type UserRepository interface {
    Create(username string, email string, hashPassword string) error
    GetByEmail(email string) (*models.User, error)
    GetById() (*models.User, error)
    GetAll() error
}

type UserRepositoryImpl struct {
    db *sql.DB
}

func NewUserRepository(_db *sql.DB) *UserRepositoryImpl {
    return &UserRepositoryImpl{db: _db}
}

func (u *UserRepositoryImpl) Create(username string, email string, hashPassword string) error {
    query := `INSERT INTO users (username, email, password) VALUES (?, ?, ?)`
    _, err := u.db.Exec(query, username, email, hashPassword)
    return err
}

func (u *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
    query := `SELECT id, email, password FROM users WHERE email=?`
    row := u.db.QueryRow(query, email)

    user := &models.User{}
    err := row.Scan(&user.Id, &user.Email, &user.Password)
    if err != nil {
        return nil, err
    }
    return user, nil
}
```

---

## Dependency Injection & App Assembly

### Concept
The application entrypoint explicitly constructs and links each layer together (Database → Repository → Service → Controller → Router → HTTP Server), ensuring easy testing and mock injection.

### Code Implementation (`app/application.go`)
```go
package app

import (
    dbConfig "AuthInGo/config/db"
    config "AuthInGo/config/env"
    "AuthInGo/controllers"
    repo "AuthInGo/db/repositories"
    router "AuthInGo/routers"
    "AuthInGo/services"
    "fmt"
    "net/http"
    "time"
)

type Application struct {
    Config Config
}

func (app *Application) Run() error {
    // 1. Initialize Database Connection
    db, err := dbConfig.SetupDB()
    if err != nil {
        return err
    }

    // 2. Dependency Injection Chain
    userRepo := repo.NewUserRepository(db)
    userService := services.NewUserService(userRepo)
    userController := controllers.NewUserController(userService)
    userRouter := router.NewUserRouter(*userController)

    // 3. Configure HTTP Server
    server := &http.Server{
        Addr:         app.Config.Addr,
        Handler:      router.SetUpRouter(userRouter),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    fmt.Println("Server running at", app.Config.Addr)
    return server.ListenAndServe()
}
```

---

## API Endpoints Summary

| Method | Endpoint | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/ping` | No | Health check ping endpoint |
| `POST` | `/signup` | No | Register new user account |
| `POST` | `/login` | No | Authenticate user & return JWT token |
| `GET` | `/profile/{id}` | **Yes (Bearer JWT)** | Fetch user profile by ID |
| `GET` | `/profiles` | No | Fetch list of users |
| `ALL` | `/fakestoreservice/*` | No | Reverse proxy to `https://fakestoreapi.com` |

---

## Node.js & OOP Developer's Guide to Go & System Design Concepts

If you are coming from Node.js (TypeScript/ES6) and building backend microservices with Object-Oriented System Design principles, Go's design philosophy is intentionally simpler and favors **composition over inheritance**.

---

### 12.1 Classes & Object-Oriented Programming

Go **does not have a `class` keyword**. Instead, you define data structures using **Structs** and attach methods to them using **Receiver Functions**.

#### Node.js / TypeScript Example
```typescript
class UserServiceImpl {
    private userRepo: UserRepository;

    constructor(userRepo: UserRepository) {
        this.userRepo = userRepo;
    }

    async createUser(email: string): Promise<void> {
        await this.userRepo.create(email);
    }
}
```

#### Equivalent Go Example (`services/userService.go`)
```go
package services

import "AuthInGo/db/repositories"

// Struct holds state (like properties in a class)
type UserServiceImpl struct {
    userRepository db.UserRepository
}

// Constructor Function (Convention in Go)
func NewUserService(userRepo db.UserRepository) *UserServiceImpl {
    return &UserServiceImpl{
        userRepository: userRepo,
    }
}

// Receiver Method: Attached to UserServiceImpl struct
// (u *UserServiceImpl) is the receiver (equivalent to 'this' or 'self')
func (u *UserServiceImpl) CreateUser(email string) error {
    return u.userRepository.Create(email)
}
```

#### Composition (Embedding) Instead of Inheritance
Go does not support class inheritance (`extends`). System design in Go uses **Struct Embedding** (composition):

```go
type BaseEntity struct {
    ID        int64
    CreatedAt time.Time
}

type User struct {
    BaseEntity // Embedded struct (inherits fields ID and CreatedAt automatically)
    Email      string
}
```

---

### 12.2 Interfaces & Polymorphism

In Node.js/TypeScript, you explicitly implement interfaces using `implements`:
```typescript
interface UserRepository {
    create(email: string): Promise<void>;
}
class MysqlRepo implements UserRepository { ... }
```

In Go, **interfaces are implemented implicitly** (Duck Typing with full compile-time safety). A struct implements an interface simply by defining matching methods — **no `implements` keyword exists**.

#### Go Implementation (`db/repositories/users.go` & `services/userService.go`)

```go
// 1. Interface definition
type UserRepository interface {
    GetByEmail(email string) (*models.User, error)
}

// 2. Struct definition
type UserRepositoryImpl struct {
    db *sql.DB
}

// 3. Implicit implementation: Defining GetByEmail satisfies UserRepository automatically
func (u *UserRepositoryImpl) GetByEmail(email string) (*models.User, error) {
    // SQL query execution...
    return &models.User{}, nil
}
```

Because `UserRepositoryImpl` has a method matching `GetByEmail(string) (*models.User, error)`, it can be passed directly anywhere a `UserRepository` interface is expected!

---

### 12.3 Enums in Go

Go does not have a native `enum` keyword like TypeScript or Java. Enums are implemented using **Custom Types + `const` + `iota`**.

#### Node.js / TypeScript Example
```typescript
enum UserRole {
    ADMIN = "ADMIN",
    USER = "USER",
    GUEST = "GUEST"
}
```

#### Equivalent Go Example
```go
package models

type UserRole int

const (
    RoleGuest UserRole = iota // 0
    RoleUser                   // 1
    RoleAdmin                  // 2
)

// Stringifier method for printing human-readable enum values
func (r UserRole) String() string {
    return [...]string{"GUEST", "USER", "ADMIN"}[r]
}
```

---

### 12.4 Visibility: Public, Private, Protected

Go does **NOT** have `public`, `private`, or `protected` keywords. Visibility is strictly controlled by **Capitalization at the Package Level**:

| Visibility Level | Go Convention | Example | Accessible From |
| :--- | :--- | :--- | :--- |
| **Public / Exported** | Starts with **Uppercase** letter | `User`, `CreateUser()`, `Email` | Any package |
| **Private / Unexported** | Starts with **Lowercase** letter | `userRepository`, `hashPassword()` | **Same package only** |
| **Protected** | **Does Not Exist** | N/A | Use composition & package encapsulation |

#### Code Example
```go
package services

type UserServiceImpl struct {
    userRepository db.UserRepository // Private field (lowercase 'u'): cannot be accessed outside package 'services'
}

func (u *UserServiceImpl) CreateUser(...) error { // Public method (uppercase 'C'): accessible anywhere
    ...
}
```

---

### 12.5 Static Members & Methods

Go structs do not have `static` fields or methods. Instead, static functionality is achieved via **Package-Level Variables & Functions**.

#### Node.js / TypeScript Example
```typescript
class PasswordUtils {
    static hash(password: string): string { ... }
}
PasswordUtils.hash("secret");
```

#### Equivalent Go Example (`utils/auth.go`)
Package-level standalone functions serve as static utility methods:

```go
package utils

import "golang.org/x/crypto/bcrypt"

// Standalone function in 'utils' package
func HashPassword(plainPassword string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
    return string(hash), err
}
```
Called in other packages as: `utils.HashPassword("secret")`.

---

### 12.6 Final, Readonly, and Immutability

Go does **not** have `final` or `readonly` keywords for struct fields. Immutability in Go system design is achieved via three primary techniques:

1. **Primitive Constants (`const`)**:
   ```go
   const MaxConnections = 100
   const AppVersion = "1.0.0"
   ```
2. **Unexported Struct Fields + Exported Getters**:
   Hide the field from external packages by keeping it lowercase, and provide only a getter method (no setter):
   ```go
   type Config struct {
       addr string // Unexported: immutable from outside package 'app'
   }

   func (c Config) Addr() string { // Getter
       return c.addr
   }
   ```
3. **Value Receivers**:
   Method receivers passed by value `(u User)` receive a copy, preventing the method from mutating the caller's struct state.

---

### 12.7 How Middlewares Work in Go vs. Node.js (Express)

#### Express.js Approach (Callback Chain)
In Express, middleware receives `(req, res, next)` and calls `next()` to advance:
```javascript
app.use((req, res, next) => {
    console.log("Pre-processing");
    next(); // Pass to next handler
    console.log("Post-processing");
});
```

#### Go Approach (Higher-Order Handler Decoration)
In Go (`net/http`), HTTP handlers implement the `http.Handler` interface (`ServeHTTP(ResponseWriter, *Request)`). 

A Go middleware is a **function that takes an `http.Handler` and returns a new `http.Handler`**:

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // 1. PRE-PROCESSING (Runs before reaching controller)
        fmt.Println("Incoming request:", r.Method, r.URL.Path)

        // 2. PASS CONTROL TO NEXT HANDLER IN CHAIN
        next.ServeHTTP(w, r)

        // 3. POST-PROCESSING (Runs after controller returns response)
        fmt.Println("Request completed")
    })
}
```

In `routers/router.go` and `routers/userRouter.go`, Chi router mounts these decorators using `r.Use(...)` or `r.With(...)`:

```go
// Global middleware
chiRouter.Use(middlewares.RateLimitMiddleware)

// Endpoint-specific middleware guard
r.With(middlewares.JWTAuthMiddleware).Get("/profile/{id}", ur.UserController.GetUserById)
```


---


---

### 12.8 Dependency Inversion Principle (DIP) & Clean Architecture

The **Dependency Inversion Principle (DIP)** states:
1. **High-level modules** (e.g. Services, Controllers) should not depend on **low-level modules** (e.g. MySQL DB drivers, specific HTTP libraries). Both should depend on **abstractions** (Interfaces).
2. **Abstractions** should not depend on details. **Details** should depend on abstractions.

#### Dependency Architecture Chain

```
[Controller Layer] ──depends on──> [UserService Interface]
                                          ▲
                                     implements
                                          │
                                 [UserServiceImpl] ──depends on──> [UserRepository Interface]
                                                                                ▲
                                                                           implements
                                                                                │
                                                                       [UserRepositoryImpl]
```

#### Code Implementation Example

1. **Service Layer Abstraction (`services/userService.go`)**:
   ```go
   type UserService interface {
       CreateUser(*dtos.CreateUserRequestDTO) error
       GetUserById(string) (*models.User, error)
       GetAllUsers() (*models.User, error)
       LoginUser(*dtos.LoginUserRequestDTO) (string, error)
   }

   type UserServiceImpl struct {
       userRepository db.UserRepository // Depend on interface, NOT concrete MySQL struct
   }

   func NewUserService(userRepo db.UserRepository) *UserServiceImpl {
       return &UserServiceImpl{userRepository: userRepo}
   }
   ```

2. **Repository Layer Abstraction (`db/repositories/users.go`)**:
   ```go
   type UserRepository interface {
       Create(username, email, hashPassword string) error
       GetByEmail(email string) (*models.User, error)
       GetById() (*models.User, error)
       GetAll() error
   }

   type UserRepositoryImpl struct {
       db *sql.DB // Concrete database driver instance
   }

   func NewUserRepository(_db *sql.DB) *UserRepositoryImpl {
       return &UserRepositoryImpl{db: _db}
   }
   ```

3. **Wiring Layers via Constructor Injection (`app/application.go`)**:
   ```go
   func (app *Application) Run() error {
       db, err := dbConfig.SetupDB()

       // Construct dependencies bottom-up and inject interfaces
       userRepo := repo.NewUserRepository(db)                     // satisfies db.UserRepository
       userService := services.NewUserService(userRepo)            // satisfies services.UserService
       userController := controllers.NewUserController(userService)
       userRouter := router.NewUserRouter(*userController)

       // launch http.Server...
   }
   ```

#### Benefits for System Design & Testing
* **Decoupled Business Logic**: Core business rules in `UserServiceImpl` have zero direct knowledge of MySQL or HTTP details.
* **Easy Mocking & Unit Testing**: You can pass a `MockUserRepository` struct to `NewUserService()` during unit tests without running a live database server.
* **Database Swapability**: Switching from MySQL to PostgreSQL or MongoDB only requires creating a new struct that satisfies `db.UserRepository` — zero lines of service layer code need to change.

---


---

### 12.9 Starter Commands: Node.js vs. Go Comparison

When starting or managing projects, Go provides direct equivalents for standard Node.js (`npm` / `yarn` / `pnpm`) terminal commands.

#### Node.js vs. Go Command Cheat Sheet

| Task | Node.js Command | Go Equivalent Command | Explanation |
| :--- | :--- | :--- | :--- |
| **Initialize new project** | `npm init -y` | `go mod init <module-name>` | Initializes a new module, creating `package.json` in Node or `go.mod` in Go. Example: `go mod init github.com/username/my-go-backend` |
| **Install package dependency** | `npm install <package>` | `go get <package-path>` | Downloads package and records it in `go.mod` & `go.sum`. Example: `go get github.com/go-chi/chi/v5` |
| **Install specific version** | `npm install <pkg>@1.2.3` | `go get <pkg>@v1.2.3` | Installs a precise version tag (e.g. `go get github.com/golang-jwt/jwt/v5@v5.3.1`). |
| **Install existing repo dependencies** | `npm install` | `go mod download` | Downloads all dependencies listed in `go.mod` to your local module cache. |
| **Clean & sync dependencies** | `npm prune` | `go mod tidy` | Scans all `.go` files, adds missing imports to `go.mod`, and removes unused packages. |
| **Run development server** | `npm run dev` (via `nodemon`) | `air` or `go run main.go` | Starts application with hot-reloading (`air`) or runs source directly (`go run main.go`). |
| **Build production binary** | `npm run build` | `go build -o bin/app main.go` | Compiles source code into a single, zero-dependency executable binary file. |
| **Run unit tests** | `npm test` | `go test -v ./...` | Runs all tests in current directory and sub-packages (`./...`). |

#### Step-by-Step: Bootstrapping a New Go Project from Scratch

1. **Create project directory & initialize module**:
   ```bash
   mkdir my-go-service && cd my-go-service
   go mod init my-go-service
   ```
2. **Install core framework dependencies**:
   ```bash
   go get github.com/go-chi/chi/v5           # HTTP Router
   go get github.com/golang-jwt/jwt/v5        # JWT Auth
   go get golang.org/x/crypto                 # Bcrypt Hashing
   go get github.com/go-sql-driver/mysql      # MySQL Driver
   go get github.com/joho/godotenv            # .env Loader
   ```
3. **Synchronize dependencies**:
   ```bash
   go mod tidy
   ```
4. **Run development mode**:
   ```bash
   go run main.go
   ```


## External Dependencies & Libraries Breakdown

Below is a detailed breakdown of every external dependency defined in `go.mod`, explaining its core responsibility and role in this application:

| Package | Purpose / Description | Role in Project |
| :--- | :--- | :--- |
| **`github.com/go-chi/chi/v5`** | Fast, lightweight HTTP router for Go | Handles HTTP routing, path parameters (`/profile/{id}`), sub-routers, and middleware composition (`r.With()`). |
| **`github.com/golang-jwt/jwt/v5`** | JSON Web Token (JWT) implementation | Signs JWT access tokens on login (`LoginUser`) and parses/verifies `Bearer` token headers in `JWTAuthMiddleware`. |
| **`golang.org/x/crypto/bcrypt`** | Cryptographic password hashing | Hashes user passwords on signup (`HashPassword`) and validates credentials on login (`CheckPasswordHash`). |
| **`github.com/go-sql-driver/mysql`** | MySQL driver for Go `database/sql` | Formats DSN connection strings, opens TCP connection pools, and executes raw SQL queries against MySQL. |
| **`github.com/lib/pq`** | PostgreSQL driver for `database/sql` | PostgreSQL database driver enabling communication with PostgreSQL instances (such as local Docker containers). |
| **`github.com/joho/godotenv`** | Environment file loader (`.env`) | Loads `.env` file variables into Go's `os.Environ` at server startup. |
| **`golang.org/x/time/rate`** | Token-bucket rate limiter | Restricts excess API traffic (`RateLimitMiddleware`) to prevent rate abuse (e.g. max 5 requests per 30 seconds). |
| **`filippo.io/edwards25519`** | Elliptic curve cryptography (Indirect) | Transitive dependency pulled in by `go-sql-driver/mysql` to support MySQL authentication plugins (`caching_sha2_password`). |
