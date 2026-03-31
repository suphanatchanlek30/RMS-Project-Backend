# RMS Project Backend

ระบบ Backend สำหรับ Restaurant Management System (RMS) พัฒนาด้วย Go, Fiber และ PostgreSQL

## ภาพรวมโปรเจกต์

โปรเจกต์นี้ใช้แนวทางแยกชั้นการทำงานแบบชัดเจน เพื่อให้ขยายระบบได้ง่ายในอนาคต

- Handler: รับคำขอ HTTP และส่งผลลัพธ์กลับ
- Service: รวมกฎธุรกิจหรือขั้นตอนการทำงาน
- Repository: ติดต่อฐานข้อมูลโดยตรง
- Model: โครงสร้างข้อมูลที่ใช้ในระบบ
- Route: ผูก URL เข้ากับ handler

## โครงสร้างโปรเจกต์ และหน้าที่

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
		health_handler.go
		menu_handler.go
		table_handler.go
	middleware/
		auth_middleware.go
	models/
		auth.go
		common.go
		menu.go
		table.go
	repositories/
		auth_repository.go
		menu_repository.go
		table_repository.go
	routes/
		routes.go
	services/
		auth_service.go
		menu_service.go
		table_service.go
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

รายละเอียดแต่ละส่วน

- [cmd/main.go](cmd/main.go): จุดเริ่มต้นของแอป โหลดค่า env, สร้างการเชื่อมต่อ DB, ตั้งค่า Fiber middleware และเปิดเซิร์ฟเวอร์
- [internal/config/env.go](internal/config/env.go): จัดการการอ่านค่าจากไฟล์ .env และอ่านค่าตัวแปรแวดล้อมพร้อม fallback
- [internal/database/postgres.go](internal/database/postgres.go): สร้าง pgx connection pool และตรวจสอบการเชื่อมต่อ PostgreSQL
- [internal/handlers/auth_handler.go](internal/handlers/auth_handler.go): endpoint สำหรับ login, me และ logout
- [internal/handlers/health_handler.go](internal/handlers/health_handler.go): endpoint ตรวจสุขภาพระบบ
- [internal/handlers/menu_handler.go](internal/handlers/menu_handler.go): endpoint ดึงเมนูสำหรับลูกค้า
- [internal/handlers/table_handler.go](internal/handlers/table_handler.go): endpoint ดึงรายการโต๊ะ
- [internal/middleware/auth_middleware.go](internal/middleware/auth_middleware.go): middleware ตรวจสอบ Bearer token
- [internal/models/auth.go](internal/models/auth.go): โครงสร้างข้อมูลของ auth เช่น login request/response
- [internal/models/common.go](internal/models/common.go): รูปแบบ response กลางของ API
- [internal/models/menu.go](internal/models/menu.go): โครงสร้างข้อมูลเมนู
- [internal/models/table.go](internal/models/table.go): โครงสร้างข้อมูลโต๊ะ
- [internal/repositories/auth_repository.go](internal/repositories/auth_repository.go): SQL สำหรับค้นหาผู้ใช้เพื่อ login และดึงข้อมูลผู้ใช้จาก token
- [internal/repositories/menu_repository.go](internal/repositories/menu_repository.go): SQL สำหรับดึงเมนูจากฐานข้อมูล
- [internal/repositories/table_repository.go](internal/repositories/table_repository.go): SQL สำหรับดึงข้อมูลโต๊ะจากฐานข้อมูล
- [internal/routes/routes.go](internal/routes/routes.go): รวม route ทั้งระบบ เช่น /health และ /api/v1/*
- [internal/services/auth_service.go](internal/services/auth_service.go): ชั้นบริการสำหรับตรวจ password, สร้าง JWT และดึงข้อมูลผู้ใช้
- [internal/services/menu_service.go](internal/services/menu_service.go): ชั้นบริการของเมนู
- [internal/services/table_service.go](internal/services/table_service.go): ชั้นบริการของโต๊ะ
- [internal/utils/jwt.go](internal/utils/jwt.go): utility สำหรับสร้างและตรวจสอบ JWT
- [internal/utils/password.go](internal/utils/password.go): utility สำหรับตรวจ bcrypt hash
- [seeds/01_schema.sql](seeds/01_schema.sql): สร้างตารางและดัชนีทั้งหมด
- [seeds/02_seed.sql](seeds/02_seed.sql): ข้อมูลตั้งต้น เช่น role, โต๊ะ, หมวดหมู่เมนู, เมนู
- [docker-compose.yml](docker-compose.yml): ตั้งค่าบริการ postgres และ api สำหรับรันทั้งระบบด้วย Docker
- [Dockerfile](Dockerfile): ขั้นตอน build และ run แอป Go ใน container
- [README_API_TEST.md](README_API_TEST.md): คู่มือทดสอบ API ทีละเส้น

## สิ่งที่ต้องมีในเครื่อง

- Go 1.25 ขึ้นไป
- Docker Desktop (แนะนำสำหรับเริ่มต้นเร็ว)

## การตั้งค่า Environment

ไฟล์ .env ที่ root ของโปรเจกต์

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

คำอธิบายค่า DB สำคัญ

- DB_HOST + DB_PORT: ใช้ตอนรัน API บนเครื่อง local
- DB_HOST_DOCKER + DB_PORT_CONTAINER: ใช้ตอนรัน API ใน Docker ให้คุยกับ service postgres ภายใน network
- DB_PORT_HOST: พอร์ตที่ expose ออกมาจาก container เพื่อให้เครื่อง local เชื่อมเข้า DB

## วิธีรันโปรเจกต์

แบบ Docker ทั้งระบบ (แนะนำ)

```bash
docker compose up --build
```

เมื่อรันแล้ว

- API: http://localhost:8080
- PostgreSQL: localhost:5435

หมายเหตุ: ไฟล์ในโฟลเดอร์ seeds จะถูกรันอัตโนมัติครั้งแรกที่สร้างฐานข้อมูล

แบบ Local Go และใช้ DB ใน Docker

1. เปิดเฉพาะฐานข้อมูล

```bash
docker compose up -d postgres
```

2. แก้ไฟล์ .env ให้ DB_HOST=localhost

3. รันแอป

```bash
go mod tidy
go run ./cmd/main.go
```

## Endpoint ที่มีตอนนี้

- GET /health
- GET /api/v1/tables
- GET /api/v1/customer/menus
- POST /api/v1/auth/login
- GET /api/v1/auth/me
- POST /api/v1/auth/logout

## ตัวอย่างเช็คระบบ

```bash
curl http://localhost:8080/health
```

ผลลัพธ์ที่คาดหวัง

```json
{
	"success": true,
	"message": "server is running",
	"data": {
		"status": "ok"
	}
}
```

## หมายเหตุเพิ่มเติม

- ถ้าปรับ schema หรือ seed แล้วต้องการเริ่มใหม่ทั้งหมด ให้ลบ volume เดิมก่อน

```bash
docker compose down -v
docker compose up --build
```

- CORS ปัจจุบันยังเปิดกว้างสำหรับการพัฒนา ควรจำกัด origin ก่อนนำขึ้น production

- สำหรับวิธีทดสอบ API แบบละเอียด ดูไฟล์ [README_API_TEST.md](README_API_TEST.md)
