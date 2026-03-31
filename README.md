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
		health_handler.go
		menu_handler.go
		table_handler.go
	models/
		common.go
		menu.go
		table.go
	repositories/
		menu_repository.go
		table_repository.go
	routes/
		routes.go
	services/
		menu_service.go
		table_service.go
seeds/
	01_schema.sql
	02_seed.sql
docker-compose.yml
Dockerfile
go.mod
README.md
```

รายละเอียดแต่ละส่วน

- [cmd/main.go](cmd/main.go): จุดเริ่มต้นของแอป โหลดค่า env, สร้างการเชื่อมต่อ DB, ตั้งค่า Fiber middleware และเปิดเซิร์ฟเวอร์
- [internal/config/env.go](internal/config/env.go): จัดการการอ่านค่าจากไฟล์ .env และอ่านค่าตัวแปรแวดล้อมพร้อม fallback
- [internal/database/postgres.go](internal/database/postgres.go): สร้าง pgx connection pool และตรวจสอบการเชื่อมต่อ PostgreSQL
- [internal/handlers/health_handler.go](internal/handlers/health_handler.go): endpoint ตรวจสุขภาพระบบ
- [internal/handlers/menu_handler.go](internal/handlers/menu_handler.go): endpoint ดึงเมนูสำหรับลูกค้า
- [internal/handlers/table_handler.go](internal/handlers/table_handler.go): endpoint ดึงรายการโต๊ะ
- [internal/models/common.go](internal/models/common.go): รูปแบบ response กลางของ API
- [internal/models/menu.go](internal/models/menu.go): โครงสร้างข้อมูลเมนู
- [internal/models/table.go](internal/models/table.go): โครงสร้างข้อมูลโต๊ะ
- [internal/repositories/menu_repository.go](internal/repositories/menu_repository.go): SQL สำหรับดึงเมนูจากฐานข้อมูล
- [internal/repositories/table_repository.go](internal/repositories/table_repository.go): SQL สำหรับดึงข้อมูลโต๊ะจากฐานข้อมูล
- [internal/routes/routes.go](internal/routes/routes.go): รวม route ทั้งระบบ เช่น /health และ /api/v1/*
- [internal/services/menu_service.go](internal/services/menu_service.go): ชั้นบริการของเมนู
- [internal/services/table_service.go](internal/services/table_service.go): ชั้นบริการของโต๊ะ
- [seeds/01_schema.sql](seeds/01_schema.sql): สร้างตารางและดัชนีทั้งหมด
- [seeds/02_seed.sql](seeds/02_seed.sql): ข้อมูลตั้งต้น เช่น role, โต๊ะ, หมวดหมู่เมนู, เมนู
- [docker-compose.yml](docker-compose.yml): ตั้งค่าบริการ postgres และ api สำหรับรันทั้งระบบด้วย Docker
- [Dockerfile](Dockerfile): ขั้นตอน build และ run แอป Go ใน container

## สิ่งที่ต้องมีในเครื่อง

- Go 1.25 ขึ้นไป
- Docker Desktop (แนะนำสำหรับเริ่มต้นเร็ว)

## การตั้งค่า Environment

ไฟล์ .env ที่ root ของโปรเจกต์

```env
APP_NAME=RMS Backend
APP_PORT=8080
APP_ENV=development

DB_HOST=postgres
DB_PORT=5435
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=rms
DB_SSLMODE=disable
```

ถ้ารัน Go แบบ local ให้เปลี่ยน DB_HOST เป็น localhost

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
