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

## วิธีเทส Customer Menus ผ่าน QR Token

ถ้าต้องการเทสเส้น `GET /api/v1/customer/menus` ให้ทำตามลำดับนี้

1. ล็อกอินด้วยบัญชี `cashier` เพื่อเอา `cashierToken`
2. เปิดโต๊ะด้วย `POST /api/v1/table-sessions/open`
3. สร้าง QR Session ด้วย `POST /api/v1/qr-sessions`
4. เอา `qrToken` จาก response ที่ได้
5. เรียก `GET /api/v1/customer/menus?qrToken={{qrToken}}`

ตัวอย่าง request ของเส้น customer menus

```http
GET {{baseUrl}}/api/v1/customer/menus?qrToken={{qrToken}}
Authorization: None
```

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงเมนูสำหรับลูกค้าสำเร็จ",
  "data": {
    "table": {
      "tableId": 1,
      "tableNumber": "A01"
    },
    "categories": [
      {
        "categoryId": 1,
        "categoryName": "อาหารจานหลัก"
      }
    ],
    "menus": [
      {
        "menuId": 101,
        "menuName": "ข้าวผัดกุ้ง",
        "price": 89.00,
        "description": "ข้าวผัดกุ้งสด",
        "menuStatus": true
      }
    ]
  }
}
```

กรณี error ที่ควรลองด้วย

- ไม่ส่ง `qrToken` -> `400 กรุณาระบุ qrToken`
- ใช้ `qrToken` ที่หมดอายุ -> `410 QR หมดอายุ`
- ใช้ `qrToken` ของ session ที่ปิดแล้ว -> `422 session ปิดแล้ว`

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
URL: `{{baseUrl}}/api/v1/customer/menus?qrToken={{qrToken}}`  
Headers: None  
Body: None

> ใช้ `qrToken` ที่ได้จากการสร้าง QR Session (หรือจาก Verify QR)

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงเมนูสำหรับลูกค้าสำเร็จ",
  "data": {
    "table": {
      "tableId": 1,
      "tableNumber": "A01"
    },
    "categories": [
      {
        "categoryId": 1,
        "categoryName": "อาหารจานหลัก"
      }
    ],
    "menus": [
      {
        "menuId": 101,
        "menuName": "ข้าวผัดกุ้ง",
        "price": 89.00,
        "description": "ข้าวผัดกุ้งสด",
        "menuStatus": true
      }
    ]
  }
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

### 1️⃣9️⃣ Close Table Session (CASHIER)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/table-sessions/1/close`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ปิดโต๊ะสำเร็จ",
  "data": {
    "sessionId": 1,
    "sessionStatus": "CLOSED",
    "endTime": "2025-08-20T14:00:00Z",
    "tableId": 1,
    "tableStatus": "AVAILABLE"
  }
}
```

### 2️⃣0️⃣ Create QR Session (CASHIER)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/qr-sessions`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "sessionId": 1
}
```

หมายเหตุ: ต้องเปิดโต๊ะใหม่ก่อน (ถ้า session 1 ถูกปิดไปแล้วจากข้อ 1️⃣9️⃣) โดยใช้ข้อ 1️⃣6️⃣ แล้วใส่ sessionId ใหม่ที่ได้

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้าง QR Session สำเร็จ",
  "data": {
    "qrSessionId": 1,
    "sessionId": 2,
    "qrCodeUrl": "http://localhost:3000/q/abcxyz123",
    "qrToken": "abcxyz123",
    "createdAt": "2025-08-20T12:00:00Z",
    "expiredAt": "2025-08-20T16:00:00Z"
  }
}
```

### 2️⃣1️⃣ Get QR Session By ID (ADMIN หรือ CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/qr-sessions/1`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูล QR Session สำเร็จ",
  "data": {
    "qrSessionId": 1,
    "sessionId": 2,
    "qrCodeUrl": "http://localhost:3000/q/abcxyz123",
    "expiredAt": "2025-08-20T16:00:00Z"
  }
}
```

### 2️⃣2️⃣ Verify QR Token (Public)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/qr/{{qrToken}}`  
Headers: None  
Body: None

