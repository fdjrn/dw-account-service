# Account Service
---
### Deskripsi
Service ini berfungsi untuk meng-handle manajemen akun wallet pengguna. Seperti:

    - Registrasi akun wallet pengguna (Regular/Merchant)
    - Detail informasi akun pengguna
    - Request Topup
    - Request Deduct
    - Request Transfer Saldo 

### Service Type
    - Message Producer
    - Message Consumer
    - RestAPI endpoint

### Subscribed/Consumed Topic
    - mdw.transaction.topup.result              ✅
    - mdw.transaction.deduct.result             ✅
    - mdw.transaction.transfer.result           ✅

### Published/Produced Topic
    - mdw.transaction.topup.request             ✅
    - mdw.transaction.deduct.request            ✅
    - mdw.transaction.transfer.request          ✅
    

### RestAPI Endpoint
    - POST | /api/v1/account/register           ✅
    - POST | /api/v1/account/unregister         ✅
    - POST | /api/v1/account/all                ✅
    - GET  | /api/v1/account/:id                ✅
    - POST | /api/v1/account/detail             ✅
    - POST | /api/v1/account/balance/inquiry    ✅
    - POST | /api/v1/merchant/members           ✅
    - POST | /api/v1/merchant/members/period    ✅
    - POST | /api/v1/merchant/balance/inquiry   ✅

### Build Docker Image
    docker build -t dw-account:1.0.0 -f Dockerfile .

### Available Environment Value:
    - DATABASE_MONGODB_URI : conncetion uri to mongodb cluster
        
        example: mongodb+srv://<user>:<password>@<cluster-host>/?retryWrites=true&w=majority

    - DATABASE_MONGODB_DB_NAME : Database Name used for parameter service

        example: dw-mdw-account

    - KAFKA_BROKERS : kafka cluster address

        example: touching-ghoul-8389-us1-kafka.upstash.io:9092

    - KAFKA_SASL_USER : kafka cluster username

    - KAFKA_SASL_PASSWORD : kafka cluster password

### Docker Run Command
    docker run -d -p 8000:8000 --name dw-account-service --env "DATABASE_MONGODB_DB_NAME=dev-mdw-account" --restart unless-stopped dw-account:1.0.0