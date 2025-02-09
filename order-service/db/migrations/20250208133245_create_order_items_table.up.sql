CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE stock_status AS ENUM ('reserved', 'deducted', 'released');
CREATE TABLE order_items (
     "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
     "order_id" uuid  NULL,
     "product_id" uuid   NULL,
     "warehouse_id" uuid   NULL,
     "quantity" INT NOT NULL,
     "price" DECIMAL(15,2) NOT NULL,
     "stock_status" stock_status DEFAULT 'released',
     "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE INDEX idx_order_items_status ON order_items(stock_status);