หมายเหตุ: ใช้ค่า `qrToken` ที่ได้จากข้อ 2️⃣0️⃣ 

Expected Response (200):

```json
{
  "success": true,
  "message": "QR ใช้งานได้",
  "data": {
    "qrSessionId": 1,
    "sessionId": 2,
    "tableId": 1,
    "tableNumber": "A01",
    "sessionStatus": "OPEN",
    "expiredAt": "2025-08-20T16:00:00Z"
  }
}
```

### 2️⃣3️⃣ Create Category (ADMIN)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/categories`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "categoryName": "เครื่องดื่ม",
  "description": "เมนูเครื่องดื่ม"
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างหมวดหมู่สำเร็จ",
  "data": {
    "categoryId": 2,
    "categoryName": "เครื่องดื่ม",
    "description": "เมนูเครื่องดื่ม",
    "createdAt": "2025-08-20T10:00:00Z"
  }
}
```

### 2️⃣4️⃣ Get All Categories (Public)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/categories`  
Headers: None  
Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงหมวดหมู่สำเร็จ",
  "data": [
    {
      "categoryId": 1,
      "categoryName": "อาหารจานหลัก",
      "description": "เมนูอาหารหลักของร้าน"
    }
  ]
}
```

### 2️⃣5️⃣ Update Category (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/categories/2`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "categoryName": "เครื่องดื่มเย็น",
  "description": "หมวดเครื่องดื่มเย็น"
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตหมวดหมู่สำเร็จ",
  "data": {
    "categoryId": 2,
    "categoryName": "เครื่องดื่มเย็น",
    "description": "หมวดเครื่องดื่มเย็น"
  }
}
```

### 2️⃣6️⃣ Create Menu (ADMIN)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/menus`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuName": "ข้าวผัดกุ้ง",
  "categoryId": 1,
  "price": 89.00,
  "description": "ข้าวผัดกุ้งสด",
  "menuStatus": true
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างเมนูสำเร็จ",
  "data": {
    "menuId": 101,
    "menuName": "ข้าวผัดกุ้ง",
    "categoryId": 1,
    "price": 89.00,
    "description": "ข้าวผัดกุ้งสด",
    "menuStatus": true,
    "createdAt": "2025-08-20T10:00:00Z"
  }
}
```

### 2️⃣7️⃣ Get All Menus (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/menus?page=1&limit=20`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Query Parameters (ทั้งหมดเป็น optional):

| Parameter    | Description              | Example     |
| ------------ | ------------------------ | ----------- |
| `categoryId` | กรองตามหมวดหมู่           | `1`         |
| `keyword`    | ค้นหาตามชื่อเมนู          | `ข้าวผัด`    |
| `status`     | กรองตามสถานะ (`true`/`false`) | `true`  |
| `page`       | หน้าที่ (default: 1)      | `1`         |
| `limit`      | จำนวนต่อหน้า (default: 20) | `20`        |

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการเมนูสำเร็จ",
  "data": {
    "items": [
      {
        "menuId": 101,
        "menuName": "ข้าวผัดกุ้ง",
        "categoryId": 1,
        "categoryName": "อาหารจานหลัก",
        "price": 89.00,
        "description": "ข้าวผัดกุ้งสด",
        "menuStatus": true
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1
    }
  }
}
```

### 2️⃣8️⃣ Get Menu By ID (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/menus/101`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลเมนูสำเร็จ",
  "data": {
    "menuId": 101,
    "menuName": "ข้าวผัดกุ้ง",
    "categoryId": 1,
    "price": 89.00,
    "description": "ข้าวผัดกุ้งสด",
    "menuStatus": true
  }
}
```

### 2️⃣9️⃣ Update Menu (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/menus/101`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuName": "ข้าวผัดกุ้งพิเศษ",
  "price": 99.00,
  "description": "เพิ่มกุ้ง"
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตเมนูสำเร็จ",
  "data": {
    "menuId": 101,
    "menuName": "ข้าวผัดกุ้งพิเศษ",
    "price": 99.00,
    "description": "เพิ่มกุ้ง"
  }
}
```

### 3️⃣0️⃣ Update Menu Status (ADMIN)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/menus/101/status`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuStatus": false
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตสถานะเมนูสำเร็จ",
  "data": {
    "menuId": 101,
    "menuStatus": false
  }
}
```

### 3️⃣1️⃣ Logout

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

### 3️⃣2️⃣ Customer Create Order (Public)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/customer/orders`  
Headers:

