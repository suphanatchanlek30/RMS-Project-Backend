# คู่มือทดสอบ API (RMS Backend)

เอกสารนี้จัดรูปแบบให้เทสใน Postman ได้ง่าย โดยใช้แพทเทิร์นเดียวกันทุกเส้น

## เตรียมระบบก่อนทดสอบ

### ตัวเลือก A: รันทั้งระบบด้วย Docker

```bash
docker compose up --build -d
```

### ตัวเลือก B: รัน DB ใน Docker และรัน API บนเครื่อง

```bash
docker compose up -d postgres
go run ./cmd/main.go
```

## ตั้งรหัสผ่านสำหรับบัญชี seed (ทำครั้งเดียว)

```bash
docker exec -it rms-postgres psql -U postgres -d rms
```

```sql
CREATE EXTENSION IF NOT EXISTS pgcrypto;

UPDATE employees
SET password_hash = crypt('Admin1234!', gen_salt('bf'))
WHERE email = 'admin@rms.com';

UPDATE employees
SET password_hash = crypt('Cashier1234!', gen_salt('bf'))
WHERE email = 'cashier@rms.com';

UPDATE employees
SET password_hash = crypt('Chef1234!', gen_salt('bf'))
WHERE email = 'chef@rms.com';
```

```sql
\q
```

## ตั้งค่า Postman Environment (แนะนำ)

สร้าง Environment แล้วใส่ตัวแปรเหล่านี้

- `baseUrl` = `http://localhost:8080`
- `adminEmail` = `admin@rms.com`
- `adminPassword` = `Admin1234!`
- `adminToken` = (ค่าว่างไว้ก่อน)
- `cashierEmail` = `cashier@rms.com`
- `cashierPassword` = `Cashier1234!`
- `cashierToken` = (ค่าว่างไว้ก่อน)

## ขั้นตอนการใช้งาน (ทดสอบลำดับนี้)

### 1️⃣ Health Check

Method: `GET`  
URL: `{{baseUrl}}/health`  
Headers: None  
Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "server is running",
  "data": {
    "status": "ok"
  }
}
```

### 2️⃣ Login (ADMIN)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/auth/login`  
Headers:

- `Content-Type: application/json`

Body:

```json
{
  "email": "{{adminEmail}}",
  "password": "{{adminPassword}}"
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "เข้าสู่ระบบสำเร็จ",
  "data": {
    "employeeId": 1,
    "employeeName": "Admin User",
    "roleId": 1,
    "roleName": "ADMIN",
    "accessToken": "<JWT_TOKEN>",
    "tokenType": "Bearer"
  }
}
```

สำคัญ: คัดลอก `data.accessToken` ไปเก็บใน `adminToken`

### 3️⃣ Get Me

