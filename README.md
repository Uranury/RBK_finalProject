# CS:GO Skin Marketplace API 🎮

A comprehensive marketplace API for buying and selling CS:GO skins, built with Go, featuring user authentication, transaction management, automated invoice generation, and email notifications.

## 🎯 Project Overview

This project is a full-featured marketplace API that simulates the CS:GO skin trading ecosystem. Users can register, deposit/withdraw funds, list skins for sale, purchase skins from other users, and receive automated invoice emails for their transactions.

### Key Features

- **User Management**: Registration, authentication, profile management
- **Skin Marketplace**: Create, list, buy, and sell CS:GO skins
- **Transaction System**: Deposit, withdraw, and track transaction history
- **Automated Invoicing**: PDF generation and email delivery for purchases
- **Background Processing**: Asynchronous task processing with Redis
- **RESTful API**: Complete API with Swagger documentation
- **Database Migrations**: Automated schema management
- **Docker Support**: Containerized deployment with optimized images

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Server    │    │   Worker        │    │   Database      │
│   (Gin)         │    │   (Asynq)       │    │   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Redis         │    │   Email Service │    │   Migrations    │
│   (Queue)       │    │   (Mailgun)     │    │   (Golang)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Tech Stack

### Backend
- **Language**: Go 1.23.1
- **Framework**: Gin (HTTP server)
- **Database**: PostgreSQL 16
- **Cache/Queue**: Redis 7
- **ORM**: SQLx
- **Authentication**: JWT
- **Background Jobs**: Asynq
- **Email Service**: Mailgun
- **PDF Generation**: gofpdf
- **API Documentation**: Swagger/OpenAPI

### DevOps & Tools
- **Containerization**: Docker & Docker Compose
- **Database Migrations**: golang-migrate
- **Configuration**: Environment variables
- **Logging**: Structured logging with slog
- **Testing**: Go testing with testify
- **Dependency Injection**: Wire

### Libraries & Dependencies
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/jmoiron/sqlx` - Database operations
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/hibiken/asynq` - Background job processing
- `github.com/mailgun/mailgun-go/v4` - Email service
- `github.com/jung-kurt/gofpdf` - PDF generation
- `github.com/golang-jwt/jwt/v4` - JWT authentication
- `github.com/google/uuid` - UUID generation
- `github.com/swaggo/gin-swagger` - API documentation
- `golang.org/x/crypto/bcrypt` - Password hashing

## 📁 Project Structure

```
finalProject/
├── cmd/                          # Application entry points
│   ├── api/                      # API server
│   │   ├── main.go              # API server entry point
│   │   └── wire.go              # Dependency injection
│   └── worker/                   # Background worker
│       ├── main.go              # Worker entry point
│       └── wire.go              # Worker dependencies
├── internal/                     # Internal application code
│   ├── auth/                    # Authentication service
│   ├── handlers/                # HTTP request handlers
│   ├── http_server/             # HTTP server setup
│   ├── middleware/              # HTTP middleware
│   ├── models/                  # Data models
│   ├── queue/                   # Background job processing
│   │   ├── handlers/            # Job handlers
│   │   └── jobs/                # Job definitions
│   ├── repositories/            # Data access layer
│   │   ├── order/               # Order repository
│   │   ├── skin/                # Skin repository
│   │   ├── transaction/         # Transaction repository
│   │   └── user/                # User repository
│   └── services/                # Business logic layer
│       ├── email_service.go     # Email service
│       ├── invoice_service.go   # PDF generation
│       ├── marketplace_service.go # Marketplace logic
│       ├── skin_service.go      # Skin management
│       ├── transaction_service.go # Transaction logic
│       └── user_service.go      # User management
├── pkg/                         # Shared packages
│   ├── apperrors/               # Error handling
│   ├── config/                  # Configuration management
│   └── db/                      # Database connection
├── migrations/                  # Database migrations
├── docs/                        # Swagger documentation
├── Dockerfile*                  # Docker configurations
├── docker-compose.yml           # Service orchestration
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23.1+ (for local development)
- Make (optional, for build scripts)

### 1. Clone the Repository
```bash
git clone <repository-url>
cd finalProject
```

### 2. Setup Environment
```bash
# Create environment file
cp .env.example .env

# Edit the environment file
nano .env
```

### 3. Build and Run
```bash
# Build and start all services
docker-compose up --build

# Or run in background
docker-compose up -d --build
```

### 4. Access the Application
- **API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Database**: localhost:5436 (PostgreSQL)
- **Redis**: localhost:6379

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `LISTEN_ADDR` | HTTP server address | `:8080` | No |
| `REDIS_ADDR` | Redis connection string | `redis:6379` | No |
| `DB_URL` | PostgreSQL connection string | `postgres://postgres:postgres@db:5432/postgres?sslmode=disable` | No |
| `MIGRATIONS_PATH` | Path to database migrations | `/app/migrations` | No |
| `JWT_SECRET` | JWT signing secret | - | **Yes** |
| `MAILGUN_DOMAIN` | Mailgun domain for emails | - | No |
| `MAILGUN_API_KEY` | Mailgun API key | - | No |
| `POSTGRES_DB` | PostgreSQL database name | `postgres` | No |
| `POSTGRES_USER` | PostgreSQL username | `postgres` | No |
| `POSTGRES_PASSWORD` | PostgreSQL password | `postgres` | No |
| `POSTGRES_PORT` | PostgreSQL port | `5432` | No |