- `Content-Type: application/json`

Body:

```json
{
  "qrToken": "abcxyz123",
  "items": [
    { "menuId": 101, "quantity": 2 },
    { "menuId": 102, "quantity": 1 }
  ]
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างคำสั่งซื้อสำเร็จ",
  "data": {
    "orderId": 9001,
    "sessionId": 1001,
    "tableId": 1,
    "orderTime": "2025-08-20T12:10:00Z",
    "orderStatus": "PENDING",
    "items": [
      {
        "orderItemId": 1,
        "menuId": 101,
        "menuName": "ข้าวผัดกุ้ง",
        "quantity": 2,
        "unitPrice": 89,
        "itemStatus": "WAITING"
      },
      {
        "orderItemId": 2,
        "menuId": 102,
        "menuName": "น้ำเปล่า",
        "quantity": 1,
        "unitPrice": 15,
        "itemStatus": "WAITING"
      }
    ]
  }
}
```

### 3️⃣3️⃣ Customer Get Orders (Public)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/customer/orders?qrToken={{qrToken}}`  
Headers: None  
Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงคำสั่งซื้อของลูกค้าสำเร็จ",
  "data": [
    {
      "orderId": 9001,
      "orderTime": "2025-08-20T12:10:00Z",
      "orderStatus": "PENDING",
      "items": [
        {
          "orderItemId": 1,
          "menuName": "ข้าวผัดกุ้ง",
          "quantity": 2,
          "unitPrice": 89,
          "itemStatus": "WAITING"
        }
      ]
    }
  ]
}
```

### 3️⃣4️⃣ Cashier Create Order

Method: `POST`  
URL: `{{baseUrl}}/api/v1/orders`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "sessionId": 1001,
  "tableId": 1,
  "createdByEmployeeId": 12,
  "items": [
    { "menuId": 101, "quantity": 1 }
  ]
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "สร้างคำสั่งซื้อสำเร็จ",
  "data": {
    "orderId": 9002,
    "sessionId": 1001,
    "tableId": 1,
    "createdByEmployeeId": 12,
    "orderTime": "2025-08-20T12:15:00Z",
    "orderStatus": "PENDING",
    "items": [
      {
        "orderItemId": 3,
        "menuId": 101,
        "menuName": "ข้าวผัดกุ้ง",
        "quantity": 1,
        "unitPrice": 89,
        "itemStatus": "WAITING"
      }
    ]
  }
}
```

### 3️⃣5️⃣ Get Order By ID

Method: `GET`  
URL: `{{baseUrl}}/api/v1/orders/9001`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูล order สำเร็จ",
  "data": {
    "orderId": 9001,
    "sessionId": 1001,
    "tableId": 1,
    "createdByEmployeeId": null,
    "orderTime": "2025-08-20T12:10:00Z",
    "orderStatus": "PENDING"
  }
}
```

### 3️⃣6️⃣ Get Orders By Session

Method: `GET`  
URL: `{{baseUrl}}/api/v1/table-sessions/1001/orders`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการ order ของโต๊ะสำเร็จ",
  "data": [
    {
      "orderId": 9001,
      "orderTime": "2025-08-20T12:10:00Z",
      "orderStatus": "PENDING"
    }
  ]
}
```

กรณี error ที่ควรลองเพิ่ม

