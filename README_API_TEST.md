# คู่มือทดสอบ API (RMS Backend)

เอกสารนี้ใช้สำหรับทดสอบ API ทีละเส้นแบบเร็ว ด้วย curl หรือ Postman

## 1) เตรียมระบบก่อนทดสอบ

### กรณีรันด้วย Docker ทั้งระบบ

```bash
docker compose up --build -d
```

### กรณีรัน DB ใน Docker และรัน API บนเครื่อง

```bash
docker compose up -d postgres
go run ./cmd/main.go
```

## 2) ตั้งรหัสผ่านสำหรับเทส Login (ทำครั้งเดียว)

หมายเหตุสำคัญ: ใน seed เดิมของ `employees.password_hash` เป็นค่า placeholder
จึงไม่สามารถใช้ login ได้ทันที ต้องตั้ง hash ใหม่ก่อน

1. เข้า PostgreSQL ใน container

```bash
docker exec -it rms-postgres psql -U postgres -d rms
```

2. รันคำสั่งด้านล่างใน psql

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

3. ออกจาก psql

```sql
\q
```

## 3) กำหนด Base URL

- Base URL: `http://localhost:8080`

ตัวอย่างตั้งตัวแปรใน PowerShell:

```powershell
$BASE_URL = "http://localhost:8080"
```

## 4) รายการ Endpoint ที่มีตอนนี้

- GET /health
- GET /api/v1/tables
- GET /api/v1/customer/menus
- GET /api/v1/roles (ต้องมี Bearer token และเป็น ADMIN)
- POST /api/v1/auth/login
- GET /api/v1/auth/me (ต้องมี Bearer token)
- POST /api/v1/auth/logout (ต้องมี Bearer token)
- POST /api/v1/employees (ต้องมี Bearer token และเป็น ADMIN)
- GET /api/v1/employees (ต้องมี Bearer token และเป็น ADMIN)
- GET /api/v1/employees/employeesid (ต้องมี Bearer token และเป็น ADMIN)
- PATCH /api/v1/employees/employeesid (ต้องมี Bearer token และเป็น ADMIN)
- PATCH /api/v1/employees/employeesid/status (ต้องมี Bearer token และเป็น ADMIN)
- GET /api/v1/tables (ต้องมี Bearer token)

## 5) ทดสอบทีละเส้น

### 5.1 Health Check

```bash
curl http://localhost:8080/health
```

คาดหวัง: success = true

### 5.2 ตารางทั้งหมด

```bash
curl http://localhost:8080/api/v1/tables
```

คาดหวัง: success = true และ data เป็นรายการโต๊ะ

### 5.3 เมนูลูกค้า

```bash
curl http://localhost:8080/api/v1/customer/menus
```

คาดหวัง: success = true และ data เป็นรายการเมนูที่ menu_status = true

### 5.4 Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@rms.com","password":"Admin1234!"}'
```

คาดหวังเมื่อสำเร็จ:
- ได้ accessToken
- tokenType เป็น Bearer
- มีข้อมูล employee และ role

### 5.5 Me (ตรวจ token)

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

คาดหวัง: success = true และได้ข้อมูลผู้ใช้จาก token

### 5.6 Logout

```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```

คาดหวัง: success = true

### 5.7 Roles (ADMIN เท่านั้น) - กรณีสำเร็จ 200

```bash
curl http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```

คาดหวัง: status 200 และ response รูปแบบใกล้เคียงนี้

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

### 5.8 Roles - กรณีไม่มี token (401)

```bash
curl http://localhost:8080/api/v1/roles
```

คาดหวัง: status 401

### 5.9 Roles - กรณีไม่ใช่ ADMIN (403)

ให้ login ด้วยบัญชี `cashier@rms.com` หรือ `chef@rms.com` แล้วใช้ token ที่ได้

```bash
curl http://localhost:8080/api/v1/roles \
  -H "Authorization: Bearer <NON_ADMIN_ACCESS_TOKEN>"
