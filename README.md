# Transactions API

A small, production-friendly REST API to manage accounts and transactions with simple business rules:

- Each cardholder (customer) has an account with their data.
- For each operation carried out by the customer, a transaction is created and associated with this account.
- Each transaction has a type (cash purchase, installment purchase, withdrawal, or payment), an amount, and a creation date.
- Purchase and withdrawal transactions are recorded with a negative value, while payment transactions are recorded with a positive value.

---

## Quickstart (Docker)

```bash
cp .env.example .env
docker compose up --build
```

API on **http://localhost:8080**.

### Health check
```bash
curl http://localhost:8080/healthc
```

### Create an account
```bash
curl -s -X POST http://localhost:8080/accounts   -H 'Content-Type: application/json'   -d '{"document_number": "12345678900"}'
```

### Get an account
```bash
curl -s http://localhost:8080/accounts/1
```

### Create a transaction
```bash
# Purchase/Withdrawl - saved as negative
curl -s -X POST http://localhost:8080/transactions   -H 'Content-Type: application/json'   -d '{"account_id":1,"operation_type_id":1,"amount":50.0}'
```

```bash
# Payment - saved as positive
curl -s -X POST http://localhost:8080/transactions   -H 'Content-Type: application/json'   -d '{"account_id":1,"operation_type_id":4,"amount":60.0}'
```

- Amount must be positive.
- For purchases/withdrawals (types 1,2,3), the API will store it as negative automatically.
- For payments (type 4), the API will store it as positive.

---