- `items` ว่าง -> `400 ข้อมูลไม่ถูกต้อง`
- `menuId` ไม่พบ -> `404 ไม่พบ QR หรือเมนู` หรือ `404 ไม่พบข้อมูลที่ต้องการ` ตามเส้นที่เรียก
- `qrToken` หมดอายุ -> `410 QR หมดอายุ`
- `session` ปิดแล้ว หรือเมนูปิดขาย -> `422 เมนูปิดขาย/โต๊ะไม่พร้อมใช้งาน`

### 3️⃣7️⃣ Get Order Items (ADMIN/CASHIER/CHEF)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/orders/9001/items`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Query Parameters: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการอาหารสำเร็จ",
  "data": [
    {
      "orderItemId": 1,
      "menuId": 101,
      "menuName": "ข้าวผัดกุ้ง",
      "quantity": 2,
      "unitPrice": 89,
      "itemStatus": "WAITING"
    }
  ]
}
```

กรณี error ที่ควรลอง:

- `orderId` ไม่ถูกต้อง -> `400 orderId ไม่ถูกต้อง`
- ไม่พบ order -> `404 ไม่พบ order`

### 3️⃣8️⃣ Update Order Item Quantity (CASHIER)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/order-items/1`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "quantity": 3
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "แก้ไขจำนวนรายการอาหารสำเร็จ",
  "data": {
    "orderItemId": 1,
    "quantity": 3
  }
}
```

กรณี error ที่ควรลอง:

- `orderItemId` ไม่ถูกต้อง -> `400 orderItemId ไม่ถูกต้อง`
- `quantity` ไม่ถูกต้องหรือไม่ส่งมา -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่พบรายการอาหาร -> `404 ไม่พบ order item`
- สถานะไม่อนุญาตให้แก้ -> `422 สถานะไม่อนุญาตให้แก้`

### 3️⃣9️⃣ Cancel Order Item (CASHIER)

Method: `DELETE`  
URL: `{{baseUrl}}/api/v1/order-items/1`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ยกเลิกรายการอาหารสำเร็จ",
  "data": {
    "orderItemId": 1,
    "itemStatus": "CANCELLED"
  }
}
```

กรณี error ที่ควรลอง:

- `orderItemId` ไม่ถูกต้อง -> `400 orderItemId ไม่ถูกต้อง`
- ไม่พบรายการอาหาร -> `404 ไม่พบรายการอาหาร`
- รายการถูกทำแล้วหรือชำระแล้ว -> `422 รายการถูกทำแล้วหรือชำระแล้ว`

### 4️⃣0️⃣ Kitchen Orders (CHEF)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/kitchen/orders?status=WAITING&page=1&limit=10`  
Headers:

- `Authorization: Bearer {{chefToken}}`

Body: None

Query Parameters (optional):

- `status` = กรองตามสถานะอาหาร เช่น `WAITING`, `PREPARING`, `COMPLETED`
- `tableId` = กรองตามโต๊ะ
- `page` = หน้าที่ (default: 1)
- `limit` = จำนวนต่อหน้า (default: 10)

หมายเหตุ: ถ้าไม่ส่ง `status` ระบบจะดึงเฉพาะรายการที่อยู่ใน `WAITING` และ `PREPARING`

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงคิวครัวสำเร็จ",
  "data": [
    {
      "orderId": 9001,
      "tableId": 1,
      "tableNumber": "A01",
      "orderTime": "2025-08-20T12:10:00Z",
      "items": [
        {
          "orderItemId": 1,
          "menuName": "ข้าวผัดกุ้ง",
          "quantity": 2,
          "itemStatus": "WAITING"
        }
      ]
    }
  ]
}
```

กรณี error ที่ควรลอง:

- ไม่มี token หรือ token ไม่ใช่ CHEF -> `401/403` จาก middleware

### 4️⃣1️⃣ Update Order Item Status (CHEF)

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/order-items/1/status`  
Headers:

- `Authorization: Bearer {{chefToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "status": "PREPARING",
  "updatedByChefId": 21
}
```

Expected Response (200):

