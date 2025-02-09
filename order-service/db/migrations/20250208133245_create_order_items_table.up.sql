CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE order_items (
     "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
     order_id uuid  NULL,
     product_id uuid   NULL,
     warehouse_id uuid   NULL,
     quantity INT NOT NULL,
     price DECIMAL(15,2) NOT NULL,
     stock_status ENUM('reserved', 'deducted', 'released') DEFAULT 'reserved',
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
);
CREATE INDEX idx_order_items_status ON order_items(stock_status);