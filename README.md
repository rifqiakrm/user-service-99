# User Service

This service provides a simple User Management API as part of a microservices architecture. It exposes endpoints to create users and retrieve users (single or paginated).

Built using Golang, Gin, and GORM (with SQLite for persistence).

---

## 🔧 Tech Stack

* **Go 1.22+**
* **Gin** – HTTP Router
* **GORM** – ORM for database operations
* **SQLite** – Embedded database
* **Testify** – Unit testing framework
* **GoMock** – Interface mocking for tests

---

## 📁 Project Structure

```
user-service/
├── db/                         # Database
│   └── database.go     
│   └── database_test.go     
├── handler/                    # HTTP handlers
│   └── user_handler.go     
│   └── user_handler_test.go     
├── model/                      # Domain models
│   └── user.go             
├── repository/                 # Database layer
│   └── user_repo.go        
│   └── user_repo_test.go      
├── service/                    # Business logic
│   └── user_service.go     
│   └── user_service_test.go       
├── mocks/                      # Generated mocks for testing
├── main.go                     # App entry point
├── go.mod
└── go.sum
```

---

## 🚀 Running the Service

```bash
go run ./cmd/main.go
```

By default, service listens on port `:8080`.

---

## 🧪 Running Tests

```bash
go test ./...
```

---

## 🧰 Generate Mocks

Mocks are generated using `go generate` and `mockgen` from GoMock.

```bash
go generate ./...
```

Make sure `mockgen` is installed:

```bash
go install github.com/golang/mock/mockgen@v1.6.0
```

---

## 📌 API Endpoints

| Method | Endpoint       | Description               |
|--------|----------------|---------------------------|
| POST   | `/users`       | Create a new user         |
| POST   | `/users/batch` | Get users by IDs          |
| GET    | `/users/:id`   | Get user by ID            |
| GET    | `/users`       | Get all users (paginated) |

### Example: Create User

```bash
curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"name":"John"}'
```

### Example: Get Users By IDs

```bash
curl --location 'localhost:6001/users/batch' \
--header 'Content-Type: application/json' \
--data '{
    "user_ids":[
        1,
        8
    ]
}'
```

### Example: Get User

```bash
curl http://localhost:8080/users/1
```

### Example: Get All Users (page=1, size=10)

```bash
curl http://localhost:8080/users?page=1&size=10
```

---

## ✅ Features

* Clean architecture separation
* Table-driven unit tests
* Mock-based service and repo testing
* SQLite in-memory support for test isolation