```json
{
  "success": true,
  "message": "อัปเดตสถานะอาหารสำเร็จ",
  "data": {
    "orderItemId": 1,
    "oldStatus": "WAITING",
    "newStatus": "PREPARING",
    "updatedTime": "2025-08-20T12:20:00Z"
  }
}
```

กรณี error ที่ควรลอง:

- ส่ง body ไม่ถูกต้อง -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่มี employeeId ใน token -> `401 token ไม่ถูกต้อง`
- ไม่พบ order item หรือสถานะไม่ถูกลำดับ -> `422 order item not found` หรือ `422 สถานะไม่ถูกต้องตามลำดับ`

### 4️⃣2️⃣ Get Order Item Status History (ADMIN/CASHIER/CHEF)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/order-items/1/history`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงประวัติสถานะสำเร็จ",
  "data": [
    {
      "statusHistoryId": 1,
      "status": "WAITING",
      "updatedByChefId": null,
      "updatedTime": "2025-08-20T12:10:00Z"
    },
    {
      "statusHistoryId": 2,
      "status": "PREPARING",
      "updatedByChefId": 21,
      "updatedTime": "2025-08-20T12:20:00Z"
    }
  ]
}
```

กรณี error ที่ควรลอง:

- ไม่พบข้อมูล -> `404 ไม่พบข้อมูล`

### 4️⃣3️⃣ Customer Order Status (Public)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/customer/order-status?qrToken={{qrToken}}`  
Headers: None  
Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงสถานะออเดอร์สำเร็จ",
  "data": {
    "tableId": 1,
    "tableNumber": "A01",
    "orders": [
      {
        "orderId": 9001,
        "orderTime": "2025-08-20T12:10:00Z",
        "items": [
          {
            "orderItemId": 1,
            "menuName": "ข้าวผัดกุ้ง",
            "quantity": 2,
            "itemStatus": "PREPARING"
          }
        ]
      }
    ]
  }
}
```

กรณี error ที่ควรลอง:

- ไม่ส่ง `qrToken` -> `400 qrToken จำเป็น`
- ไม่พบ QR หรือคำสั่งซื้อ -> `404 ไม่พบ QR หรือคำสั่งซื้อ`
- QR หมดอายุ -> `410 QR หมดอายุ`
- session ปิดแล้ว -> `422 session ปิดแล้ว`

### 4️⃣4️⃣ Get Session Bill (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/table-sessions/1001/bill`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "คำนวณบิลสำเร็จ",
  "data": {
    "sessionId": 1001,
    "tableId": 1,
    "tableNumber": "A01",
    "items": [
      {
        "orderItemId": 1,
        "menuName": "ข้าวผัดกุ้ง",
        "quantity": 2,
        "unitPrice": 89,
        "lineTotal": 178
      },
      {
        "orderItemId": 2,
        "menuName": "น้ำเปล่า",
        "quantity": 1,
        "unitPrice": 15,
        "lineTotal": 15
      }
    ],
    "subtotal": 193,
    "serviceCharge": 0,
    "vat": 0,
    "totalAmount": 193
  }
}
```

กรณี error ที่ควรลอง:

- `sessionId` ไม่ถูกต้อง -> `400 sessionId ไม่ถูกต้อง`
- ไม่พบ session -> `404 ไม่พบ session`
- session ไม่พร้อมคิดเงิน -> `422 session ไม่พร้อมคิดเงิน`

### 4️⃣5️⃣ Create Payment (CASHIER)

Method: `POST`  
URL: `{{baseUrl}}/api/v1/payments`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "sessionId": 1001,
  "paymentMethodId": 1,
  "receivedAmount": 200
}
```

Expected Response (201):

```json
{
  "success": true,
  "message": "ชำระเงินสำเร็จ",
  "data": {
    "paymentId": 3001,
    "sessionId": 1001,
    "paymentMethodId": 1,
    "paymentMethodName": "CASH",
    "totalAmount": 193,
    "receivedAmount": 200,
    "changeAmount": 7,
    "paymentTime": "2025-08-20T13:20:00Z",
    "paymentStatus": "PAID"
  }
}
```

