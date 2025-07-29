
CREATE TABLE "clinic_images" (
  "id" bigserial PRIMARY KEY,
  "clinic_id" bigint REFERENCES clinics(id) ON DELETE CASCADE,
  "img_url" varchar(255) NOT NULL
);
