# CS:GO Skin Marketplace API 🎮

[![Go Version](https://img.shields.io/badge/Go-1.23.1+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue.svg)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-7-red.svg)](https://redis.io)

[![GitHub](https://img.shields.io/badge/GitHub-Profile-black?style=for-the-badge&logo=github)](https://github.com/Uranury)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-blue?style=for-the-badge&logo=linkedin)](https://linkedin.com/in/alibi-ulanuly-37700330b)

A comprehensive marketplace API for buying and selling CS:GO skins, built with Go, featuring user authentication, transaction management, automated invoice generation, and email notifications.

## 🎯 Quick Start

```bash
# Clone and setup
git clone <repository-url>
cd finalProject

# Minimal setup (only 2 required variables!)
cp env.example .env
echo "JWT_SECRET=your-secret-key-here" >> .env
echo "MIGRATIONS_PATH=./migrations" >> .env

# Run with Docker
docker-compose up --build
```

**Access:** http://localhost:8080 | **API Docs:** http://localhost:8080/swagger/index.html

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Server    │    │   Worker        │    │   Database      │
│   (Gin)         │    │   (Asynq)       │    │   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Redis         │    │   Email Service │    │   Migrations    │
│   (Queue)       │    │   (Mailgun)     │    │   (Golang)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Tech Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.23.1 |
| **Framework** | Gin (HTTP server) |
| **Database** | PostgreSQL 16 |
| **Cache/Queue** | Redis 7 |
| **ORM** | SQLx |
| **Auth** | JWT |
| **Background Jobs** | Asynq |
| **Email** | Mailgun |
| **PDF** | gofpdf |
| **Docs** | Swagger/OpenAPI |
| **Containerization** | Docker & Docker Compose |

## 📁 Project Structure

```
finalProject/
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   └── worker/            # Background worker
├── internal/              # Application code
│   ├── auth/              # Authentication
│   ├── handlers/          # HTTP handlers
│   ├── models/            # Data models
│   ├── repositories/      # Data access layer
│   ├── services/          # Business logic
│   └── queue/             # Background jobs
├── pkg/                   # Shared packages
├── migrations/            # Database migrations
├── docs/                  # Swagger documentation
└── Dockerfile*            # Container configs
```

## ⚙️ Configuration

### Required Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing secret | `your-secret-key-here` |
| `MIGRATIONS_PATH` | Path to migrations | `./migrations` |

### Optional Variables (with defaults)

| Variable | Default | Description |
|----------|---------|-------------|
| `LISTEN_ADDR` | `:8080` | HTTP server address |
| `REDIS_ADDR` | `:6379` | Redis connection |
| `DB_URL` | `postgres://postgres:postgres@db:5432/postgres?sslmode=disable` | Database URL |
| `MAILGUN_DOMAIN` | - | Email domain (optional) |
| `MAILGUN_API_KEY` | - | Email API key (optional) |

## 📚 API Endpoints

**Full Documentation:** http://localhost:8080/swagger/index.html

### Key Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/signup` | Register user |
| `POST` | `/login` | Authenticate user |
| `GET` | `/profile` | Get user profile |
| `GET` | `/marketplace/skins` | List available skins |
| `POST` | `/marketplace/purchase` | Purchase skin |
| `POST` | `/marketplace/sell` | List skin for sale |
| `POST` | `/transactions/deposit` | Deposit funds |
| `POST` | `/transactions/withdraw` | Withdraw funds |

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific service tests
go test ./internal/services -v
```

## 🐳 Docker Commands

```bash
# Build and start
docker-compose up --build

# Start in background
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

## 🔧 Development

```bash
# Setup development environment
make dev-setup

# Run locally (requires PostgreSQL & Redis)
make run

# Database migrations
make migrate-up
make migrate-down

# Code generation
make wire
make swagger
```

## 🔒 Security Features

- **JWT Authentication** with secure token management
- **Password Hashing** using bcrypt
- **Input Validation** and sanitization
- **SQL Injection Protection** with parameterized queries
- **Non-root Containers** for security hardening
- **Environment Variables** for secure configuration

## 📊 Performance Optimizations

- **Multi-stage Docker builds** (80% image size reduction)
- **Database connection pooling**
- **Redis caching** and background job queue
- **Asynchronous processing** for emails and PDFs
- **Structured logging** with slog
- **Alpine-based images** for smaller footprint

## 🆘 Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| **Email not sending** | Check Mailgun credentials in `.env` |
| **Database connection** | Verify PostgreSQL is running |
| **Build issues** | Clear Docker cache: `docker system prune -a` |

### Getting Help

- **API Documentation:** http://localhost:8080/swagger/index.html
- **Application Logs:** `docker-compose logs`
- **Worker Logs:** `docker-compose logs worker`

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License.

---

**Built with joy using Go (while being sick) and modern DevOps practices**

[![GitHub](https://img.shields.io/badge/GitHub-Profile-black?style=for-the-badge&logo=github)](https://github.com/Uranury)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-blue?style=for-the-badge&logo=linkedin)](https://linkedin.com/in/alibi-ulanuly-37700330b)
