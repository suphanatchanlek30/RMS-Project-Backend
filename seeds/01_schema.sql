CREATE TABLE roles (
    role_id SERIAL PRIMARY KEY,
    role_name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE employees (
    employee_id SERIAL PRIMARY KEY,
    employee_name VARCHAR(255) NOT NULL,
    role_id INT NOT NULL REFERENCES roles(role_id),
    phone_number VARCHAR(20),
    email VARCHAR(255) NOT NULL UNIQUE,
    hire_date DATE NOT NULL,
    employee_status BOOLEAN NOT NULL DEFAULT TRUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE restaurant_tables (
    table_id SERIAL PRIMARY KEY,
    table_number VARCHAR(20) NOT NULL UNIQUE,
    capacity INT NOT NULL CHECK (capacity > 0),
    table_status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE'
        CHECK (table_status IN ('AVAILABLE', 'OCCUPIED', 'RESERVED', 'UNAVAILABLE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE table_sessions (
    session_id SERIAL PRIMARY KEY,
    table_id INT NOT NULL REFERENCES restaurant_tables(table_id),
    start_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMPTZ NULL,
    session_details TEXT,
    session_status VARCHAR(20) NOT NULL DEFAULT 'OPEN'
        CHECK (session_status IN ('OPEN', 'CLOSED'))
);

CREATE TABLE qr_sessions (
    qr_session_id SERIAL PRIMARY KEY,
    session_id INT NOT NULL UNIQUE REFERENCES table_sessions(session_id),
    qr_code_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE menu_categories (
    category_id SERIAL PRIMARY KEY,
    category_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menus (
    menu_id SERIAL PRIMARY KEY,
    menu_name VARCHAR(255) NOT NULL,
    category_id INT NOT NULL REFERENCES menu_categories(category_id),
    price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
    description TEXT,
    menu_status BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ingredients (
    ingredient_id SERIAL PRIMARY KEY,
    ingredient_name VARCHAR(255) NOT NULL UNIQUE,
    unit VARCHAR(50) NOT NULL,
    stock_quantity NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE menu_ingredients (
    menu_ingredient_id SERIAL PRIMARY KEY,
    menu_id INT NOT NULL REFERENCES menus(menu_id) ON DELETE CASCADE,
    ingredient_id INT NOT NULL REFERENCES ingredients(ingredient_id) ON DELETE CASCADE,
    quantity_used NUMERIC(10,2) NOT NULL CHECK (quantity_used > 0),
    UNIQUE (menu_id, ingredient_id)
);

CREATE TABLE customer_orders (
    order_id SERIAL PRIMARY KEY,
    session_id INT NOT NULL REFERENCES table_sessions(session_id),
    table_id INT NOT NULL REFERENCES restaurant_tables(table_id),
    order_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    order_status VARCHAR(20) NOT NULL DEFAULT 'PENDING'
        CHECK (order_status IN ('PENDING', 'PREPARING', 'COMPLETED', 'CANCELLED')),
    created_by_employee_id INT NULL REFERENCES employees(employee_id)
);

CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES customer_orders(order_id) ON DELETE CASCADE,
    menu_id INT NOT NULL REFERENCES menus(menu_id),
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price NUMERIC(10,2) NOT NULL CHECK (unit_price >= 0),
    item_status VARCHAR(20) NOT NULL DEFAULT 'WAITING'
        CHECK (item_status IN ('WAITING', 'PREPARING', 'COMPLETED', 'CANCELLED'))
);

CREATE TABLE order_status_history (
    status_history_id SERIAL PRIMARY KEY,
    order_item_id INT NOT NULL REFERENCES order_items(order_item_id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('WAITING', 'PREPARING', 'COMPLETED', 'CANCELLED')),
    updated_by_chef_id INT NULL REFERENCES employees(employee_id),
    updated_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payment_methods (
    payment_method_id SERIAL PRIMARY KEY,
    method_name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE payments (
    payment_id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES customer_orders(order_id),
    payment_method_id INT NOT NULL REFERENCES payment_methods(payment_method_id),
    total_amount NUMERIC(10,2) NOT NULL CHECK (total_amount >= 0),
    payment_time TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'PAID'
        CHECK (payment_status IN ('PENDING', 'PAID', 'FAILED', 'REFUNDED'))
);

CREATE TABLE receipts (
    receipt_id SERIAL PRIMARY KEY,
    payment_id INT NOT NULL UNIQUE REFERENCES payments(payment_id),
    receipt_number VARCHAR(100) NOT NULL UNIQUE,
    issue_date TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    total_amount NUMERIC(10,2) NOT NULL CHECK (total_amount >= 0)
);

CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_restaurant_tables_status ON restaurant_tables(table_status);
CREATE INDEX idx_table_sessions_table_id ON table_sessions(table_id);
CREATE INDEX idx_customer_orders_session_id ON customer_orders(session_id);
CREATE INDEX idx_customer_orders_table_id ON customer_orders(table_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_status ON order_items(item_status);
CREATE INDEX idx_payments_order_id ON payments(order_id);