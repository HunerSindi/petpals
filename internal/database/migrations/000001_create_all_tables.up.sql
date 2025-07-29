
CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(100),
  "last_name" varchar(100),
  "phone" varchar(50),
  "email" varchar(150) UNIQUE,
  "password" varchar(255),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "pets" (
  "id" bigserial PRIMARY KEY,
  "uuid" char(36),
  "user_id" bigint REFERENCES users(id),
  "name" varchar(100),
  "type" varchar(50),
  "birth_date" date,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "categories" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(100),
  "img_url" varchar(255)
);

CREATE TABLE "products" (
  "id" bigserial PRIMARY KEY,
  "category_id" bigint REFERENCES categories(id),
  "name" varchar(150),
  "description" text,
  "price" decimal(10,2),
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "product_images" (
  "id" bigserial PRIMARY KEY,
  "product_id" bigint REFERENCES products(id),
  "img_url" varchar(255),
  "is_primary" boolean DEFAULT false
);

CREATE TABLE "clinics" (
  "id" bigserial PRIMARY KEY,
  "first_name" varchar(100),
  "last_name" varchar(100),
  "clinic_name" varchar(150),
  "email" varchar(150) UNIQUE,
  "password" varchar(255),
  "open_time" time,
  "close_time" time,
  "description" text,
  "created_at" timestamp DEFAULT (now())
);

CREATE TABLE "clinic_locations" (
  "id" bigserial PRIMARY KEY,
  "clinic_id" bigint REFERENCES clinics(id),
  "address" text,
  "city" varchar(100),
  "phone" varchar(20)
);

CREATE TYPE appointment_status AS ENUM ('pending', 'confirmed', 'cancelled');
CREATE TABLE "appointments" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "clinic_id" bigint REFERENCES clinics(id),
  "pet_id" bigint REFERENCES pets(id),
  "appointment_date" date,
  "appointment_time" time,
  "status" appointment_status DEFAULT 'pending',
  "created_at" timestamp DEFAULT (now())
);

CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled');
CREATE TABLE "orders" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "total_amount" decimal(10,2),
  "status" order_status DEFAULT 'pending',
  "delivery_address" text,
  "order_date" timestamp DEFAULT (now()),
  "delivered_at" timestamp
);

CREATE TABLE "order_items" (
  "id" bigserial PRIMARY KEY,
  "order_id" bigint REFERENCES orders(id),
  "product_id" bigint REFERENCES products(id),
  "quantity" int,
  "price" decimal(10,2)
);

CREATE TABLE "user_addresses" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint REFERENCES users(id),
  "address_line1" varchar(255),
  "address_line2" varchar(255),
  "city" varchar(100),
  "state" varchar(100),
  "postal_code" varchar(20),
  "is_default" boolean DEFAULT false
);