กรณี error ที่ควรลอง:

- body ไม่ถูกต้อง -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่พบ session -> `404 ไม่พบ session`
- ไม่พบ payment method -> `404 ไม่พบ payment method`
- จ่ายซ้ำ session เดิม -> `409 จ่ายแล้ว`
- receivedAmount ไม่พอ -> `422 receivedAmount ไม่พอ`

### 4️⃣6️⃣ Get Payment By ID (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/payments/3001`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลการชำระเงินสำเร็จ",
  "data": {
    "paymentId": 3001,
    "sessionId": 1001,
    "paymentMethodId": 1,
    "paymentMethodName": "CASH",
    "totalAmount": 193,
    "paymentTime": "2025-08-20T13:20:00Z",
    "paymentStatus": "PAID"
  }
}
```

กรณี error ที่ควรลอง:

- `paymentId` ไม่ถูกต้อง -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่พบ payment -> `404 ไม่พบ payment`

### 4️⃣7️⃣ Get All Payments (ADMIN)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/payments?dateFrom=2025-08-01T00:00:00Z&dateTo=2025-08-31T23:59:59Z&paymentMethodId=1&status=PAID&page=1&limit=20`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Query Parameters (ทั้งหมดเป็น optional):

- `dateFrom` = วันเวลาเริ่มต้นของช่วงค้นหา
- `dateTo` = วันเวลาสิ้นสุดของช่วงค้นหา
- `paymentMethodId` = กรองตามวิธีการชำระ
- `status` = กรองตามสถานะการชำระ เช่น `PAID`
- `page` = หน้าที่ (default: 1)
- `limit` = จำนวนต่อหน้า (default: 20)

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงรายการการชำระเงินสำเร็จ",
  "data": {
    "items": [
      {
        "paymentId": 3001,
        "sessionId": 1001,
        "paymentMethodName": "CASH",
        "totalAmount": 193,
        "paymentStatus": "PAID",
        "paymentTime": "2025-08-20T13:20:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 1
    }
  }
}
```

กรณี error ที่ควรลอง:

- ไม่มี token -> `401`
- token ไม่ใช่ ADMIN -> `403`

### 4️⃣8️⃣ Get All Payment Methods (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/payment-methods`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงวิธีชำระเงินสำเร็จ",
  "data": [
    {
      "paymentMethodId": 1,
      "methodName": "CASH"
    },
    {
      "paymentMethodId": 2,
      "methodName": "QR"
    }
  ]
}
```

กรณี error ที่ควรลอง:

- ไม่มี token -> `401`
- token ไม่ใช่ ADMIN หรือ CASHIER -> `403`

### 4️⃣9️⃣ Get Receipt By Payment ID (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/payments/3001/receipt`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลใบเสร็จสำเร็จ",
  "data": {
    "receiptId": 4001,
    "receiptNumber": "RCT-20250820-0001",
    "issueDate": "2025-08-20T13:20:10Z",
    "totalAmount": 193,
    "payment": {
      "paymentId": 3001,
      "paymentMethodName": "CASH",
      "paymentTime": "2025-08-20T13:20:00Z"
    },
    "table": {
      "tableId": 1,
      "tableNumber": "A01"
    },
    "items": [
      {
        "menuName": "ข้าวผัดกุ้ง",
        "quantity": 2,
        "unitPrice": 89,
        "lineTotal": 178
      },
      {
        "menuName": "น้ำเปล่า",
        "quantity": 1,
        "unitPrice": 15,
        "lineTotal": 15
      }
    ]
  }
}
```

กรณี error ที่ควรลอง:

- `paymentId` ไม่ถูกต้อง -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่พบ receipt หรือ payment -> `404 ไม่พบ receipt หรือ payment`

### 5️⃣0️⃣ Get Receipt By Receipt ID (ADMIN/CASHIER)

Method: `GET`  
URL: `{{baseUrl}}/api/v1/receipts/4001`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Body: None

Expected Response (200):

```json
{
  "success": true,
  "message": "ดึงข้อมูลใบเสร็จสำเร็จ",
  "data": {
    "receiptId": 4001,
    "receiptNumber": "RCT-20250820-0001",
    "issueDate": "2025-08-20T13:20:10Z",
    "totalAmount": 193,
    "payment": {
      "paymentId": 3001,
      "paymentMethodName": "CASH",
      "paymentTime": "2025-08-20T13:20:00Z"
    },
    "table": {
      "tableId": 1,
      "tableNumber": "A01"
    },
    "items": [
      {
        "menuName": "ข้าวผัดกุ้ง",
        "quantity": 2,
        "unitPrice": 89,
        "lineTotal": 178
      },
      {
        "menuName": "น้ำเปล่า",
        "quantity": 1,
        "unitPrice": 15,
        "lineTotal": 15
      }
    ]
  }
}
```

กรณี error ที่ควรลอง:

- `receiptId` ไม่ถูกต้อง -> `400 ข้อมูลไม่ถูกต้อง`
- ไม่พบ receipt -> `404 ไม่พบ receipt`

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

### K) Close Session ไม่พบ session

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/table-sessions/99999/close`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบ session"
}
```

### L) Close Session ที่ปิดไปแล้ว

ปิด session เดิมซ้ำอีกครั้ง

Expected Response (422):

```json
{
  "success": false,
  "message": "session ปิดไปแล้ว"
}
```

### M) Close Session ที่ยังมีบิลค้างชำระ

เปิดโต๊ะใหม่ สร้าง order ที่สถานะ PENDING แล้วลองปิด

Expected Response (409):

```json
{
  "success": false,
  "message": "ยังมีบิลค้างชำระ"
}
```

### N) Create QR Session ไม่พบ session

Method: `POST`  
URL: `{{baseUrl}}/api/v1/qr-sessions`  
Headers:

- `Authorization: Bearer {{cashierToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "sessionId": 99999
}
```

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบ session"
}
```

### O) Create QR Session ซ้ำ (มี QR active อยู่แล้ว)

สร้าง QR Session ด้วย sessionId เดิมซ้ำอีกครั้ง

Expected Response (409):

```json
{
  "success": false,
  "message": "มี QR active อยู่แล้ว"
}
```

### P) Verify QR token ไม่พบ

Method: `GET`  
URL: `{{baseUrl}}/api/v1/qr/invalidtoken999`  
Headers: None

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบ QR"
}
```

### Q) Verify QR ที่หมดอายุ

รอให้ QR หมดอายุแล้วลองเรียกอีกครั้ง

Expected Response (410):

```json
{
  "success": false,
  "message": "QR หมดอายุ"
}
```

### R) Verify QR ที่ session ปิดแล้ว

ปิด session แล้วลองเรียก QR token เดิม

Expected Response (422):

```json
{
  "success": false,
  "message": "session ปิดแล้ว"
}
```

### S) Get QR Session By ID ไม่พบ

Method: `GET`  
URL: `{{baseUrl}}/api/v1/qr-sessions/99999`  
Headers:

- `Authorization: Bearer {{cashierToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบ QR Session"
}
```

### T) Create Category ชื่อซ้ำ

สร้างหมวดหมู่ด้วยชื่อเดิมซ้ำอีกครั้ง

Expected Response (409):

```json
{
  "success": false,
  "message": "ชื่อหมวดหมู่ซ้ำ"
}
```

### U) Update Category ไม่พบ

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/categories/99999`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "categoryName": "ทดสอบ",
  "description": "ทดสอบ"
}
```

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบหมวดหมู่"
}
```

