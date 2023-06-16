# Account Service
---
### Deskripsi
Service ini berfungsi untuk meng-handle manajemen akun wallet pengguna. Seperti:

    - Registrasi akun wallet pengguna
    - Detail informasi akun pengguna
    - Request Topup
    - Request Deduct
    - Request Transfer Saldo 

### Service Type
    - Message Producer
    - Message Consumer
    - RestAPI endpoint

### Subscribed/Consumed Topic
    - mdw.transaction.topup.created
    - mdw.transaction.deduct.created
    - mdw.transaction.transfer.created

### Published/Produced Topic
    - mdw.transaction.topup.requested
    - mdw.transaction.deduct.requested
    - mdw.transaction.transfer.requested
    

### RestAPI Endpoint
    - /api/v1/account/all
    - /api/v1/account/all
    - /api/v1/account/:id
    - /api/v1/account/uid/:uid
    - /api/v1/account/register
    - /api/v1/account/unregister
    - /api/v1/account/balance/inquiry/:uid
    - /api/v1/account/balance/topup
    - /api/v1/account/balance/deduct
    - /api/v1/account/balance/transfer
