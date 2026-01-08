# Transaction Upload & Processing (In-Memory)

This service provides an asynchronous CSV upload to ingest transactions, process them using a worker pool, and expose query APIs with cursor-based pagination.

---

## ✨ Features

- CSV upload via HTTP API
- Asynchronous processing using in-memory queue + worker pool
- Thread-safe in-memory storage
- Cursor-based pagination (timestamp cursor)
- Idempotent background consumers
- Graceful shutdown support
- Unit tests, integration tests, and race-safe (`go test -race`)

---

## 🏗 Architecture Overview
cmd/\
└ api/\
└─ main.go\
internal/\
└ config/          # Configuration loading & setup\
└ contract/        # Request/response DTOs (API contracts)\
└ deliveries/      # Transport layer\
└─ http/        # HTTP handlers\
└─ consumer/    # Event/message consumers\
└─ worker/      # Background workers\
└ models/          # Domain models\
└ pkg/             # Shared utilities (logger, helpers, etc.)\
└ repositories/    # Data access layer\
└─ mock/        # Repository mocks\
└ services/        # Business logic\
└─ mock/        # Service mocks


---

## 📂 Main Components
### 1. Upload API

**Endpoint**\
POST /statements/upload

**Curl**
```bash
curl --location 'localhost:8080/statements' \
--form 'file=@"/Users/rafelychandrarizkilillah/go/src/go-flip-life-style-products/example.csv"'
```
**Response**
```json
{
    "upload_id": "6769ca92-0546-427d-b747-5f2a3731e37c"
}
```
**Logic**
- in service open and read the file
- saves the file to the OS with temp directory
- enqueues a background job with payload path the files, and upload_id
- if any data FAILED, publish to event
- returning the upload_id

### 2. Get Balance

**Endpoint**\
GET /balance

**Query Param**\
?upload_id=

**Curl**
```bash
curl --location 'localhost:8080/balance?upload_id=6769ca92-0546-427d-b747-5f2a3731e37c'
```
**Response**
```json
{
  "uploadID": "683f1833-2086-4e7c-8aa3-0740a93528ab",
  "balance": "-320000"
}
```
**Logic**
- get the list data by upload id from the repository
- after got the list data, service will calculate the balance within status SUCCESS and credit for add, debit for sub

### 2. Get List Issues Transactions

**Endpoint**\
GET /transactions/issues

**Query Param**\
?upload_id=\
limit=    # default 10\
next_cursor=\
prev_cursor=\
status= # default FAILED, PENDING\
transaction_type= \

**Curl**
```bash
curl --location 'localhost:8080/transactions/issues?upload_id=683f1833-2086-4e7c-8aa3-0740a93528ab'
```
**Response**
```json
{
  "kind": "collection",
  "contents": [
    {
      "uploadID": "847b3b8a-56d8-4e8a-bd1e-020fc06bf477",
      "timestamp": 1674509300,
      "counterParty": "NATALIE PARK",
      "type": "CREDIT",
      "amount": "800000",
      "status": "PENDING",
      "description": "bonus"
    },
    {
      "uploadID": "847b3b8a-56d8-4e8a-bd1e-020fc06bf477",
      "timestamp": 1674508858,
      "counterParty": "DAVID WONG",
      "type": "DEBIT",
      "amount": "100000",
      "status": "FAILED",
      "description": "transportation"
    },
    {
      "uploadID": "847b3b8a-56d8-4e8a-bd1e-020fc06bf477",
      "timestamp": 1674508555,
      "counterParty": "MICHAEL LEE",
      "type": "CREDIT",
      "amount": "150000",
      "status": "PENDING",
      "description": "refund"
    },
    {
      "uploadID": "847b3b8a-56d8-4e8a-bd1e-020fc06bf477",
      "timestamp": 1674508250,
      "counterParty": "ALICE TAN",
      "type": "DEBIT",
      "amount": "75000",
      "status": "FAILED",
      "description": "parking"
    }
  ],
  "pagination": {
    "prev": "",
    "next": "",
    "totalEntries": 4
  }
}
```
**Logic**
- mapping filtering from query param
- get the list data by upload id and filter from the repository
- get next cursor if there's any more pages


### 4. Queueing Store Data

**Logic**
- enqueue the job to the queue in service statements
- in worker pool, dequeue the job from the queue
- process the job by reading the file and store the data to the repository

### 5. Event Consumer Reconciliation Data

**Logic**
- publish event in service statements if the status FAILED
- in worker pool, consume the event
- process the event, if already process will return nil because we have idempotency for protect duplicate recon
- process the event max back off is 3x

Trade-offs in this project
- hard to scale the system because we use in-memory storage
- hard to see the log because api, consumer, and worker running in the same binary

🚀 How to Run\
I already set example.csv in this project, you can use it for testing\
You need a config.yaml
```yaml
app:
  name: go-flip-life-style-products
  port: 8080
  graceful_timeout: 5s
  time_out_api: 5s
worker:
  upload_worker:
    size: 5
consumer:
  reconciliation_consumer:
    size: 5
```
you can set the config.yaml in the root directory, like changing size binary for upload_worker, or change graceful_timeout
```bash
using make File
make run-api --> for running the system
make mock-gen --> for generate mock file using go.ubser
make test-cover --> for running unit test entire service
```

📌 Design Notes
- decimal is used for monetary values to avoid floating-point errors 
- cursor pagination chosen over offset pagination for stability 
- interfaces everywhere for testability

👤 Author\
Rafely Chandra Rizkilillah
