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
    - mdw.transaction.topup.result
    - mdw.transaction.deduct.result
    - mdw.transaction.transfer.result

### Published/Produced Topic
    - mdw.transaction.topup.request
    - mdw.transaction.deduct.request
    - mdw.transaction.transfer.request
    

### RestAPI Endpoint
    - POST | /api/v1/account/register
    - POST | /api/v1/account/unregister
    - POST | /api/v1/account/all
    - GET  | /api/v1/account/:id
    - POST | /api/v1/account/detail
    - POST | /api/v1/account/balance/inquiry
    - POST | /api/v1/merchant/members
    - POST | /api/v1/merchant/members/period
    - POST | /api/v1/merchant/balance/inquiry