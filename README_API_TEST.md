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
- POST /api/v1/auth/login
- GET /api/v1/auth/me (ต้องมี Bearer token)
- POST /api/v1/auth/logout (ต้องมี Bearer token)

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
