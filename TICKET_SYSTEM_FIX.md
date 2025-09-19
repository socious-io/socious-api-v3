# Ticket Management System Fix Guide

## Executive Summary

The ticket management system has two critical issues preventing proper ticket generation and tracking:

1. **Missing Tickets**: Only 134 out of 202 sold tickets are being captured because the system only fetches **active** payment links, missing tickets purchased through now-inactive links.

2. **Incorrect Name Capture**: Tickets are using inconsistent Stripe form data instead of the standardized Custom Field 1 (Full Name) and Custom Field 2 (Organization).

## Current System Architecture

### Ticket Processing Flow
```
Stripe Payment Links → Checkout Sessions → Customer Data → NATS Queue → PDF Generation → Email Delivery
```

### Key Components
- **stripe.go**: Fetches payment links and checkout sessions from Stripe
- **main.go**: Orchestrates the ticket processing pipeline
- **pdf.go**: Generates PDF tickets with customer information
- **email.go**: Sends tickets via SendGrid
- **csv.go**: Handles manual CSV imports as a fallback

## Problem Analysis

### Issue 1: Missing Tickets (68 tickets not captured)

**Root Cause**: Line 14 in `stripe.go`
```go
params := &stripe.PaymentLinkListParams{
    Active: stripe.Bool(true),  // ← Only fetches active links
}
```

**Impact**:
- Payment links that have been deactivated, archived, or expired after purchases were made are not included
- Results in ~33% of tickets not being processed automatically

### Issue 2: Inconsistent Name Capture

**Current Implementation** (Lines 55-68 in `stripe.go`):
```go
// Getting name from inconsistent sources
if s.CustomerDetails != nil {
    customer.Name = s.CustomerDetails.Name  // ← Uses Stripe's default form
    customer.Email = s.CustomerDetails.Email
}

// Only captures "companyname" custom field
for _, field := range s.CustomFields {
    if field.Key == "companyname" {
        customer.Company = field.Text.Value
    }
}
```

**Problems**:
- Uses Stripe's default `CustomerDetails.Name` which may be incomplete or formatted inconsistently
- Doesn't capture Custom Field 1 (Full Name)
- Doesn't capture Custom Field 2 (Organization)
- Hardcoded to look for "companyname" field which may not exist

## Complete Solution

### Step 1: Fix Missing Tickets Issue

**File**: `src/apps/tickets/stripe.go`

**Changes Required**:
```go
// Line 12-31: Update fetchPaymentLinks function
func fetchPaymentLinks() {
    // Remove the Active filter to get ALL payment links
    params := &stripe.PaymentLinkListParams{}

    // Increase pagination limit for better performance
    params.Limit = stripe.Int64(100)

    // List ALL payment links (active and inactive)
    i := paymentlink.List(params)
    for i.Next() {
        link := i.PaymentLink()
        fmt.Printf("Payment Link ID: %s, URL: %s, Active: %v\n",
            link.ID, link.URL, link.Active)
        fetchSuccessfulPaymentsForLink(link.ID)
    }

    if err := i.Err(); err != nil {
        log.Printf("Error listing payment links: %v\n", err)
    }
}
```

### Step 2: Fix Name Capture from Custom Fields

**File**: `src/apps/tickets/stripe.go`

