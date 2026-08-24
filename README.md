# go-rest

A lightweight Go application framework providing reusable infrastructure for building REST APIs and backend services.

## ⚙️ Installation

**go-rest requires Go 1.25 or higher.**

If you need to install or upgrade Go, visit the [official Go download page](https://go.dev/dl/).

Create a new directory for your project and initialize it with Go Modules:

```bash
mkdir my-rest-api
cd my-rest-api

go mod init github.com/your/repo
```

Then install `go-rest:

```bash
go get github.com/elsyahtech/go-rest
```

For more information about Go Modules, see the [Using Go Modules](https://go.dev/blog/using-go-modules) documentation.

---

## 🚀 Quickstart

The easiest way to get started is to use the `golang-rest-api-starter` repository.

```bash
git clone github.com/elmansyah/golang-rest-api-starter
```

The starter provides a working example of an application built with `go-rest`, including:

* Configuration
* Database connection
* Database migration
* Handler
* Library
* Entity
* Model
* Module
* Filter / Authentication
* View / Response

The starter is intended as a **reference implementation and learning example**.

Developers are free to modify the application layer according to their own requirements.

---

## 📁 Project Structure

A typical application built with `go-rest` can be organized as follows:

```text
.
├── app/
│   ├── config/
│   │   ├── app.go
│   │   ├── config.go
│   │   ├── cookies.go
│   │   ├── database.go
│   │   ├── filter.go
│   │   ├── migration.go
│   │   ├── module.go
│   │   ├── router.go
│   │   ├── server.go
│   │   └── token.go
│   │
│   ├── database/
│   │   ├── migrations
│   │   └── seeds
│   │
│   ├── entities/
│   ├── filters/
│   ├── handlers/
│   ├── helpers/
│   ├── libraries/
│   ├── models/
│   ├── modules/
│   └── views/
│
├── public/
│   └── upload/
│
├── test/
│
├── writable/
│   ├── info.log
│   └── error.log
│
└── main.go
```

---

## 🧩 Application Architecture

The application layer is intentionally kept simple:

```text
HTTP Request
     │
     ▼
  Handler
     │
     ▼
  Library
     │
     ▼
   Model
     │
     ▼
  Database
```

Each layer has a specific responsibility.

### Handler

Responsible for HTTP-related operations:

* Parse HTTP requests
* Validate request format
* Call the Library
* Transform data into an HTTP response
* Return HTTP status codes and JSON responses

### Library

Responsible for application and business logic:

* Business rules
* Data processing
* Validation of business conditions
* Data transformation
* Combining data from multiple models
* Calling external services
* Application-specific workflows

### Model

Responsible for database operations:

* SQL queries
* MongoDB queries
* Database-specific operations
* Reading and writing database entities

### Entity

Responsible for application data structures shared between layers.

### Module

Responsible for registering HTTP routes and connecting handlers to endpoints.

### Filter

Responsible for request-level filtering such as:

* Authentication
* Authorization
* JWT
* OAuth
* Other application-specific middleware

### View

Responsible for transforming internal entities into API responses.

### Helper

Contains reusable application-level helper functions.

---

## 🛢️ Database Support

`go-rest` provides infrastructure for multiple database drivers.

Currently supported database types include:

* MySQL
* PostgreSQL
* MSSQL
* SQLite
* MongoDB

The framework provides the database infrastructure and connection lifecycle.

Database-specific queries remain in the application Model layer.

For example:

```text
./app/models/
├── mysql/
├── postgresql/
├── mssql/
├── sqlite/
└── mongodb/
```

This allows developers to use database-specific query capabilities without changing the application Handler or Library flow.

---

## 🔄 Multi-Database Architecture

The same application Entity can be used across different database implementations.

For example:

```go
type User struct {
    ID        string     `gorm:"primaryKey;column:id" bson:"_id,omitempty"`
    Email     string     `gorm:"column:email" bson:"email"`
    Password  string     `gorm:"column:password" bson:"password"`
    IsActive  bool       `gorm:"column:is_active" bson:"is_active"`
    IsDeleted bool       `gorm:"column:is_deleted" bson:"is_deleted"`
    CreatedAt time.Time  `gorm:"column:created_at" bson:"created_at"`
    UpdatedAt *time.Time `gorm:"column:updated_at" bson:"updated_at"`
    DeletedAt *time.Time `gorm:"column:deleted_at" bson:"deleted_at"`
}
```

SQL-based databases can use SQL queries:

```go
rows, err := db.QuerySQL(ctx, query)
```

MongoDB can use MongoDB-native operations:

```go
collection := mongoConn.Collection("users")

cursor, err := collection.Find(ctx, bson.M{})
```

The framework does not force developers to use the same query implementation for every database.

Use the query model that is appropriate for the selected database.

---

## 🗃️ Database Migration

`go-rest` provides database migration infrastructure.

Migrations are executed during application startup according to the configured database driver.

The framework supports database-specific migration implementations.

For SQL databases, migrations can use SQL.

For MongoDB, migrations can use native Go migration implementations.

Example:

```go
func GetMongoDBMigrations() []sysconfig.MongoDBMigration {
    return []sysconfig.MongoDBMigration{
        CreateUsersCollectionMigration{},
        CreateProductsCollectionMigration{},
    }
}
```

Developers only need to register their migrations.

The framework handles migration execution and migration history.

---

## 🔐 Authentication

Routes can use authentication filters such as JWT or OAuth depending on the application's configuration.

Example:

```go
userGroup.Get(
    "/",
    filters.JWTAuth(),
    handler.GetAll,
)
```

Authentication runs before the Handler:

```text
HTTP Request
     │
     ▼
JWT / OAuth Filter
     │
     ├── Authentication failed
     │        ↓
     │      STOP
     │
     └── Authentication successful
              ↓
           Handler
```

Public routes can omit authentication middleware when appropriate:

```go
userGroup.Get("/:id", handler.GetByID)
userGroup.Post("/", handler.Create)
```

Developers can add additional application-specific filters inside `./app`.

---

## 🛠️ Developer Guide

The `./app` directory is the application development area.

Developers are free to:

* Add handlers
* Add libraries
* Add models
* Add entities
* Add modules
* Add routes
* Add filters
* Add helpers
* Customize views
* Implement business logic
* Implement database-specific queries

The core infrastructure is provided by `go-rest`.

```text
go-rest
   │
   ├── Core infrastructure
   │
   ▼
Application
   │
   └── ./app
       ├── handlers
       ├── libraries
       ├── models
       ├── entities
       ├── modules
       ├── filters
       ├── helpers
       └── views
```

The application layer should be adapted to the requirements of each project.

---

## 🧠 Design Philosophy

`go-rest` aims to provide a simple and predictable foundation for Go applications.

The framework handles common infrastructure concerns while leaving application-specific decisions to the developer.

### Core

The framework provides:

* Application lifecycle
* Configuration
* Database connections
* Database migrations
* Logging
* Routing infrastructure
* Security infrastructure
* Server lifecycle

### Application

The developer provides:

* Business logic
* HTTP handlers
* Database queries
* Application models
* Routes
* Application-specific filters
* Application-specific helpers
* API responses

In simple terms:

```text
go-rest
    ↓
Core infrastructure

./app
    ↓
Your application
```

---

## 🚀 Application Lifecycle

A typical application startup flow is:

```text
main()
   │
   ▼
AppRun()
   │
   ├── Load configuration
   ├── Load timezone
   ├── Run testing / application mode
   ├── Initialize logger
   ├── Connect database
   ├── Run database migrations
   ├── Initialize router
   └── Start HTTP server
```

The application entry point can remain intentionally simple:

```go
func main() {
    if err := bootstrap.AppRun(); err != nil {
        log.Printf("application stopped: %v", err)
        os.Exit(1)
    }
}
```

Application initialization is handled by the framework bootstrap process.

---

## 📚 Example Project

For a complete working example, see:

```text
github.com/elmansyah/golang-rest-api-starter
```

The starter demonstrates how to build an application using `go-rest`, including:

```text
Handler
   ↓
Library
   ↓
Model
   ↓
Database
```

It is designed to be both a starting point and a learning reference for developers.

---

## 🧪 Testing

Testing support is part of the application development workflow.

A recommended test structure is:

```text
test/
├── unit/
├── integration/
└── ...
```

Use mocks or other test doubles when testing application components independently from real infrastructure.

---

## 📦 Requirements

* Go 1.25+
* A supported database driver when database access is required

Supported databases:

* MySQL
* PostgreSQL
* MSSQL
* SQLite
* MongoDB

---

## 🤝 Contributing

Contributions are welcome.

Before contributing:

1. Read the project documentation.
2. Understand the separation between core infrastructure and application code.
3. Keep changes focused.
4. Add or update tests where appropriate.
5. Make sure the project builds successfully.
6. Make sure existing functionality is not unintentionally broken.

---

## 📄 License

This project is licensed under the MIT License.

See the `LICENSE` file for details.

---

## 👤 Author

**Elsyah Technology Indonesia**

Repository:

```text
github.com/elsyahtech/go-rest
```

Starter:

```text
github.com/elmansyah/golang-rest-api-starter
```
