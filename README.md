# 🏗️ Scalable eCommerce Platform (Go + Chi + SQLC + PGX + Stripe)

A **high-performance, modular, and scalable eCommerce backend** built with **Golang**, following industry best practices and inspired by the [golang-standards/project-layout](https://github.com/golang-standards/project-layout).  
The project is designed for **real-world scalability**, **robust domain isolation**, and **seamless third-party integrations** — including full **Stripe payment processing**.

---

## 🚀 Overview

This backend serves as the foundation for a **production-grade eCommerce system**, providing:

- Domain-driven modular architecture
- Optimized PostgreSQL interactions using `pgx` + `sqlc`
- Secure Stripe integration for real payments
- Seamless migration handling with `golang-migrate`
- Clean routing and middleware powered by `chi`
- Extensible, testable, and horizontally scalable design

---

## 🧱 Key Technologies

| Layer                 | Technology                                                  | Purpose                                                 |
| --------------------- | ----------------------------------------------------------- | ------------------------------------------------------- |
| **Router**            | [Chi](https://github.com/go-chi/chi)                        | Lightweight, idiomatic HTTP routing                     |
| **Database**          | [PGX](https://github.com/jackc/pgx)                         | High-performance PostgreSQL driver                      |
| **Query Generation**  | [SQLC](https://sqlc.dev)                                    | Type-safe SQL query to Go code generator                |
| **Migration**         | [golang-migrate](https://github.com/golang-migrate/migrate) | Database schema migration                               |
| **Validation**        | [go-playground](https://github.com/go-playground/validator) | Request payload validation                              |
| **Payments**          | [Stripe](https://stripe.com/docs/api)                       | Secure payment integration                              |
| **Config Management** | `.env` + `configs` package                                  | Environment-based configuration setup                   |
| **Architecture**      | Modular (Domain-Driven)                                     | Each domain encapsulates handler, service, repo, routes |

---

## 🗂️ Project Structure

```bash
├── cmd/api
│   ├── main.go              # Application entry point
│   └── bind.go              # Server bootstrap and dependency binding
│
├── configs
│   ├── configs.go           # Configuration loader
│   ├── keys.go              # Env key constants
│   └── types.go             # Config schema
│
├── internal
│   ├── app/api/router
│   │   └── router.go        # Chi route setup
│   ├── domain
│   │   ├── product/         # Product module (handler, service, repo, routes)
│   │   ├── order/           # Order module
│   │   ├── user/            # User module
│   │   └── payment/         # Stripe integration module
│   │       ├── handler.go   # Webhook and payment endpoints
│   │       ├── service.go   # Stripe payment orchestration logic
│   │       ├── repository.go # Payment transaction persistence
│   │       └── routes.go    # Route registration
│   │       └── types.go     # Entity structs
│   │       └── dto.go       # Request payload (dto) structs
│   │   └── stripe/          # Stripe payment integration
│   │       ├── client.go    # Stripe client: handles webhooks, payment intents, refunds
│   │       └── types.go     # Strongly typed Stripe-related structs
│   │   - (All the rest domains)
│   └── infra/db
│       └── postgres.go      # Database connection via pgxpool
│
├── pkg
│   ├── jwt/                 # JWT token handling
│   ├── password/            # Password hashing utilities
│   ├── logger/              # Structured logging
│   ├── middleware/          # Role based access control middleware
│   ├── response/            # Unified API response format
│   ├── validator/           # Input validation
│   ├── pagination/          # Pagination helpers
│   ├── idgen/               # UUID/random ID generation
│   ├── httputil/            # Common HTTP helpers
│   └── uploader/            # File upload utility
│
├── scripts/
│   ├── sql/
│   │   ├── migrations/      # SQL migration files for golang-migrate
│   │   └── queries/         # SQLC query definitions
│
├── docker-compose.yaml       # Development services setup
├── .env                      # Environment configuration
└── .gitignore
```

💳 Stripe Payment Integration

The payment domain encapsulates all Stripe-related functionality.

Features

- Create payment intents for checkout

- Handle Stripe webhooks for confirmation events

- Store and update transaction records

- Graceful failure handling with rollback logic

- Secure secret key management via .env

🧩 Architectural Principles

- Modular Domains — Each feature area (product, order, payment, user, etc.) is fully self-contained.

- No Global Dependencies — Dependency injection ensures testability and scalability.

- SQLC + PGX — Compile-time validated SQL and optimized DB access.

- High Testability — Each domain can be tested independently.

Environment Isolation — .env-driven configuration for local, staging, and production environments.

⚙️ Setup & Installation

1️⃣ Prerequisites

- Go 1.22+

- PostgreSQL 15+

- Docker & Docker Compose

- Stripe account and API keys

- golang-migrate CLI

2️⃣ Clone the repository

```bash
git clone https://github.com/natnael772/ecommerce-golang.git
cd ecommerce-golang
```

3️⃣ Configure Environment
Copy and update your environment file:

```bash
cp .env.example .env
```

4️⃣ Install Go Dependencies

Before running any Go code, install all required packages:

```bash
go mod tidy
```

5️⃣ Run Services

```bash
docker compose up -d
```

6️⃣ Run Database Migrations

```bash
migrate -path scripts/sql/migrations -database "postgres://user:password@localhost:5432/ecommerce?sslmode=disable" up
```

7️⃣ Generate SQLC Code
Generate Go code from SQL queries:

```bash
sqlc generate
```

8️⃣ Start the API

```bash
go run cmd/api/main.go
```

🧠 Development Philosophy

This project embodies:

- Performance-first design — Leveraging PGX and SQLC.

- Domain encapsulation — Modules own their data, logic, and interfaces.

- Security by design — JWT, password hashing, Stripe webhook validation.

- Scalability — Horizontally extendable via independent domain modules.

- Clean abstractions — Repository, Service, Handler pattern.

🧪 Testing

Run all unit and integration tests:

```bash
go test ./... -v
```

Each domain is independently testable — enabling isolated business logic validation.

🐳 Docker Support

To run the entire system in Docker:

```bash
docker compose up --build
```

This spins up:

- The API service

- PostgreSQL

- Optional local Stripe CLI webhook listener (recommended for dev)

📜 License

This project is licensed under the MIT License.
