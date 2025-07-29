
CREATE TABLE "admin_users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(100) UNIQUE NOT NULL,
  "password" varchar(255) NOT NULL
);
