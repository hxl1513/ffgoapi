# Examples

The `examples/` directory contains a complete, runnable example that demonstrates the framework's key concepts in the context of a user resource with full CRUD operations.

## User CRUD Example

**Location:** `examples/users/`

### What It Demonstrates

| Concept | File | Description |
|---|---|---|
| Data model | `user_model.go` | `User` struct embedding `gorm.Model` for database persistence. |
| Repository pattern | `user_repository.go` | Database operations (Create, FindByID, FindAll, Update, Delete) using GORM. |
| Service layer | `user_service.go` | Business logic wrapping the repository; exposes a `BaseService` for transactions and caching. |
| Controller layer | `user_controller.go` | HTTP handlers using `BaseController` helpers (`Success`, `BadRequest`, `NotFound`, `Created`). |
| DTOs | `user_dto.go` | Request and response types with validation tags. |

### Running the Example

```bash
cd server/examples/users
go run .
```

The server starts on `http://localhost:3000`. All endpoints are prefixed under `/users`.

### Available Endpoints

| Method | Path | Description |
|---|---|---|
| `POST` | `/users` | Create a new user |
| `GET` | `/users/:id` | Retrieve a user by ID |
| `GET` | `/users` | List all users |
| `PUT` | `/users/:id` | Update an existing user |
| `DELETE` | `/users/:id` | Delete a user |

### Example Requests

**Create a user:**

```bash
curl -X POST http://localhost:3000/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'
```

**List all users:**

```bash
curl http://localhost:3000/users
```

**Get a single user:**

```bash
curl http://localhost:3000/users/1
```

### Environment Variables

The example uses SQLite by default — no external services are required. To switch to MySQL or PostgreSQL, set the following environment variables before starting the example:

```bash
# MySQL example
export DB_TYPE=mysql
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=
export DB_NAME=myapp

# PostgreSQL example
export DB_TYPE=postgres
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=
export DB_NAME=myapp
```

### Adding Your Own Example

To create a new example resource:

1. Create a subdirectory under `server/examples/`.
2. Define a `User`-style struct in the `model` package.
3. Implement a `Repository` that wraps GORM calls.
4. Implement a `Service` that embeds `BaseService` and holds the repository.
5. Implement a `Controller` that embeds `BaseController`.
6. Wire everything up in `main.go` with `srv.RegisterRoutes()`.

Refer to the framework design document for details on the expected layered structure.
