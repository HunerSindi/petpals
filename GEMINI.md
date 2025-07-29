Database Schema

users table
Field	Type	Notes
id	BIGINT (PK)	auto increment
first_name	VARCHAR(100)	
last_name	VARCHAR(100)	
Phone	VARCHAR(50)	
email	VARCHAR(150)	unique
password	VARCHAR(255)	hashed
created_at	TIMESTAMP	default CURRENT_TIMESTAMP

pets table
Field	Type	Notes
id	BIGINT (PK)	auto increment
uuid	CHAR(36)	provided by user
user_id	BIGINT (FK)	references users(id)
name	VARCHAR(100)	
type	VARCHAR(50)	e.g. Dog, Cat
birth_date	DATE	date of birth
created_at	TIMESTAMP	default CURRENT_TIMESTAMP

categories table
Field	Type	Notes
id	BIGINT (PK)	auto increment
name	VARCHAR(100)	e.g. dog, cat
img_url	VARCHAR(255)	optional category image

products table
Field	Type	Notes
id	BIGINT (PK)	auto increment
category_id	BIGINT (FK)	references categories(id)
name	VARCHAR(150)	
description	TEXT	
price	DECIMAL(10,2)	
created_at	TIMESTAMP	default CURRENT_TIMESTAMP

product_images table
Field	Type	Notes
id	BIGINT (PK)	auto increment
product_id	BIGINT (FK)	references products(id)
img_url	VARCHAR(255)	
is_primary	BOOLEAN	default FALSE

clinics table
Field	Type	Notes
id	BIGINT (PK)	auto increment
first_name	VARCHAR(100)	
last_name	VARCHAR(100)	
clinic_name	VARCHAR(150)	
email	VARCHAR(150)	unique
password	VARCHAR(255)	hashed
open_time	TIME	
close_time	TIME	
description	TEXT	
created_at	TIMESTAMP	default CURRENT_TIMESTAMP

clinic_locations table
Field	Type	Notes
id	BIGINT (PK)	auto increment
clinic_id	BIGINT (FK)	references clinics(id)
address	TEXT	
city	VARCHAR(100)	
phone	VARCHAR(20)	

appointments table
Field	Type	Notes
id	BIGINT (PK)	auto increment
user_id	BIGINT (FK)	references users(id)
clinic_id	BIGINT (FK)	references clinics(id)
pet_id	BIGINT (FK)	references pets(id)
appointment_date	DATE	
appointment_time	TIME	
status	ENUM('pending','confirmed','cancelled')	default 'pending'
created_at	TIMESTAMP	default CURRENT_TIMESTAMP

orders table
Field	Type	Notes
id	BIGINT (PK)	auto increment
user_id	BIGINT (FK)	references users(id)
total_amount	DECIMAL(10,2)	
status	ENUM('pending','confirmed','shipped','delivered','cancelled')	default 'pending'
delivery_address	TEXT	
order_date	TIMESTAMP	default CURRENT_TIMESTAMP
delivered_at	TIMESTAMP	nullable

order_items table
Field	Type	Notes
id	BIGINT (PK)	auto increment
order_id	BIGINT (FK)	references orders(id)
product_id	BIGINT (FK)	references products(id)
quantity	INT	
price	DECIMAL(10,2)	price at time of order

user_addresses table
Field	Type	Notes
id	BIGINT (PK)	auto increment
user_id	BIGINT (FK)	references users(id)
address_line1	VARCHAR(255)	
address_line2	VARCHAR(255)	nullable
city	VARCHAR(100)	
state	VARCHAR(100)	
postal_code	VARCHAR(20)	
is_default	BOOLEAN	default FALSE

End point api 
User App APIs
Authentication:
* POST /api/auth/register - User registration
* POST /api/auth/login - User login
* POST /api/auth/logout - User logout
* GET /api/auth/profile - Get user profile
* PUT /api/auth/profile - Update user profile

Home/Categories & Products:
* GET /api/categories - Get all categories
* GET /api/categories/{id}/products - Get products by category ( pagination ) 
* GET /api/products/{id} - Get single product details
* GET /api/products/{id}/images - Get product images

Orders/Accessories:
* POST /api/orders - Create new order
* GET /api/orders - Get user's orders
* GET /api/orders/{id} - Get specific order details
* PUT /api/orders/{id}/cancel - Cancel order

Clinics:
* GET /api/clinics - Get all clinics
* GET /api/clinics/{id} - Get clinic details
* GET /api/clinics/{id}/available-slots - Get available appointment slots
Appointments:
* POST /api/appointments - Book appointment
* GET /api/appointments - Get user's appointments
* PUT /api/appointments/{id} - Update appointment
* DELETE /api/appointments/{id} - Cancel appointment
Pets:
* GET /api/pets - Get user's pets
* POST /api/pets - Add new pet
* GET /api/pets/{id} - Get pet details
* PUT /api/pets/{id} - Update pet
* DELETE /api/pets/{id} - Delete pet
Addresses:
* GET /api/addresses - Get user addresses
* POST /api/addresses - Add new address
* PUT /api/addresses/{id} - Update address
* DELETE /api/addresses/{id} - Delete address
Clinic App APIs
Authentication:
* POST /api/clinic/auth/login - Clinic login
* GET /api/clinic/auth/profile - Get clinic profile
* PUT /api/clinic/auth/profile - Update clinic profile
Appointments Management:
* GET /api/clinic/appointments - Get clinic's appointments
* PUT /api/clinic/appointments/{id}/confirm - Confirm appointment
* PUT /api/clinic/appointments/{id}/cancel - Cancel appointment
* GET /api/clinic/appointments/calendar - Get calendar view
Schedule Management:
* GET /api/clinic/schedule - Get clinic schedule
* PUT /api/clinic/schedule - Update working hours
Admin APIs (For You)
Clinic Management:
* POST /api/admin/clinics - Register new clinic
* GET /api/admin/clinics - Get all clinics
* PUT /api/admin/clinics/{id} - Update clinic
* DELETE /api/admin/clinics/{id} - Delete clinic
Product Management:
* POST /api/admin/products - Add new product
* PUT /api/admin/products/{id} - Update product
* DELETE /api/admin/products/{id} - Delete product
* POST /api/admin/products/{id}/images - Add product images
Order Management:
* GET /api/admin/orders - Get all orders
* PUT /api/admin/orders/{id}/status - Update order status
