# Quick Fix Patch for Ticket System

## Critical Changes Required

### 1. Fix stripe.go (Line 14)
**REMOVE**:
```go
params := &stripe.PaymentLinkListParams{
    Active: stripe.Bool(true),
}
```

**REPLACE WITH**:
```go
params := &stripe.PaymentLinkListParams{}
```

### 2. Fix stripe.go (Lines 55-69)
**REMOVE**:
```go
if s.CustomerDetails != nil {
    customer.Name = s.CustomerDetails.Name
    customer.Email = s.CustomerDetails.Email
}

if s.Customer != nil {
    customer.Name = s.Customer.Name
    customer.Email = s.Customer.Email
}

for _, field := range s.CustomFields {
    if field.Key == "companyname" {
        customer.Company = field.Text.Value
    }
}
```

**REPLACE WITH**:
```go
// Initialize email
if s.CustomerEmail != "" {
    customer.Email = s.CustomerEmail
} else if s.CustomerDetails != nil {
    customer.Email = s.CustomerDetails.Email
}

// Capture from custom fields with priority
var customFullName string
var customOrganization string

for _, field := range s.CustomFields {
    // Check for Full Name (Custom Field 1)
    if field.Key == "fullname" || field.Key == "full_name" || field.Key == "name" {
        customFullName = field.Text.Value
    }
    // Check for Organization (Custom Field 2)
    if field.Key == "organization" || field.Key == "company" || field.Key == "companyname" || field.Key == "org" {
        customOrganization = field.Text.Value
    }
}

// Set name with priority to custom field
if customFullName != "" {
    customer.Name = customFullName
} else if s.CustomerDetails != nil && s.CustomerDetails.Name != "" {
    customer.Name = s.CustomerDetails.Name
} else if s.Customer != nil && s.Customer.Name != "" {
    customer.Name = s.Customer.Name
}

// Set organization
if customOrganization != "" {
    customer.Company = customOrganization
}
```

## Testing Commands

### 1. Run with Debug Output
```bash
go run cmd/tickets/main.go -c config.yml -m producer -t templates/ticket.pdf -o generated_tickets/
```

### 2. Process Missing Tickets from CSV
```bash
go run cmd/tickets/main.go -c config.yml -m csv -csv missing_tickets.csv -t templates/ticket.pdf -o generated_tickets/
```

### 3. Send Individual Ticket
```bash
go run cmd/tickets/main.go -c config.yml -m publish-customer \
  -email "user@example.com" \
  -name "John Doe" \
  -company "Acme Corp" \
  -type "Standard" \
  -t templates/ticket.pdf \
  -o generated_tickets/
```

## Verification

After applying the fix, run this SQL to verify all tickets are captured:
```sql
SELECT COUNT(DISTINCT email) FROM users WHERE 'event_uuid' = ANY(events);
```

This should return **202** (matching Stripe records).

## Emergency CSV Format

If manual recovery is needed:
```csv
email,firstname,lastname,company,ticket_type
patrick@tobler.de,Patrick,Tobler,CEO,Standard
kitwillow91@gmail.com,Mubarak,Oladimeji,Student,Standard
```

---
**Deploy immediately to prevent further ticket generation issues.**