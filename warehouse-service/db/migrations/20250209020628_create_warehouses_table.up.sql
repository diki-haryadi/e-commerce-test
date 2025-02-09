CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE warehouses (
    "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    "shop_id"  uuid NULL,
    "name" VARCHAR(255) NOT NULL,
    "status" bool NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);