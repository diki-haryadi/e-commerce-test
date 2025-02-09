CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'cancelled', 'completed');
CREATE TABLE orders (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "user_id" uuid  NULL,
    "status" order_status DEFAULT 'pending',
    "total_amount" DECIMAL(15,2) NOT NULL,
    "payment_deadline" TIMESTAMP,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX idx_order_status ON orders(status);