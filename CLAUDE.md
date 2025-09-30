# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Building and Running
- **Main API Server**: `go run cmd/app/main.go` - Starts the main API server on port 3000
- **Background Worker**: `go run cmd/worker/main.go` - Starts the background job processor 
- **Ticket Service**: `go run cmd/tickets/main.go` - Starts the report/PDF generation service
- **Hot Reload Development**: `air` - Uses Air for hot reloading during development
- **Build Binary**: `go build -o ./bin/app cmd/app/main.go` - Builds the main API server binary

### Database Operations
- **Run Migrations**: `go run cmd/migrate/main.go up` - Apply all pending migrations
- **Rollback Migration**: `go run cmd/migrate/main.go down 1` - Rollback the last migration
- **Check Migration Status**: `go run cmd/migrate/main.go version` - Show current migration version
- **Create New Migration**: `migrate create -ext sql -dir src/sql/migrations -seq migration_name`

### Testing
- **Run All Tests**: `go test ./...` - Execute all tests in the project
- **Run Tests with Coverage**: `go test -cover ./...` - Execute tests and show coverage
- **Generate Coverage Report**: `go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out`
- **Run BDD Tests**: `ginkgo tests/` - Run Ginkgo BDD test suite
- **Test Database**: Use `test.config.yml` configuration for separate test database

### Docker Development
- **Start Dependencies**: `docker-compose up -d` - Start PostgreSQL and NATS
- **Build Docker Image**: `docker build -t socious-api:latest .`
- **Full Stack**: `docker-compose up` - Start all services including API

### Configuration
- **Development Config**: `cp .tmp.config.yml config.yml` - Copy template configuration
- **Test Config**: Use `test.config.yml` for running tests

## Architecture Overview

### Service Architecture
This is a Go-based API built with Gin framework that follows a layered architecture:

- **cmd/**: Entry points for different services (app, worker, tickets, migrate)
- **src/apps/**: Application layer with business logic
- **src/sql/**: Database queries and migrations organized by entity

### Key Components

#### Multi-Service Architecture
1. **Main API Service** (`cmd/app`): HTTP API server handling client requests
2. **Worker Service** (`cmd/worker`): Background job processor for emails, notifications
3. **Ticket Service** (`cmd/tickets`): Report generation, PDF creation, CSV exports
4. **Migration Service** (`cmd/migrate`): Database schema management

#### Core Business Entities
- **Dual Identity System**: Users can act as themselves or on behalf of organizations
- **Projects**: Social impact initiatives with budgets and requirements
- **Contracts**: Work agreements with milestone tracking and escrow
- **Payment Processing**: Unified fiat (Stripe) and blockchain payment support
- **Impact Points**: Quantified social impact tracking system
- **Event System**: Comprehensive audit trail for all activities

#### Technology Stack Integration
- **Database**: PostgreSQL with sqlx for query execution
- **Message Queue**: NATS for async processing between services  
- **Authentication**: Socious ID (GoAccount SDK) with JWT tokens
- **Payment**: GoPay unified interface supporting Stripe and multiple blockchains
- **Email**: SendGrid integration via background workers
- **File Storage**: S3 with CDN support for media management

### Data Layer Organization

#### SQL Structure (`src/sql/`)
Each entity has its own directory containing:
- `fetch.sql` - Retrieval queries
- `upsert.sql` - Insert/update operations  
- `delete.sql` - Deletion operations
- Entity-specific complex queries

Key entities: users, organizations, identities, projects, contracts, wallets, media, events, feedbacks

#### Authentication Flow
1. Client initiates auth via `/auth` endpoint
2. OAuth flow with Socious ID service
3. Code exchange for JWT tokens at `/auth/session`
4. JWT validation on protected routes
5. Identity switching via `Current-Identity` header for org actions

#### Payment Architecture
- **Fiat Payments**: Stripe integration with multi-region support
- **Blockchain**: Multi-chain support (Cardano, EVM-compatible) via GoPay
- **Escrow**: Secure fund holding for contract milestones
- **Webhooks**: Payment confirmation processing

### Development Patterns

#### Request Flow
1. HTTP request hits Gin router in `src/apps/views/`
2. Middleware chain: security headers, auth, sanitization  
3. Handler validates input and calls model methods
4. Model executes SQL queries from `src/sql/`
5. Response formatted and returned to client

#### Background Processing
- NATS message queue for async operations
- Worker pools process email, notifications, analytics
- Event-driven architecture for cross-service communication

#### Security Implementation  
- XSS protection via bluemonday sanitization
- Parameterized queries prevent SQL injection
- JWT-based stateless authentication
- CORS configuration for web client access
- Secure headers middleware (HSTS, CSP, X-Frame-Options)

### Configuration Management

The system uses YAML configuration files:
- `config.yml` - Main configuration (copy from `.tmp.config.yml`)
- `test.config.yml` - Test environment settings
- Environment variables override YAML settings

Critical configuration sections:
- Database connection and SQL directories
- NATS message broker settings
- S3/CDN for file storage
- Payment provider credentials (Stripe, blockchain APIs)
- Socious ID authentication settings
- CORS origins for web clients

### Testing Strategy

Uses Ginkgo BDD framework with separate test database:
- **Unit Tests**: Model and utility function testing
- **Integration Tests**: API endpoint testing
- **Contract Tests**: External service mocking
- **E2E Tests**: Complete workflow validation

Test organization mirrors source structure with dedicated test database configuration.