### Example .env File
```env
# Application Configuration
LISTEN_ADDR=:8080
REDIS_ADDR=redis:6379

# Database Configuration
DB_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
MIGRATIONS_PATH=/app/migrations

# PostgreSQL Configuration
POSTGRES_DB=postgres
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_PORT=5432

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here

# Mailgun Configuration (optional)
MAILGUN_DOMAIN=your-mailgun-domain.com
MAILGUN_API_KEY=your-mailgun-api-key
```

## 📚 API Documentation

The API is fully documented with Swagger/OpenAPI. Once the application is running, visit:

**http://localhost:8080/swagger/index.html**

### Key Endpoints

#### Authentication
- `POST /signup` - Register a new user
- `POST /login` - Authenticate user
- `GET /profile` - Get user profile (authenticated)

#### Marketplace
- `GET /marketplace/skins` - List available skins
- `GET /marketplace/skins/mine` - Get user's skins (authenticated)
- `POST /marketplace/purchase` - Purchase a skin (authenticated)
- `POST /marketplace/sell` - List a skin for sale (authenticated)
- `DELETE /marketplace/skins/{skin_id}` - Remove skin from listing (authenticated)

#### Transactions
- `POST /transactions/deposit` - Deposit funds (authenticated)
- `POST /transactions/withdraw` - Withdraw funds (authenticated)
- `GET /transactions/history` - Get transaction history (authenticated)

#### Skins
- `POST /skins` - Create a new skin (authenticated)
- `GET /guns` - Get available guns
- `GET /wears` - Get available wear levels

## 🧪 Testing

### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./internal/services -v

```

### Test Coverage
The project includes comprehensive unit tests for:
- User service (authentication, registration, profile management)
- Skin service (creation, validation, retrieval)
- Email service (configuration and structure)
- Business logic validation

### Test Structure
```
internal/services/
├── user_service_test.go      # User service tests
├── skin_service_test.go      # Skin service tests
└── email_service_test.go     # Email service tests
```

## 🐳 Docker Optimization

The project includes optimized Docker images with:

- **Multi-stage builds** for smaller production images
- **Build caching** for faster rebuilds
- **Selective file copying** to minimize image size
- **Security hardening** with non-root users
- **Alpine-based images** for smaller footprint

### Image Sizes
- **API Image**: ~20MB (80% reduction from original)
- **Worker Image**: ~15MB (85% reduction from original)
- **Database Images**: ~50% reduction with Alpine variants

## 🔧 Development

### Local Development Setup
```bash
# Install dependencies
go mod download

# Run migrations
make migrate-up

# Start services
docker-compose up db redis

# Run API server
go run cmd/api/main.go

# Run worker
go run cmd/worker/main.go
```

### Database Migrations
```bash
# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down

# Create new migration
make migrate-create name=migration_name
```

### Code Generation
```bash

# Generate Swagger documentation
make swagger
```

## 📊 Performance

### Optimizations Applied
- **Database Connection Pooling**: Optimized connection settings
- **Redis Caching**: Background job queue and caching
- **Asynchronous Processing**: Email sending and PDF generation
- **Optimized Docker Images**: Multi-stage builds and Alpine base
- **Structured Logging**: Efficient logging with slog

### Monitoring
- **Health Checks**: Database and Redis health monitoring
- **Structured Logging**: JSON-formatted logs for easy parsing
- **Error Tracking**: Comprehensive error handling and logging

## 🔒 Security

### Security Features
- **JWT Authentication**: Secure token-based authentication
- **Password Hashing**: bcrypt for password security
- **Input Validation**: Comprehensive request validation
- **SQL Injection Protection**: Parameterized queries with SQLx
- **Non-root Containers**: Docker security hardening
- **Environment Variables**: Secure configuration management

### Best Practices
- Input sanitization and validation
- Proper error handling without information leakage
- Secure password storage with bcrypt
- JWT token expiration and validation
- Database transaction management

## 🤝 Contributing

### Development Workflow
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Code Style
- Follow Go conventions and idioms
- Use `gofmt` for code formatting
- Write comprehensive tests
- Add comments for complex logic
- Use meaningful variable and function names

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

### Common Issues

#### Email Not Sending
- Check Mailgun credentials in environment variables
- Verify domain configuration
- Check worker logs for email processing errors

#### Database Connection Issues
- Ensure PostgreSQL is running
- Check database credentials
- Verify network connectivity between services

### Getting Help
- Check the API documentation at `/swagger/index.html`
- Review application logs: `docker-compose logs`
- Check worker logs: `docker-compose logs worker`
- Verify environment configuration

### Technical Improvements
- **GraphQL API**: Alternative to REST endpoints
- **Microservices**: Split into smaller, focused services
- **Event Sourcing**: Complete audit trail of all changes
- **Caching Layer**: Redis caching for frequently accessed data
- **Load Balancing**: Horizontal scaling support
- **Monitoring**: Prometheus metrics and Grafana dashboards

---

**Built with ❤️ using Go and modern DevOps practices**
