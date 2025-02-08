# OAuth2 Server Documentation

A comprehensive OAuth2 authorization server implementation supporting multiple grant types and complete user, client, and token management.

## Features

### OAuth2 Grant Types Support
- Authorization Code Grant
- Client Credentials Grant
- Password Grant (Resource Owner Password Credentials)
- Refresh Token Grant

### Client Management
- Client registration and authentication
- Multiple redirect URIs support
- Grant type restrictions
- Scope-based access control
- Confidential and public client support

### User Management
- User registration and authentication
- Role-based access control
- Password reset functionality
- Email verification
- Profile management

### Token Management
- Access token generation and validation
- Refresh token handling
- Token introspection
- Token revocation (single and bulk)
- Active session management

### Security Features
- Rate limiting
- Audit logging
- Session management
- Basic authentication support
- Scope-based authorization
- Token introspection

### Consent Management
- User consent tracking
- Granular permission control
- Consent revocation
- Consent history

## API Documentation

### OAuth2 Endpoints

#### Token Endpoint
```http
POST /api/v1/oauth/tokens
Authorization: Basic {client_credentials}
Content-Type: application/x-www-form-urlencoded
```

Supported grant types:
1. Authorization Code
```
grant_type=authorization_code
code={authorization_code}
redirect_uri={redirect_uri}
```

2. Client Credentials
```
grant_type=client_credentials
scope={scope}
```

3. Password
```
grant_type=password
username={username}
password={password}
scope={scope}
```

4. Refresh Token
```
grant_type=refresh_token
refresh_token={refresh_token}
```

#### Token Introspection
```http
POST /api/v1/oauth/introspect
Authorization: Basic {client_credentials}
Content-Type: application/x-www-form-urlencoded

token={token}
token_type_hint={access_token|refresh_token}
```

### Client Management API

#### Register New Client
```http
POST /api/v1/oauth/clients
Content-Type: application/json

{
    "name": "My Application",
    "redirect_uris": ["https://app.example.com/callback"],
    "grant_types": ["authorization_code", "refresh_token"],
    "scope": "read write",
    "confidential": true
}
```

#### List Clients
```http
GET /api/v1/oauth/clients
```

#### Update Client
```http
PUT /api/v1/oauth/clients/{client_id}
```

#### Delete Client
```http
DELETE /api/v1/oauth/clients/{client_id}
```

### User Management API

#### Register User
```http
POST /api/v1/users
Content-Type: application/json

{
    "username": "john.doe@example.com",
    "password": "secure_password",
    "name": "John Doe",
    "roles": ["user"]
}
```

#### User Operations
```http
PUT /api/v1/users/profile          # Update user profile
POST /api/v1/users/reset-password  # Reset password
POST /api/v1/users/verify-email    # Verify email
GET /api/v1/users                  # List users (admin only)
PUT /api/v1/users/{user_id}        # Update user (admin only)
DELETE /api/v1/users/{user_id}     # Delete user (admin only)
```

### Scope Management API

#### Create Scope
```http
POST /api/v1/oauth/scopes
Content-Type: application/json

{
    "name": "read_profile",
    "description": "Read user profile information"
}
```

#### Scope Operations
```http
GET /api/v1/oauth/scopes           # List scopes
PUT /api/v1/oauth/scopes/{scope_id}      # Update scope
DELETE /api/v1/oauth/scopes/{scope_id}    # Delete scope
```

### Consent Management API

#### Create Consent
```http
POST /api/v1/oauth/consents
Content-Type: application/json

{
    "client_id": "client_id",
    "scopes": ["read", "write"]
}
```

#### Consent Operations
```http
GET /api/v1/oauth/consents                # List consents
DELETE /api/v1/oauth/consents/{consent_id} # Revoke consent
```

### Security & Monitoring API

#### Token Management
```http
GET /api/v1/oauth/tokens          # List active tokens
POST /api/v1/oauth/tokens/revoke  # Revoke specific token
POST /api/v1/oauth/tokens/bulk-revoke # Bulk revoke tokens
```

#### Monitoring
```http
GET /api/v1/audit-logs           # View audit logs
GET /api/v1/oauth/rate-limits    # Check rate limit status
GET /api/v1/oauth/sessions       # List active sessions
DELETE /api/v1/oauth/sessions/{session_id} # End session
```

## Authentication

Most endpoints require authentication using HTTP Basic Authentication with client credentials:
```
Authorization: Basic base64(client_id:client_secret)
```

## Error Handling

The API uses standard HTTP status codes and returns errors in the following format:
```json
{
    "error": "error_code",
    "error_description": "Detailed error message",
    "error_uri": "https://documentation/errors/error_code"
}
```

## Rate Limiting

The API implements rate limiting per client. Current limits can be checked via the rate-limiting endpoint.

## Security Considerations

1. Always use HTTPS in production
2. Implement proper password hashing
3. Store client secrets securely
4. Implement token encryption
5. Set up proper CORS configuration
6. Enable audit logging
7. Implement IP whitelisting where appropriate

## Database Schema

The server requires the following database tables:
- users
- clients
- access_tokens
- refresh_tokens
- authorization_codes
- scopes
- client_scopes
- user_consents
- audit_logs