**Changes Required**:
```go
// Line 33-84: Update fetchSuccessfulPaymentsForLink function
func fetchSuccessfulPaymentsForLink(paymentLinkID string) {
    // Create parameters for listing checkout sessions
    params := &stripe.CheckoutSessionListParams{}
    params.Filters.AddFilter("payment_link", "", paymentLinkID)
    params.Filters.AddFilter("status", "", "complete")

    // Expand to get custom fields
    params.AddExpand("data.customer")
    params.AddExpand("data.customer_details")
    params.AddExpand("data.custom_fields")  // Add this expansion

    // List checkout sessions for the payment link
    i := session.List(params)
    for i.Next() {
        s := i.CheckoutSession()
        fmt.Printf("Session ID: %s\n", s.ID)
        customer := new(Customer)

        // Initialize with email first
        if s.CustomerEmail != "" {
            customer.Email = s.CustomerEmail
        } else if s.CustomerDetails != nil {
            customer.Email = s.CustomerDetails.Email
        }

        // NEW: Properly capture from custom fields
        var customFullName string
        var customOrganization string

        // Process all custom fields
        for _, field := range s.CustomFields {
            fmt.Printf("Custom Field - Key: %s, Label: %s, Value: %s\n",
                field.Key, field.Label.Custom, field.Text.Value)

            // Check for Full Name field (Custom Field 1)
            // Match by position, label, or key patterns
            if field.Label.Custom == "Full Name" ||
               field.Key == "fullname" ||
               field.Key == "full_name" ||
               field.Key == "name" {
                customFullName = field.Text.Value
            }

            // Check for Organization field (Custom Field 2)
            // Match by position, label, or key patterns
            if field.Label.Custom == "Organization" ||
               field.Key == "organization" ||
               field.Key == "company" ||
               field.Key == "companyname" ||
               field.Key == "org" {
                customOrganization = field.Text.Value
            }
        }

        // Set customer data with priority to custom fields
        if customFullName != "" {
            customer.Name = customFullName
        } else if s.CustomerDetails != nil && s.CustomerDetails.Name != "" {
            // Fallback to CustomerDetails if custom field not found
            customer.Name = s.CustomerDetails.Name
        } else if s.Customer != nil && s.Customer.Name != "" {
            // Final fallback to Customer object
            customer.Name = s.Customer.Name
        }

        // Set organization/company
        if customOrganization != "" {
            customer.Company = customOrganization
        }

        // Set ticket type
        customer.TicketType = linkType(paymentLinkID)

        // Log the extracted data for debugging
        fmt.Printf("Extracted Customer Data:\n")
        fmt.Printf("  Name: %s\n", customer.Name)
        fmt.Printf("  Email: %s\n", customer.Email)
        fmt.Printf("  Company: %s\n", customer.Company)
        fmt.Printf("  Ticket Type: %s\n", customer.TicketType)

        // Only publish if we have minimum required data
        if customer.Email != "" && customer.Name != "" {
            publish(consumerTitle(CUSTOMER), customer)
        } else {
            log.Printf("WARNING: Missing required data for session %s (Name: %s, Email: %s)\n",
                s.ID, customer.Name, customer.Email)
        }
    }

    if err := i.Err(); err != nil {
        log.Printf("Error listing checkout sessions: %v\n", err)
    }
}
```

### Step 3: Add Debug Mode for Custom Field Detection

Create a utility function to help identify the correct custom field keys:

**File**: `src/apps/tickets/stripe.go`

**Add this function**:
```go
// Add this debug function to identify custom field structure
func debugCustomFields(s *stripe.CheckoutSession) {
    fmt.Println("=== DEBUG: Custom Fields Structure ===")
    fmt.Printf("Session ID: %s\n", s.ID)

    if s.CustomFields == nil || len(s.CustomFields) == 0 {
        fmt.Println("No custom fields found")
        return
    }

    for i, field := range s.CustomFields {
        fmt.Printf("Field %d:\n", i+1)
        fmt.Printf("  Key: %s\n", field.Key)
        if field.Label != nil {
            fmt.Printf("  Label (Custom): %s\n", field.Label.Custom)
            fmt.Printf("  Label (Type): %s\n", field.Label.Type)
        }
        if field.Text != nil {
            fmt.Printf("  Value: %s\n", field.Text.Value)
        }
        if field.Dropdown != nil {
            fmt.Printf("  Dropdown Value: %s\n", field.Dropdown.Value)
        }
        if field.Numeric != nil {
            fmt.Printf("  Numeric Value: %s\n", field.Numeric.Value)
        }
        fmt.Printf("  Type: %s\n", field.Type)
        fmt.Printf("  Optional: %v\n", field.Optional)
        fmt.Println("---")
    }
    fmt.Println("=====================================")
}
```

### Step 4: Add Venue Map Support

Since they mentioned including a new venue map with tickets, update the email template:

**File**: `src/apps/tickets/email.go`

```go
func sendEmail(apikey, email, name, ticketPath, venueMapPath string) {
    from := mail.NewEmail("Socious", "info@socious.io")
    to := mail.NewEmail(name, email)

    message := mail.NewV3Mail()
    message.SetFrom(from)
    message.SetTemplateID(EMAIL_TEMPLATE)

    personalization := mail.NewPersonalization()
    personalization.AddTos(to)
    personalization.SetDynamicTemplateData("name", name)

    message.AddPersonalizations(personalization)

    // Attach ticket PDF
    ticketBytes, err := os.ReadFile(ticketPath)
    if err != nil {
        log.Printf("Failed to read ticket file: %v \n", err)
        return
    }

    ticketAttachment := mail.NewAttachment()
    ticketAttachment.Content = base64.StdEncoding.EncodeToString(ticketBytes)
    ticketAttachment.Type = "application/pdf"
    ticketAttachment.Filename = "ticket.pdf"
    ticketAttachment.Disposition = "attachment"
    message.AddAttachment(ticketAttachment)

    // Attach venue map if provided
    if venueMapPath != "" {
        mapBytes, err := os.ReadFile(venueMapPath)
        if err != nil {
            log.Printf("Failed to read venue map: %v \n", err)
        } else {
            mapAttachment := mail.NewAttachment()
            mapAttachment.Content = base64.StdEncoding.EncodeToString(mapBytes)
            mapAttachment.Type = "application/pdf"
            mapAttachment.Filename = "venue_map.pdf"
            mapAttachment.Disposition = "attachment"
            message.AddAttachment(mapAttachment)
        }
    }

    client := sendgrid.NewSendClient(apikey)
    response, err := client.Send(message)
    if err != nil {
        log.Printf("Send error: %v \n", err)
    } else {
        log.Printf("Email sent to %s - Status: %d\n", email, response.StatusCode)
    }
}
```