### V) Update Category ชื่อซ้ำ

แก้ไขชื่อหมวดหมู่เป็นชื่อที่มีอยู่แล้ว

Expected Response (409):

```json
{
  "success": false,
  "message": "ชื่อหมวดหมู่ซ้ำ"
}
```

### W) Create Menu category ไม่พบ

Method: `POST`  
URL: `{{baseUrl}}/api/v1/menus`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuName": "ทดสอบ",
  "categoryId": 99999,
  "price": 50.00,
  "description": "ทดสอบ",
  "menuStatus": true
}
```

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบหมวดหมู่"
}
```

### X) Create Menu ชื่อซ้ำ

สร้างเมนูด้วยชื่อเดิมซ้ำอีกครั้ง

Expected Response (409):

```json
{
  "success": false,
  "message": "ชื่อเมนูซ้ำ"
}
```

### Y) Get Menu By ID ไม่พบ

Method: `GET`  
URL: `{{baseUrl}}/api/v1/menus/99999`  
Headers:

- `Authorization: Bearer {{adminToken}}`

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบเมนู"
}
```

### Z) Update Menu ไม่พบ

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/menus/99999`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuName": "ทดสอบ",
  "price": 50.00,
  "description": "ทดสอบ"
}
```

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบเมนู"
}
```

