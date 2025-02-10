# E-Commerce Microservices Platform

A modern e-commerce platform built with microservices architecture using Go, PostgreSQL, and Docker.

## Services Overview

### 1. Authentication Service
- User registration and authentication
- JWT token management
- User profile management

**Endpoints:**
- POST `/api/v1/auth/signup` - Register new user
- POST `/api/v1/auth/signin` - User login
- POST `/api/v1/auth/refresh` - Refresh access token
- GET `/api/v1/me` - Get user profile

### 2. Product Service
- Product catalog management
- Product search and filtering
- Product category management

**Endpoints:**
- GET `/api/v1/product` - List all products
- GET `/api/v1/product/{id}` - Get product details

### 3. Order Service
- Order processing
- Order history
- Payment integration

**Endpoints:**
- POST `/api/v1/orders/checkout` - Create new order

### 4. Shop Service
- Shop management
- Shop analytics
- Shop settings

**Endpoints:**
- POST `/api/v1/shop` - Create new shop
- GET `/api/v1/shop/{id}/warehouses` - Get shop warehouses

### 5. Warehouse Service
- Inventory management
- Stock transfer
- Warehouse operations

**Endpoints:**
- POST `/api/v1/warehouse` - Create warehouse
- PATCH `/api/v1/warehouse/{id}/status` - Update warehouse status
- POST `/api/v1/warehouse/transfer` - Transfer stock
- GET `/api/v1/warehouse/{id}/stock` - Get warehouse stock

## Technology Stack

### Core Technologies
- **Go (1.21+)** - Main programming language
- **PostgreSQL (14+)** - Primary database
- **Docker** - Containerization
- **Docker Compose** - Container orchestration
- **KrakenD** - API Gateway
- **JWT** - Authentication
- **OpenAPI/Swagger** - API documentation
- **Testify** - Testing framework

### Go Libraries
- **Echo** - Web framework
- **GORM** - ORM
- **go-migrate** - Database migrations
- **testify** - Testing
- **zap** - Logging
- **viper** - Configuration
- **validator** - Request validation

## Project Structure

```
.
├── auth-service/
│   ├── cmd/
│   ├── internal/
│   │   ├── controller/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
│   ├── migrations/
│   ├── tests/
│   └── Dockerfile
├── product-service/
│   ├── [similar structure]
├── order-service/
│   ├── [similar structure]
├── shop-service/
│   ├── [similar structure]
├── warehouse-service/
│   ├── [similar structure]
├── deploy/
│   ├── docker-compose.yml
│   └── krakend/
├── docs/
│   ├── openapi.json
│   └── diagrams/
├── scripts/
│   └── migrations/
├── Makefile
└── README.md
```

## Database Schema

Each service has its own database with the following key tables:

### Auth Service
```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Product Service
```sql
-- Products table
CREATE TABLE products (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Shop Service
```sql
-- Shops table
CREATE TABLE shops (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

### Warehouse Service
```sql
-- Warehouses table
CREATE TABLE warehouses (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    shop_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Stock table
CREATE TABLE stocks (
    id UUID PRIMARY KEY,
    warehouse_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

## Getting Started

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Make
- PostgreSQL 14+ (for local development)

### Setup
1. Clone the repository:
```bash
git clone https://github.com/your-org/ecommerce-platform.git
cd ecommerce-platform
```

2. Start the services:
```bash
make up
```

3. Run migrations:
```bash
make migrate-up
```

4. Seed the database:
```bash
make seed
```

### Available Make Commands
```
# Start all services
up:
    docker-compose up -d

# Stop all services
down:
    docker-compose down

# Run migrations
migrate-up:
    docker-compose exec auth-service go run cmd/migrate/main.go up
    docker-compose exec product-service go run cmd/migrate/main.go up
    docker-compose exec order-service go run cmd/migrate/main.go up
    docker-compose exec shop-service go run cmd/migrate/main.go up
    docker-compose exec warehouse-service go run cmd/migrate/main.go up

# Run database seeders
seed:
    docker-compose exec auth-service go run cmd/seeder/main.go
    docker-compose exec product-service go run cmd/seeder/main.go
    docker-compose exec shop-service go run cmd/seeder/main.go

# Run tests
test:
    go test ./... -v

# Run integration tests
test-integration:
    go test ./... -tags=integration -v

# Generate OpenAPI documentation
gen-docs:
    swag init -g cmd/main.go
```

## Testing

### Unit Tests
Run unit tests with:
```bash
make test
```

### Integration Tests
The project uses test suites for organized testing. Run integration tests with:
```bash
make test-integration
```

## API Documentation
The API is documented using OpenAPI. Access the documentation at:
[openapi.yml](openapi/openapi.yml)
[docs.md](docs.md)

## Monitoring and Logging

- Each service uses structured logging with Zap

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.