# API Documentation

Base URL: `http://localhost:8080/api/v1`

## Authentication

### Sign Up
```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jhon@gmail.com",
    "password": "your_password"
  }'
```

### Sign In
```bash
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jhon@gmail.com",
    "password": "your_password"
  }'
```

### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your_refresh_token"
  }'
```

## User

### Get Profile
```bash
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzkyMDY4MDcsIm5iZiI6MTczOTIwMzIwNywiaWF0IjoxNzM5MjAzMjA3LCJzdWIiOiI1ZGVlOTc3YS0xNzY1LTQ2M2EtYTgzMC00NTU3ODJhYzUwN2UiLCJpc3MiOiJhdXRoLXNlcnZpY2UiLCJhdWQiOlsiYXV0aC1zZXJ2aWNlIiwicHJvZHVjdC1zZXJ2aWNlIiwib3JkZXItc2VydmljZSIsIndhcmVob3VzZS1zZXJ2aWNlIl0sInNjb3BlIjoidXNlcjpyZWFkIHVzZXI6d3JpdGUifQ.18dCt6Xo_GER-LTaKWcT1u2jQWEHDdjiaOVVLj4To40"
```

## Orders

### Checkout
```bash
curl -X POST http://localhost:8080/api/v1/orders/checkout \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": "product_id_1",
        "quantity": 2
      },
      {
        "product_id": "product_id_2",
        "quantity": 1
      }
    ]
  }'
```

## Products

### Get All Products
```bash
curl -X GET http://localhost:8080/api/v1/product \
  -H "Authorization: Bearer your_access_token"
```

### Get Product by ID
```bash
curl -X GET http://localhost:8080/api/v1/product/product_id \
  -H "Authorization: Bearer your_access_token"
```

## Shops

### Create Shop
```bash
curl -X POST http://localhost:8080/api/v1/shop \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My Shop",
    "description": "My awesome shop description"
  }'
```

### Get Shop Warehouses
```bash
curl -X GET http://localhost:8080/api/v1/shop/shop_id/warehouses \
  -H "Authorization: Bearer your_access_token"
```

## Warehouses

### Create Warehouse
```bash
curl -X POST http://localhost:8080/api/v1/warehouse \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Main Warehouse",
    "location": "123 Street, City",
    "status": "active"
  }'
```

### Update Warehouse Status
```bash
curl -X PATCH http://localhost:8080/api/v1/warehouse/warehouse_id/status \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "inactive"
  }'
```

### Transfer Stock
```bash
curl -X POST http://localhost:8080/api/v1/warehouse/transfer \
  -H "Authorization: Bearer your_access_token" \
  -H "Content-Type: application/json" \
  -d '{
    "source_warehouse_id": "source_id",
    "destination_warehouse_id": "destination_id",
    "products": [
      {
        "product_id": "product_id",
        "quantity": 5
      }
    ]
  }'
```

### Get Warehouse Stock
```bash
curl -X GET http://localhost:8080/api/v1/warehouse/warehouse_id/stock \
  -H "Authorization: Bearer your_access_token"
```

## Response Examples

### Authentication Response
```json
{
  "access_token": "eyJhbGciOiJS...",
  "refresh_token": "eyJhbGciOiJS...",
  "token_type": "Bearer"
}
```

### User Profile Response
```json
{
  "id": "user_id",
  "email": "user@example.com",
  "name": "John Doe"
}
```

### Product Response
```json
{
  "id": "product_id",
  "name": "Product Name",
  "price": 29.99,
  "description": "Product description"
}
```

### Shop Response
```json
{
  "id": "shop_id",
  "name": "Shop Name",
  "description": "Shop description"
}
```

### Warehouse Response
```json
{
  "id": "warehouse_id",
  "name": "Warehouse Name",
  "location": "Warehouse Location",
  "status": "active"
}
```

### Stock Response
```json
[
  {
    "product_id": "product_id",
    "quantity": 100
  }
]
```

## Error Responses

The API returns standard HTTP status codes:

- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 500: Internal Server Error

Error Response Format:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message description"
  }
}
```