```

คาดหวัง: status 403

### 5.10 Employees (ADMIN เท่านั้น) - กรณีสำเร็จ (201)

ให้ login ด้วยบัญชี admin แล้วใช้ token ที่ได้

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeName": "สมชาย ใจดี",
    "roleId": 2,
    "phoneNumber": "0812345678",
    "email": "cashier1@rms.com",
    "hireDate": "2025-08-20",
    "password": "12345678"
  }'
```

คาดหวัง: status 201

```json
{
  "success": true,
  "message": "สร้างพนักงานสำเร็จ",
  "data": {
    "employeeId": 15,
    "employeeName": "สมชาย ใจดี",
    "roleId": 2,
    "phoneNumber": "0812345678",
    "email": "cashier1@rms.com",
    "hireDate": "2025-08-20",
    "employeeStatus": true
  }
}
```
### 5.11 Employees - กรณี email ซ้ำ (409)

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeName": "สมชาย ใจดี",
    "roleId": 2,
    "phoneNumber": "0812345678",
    "email": "cashier1@rms.com",
    "hireDate": "2025-08-20",
    "password": "12345678"
  }'
```

คาดหวัง: status 409

### 5.12 Employees - กรณีข้อมูลไม่ครบ (400)

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeName": "",
    "roleId": 2
  }'
```

คาดหวัง: status 400

### 5.13 Employees - กรณี role ไม่พบ (404)

```bash
curl -X POST http://localhost:8080/api/v1/employees \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "employeeName": "ทดสอบ",
    "roleId": 999,
    "phoneNumber": "0812345678",
    "email": "test999@rms.com",
    "hireDate": "2025-08-20",
    "password": "12345678"
  }'
```

คาดหวัง: status 404

### 5.14 Employees - ดูรายชื่อพนักงานทั้งหมด กรณีสำเร็จ (200)

```bash
curl "http://localhost:8080/api/v1/employees?page=1&limit=10" \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```
```json
  {
    "data": {
        "items": [
            {
                "employeeId": 1,
                "employeeName": "Admin User",
                "roleId": 1,
                "roleName": "ADMIN",
                "phoneNumber": "0811111111",
                "email": "admin@rms.com",
                "hireDate": "",
                "employeeStatus": true
            },
            {
                "employeeId": 2,
                "employeeName": "Cashier User",
                "roleId": 2,
                "roleName": "CASHIER",
                "phoneNumber": "0822222222",
                "email": "cashier@rms.com",
                "hireDate": "",
                "employeeStatus": true
            },
        ],
        "pagination": {
            "limit": 20,
            "page": 1,
            "total": 6
        }
    },
    "message": "ดึงรายการพนักงานสำเร็จ",
    "success": true
  }
```

คาดหวัง: status 200

### 5.15 Employee by ID - ดูข้อมูลพนักงานรายคน (200)

```bash
curl "http://localhost:8080/api/v1/employees?page=1&limit=10" \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```
```json
  {
      "data": {
         "employeeId": 1,
         "employeeName": "Admin User",
         "roleId": 1,
         "roleName": "ADMIN",
          "phoneNumber": "0811111111",
         "email": "admin@rms.com",
         "hireDate": "2025-01-01",
         "employeeStatus": true
      },
     "message": "ดึงข้อมูลพนักงานสำเร็จ",
      "success": true
  }
```
คาดหวัง: status 200

### 5.16 Employee by ID - ดูข้อมูลพนักงานรายคน ไม่พบพนักงาน (404)

```bash
curl "http://localhost:8080/api/v1/employees?page=1&limit=10" \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```
```json
  {
      "message": "ไม่พบพนักงาน",
      "success": false
  }
```
คาดหวัง: status 404

### 5.17 Employee by ID - แก้ไขข้อมูลพนักงาน (200)

