CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
   "id" uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
   "name" VARCHAR(254) UNIQUE NOT NULL,
   "description" text NOT NULL,
   "category_id" uuid DEFAULT uuid_generate_v4(),
   "price" int NOT NULL,
   "stock" int NOT NULL,
   "status" bool NOT NULL,
   "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   "deleted_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
