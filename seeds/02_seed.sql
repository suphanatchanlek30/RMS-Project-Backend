INSERT INTO roles (role_name) VALUES
('ADMIN'),
('CASHIER'),
('CHEF')
ON CONFLICT (role_name) DO NOTHING;

INSERT INTO employees (
    employee_name,
    role_id,
    phone_number,
    email,
    hire_date,
    employee_status,
    password_hash
) VALUES
('Admin User', 1, '0811111111', 'admin@rms.com', '2025-01-01', TRUE, '$2a$10$example_admin_hash'),
('Cashier User', 2, '0822222222', 'cashier@rms.com', '2025-01-02', TRUE, '$2a$10$example_cashier_hash'),
('Chef User', 3, '0833333333', 'chef@rms.com', '2025-01-03', TRUE, '$2a$10$example_chef_hash')
ON CONFLICT (email) DO NOTHING;

INSERT INTO restaurant_tables (table_number, capacity, table_status) VALUES
('A01', 4, 'AVAILABLE'),
('A02', 2, 'AVAILABLE'),
('A03', 6, 'AVAILABLE')
ON CONFLICT (table_number) DO NOTHING;

INSERT INTO menu_categories (category_name, description) VALUES
('อาหารจานหลัก', 'เมนูอาหารหลักของร้าน'),
('เครื่องดื่ม', 'เมนูเครื่องดื่ม'),
('ของหวาน', 'เมนูของหวาน')
ON CONFLICT (category_name) DO NOTHING;

INSERT INTO menus (menu_name, category_id, price, description, menu_status) VALUES
('ข้าวผัดกุ้ง', 1, 89.00, 'ข้าวผัดกุ้งสด', TRUE),
('ผัดไทยกุ้งสด', 1, 99.00, 'ผัดไทยกุ้งสดรสเข้มข้น', TRUE),
('น้ำเปล่า', 2, 15.00, 'น้ำดื่มบรรจุขวด', TRUE),
('ชาเย็น', 2, 35.00, 'ชาเย็นหวานมัน', TRUE),
('ไอศกรีมวานิลลา', 3, 45.00, 'ไอศกรีมวานิลลา 1 scoop', TRUE);

INSERT INTO ingredients (ingredient_name, unit, stock_quantity) VALUES
('ข้าวสวย', 'จาน', 100),
('กุ้ง', 'ตัว', 200),
('เส้นผัดไทย', 'ห่อ', 50),
('น้ำเปล่า', 'ขวด', 300),
('ชา', 'ถุง', 40),
('นมข้น', 'กระป๋อง', 30)
ON CONFLICT (ingredient_name) DO NOTHING;

INSERT INTO menu_ingredients (menu_id, ingredient_id, quantity_used) VALUES
(1, 1, 1),
(1, 2, 5),
(2, 2, 4),
(2, 3, 1),
(3, 4, 1),
(4, 5, 1),
(4, 6, 1)
ON CONFLICT (menu_id, ingredient_id) DO NOTHING;