Method: `GET`  
URL: `{{baseUrl}}/api/v1/auth/me`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลผู้ใช้สำเร็จ",
  "data": {
    "employeeId": 1,
    "employeeName": "Admin User",
    "roleId": 1,
    "roleName": "ADMIN"
  }
}
```

### 4️⃣ Roles (ADMIN เท่านั้น)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/roles`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการ role สำเร็จ",
  "data": [
    { "roleId": 1, "roleName": "ADMIN" },
    { "roleId": 2, "roleName": "CASHIER" },
    { "roleId": 3, "roleName": "CHEF" }
  ]
}
```

### 5️⃣ Customer Menus (Public)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/customer/menus`  
Headers: None  
Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "fetch customer menus success",
  "data": [
    {
      "menuId": 1,
      "menuName": "..."
    }
  ]
}
```

### 6️⃣ Create Employee (ADMIN)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/employees`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "employeeName": "สมชาย ใจดี",
  "roleId": 2,
  "phoneNumber": "0812345678",
  "email": "cashier1@rms.com",
  "hireDate": "2025-08-20",
  "password": "12345678"
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างพนักงานสำเร็จ",
  "data": {
    "employeeId": 15,
    "employeeName": "สมชาย ใจดี",
    "roleId": 2,
    "roleName": "CASHIER",
    "phoneNumber": "0812345678",
    "email": "cashier1@rms.com",
    "hireDate": "2025-08-20",
    "employeeStatus": true
  }
}
```

### 7️⃣ Get Employees (ADMIN)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/employees?page=1&limit=10`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการพนักงานสำเร็จ",
  "data": {
    "items": [
      {
        "employeeId": 1,
        "employeeName": "Admin User"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 6
    }
  }
}
```

### 8️⃣ Get Employee By ID (ADMIN)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/employees/1`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลพนักงานสำเร็จ",
  "data": {
    "employeeId": 1,
    "employeeName": "Admin User",
    "roleId": 1,
    "roleName": "ADMIN",
    "phoneNumber": "0811111111",
    "email": "admin@rms.com",
    "hireDate": "",
    "employeeStatus": true
  }
}
```

### 9️⃣ Update Employee By ID (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/employees/2`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "employeeName": "สมชาย ใจดี",
  "roleId": 2,
  "phoneNumber": "0812345678"
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตข้อมูลพนักงานสำเร็จ",
  "data": {
    "employeeId": 2,
    "employeeName": "สมชาย ใจดี",
    "roleId": 2,
    "roleName": "CASHIER",
    "phoneNumber": "0812345678",
    "email": "cashier@rms.com",
    "hireDate": "",
    "employeeStatus": true
  }
}
```

### 🔟 Update Employee Status (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/employees/2/status`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "employeeStatus": false
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตสถานะพนักงานสำเร็จ",
  "data": {
    "employeeId": 2,
    "employeeStatus": false
  }
}
```

### 1️⃣1️⃣ Get Tables (ADMIN หรือ CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/tables?status=AVAILABLE&page=1&limit=5`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการโต๊ะสำเร็จ",
  "data": [
    {
      "tableId": 1,
      "tableNumber": "A01",
      "capacity": 4,
      "tableStatus": "AVAILABLE"
    }
  ]
}
```

### 1️⃣2️⃣ Get Table By ID (ADMIN หรือ CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/tables/1`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลโต๊ะสำเร็จ",
  "data": {
    "tableId": 1,
    "tableNumber": "A01",
    "capacity": 4,
    "tableStatus": "AVAILABLE"
  }
}
```

### 1️⃣3️⃣ Create Table (ADMIN)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/tables`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "tableNumber": "A07",
  "capacity": 5
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างโต๊ะสำเร็จ",
  "data": {
    "tableId": 4,
    "tableNumber": "A07",
    "capacity": 5,
    "tableStatus": "AVAILABLE"
  }
}
```

### 1️⃣4️⃣ Update Table (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/tables/4`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "tableNumber": "A04",
  "capacity": 8
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตข้อมูลโต๊ะสำเร็จ",
  "data": {
    "tableId": 4,
    "tableNumber": "A04",
    "capacity": 8,
    "tableStatus": "AVAILABLE"
  }
}
```

### 1️⃣5️⃣ Login (CASHIER)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/auth/login`  
Headers:

- `Content-Type: application/json`

Body:

```json
{
  "email": "{{cashierEmail}}",
  "password": "{{cashierPassword}}"
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "เข้าสู่ระบบสำเร็จ",
  "data": {
    "employeeId": 2,
    "employeeName": "Cashier User",
    "roleId": 2,
    "roleName": "CASHIER",
    "accessToken": "<JWT_TOKEN>",
    "tokenType": "Bearer"
  }
}
```

สำคัญ: คัดลอก `data.accessToken` ไปเก็บใน `cashierToken`

