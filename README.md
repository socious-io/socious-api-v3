# Socious API V3

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Core Components](#core-components)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Documentation](#api-documentation)
- [Database](#database)
- [Authentication & Authorization](#authentication--authorization)
- [Payment Integration](#payment-integration)
- [Message Queue System](#message-queue-system)
- [Testing](#testing)
- [Deployment](#deployment)
- [Development Guidelines](#development-guidelines)

## Overview

Socious API V3 is a robust, scalable backend service for the Socious platform - a social impact ecosystem that connects organizations, individuals, and projects focused on social good. The API provides comprehensive functionality for user management, organization management, project creation, contract handling, payment processing, and impact tracking.

### Key Features
- **Identity Management**: Dual identity system supporting both users and organizations
- **Project Management**: Create, manage, and track social impact projects
- **Contract System**: Handle agreements between parties with escrow support
- **Payment Processing**: Integrated fiat (Stripe) and blockchain payment capabilities
- **Impact Points**: Track and measure social impact contributions
- **Feedback System**: Comprehensive review and rating system
- **Event Tracking**: Audit trail for all system activities
- **Media Management**: S3-based file storage and CDN integration
- **Verification Credentials**: Support for verifiable credentials and attestations

## Architecture

The system follows a layered architecture pattern:

```
┌─────────────────────────────────────────────────┐
│                   Client Layer                   │
│         (Web App, Mobile App, Third-party)       │
└─────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────┐
│                    API Layer                     │
│               (Gin HTTP Framework)               │
│                                                  │
│  ┌──────────────┐  ┌──────────────┐            │
│  │   Routes     │  │ Middlewares  │            │
│  └──────────────┘  └──────────────┘            │
└─────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────┐
│                 Business Logic                   │
│                                                  │
│  ┌──────────────┐  ┌──────────────┐            │
│  │    Models    │  │   Services   │            │
│  └──────────────┘  └──────────────┘            │
└─────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────┐
│                  Data Layer                      │
│                                                  │
│  ┌──────────────┐  ┌──────────────┐            │
│  │  PostgreSQL  │  │     NATS     │            │
│  └──────────────┘  └──────────────┘            │
└─────────────────────────────────────────────────┘
```

### Service Components

- **Main API Service**: Handles HTTP requests and business logic
- **Worker Service**: Processes background jobs (email, notifications, etc.)
- **Ticket Service**: Generates reports, PDFs, and handles bulk operations
- **Migration Service**: Database schema management

## Technology Stack

### Core Technologies
- **Language**: Go 1.24.2
- **Web Framework**: Gin v1.10.0
- **Database**: PostgreSQL 14.4
- **Message Queue**: NATS
- **Container**: Docker & Docker Compose

### Key Dependencies
- **Authentication**: `socious-io/goaccount` - Socious ID SDK for auth management
- **Payment Processing**:
  - `socious-io/gopay` v0.0.37 - Unified payment interface
  - `stripe/stripe-go` v81.1.1 - Stripe integration
- **Database**:
  - `jmoiron/sqlx` v1.4.0 - SQL extensions
  - `lib/pq` v1.10.9 - PostgreSQL driver
  - `golang-migrate/migrate` v4.17.1 - Database migrations
- **Messaging**:
  - `socious-io/gomq` v0.2.2 - Message queue abstraction
  - `nats-io/nats.go` v1.39.1 - NATS client
- **Email**: `sendgrid/sendgrid-go` v3.16.1 - Email service
- **Security**:
  - `microcosm-cc/bluemonday` v1.0.27 - XSS sanitization
  - `unrolled/secure` v1.17.0 - HTTP security headers
- **Testing**:
  - `onsi/ginkgo` v2.20.2 - BDD testing framework
  - `onsi/gomega` v1.34.2 - Matcher library
- **Utilities**:
  - `pdfcpu/pdfcpu` v0.11.0 - PDF generation
  - `skip2/go-qrcode` - QR code generation

## Project Structure

```
socious-api-v3/
├── cmd/                        # Application entry points
│   ├── app/                    # Main API server
│   │   └── main.go
│   ├── migrate/                # Database migration tool
│   │   └── main.go
│   ├── tickets/                # Ticket generation service
│   │   └── main.go
│   └── worker/                 # Background job processor
│       └── main.go
├── src/
│   ├── apps/                   # Application layer
│   │   ├── apps.go            # App initialization
│   │   ├── lib/               # Shared libraries
│   │   │   └── payment.go     # Payment utilities
│   │   ├── models/            # Data models
│   │   │   ├── user.go
│   │   │   ├── organization.go
│   │   │   ├── project.go
│   │   │   ├── contract.go
│   │   │   ├── identity.go
│   │   │   ├── media.go
│   │   │   ├── wallet.go
│   │   │   ├── card.go
│   │   │   ├── event.go
│   │   │   ├── feedback.go
│   │   │   ├── impact_points.go
│   │   │   ├── oauth_connect.go
│   │   │   ├── referring.go
│   │   │   ├── verification_credential.go
│   │   │   └── enum.go
│   │   ├── tickets/           # Ticket generation
│   │   │   ├── main.go
│   │   │   ├── pdf.go
│   │   │   ├── csv.go
│   │   │   ├── email.go
│   │   │   └── stripe.go
│   │   ├── utils/             # Utility functions
│   │   │   ├── utils.go
│   │   │   ├── security.go
│   │   │   └── sdg.go        # SDG-related utilities
│   │   ├── views/             # HTTP handlers
│   │   │   ├── views.go      # Route registration
│   │   │   ├── auth.go       # Authentication endpoints
│   │   │   ├── users.go      # User endpoints
│   │   │   ├── organizations.go
│   │   │   ├── projects.go
│   │   │   ├── contracts.go
│   │   │   ├── identities.go
│   │   │   ├── forms.go      # Request/response forms
│   │   │   └── middlewares.go # HTTP middlewares
│   │   └── workers/           # Background workers
│   │       ├── workers.go
│   │       ├── events.go
│   │       └── forms.go
│   ├── config/                # Configuration
│   │   └── config.go
│   └── sql/                   # SQL queries and migrations
│       ├── migrations/        # Database migrations
│       ├── users/            # User-related queries
│       ├── organizations/
│       ├── projects/
│       ├── contracts/
│       ├── cards/
│       ├── events/
│       ├── feedbacks/
│       ├── identities/
│       ├── media/
│       ├── oauth_connects/
│       ├── referrings/
│       └── wallets/
├── tests/                     # Test files
│   ├── main_test.go
│   ├── auth_test.go
│   ├── user_test.go
│   ├── project_test.go
│   ├── contract_test.go
│   └── data_test.go
├── docs/                      # Documentation
│   └── swagger.yaml          # OpenAPI specification
├── .air.toml                 # Hot reload configuration
├── .gitignore
├── .tmp.config.yml          # Template configuration
├── docker-compose.yml       # Docker services
├── Dockerfile              # Container definition
├── go.mod                  # Go dependencies
└── go.sum                  # Dependency checksums
```

## Core Components

### 1. Identity System
The platform uses a dual identity system where users can act as themselves or on behalf of organizations:
- **User Identity**: Personal profile with skills, bio, impact points
- **Organization Identity**: Company/NGO profile with mission, verification status
- Identity switching via `current-identity` header

### 2. Authentication & Authorization
- JWT-based authentication via Socious ID (GoAccount SDK)
- OAuth2 flow support with session management
- Role-based access control for organizations
- Middleware-enforced security on protected routes

### 3. Data Models

#### Primary Entities
- **Users**: Individual platform members with profiles, skills, and impact tracking
- **Organizations**: Verified entities that can post projects and hire
- **Projects**: Social impact initiatives with budgets and requirements
- **Contracts**: Agreements between parties with milestone tracking
- **Identities**: Unified interface for users and organizations
- **Wallets**: Payment method storage (fiat and crypto)
- **Media**: File attachments and profile images
- **Events**: Audit log of all system activities
- **Feedback**: Reviews and ratings between parties
- **Cards**: Referral and achievement cards
- **Impact Points**: Quantified social impact metrics

### 4. Payment System
Integrated payment processing supporting multiple methods:
- **Fiat Payments**: Stripe integration (multiple regions)
- **Blockchain Payments**: Support for multiple chains (Cardano, EVM-compatible)
- **Escrow System**: Secure fund holding for contracts
- **Multi-currency Support**: Handle various fiat and cryptocurrencies

### 5. Message Queue System
NATS-based asynchronous processing:
- Email notifications via SendGrid
- Background job processing
- Event streaming
- Worker pool management

### 6. Security Features
- XSS protection via bluemonday sanitization
- CSRF protection
- SQL injection prevention via parameterized queries
- Rate limiting
- Secure headers (HSTS, CSP, X-Frame-Options)
- Input validation and sanitization

## Getting Started

### Prerequisites
- Go 1.24.2 or higher
- PostgreSQL 14.4 or higher
- NATS server
- Docker & Docker Compose (optional)

### Quick Start

1. **Clone the repository**
```bash
git clone https://github.com/socious-io/socious-api-v3.git
cd socious-api-v3
```

2. **Set up configuration**
```bash
cp .tmp.config.yml config.yml
# Edit config.yml with your settings
```

3. **Start dependencies**
```bash
docker-compose up -d
```

4. **Install Go dependencies**
```bash
go mod download
```

5. **Run database migrations**
```bash
go run cmd/migrate/main.go
```

6. **Start the API server**
```bash
go run cmd/app/main.go
```

7. **Start the worker (in another terminal)**
```bash
go run cmd/worker/main.go
```

The API will be available at `http://localhost:3000`

### Development Mode
For hot reload during development:
```bash
air
```

## Configuration

Configuration is managed via YAML files. Create a `config.yml` from the template:

```yaml
env: development           # Environment: development, staging, production
port: 3000                # API server port
host: "http://localhost:3000"
secret: this-is-secret    # JWT secret
debug: true               # Debug mode

database:
  url: postgresql://user:pass@localhost:5432/socious?sslmode=disable
  sqldir: src/sql         # SQL queries directory
  migrations: src/sql/migrations

s3:
  access_key_id: your-key
  secret_access_key: your-secret
  default_region: ap-northeast-1
  bucket: your-bucket
  cdn_url: https://cdn.example.com

cors:
  origins:
    - 'http://localhost:3001'
    - 'https://app.socious.io'

nats:
  url: nats://127.0.0.1:4222
  secret: nats-secret

goaccounts:
  host: https://id.socious.io
  id: your-app-id
  secret: your-app-secret
  jwt:
    secret: jwt-secret
    duration: 24h

payment:
  fiats:
    - name: STRIPE
      apikey: sk_test_xxx
      service: STRIPE
      callback: /payment/callback
    - name: STRIPE_JP
      apikey: sk_test_xxx
      service: STRIPE
      callback: /payment/callback
  chains:
    - name: Cardano
      explorer: https://cardanoscan.io
      contractaddress: addr1xxx
      tokens:
        - name: ADA
          symbol: ADA
          address: native
          decimals: 6
      type: CARDANO
      mode: MAINNET
      apikey: your-api-key

sendgrid:
  apikey: your-sendgrid-key
  url: https://api.sendgrid.com
  templates:
    welcome: d-xxx
    verification: d-xxx
```

### Environment Variables
The system also supports environment variable overrides:
- `DATABASE_URL`: Override database connection
- `PORT`: Override server port
- `JWT_SECRET`: Override JWT secret
- `SENDGRID_API_KEY`: SendGrid API key
- `STRIPE_API_KEY`: Stripe API key

## API Documentation

### Swagger/OpenAPI
Interactive API documentation is available at:
- Development: `http://localhost:3000/docs`
- Production: `https://api.socious.io/docs`

### Base URLs
- Development: `https://dev.socious.io/api/v3/`
- Production: `https://socious.io/api/v3/`
- Local: `http://localhost:3000/`

### Authentication
Most endpoints require authentication via Bearer token:
```
Authorization: Bearer <token>
```

For organization actions, include the identity header:
```
Current-Identity: <organization-uuid>
```

### Main Endpoint Groups

#### Authentication (`/auth`)
- `POST /auth` - Start authentication flow
- `POST /auth/session` - Exchange code for tokens
- `POST /auth/refresh` - Refresh access token

#### Users (`/users`)
- `GET /users` - List users
- `GET /users/:id` - Get user details
- `PUT /users/:id` - Update user profile
- `DELETE /users/:id` - Delete user

#### Organizations (`/organizations`)
- `GET /organizations` - List organizations
- `POST /organizations` - Create organization
- `GET /organizations/:id` - Get organization details
- `PUT /organizations/:id` - Update organization
- `DELETE /organizations/:id` - Delete organization
- `POST /organizations/:id/members` - Add member
- `DELETE /organizations/:id/members/:userId` - Remove member

#### Projects (`/projects`)
- `GET /projects` - List projects (with filters)
- `POST /projects` - Create project
- `GET /projects/:id` - Get project details
- `PUT /projects/:id` - Update project
- `DELETE /projects/:id` - Delete project
- `POST /projects/:id/apply` - Apply to project
- `GET /projects/:id/applicants` - List applicants

#### Contracts (`/contracts`)
- `GET /contracts` - List contracts
- `POST /contracts` - Create contract
- `GET /contracts/:id` - Get contract details
- `PUT /contracts/:id` - Update contract
- `POST /contracts/:id/accept` - Accept contract
- `POST /contracts/:id/complete` - Mark as complete
- `POST /contracts/:id/dispute` - Raise dispute

#### Identities (`/identities`)
- `GET /identities/:id` - Get identity details
- `GET /identities/:id/projects` - List identity's projects
- `GET /identities/:id/contracts` - List identity's contracts

### Response Format
All API responses follow a consistent format:

Success response:
```json
{
  "data": {...},
  "message": "Success"
}
```

Error response:
```json
{
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

Paginated response:
```json
{
  "results": [...],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

## Database

### Schema Management
Database migrations are managed using golang-migrate:

```bash
# Create a new migration
migrate create -ext sql -dir src/sql/migrations -seq create_new_table

# Run migrations
go run cmd/migrate/main.go up

# Rollback
go run cmd/migrate/main.go down 1

# Check migration status
go run cmd/migrate/main.go version
```

### Key Tables
- `users` - User accounts and profiles
- `organizations` - Organization entities
- `identities` - Unified identity view
- `projects` - Project listings
- `contracts` - Work agreements
- `wallets` - Payment methods
- `media` - File storage references
- `events` - Audit trail
- `feedbacks` - Reviews and ratings
- `oauth_connects` - OAuth provider connections
- `impact_points_history` - Impact tracking

### Query Organization
SQL queries are organized by domain in `src/sql/`:
- Each entity has its own directory
- Common operations: `fetch.sql`, `upsert.sql`, `delete.sql`
- Complex queries are in separate files

### Database Best Practices
1. Use transactions for multi-table operations
2. Add appropriate indexes for frequently queried columns
3. Use parameterized queries to prevent SQL injection
4. Implement soft deletes where appropriate
5. Maintain referential integrity with foreign keys

## Authentication & Authorization

### Socious ID Integration
The platform uses Socious ID (GoAccount) for authentication:

1. **Session Creation**: Client initiates auth flow
2. **OAuth Flow**: User authenticates with Socious ID
3. **Token Exchange**: Code exchanged for JWT tokens
4. **Token Validation**: Each request validates JWT
5. **Identity Context**: User/org context determined per request

### JWT Token Structure
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "type": "user|organization",
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Middleware Chain
1. `SecureHeaders` - Add security headers
2. `SecureRequest` - Sanitize input
3. `LoginRequired` - Enforce authentication
4. `LoginOptional` - Optional auth
5. `AccountCenterRequired` - Admin endpoints

## Payment Integration

### Stripe Integration
Fiat payment processing via Stripe:
- Customer management
- Payment method storage
- Subscription handling
- Webhook processing
- Multi-region support (US, JP, EU)

### Blockchain Integration
Cryptocurrency payments via GoPay:
- Multi-chain support (Cardano, Ethereum, etc.)
- Wallet management
- Transaction tracking
- Smart contract interaction
- Token support

### Escrow System
Secure fund management for contracts:
1. Funds locked on contract creation
2. Released on milestone completion
3. Dispute resolution process
4. Automatic refunds on cancellation

## Message Queue System

### NATS Configuration
Message broker for asynchronous processing:
- Pub/Sub pattern for events
- Request/Reply for RPC
- Queue groups for load balancing
- Persistent subscriptions

### Worker Types
1. **Email Worker**: SendGrid email delivery
2. **Notification Worker**: Push notifications
3. **Analytics Worker**: Event processing
4. **Payment Worker**: Transaction processing

### Message Format
```json
{
  "id": "msg-uuid",
  "type": "email.send",
  "payload": {...},
  "timestamp": "2024-01-01T00:00:00Z",
  "retries": 0
}
```

## Testing

### Test Structure
Tests use Ginkgo BDD framework:

```bash
# Run all tests
go test ./...

# Run specific test suite
ginkgo tests/

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Categories
1. **Unit Tests**: Model and utility functions
2. **Integration Tests**: API endpoints
3. **Contract Tests**: External service mocks
4. **E2E Tests**: Full workflow validation

### Test Database
Use separate test database:
```yaml
# test.config.yml
database:
  url: postgresql://test:test@localhost:5432/socious_test
```

## Deployment

### Docker Deployment

Build and run with Docker:
```bash
# Build image
docker build -t socious-api:latest .

# Run container
docker run -p 3000:3000 \
  -e DATABASE_URL=postgresql://... \
  -e JWT_SECRET=... \
  socious-api:latest
```

### Docker Compose
Full stack deployment:
```bash
docker-compose up -d
```

Services:
- `api` - Main API server
- `worker` - Background processor
- `postgres` - PostgreSQL database
- `nats` - Message broker

### Production Considerations

1. **Environment Variables**: Use secrets management (AWS Secrets Manager, etc.)
2. **Database**: Use managed PostgreSQL (RDS, Cloud SQL)
3. **File Storage**: Configure S3 with CloudFront CDN
4. **Message Queue**: Use managed NATS or alternative (SQS, Pub/Sub)
5. **Monitoring**: Implement logging, metrics, and tracing
6. **Load Balancing**: Use ALB/NLB for high availability
7. **Auto-scaling**: Configure based on CPU/memory metrics
8. **Backup**: Regular database backups and point-in-time recovery

### Health Checks
- `GET /health` - Basic health check
- `GET /ready` - Readiness probe (checks DB, NATS)

## Development Guidelines

### Code Organization
1. **Models**: One file per entity in `src/apps/models/`
2. **Views**: Group related endpoints in `src/apps/views/`
3. **SQL**: Organize queries by entity in `src/sql/`
4. **Utils**: Reusable functions in `src/apps/utils/`

### Best Practices
1. **Error Handling**: Always return meaningful error messages
2. **Logging**: Use structured logging with context
3. **Validation**: Validate input at handler level
4. **Transactions**: Use for multi-step operations
5. **Testing**: Write tests for new features
6. **Documentation**: Update Swagger for API changes

### Git Workflow
1. Create feature branch from `main`
2. Make changes with descriptive commits
3. Write/update tests
4. Update documentation
5. Create pull request
6. Code review and merge

### Code Style
- Follow Go best practices and idioms
- Use `gofmt` for formatting
- Use `golint` for linting
- Keep functions small and focused
- Add comments for complex logic

## Troubleshooting

### Common Issues

1. **Database Connection Failed**
   - Check PostgreSQL is running
   - Verify connection string
   - Check network/firewall settings

2. **NATS Connection Failed**
   - Ensure NATS server is running
   - Check authentication credentials
   - Verify network connectivity

3. **Authentication Errors**
   - Verify JWT secret matches
   - Check token expiration
   - Ensure Socious ID is configured

4. **Payment Processing Issues**
   - Verify Stripe/blockchain API keys
   - Check webhook configuration
   - Review payment logs

### Debug Mode
Enable debug mode in config:
```yaml
debug: true
```

This enables:
- Detailed error messages
- Request/response logging
- SQL query logging
- Performance metrics

## Support

For issues and questions:
- GitHub Issues: [socious-io/socious-api-v3/issues](https://github.com/socious-io/socious-api-v3/issues)
- Documentation: [docs.socious.io](https://docs.socious.io)
- Email: support@socious.io

## License

Copyright (c) 2024 Socious Inc. All rights reserved.

---

*This documentation is maintained by the Socious development team. Last updated: 2024*