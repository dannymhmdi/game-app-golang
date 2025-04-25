# payment proccess steps:
sequenceDiagram
Client->>+API: POST /payments/token
API->>+Redis: Store token (TTL=30m)
Redis-->>-API: OK
API-->>-Client: {token: "pay_xyz123"}

    Client->>+API: POST /payments/initiate (token)
    API->>+Redis: Validate token
    Redis-->>-API: OK
    API->>+MySQL: Create payment record
    MySQL-->>-API: OK
    API->>+Payment Processor: Charge card
    Payment Processor-->>-API: Success
    API->>+MySQL: Update status=SUCCESS
    MySQL-->>-API: OK
    API-->>-Client: {status: "SUCCESS"}
    API->>+Webhook: Notify merchant
## acording to this diagram :
1) first client request to POST /payments/token to generate a token for payment process
2) store token in redis and send back token to client
3) client requests to POST /payments/initiate and send token by request
4) token in request body checks with token stored in redis
5) after validation a record insert in payment db
6) The API sends the payment details (e.g., card info, amount) to a payment processor (e.g., Stripe, PayPal).

 **A payment processor is a third-party service that handles the actual transfer of money between a customer and a merchant. It securely communicates between banks, credit card networks, and your system to approve/decline transactions.**
 
7) Processor Responds to API /payments/initiate (Payment Processor → API)

If successful:
```go
{ "status": "SUCCESS", "transaction_id": "txn_123" }
```
if failed:

```go
{ "status": "FAILED", "error": "Insufficient funds" }
```
8) Update database : The API updates the payment status in MySQL

```go
UPDATE payments SET status = 'SUCCESS' WHERE id = 'pay_123';
```
9) Notify Client (API → Client): the current API /payments/initiate replies to the client:

```go
{ "status": "SUCCESS", "payment_id": "pay_123" }
```
it can be done by publish an event and notification service listen this event to send message to client

10) Notify Merchant via Webhook (API → Merchant):The API sends an async notification (e.g., HTTP POST) to the merchant’s server:

```go
POST https://merchant.com/webhook
{ "event": "payment_success", "payment_id": "pay_123", "amount": 100 }
```