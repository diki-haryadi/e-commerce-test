CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE stock_transfer_status  AS ENUM ('pending', 'completed', 'cancelled');
CREATE TABLE stock_transfers (
     "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
     "from_warehouse_id" uuid NOT NULL,
     "to_warehouse_id" uuid NOT NULL,
     "product_id" uuid NOT NULL,
     "quantity" INT NOT NULL,
     "status" stock_transfer_status DEFAULT 'pending',
     "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
CREATE INDEX idx_stock_transfers_status ON stock_transfers(status);