```bash
curl curl -X PATCH http://localhost:8080/api/v1/employees/2 \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```
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
        "email": "",
        "hireDate": "",
        "employeeStatus": false
    }
  }
```
คาดหวัง: status 200

### 5.18 Employee by ID - เปิด/ปิดการใช้งานพนักงาน (200)

```bash
curl curl -X PATCH http://localhost:8080/api/v1/employees/2 \
  -H "Authorization: Bearer <ADMIN_ACCESS_TOKEN>"
```
```json
  {
    "success": true,
    "message": "อัปเดตสถานะพนักงานสำเร็จ",
    "data": {
        "employeeId": 11,
        "employeeStatus": false
    }
  }
```
คาดหวัง: status 200

### 5.19 Table - ดูรายการโต๊ะทั้งหมด (200)

```bash
curl -X GET "http://localhost:8080/api/v1/tables?status=AVAILABLE&page=1&limit=5" \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
```json
  {
    "success": true,
    "message": "ดึงรายการโต๊ะสำเร็จ",
    "data": [
        {
            "tableId": 1,
            "tableNumber": "A01",
            "capacity": 4,
            "tableStatus": "AVAILABLE",
            "createdAt": "2026-04-11T05:33:45.291484Z"
        },
        {
            "tableId": 2,
            "tableNumber": "A02",
            "capacity": 2,
            "tableStatus": "AVAILABLE",
            "createdAt": "2026-04-11T05:33:45.291484Z"
        }
    ]
  }
```
คาดหวัง: status 200

### 5.20 Table by ID - ดูข้อมูลโต๊ะรายตัว (200)

```bash
curl -X GET http://localhost:8080/api/v1/tables/1 \
  -H "Authorization: Bearer <ACCESS_TOKEN>"
```
```json
  {
    "success": true,
    "message": "ดึงข้อมูลโต๊ะสำเร็จ",
    "data": {
        "tableId": 1,
        "tableNumber": "A01",
        "capacity": 4,
        "tableStatus": "AVAILABLE",
        "createdAt": "2026-04-11T05:33:45.291484Z"
    }
}
```
คาดหวัง: status 200

## 6) วิธีทดสอบใน Postman

1. สร้าง Environment แล้วใส่ตัวแปร `baseUrl = http://localhost:8080`
2. สร้าง Request ตาม endpoint ด้านบน
3. หลัง login สำเร็จ ให้เก็บ token ลงตัวแปร `token`
4. ใส่ Header ในเส้น protected:
   - Authorization: Bearer {{token}}

## 7) ปัญหาที่พบบ่อย

### Login ไม่ผ่านทั้งที่ email ถูกต้อง

สาเหตุที่พบบ่อย: ยังไม่ได้ทำขั้นตอนตั้งรหัสผ่านในหัวข้อ "ตั้งรหัสผ่านสำหรับเทส Login"

ไฟล์ seed ที่เกี่ยวข้อง:
- `seeds/02_seed.sql`

### ต่อ DB ได้ แต่ API query ไม่ได้

ตรวจสอบว่า:
- DB container ทำงานอยู่
- API ใช้ค่า DB_HOST/DB_PORT ถูกต้องตามโหมดที่รัน
- สำหรับ local machine: ใช้ localhost:5435
- สำหรับ container-to-container: ใช้ postgres:5432

## 8) คำสั่งช่วยตรวจสถานะ

```bash
docker compose ps
docker compose logs -f postgres
```

```bash
go build ./...
```

## 9) ตัวอย่างลำดับเทสแบบเร็ว

1. GET /health
2. POST /api/v1/auth/login ด้วย admin@rms.com / Admin1234!
3. คัดลอก accessToken
4. GET /api/v1/auth/me พร้อม Authorization: Bearer <token>
5. GET /api/v1/tables และ GET /api/v1/customer/menus
6. GET /api/v1/roles ด้วย token ของ ADMIN