### AA) Update Menu ชื่อซ้ำ

แก้ไขชื่อเมนูเป็นชื่อที่มีอยู่แล้ว

Expected Response (409):

```json
{
  "success": false,
  "message": "ชื่อเมนูซ้ำ"
}
```

### AB) Update Menu Status ไม่พบ

Method: `PATCH`  
URL: `{{baseUrl}}/api/v1/menus/99999/status`  
Headers:

- `Authorization: Bearer {{adminToken}}`
- `Content-Type: application/json`

Body:

```json
{
  "menuStatus": false
}
```

Expected Response (404):

```json
{
  "success": false,
  "message": "ไม่พบเมนู"
}
```

### AC) Customer Menus ไม่ส่ง qrToken

Method: `GET`  
URL: `{{baseUrl}}/api/v1/customer/menus`  
Headers: None

Expected Response (400):

```json
{
  "success": false,
  "message": "กรุณาระบุ qrToken"
}
```

### AD) Customer Menus QR หมดอายุ

ใช้ qrToken ที่หมดอายุแล้ว

Expected Response (410):

```json
{
  "success": false,
  "message": "QR หมดอายุ"
}
```

### AE) Customer Menus session ปิดแล้ว

ใช้ qrToken ของ session ที่ปิดแล้ว

Expected Response (422):

```json
{
  "success": false,
  "message": "session ปิดแล้ว"
}
```

## สรุป Endpoint ทั้งหมดในระบบปัจจุบัน

- `GET /health`
- `POST /api/v1/auth/login`
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`
- `GET /api/v1/customer/menus`
- `GET /api/v1/customer/order-status`
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
- `GET /api/v1/table-sessions/:sessionId/bill`
- `GET /api/v1/tables/:tableId/current-session`
- `PATCH /api/v1/table-sessions/:sessionId/close`
- `POST /api/v1/qr-sessions`
- `GET /api/v1/qr-sessions/:qrSessionId`
- `GET /api/v1/qr/:token`
- `POST /api/v1/categories`
- `GET /api/v1/categories`
- `PATCH /api/v1/categories/:categoryId`
- `POST /api/v1/menus`
- `GET /api/v1/menus`
- `GET /api/v1/menus/:menuId`
- `PATCH /api/v1/menus/:menuId`
- `PATCH /api/v1/menus/:menuId/status`
- `GET /api/v1/orders/:orderId/items`
- `PATCH /api/v1/order-items/:orderItemId`
- `DELETE /api/v1/order-items/:orderItemId`
- `GET /api/v1/kitchen/orders`
- `PATCH /api/v1/order-items/:orderItemId/status`
- `GET /api/v1/order-items/:orderItemId/history`
- `POST /api/v1/payments`
- `GET /api/v1/payments/:paymentId`
- `GET /api/v1/payments`
- `GET /api/v1/payment-methods`
- `GET /api/v1/payments/:paymentId/receipt`
- `GET /api/v1/receipts/:receiptId`

## คำสั่งช่วยตรวจสถานะ

```bash
docker compose ps
docker compose logs -f postgres
go build ./...
```
