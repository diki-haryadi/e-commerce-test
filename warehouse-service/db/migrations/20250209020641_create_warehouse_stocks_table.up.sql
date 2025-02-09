CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE warehouse_stocks (
      "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
      "warehouse_id" uuid NOT NULL,
      "product_id" uuid NOT NULL,
      "quantity" INT NOT NULL DEFAULT 0,
      "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
      "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);