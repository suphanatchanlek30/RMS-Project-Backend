# RMS Project Backend

โปรเจกต์ Backend ของระบบจัดการร้านอาหาร (Restaurant Management System) พัฒนาด้วย Go + Fiber + PostgreSQL + Docker

## มีอะไรในโปรเจกต์นี้

- โครงสร้างแยกชั้นชัดเจน: Handler, Service, Repository, Model
- ระบบยืนยันตัวตนด้วย JWT (Login, Me, Logout)
- ระบบสิทธิ์การใช้งานตาม Role (ADMIN, CASHIER, CHEF)
- API จัดการพนักงาน (สร้าง, ค้นหา, แก้ไข, ปิด/เปิดใช้งาน)
- API จัดการโต๊ะ (ดูรายการ, ดูรายตัว, สร้าง, แก้ไข)
- API เมนูสำหรับลูกค้า
- API คำสั่งซื้อของลูกค้าผ่าน QR และคำสั่งซื้อของ cashier
- SQL seed สำหรับสร้าง schema และข้อมูลตั้งต้น
- รองรับการรันแบบ Docker ทั้งระบบ หรือรัน Go local + DB ใน Docker

## เทคโนโลยีที่ใช้

- Go 1.25
- Fiber v2
- PostgreSQL 17
- pgx/v5
- JWT (golang-jwt)
- Docker Compose

## โครงสร้างโปรเจกต์

```text
cmd/
  main.go
internal/
  config/
    env.go
  database/
    postgres.go
  handlers/
    auth_handler.go
    employee_handler.go
    health_handler.go
    menu_handler.go
    role_handler.go
    table_handler.go
    table_session_handler.go
    qr_session_handler.go
    order_handler.go
    category_handler.go
  middleware/
    auth_middleware.go
  models/
    auth.go
    common.go
    employee.go
    menu.go
    role.go
    table.go
    table_session.go
    qr_session.go
    order.go
    category.go
  repositories/
    auth_repository.go
    employee_repository.go
    menu_repository.go
    role_repository.go
    table_repository.go
    table_session_repository.go
    qr_session_repository.go
    order_repository.go
    category_repository.go
  routes/
    routes.go
  services/
    auth_service.go
    employee_service.go
    menu_service.go
    role_service.go
    table_service.go
    table_session_service.go
    qr_session_service.go
    order_service.go
    category_service.go
  utils/
    jwt.go
    password.go
seeds/
  01_schema.sql
  02_seed.sql
docker-compose.yml
Dockerfile
go.mod
README.md
README_API_TEST.md
```

## สถาปัตยกรรมโดยย่อ

- `Handler`: รับ HTTP request, validate input, ส่ง response
- `Service`: รวม business logic
- `Repository`: คุยกับฐานข้อมูลโดยตรง
- `Model`: โครงสร้าง request/response/entity
- `Route`: ผูก endpoint กับ handler และ middleware

## สิ่งที่ต้องมีในเครื่อง

- Go 1.25 ขึ้นไป
- Docker Desktop

## การตั้งค่าไฟล์ .env

สร้างไฟล์ `.env` ที่ root โปรเจกต์

```env
APP_NAME=RMS Backend
APP_PORT=8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5435
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=rms
DB_SSLMODE=disable
DB_MAX_CONNS=10

DB_HOST_DOCKER=postgres
DB_PORT_CONTAINER=5432
DB_PORT_HOST=5435

JWT_SECRET=super-secret-rms-key
JWT_EXPIRES_IN_SECONDS=3600
```

ความหมายค่า DB ที่ใช้บ่อย

- `DB_HOST` + `DB_PORT`: ใช้ตอนรัน API บนเครื่อง local
- `DB_HOST_DOCKER` + `DB_PORT_CONTAINER`: ใช้ตอนรัน API ใน Docker
- `DB_PORT_HOST`: พอร์ต DB ที่ expose ออกมาให้ local ต่อเข้าได้

## วิธีติดตั้งและรัน (Setup)

### 1) รันทั้งระบบด้วย Docker (แนะนำ)

```bash
docker compose up --build -d
```

ระบบจะพร้อมใช้งานที่

- API: `http://localhost:8080`
- PostgreSQL: `localhost:5435`

### 2) รัน Go local + DB ใน Docker

```bash
docker compose up -d postgres
go mod tidy
go run ./cmd/main.go
```

## วิธีเช็กว่า API ทำงาน

```bash
curl http://localhost:8080/health
```

Expected Response:

```json
{
  "success": true,
  "message": "server is running",
  "data": {
    "status": "ok"
  }
}
```

## คำสั่งที่ใช้บ่อย

รันและดูสถานะ container:

```bash
docker compose ps
docker compose logs -f postgres
```

เช็กว่าโค้ด build ได้:

```bash
go build ./...
```

รีเซ็ตฐานข้อมูลทั้งหมด (ลบ volume เดิม):

```bash
docker compose down -v
docker compose up --build -d
```

## Endpoint ที่มีในระบบตอนนี้

### Public

- `GET /health`
- `GET /api/v1/customer/menus?qrToken=xxx`
- `GET /api/v1/customer/orders?qrToken=xxx`
- `POST /api/v1/customer/orders`
- `POST /api/v1/auth/login`
- `GET /api/v1/qr/:token`
- `GET /api/v1/categories`

### ต้องมี Bearer Token

- `GET /api/v1/auth/me`
- `POST /api/v1/auth/logout`

### ADMIN เท่านั้น

- `GET /api/v1/roles`
- `GET /api/v1/dashboard/summary`
- `POST /api/v1/employees`
- `GET /api/v1/employees`
- `GET /api/v1/employees/:employeeId`
- `PATCH /api/v1/employees/:employeeId`
- `PATCH /api/v1/employees/:employeeId/status`
- `POST /api/v1/tables`
- `PATCH /api/v1/tables/:tableId`
- `POST /api/v1/categories`
- `PATCH /api/v1/categories/:categoryId`
- `POST /api/v1/menus`
- `PATCH /api/v1/menus/:menuId`
- `PATCH /api/v1/menus/:menuId/status`

### ADMIN หรือ CASHIER

- `GET /api/v1/tables`
- `GET /api/v1/tables/:tableId`
- `GET /api/v1/tables/:tableId/current-session`
- `GET /api/v1/table-sessions/:sessionId`
- `GET /api/v1/qr-sessions/:qrSessionId`
- `GET /api/v1/menus`
- `GET /api/v1/menus/:menuId`
- `GET /api/v1/table-sessions/:sessionId/orders`

### ADMIN, CASHIER หรือ CHEF

- `GET /api/v1/orders/:orderId`

### CASHIER เท่านั้น

- `POST /api/v1/table-sessions/open`
- `PATCH /api/v1/table-sessions/:sessionId/close`
- `POST /api/v1/qr-sessions`
- `POST /api/v1/orders`

หมายเหตุ: ต้องตั้ง `password_hash` ให้บัญชี seed ก่อนตามขั้นตอนใน [README_API_TEST.md](README_API_TEST.md)

## คู่มือทดสอบ API แบบละเอียด

ดูที่ไฟล์ [README_API_TEST.md](README_API_TEST.md)

ไฟล์นี้จัดเป็นลำดับทดสอบสำหรับ Postman พร้อมตัวอย่าง Method, URL, Headers, Body และ Expected Response ของทุก endpoint