## Implementation Steps

### 1. Immediate Fix (Production Hotfix)
```bash
# 1. Update stripe.go with the fixes
# 2. Test with debug mode enabled
go run cmd/tickets/main.go \
    -c config.yml \
    -m producer \
    -t templates/ticket.pdf \
    -o generated_tickets/ \
    -ak "sendgrid_api_key"

# 3. Verify all 202 tickets are captured
# 4. Deploy the fix
```

### 2. Manual Recovery for Missing Tickets
```bash
# Use CSV import for immediate recovery
go run cmd/tickets/main.go \
    -c config.yml \
    -m csv \
    -csv missing_tickets.csv \
    -t templates/ticket.pdf \
    -o generated_tickets/
```

CSV Format for recovery:
```csv
email,firstname,lastname,company,ticket_type
user@example.com,John,Doe,Acme Corp,Standard
```

### 3. Testing Checklist
- [ ] Verify inactive payment links are fetched
- [ ] Confirm Custom Field 1 (Full Name) is captured
- [ ] Confirm Custom Field 2 (Organization) is captured
- [ ] Test with missing custom fields (fallback to defaults)
- [ ] Verify venue map attachment
- [ ] Check Metabase count matches Stripe (202 tickets)
- [ ] Test email delivery with attachments

## Monitoring & Validation

### SQL Query to Verify Ticket Count
```sql
-- Check users with event registrations
SELECT
    COUNT(DISTINCT email) as unique_attendees,
    COUNT(*) as total_tickets,
    array_agg(DISTINCT events) as event_ids
FROM users
WHERE events && ARRAY['current_event_uuid']::uuid[];
```

### Stripe Validation Script
```go
// Add this as a validation mode
func validateTicketCount() {
    totalSessions := 0
    params := &stripe.PaymentLinkListParams{}
    params.Limit = stripe.Int64(100)

    i := paymentlink.List(params)
    for i.Next() {
        link := i.PaymentLink()

        sessionParams := &stripe.CheckoutSessionListParams{}
        sessionParams.Filters.AddFilter("payment_link", "", link.ID)
        sessionParams.Filters.AddFilter("status", "", "complete")

        j := session.List(sessionParams)
        for j.Next() {
            totalSessions++
        }
    }

    fmt.Printf("Total successful payments found: %d\n", totalSessions)
}
```

## Preventive Measures

### 1. Add Automated Monitoring
- Set up daily reconciliation between Stripe and database
- Alert when discrepancy > 5%
- Log all skipped sessions with reasons

### 2. Improve Error Handling
- Add retry mechanism for failed ticket generations
- Queue failed tickets for manual review
- Send admin notifications for processing failures

### 3. Configuration Management
- Make custom field keys configurable
- Add field mapping configuration:
```yaml
stripe:
  custom_fields:
    full_name:
      - "fullname"
      - "full_name"
      - "name"
    organization:
      - "organization"
      - "company"
      - "companyname"
      - "org"
```

## Rollback Plan

If issues arise after deployment:

1. **Revert Code**:
```bash
git revert [commit-hash]
```

2. **Use CSV Fallback**: Process all tickets manually via CSV import

3. **Direct Database Update**: If needed, directly update user events:
```sql
UPDATE users
SET events = array_append(events, 'event_uuid'::uuid)
WHERE email IN (SELECT email FROM ticket_csv_import);
```

## Contact for Support

- **Primary Developer**: Ehsan Mahmoudi
- **Team Lead**: Joshua / Seira
- **Escalation**: Mo (for complex debugging)

## Appendix: Common Issues

### Issue: "Customer not found in list"
**Solution**: Check if payment link was archived. Remove Active filter.

### Issue: "Name showing as blank"
**Solution**: Custom fields not properly configured in Stripe payment link.

### Issue: "Metabase count mismatch"
**Solution**: Run validation script to identify missing sessions.

---

*Last Updated: September 2025*
*Version: 1.0*