### 1️⃣6️⃣ Open Table Session (CASHIER)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/table-sessions/open`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "tableId": 1,
  "employeeId": 2
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "เปิดโต๊ะสำเร็จ",
  "data": {
    "sessionId": 1,
    "tableId": 1,
    "tableNumber": "A01",
    "startTime": "2025-08-20T12:00:00Z",
    "sessionStatus": "OPEN"
  }
}
```


### 1️⃣7️⃣ Get Table Session By ID (ADMIN หรือ CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/table-sessions/1`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูล session สำเร็จ",
  "data": {
    "sessionId": 1,
    "tableId": 1,
    "tableNumber": "A01",
    "startTime": "2025-08-20T12:00:00Z",
    "endTime": null,
    "sessionStatus": "OPEN"
  }
}
```

### 1️⃣8️⃣ Get Current Session By Table ID (ADMIN หรือ CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/tables/1/current-session`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึง session ปัจจุบันสำเร็จ",
  "data": {
    "sessionId": 1,
    "tableId": 1,
    "sessionStatus": "OPEN",
    "startTime": "2025-08-20T12:00:00Z"
  }
}
```

### 1️⃣9️⃣ Logout

Method: `POST`  
URL: `{{baseUrl}}/api/v1/auth/logout`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ออกจากระบบสำเร็จ",
  "data": null
}
```

## Negative Test ที่ควรลองเพิ่ม

### A) Roles ไม่มี token

Method: `GET`  
URL: `{{baseUrl}}/api/v1/roles`  
Headers: None

Expected: `401`

### B) Roles ใช้ token ที่ไม่ใช่ ADMIN

1. Login ด้วย `cashier@rms.com / Cashier1234!` แล้วเก็บ token ใน `cashierToken`
2. เรียก `GET {{baseUrl}}/api/v1/roles` ด้วย `Authorization: Bearer {{cashierToken}}`

Expected: `403`

### C) Employee By ID ไม่พบข้อมูล

Method: `GET`  
URL: `{{baseUrl}}/api/v1/employees/99999`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบพนักงาน"
}
```

### D) Create Employee email ซ้ำ

ใช้ body เดิมของข้อ 6 อีกครั้ง

Expected: `409`

### E) Create Employee role ไม่พบ

ตั้ง `roleId` เป็น `999`

Expected: `404`

### F) Open Table Session ไม่พบโต๊ะ

Method: `POST`  
URL: `{{baseUrl}}/api/v1/table-sessions/open`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "tableId": 99999,
  "employeeId": 2
}
```

Expected: `404`

### G) Open Table Session โต๊ะกำลังใช้งานอยู่

เปิดโต๊ะเดิมที่เปิดไปแล้วซ้ำอีกครั้ง

Expected: `409`

### H) Open Table Session ใช้ token ADMIN

Method: `POST`  
URL: `{{baseUrl}}/api/v1/table-sessions/open`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "tableId": 2,
  "employeeId": 1
}
```

Expected: `403`

### I) Get Table Session By ID ไม่พบ session

Method: `GET`  
URL: `{{baseUrl}}/api/v1/table-sessions/99999`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบ session"
}
```

### J) Get Current Session โต๊ะไม่มี session เปิดอยู่

Method: `GET`  
URL: `{{baseUrl}}/api/v1/tables/3/current-session`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่มี session ที่เปิดอยู่"
}
```

## สรุป Endpoint ทั้งหมดในระบบปัจจุบัน

- `GET /health`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`
- `GET /api/v1/customer/menus`
- `GET /api/v1/roles`
- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:employeeId`
- `PATCH /api/v1/employees/:employeeId`
- `PATCH /api/v1/employees/:employeeId/status`
- `GET /api/v1/tables`
- `GET /api/v1/tables/:tableId`
- `POST /api/v1/tables`
- `PATCH /api/v1/tables/:tableId`
- `POST /api/v1/table-sessions/open`
- `GET /api/v1/table-sessions/:sessionId`
- `GET /api/v1/tables/:tableId/current-session`

## คำสั่งช่วยตรวจสถานะ

```bash
docker compose ps
docker compose logs -f postgres
go build ./...
```
