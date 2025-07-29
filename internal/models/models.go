
package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Pet struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

type PetResponse struct {
	ID        int64     `json:"id"`
	UUID      uuid.UUID `json:"uuid"`
	UserID    int64     `json:"user_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	ImgUrl  string `json:"img_url"`
}

type CategoryResponse struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	ImgUrl  string `json:"img_url"`
}

type Product struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProductResponse struct {
	ID          int64     `json:"id"`
	CategoryID  int64     `json:"category_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	Images      []string  `json:"images"`
}

type ProductImage struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	ImgUrl    string `json:"img_url"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductImageResponse struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	ImgUrl    string `json:"img_url"`
	IsPrimary bool   `json:"is_primary"`
}

type Clinic struct {
	ID         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ClinicName string    `json:"clinic_name"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	OpenTime   time.Time `json:"open_time"`
	CloseTime  time.Time `json:"close_time"`
	Description string    `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
}

type ClinicResponse struct {
	ID         int64     `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	ClinicName string    `json:"clinic_name"`
	Email      string    `json:"email"`
	OpenTime   time.Time `json:"open_time"`
	CloseTime  time.Time `json:"close_time"`
	Description string    `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
	Images     []string  `json:"images"`
}

type ClinicImage struct {
	ID       int64  `json:"id"`
	ClinicID int64  `json:"clinic_id"`
	ImgUrl   string `json:"img_url"`
}

type ClinicImageResponse struct {
	ID       int64  `json:"id"`
	ClinicID int64  `json:"clinic_id"`
	ImgUrl   string `json:"img_url"`
}

type ClinicLocation struct {
	ID       int64  `json:"id"`
	ClinicID int64  `json:"clinic_id"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Phone    string `json:"phone"`
}

type ClinicLocationResponse struct {
	ID       int64  `json:"id"`
	ClinicID int64  `json:"clinic_id"`
	Address  string `json:"address"`
	City     string `json:"city"`
	Phone    string `json:"phone"`
}

type AppointmentStatus string

const (
	AppointmentStatusPending   AppointmentStatus = "pending"
	AppointmentStatusConfirmed AppointmentStatus = "confirmed"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type Appointment struct {
	ID             int64             `json:"id"`
	UserID         int64             `json:"user_id"`
	ClinicID       int64             `json:"clinic_id"`
	PetID          int64             `json:"pet_id"`
	AppointmentDate time.Time         `json:"appointment_date"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
}

type AppointmentResponse struct {
	ID             int64             `json:"id"`
	UserID         int64             `json:"user_id"`
	ClinicID       int64             `json:"clinic_id"`
	PetID          int64             `json:"pet_id"`
	AppointmentDate time.Time         `json:"appointment_date"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
}

type AppointmentDetailsResponse struct {
	ID             int64             `json:"id"`
	AppointmentDate time.Time         `json:"appointment_date"`
	AppointmentTime time.Time         `json:"appointment_time"`
	Status          AppointmentStatus `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	User            UserResponse      `json:"user"`
	Clinic          ClinicResponse    `json:"clinic"`
	Pet             PetResponse       `json:"pet"`
}

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusShipped  OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)


type Order struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	TotalAmount    float64     `json:"total_amount"`
	Status         OrderStatus `json:"status"`
	DeliveryAddress string      `json:"delivery_address"`
	OrderDate      time.Time   `json:"order_date"`
	DeliveredAt    *time.Time  `json:"delivered_at"` // Nullable
}

type OrderResponse struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	TotalAmount    float64     `json:"total_amount"`
	Status         OrderStatus `json:"status"`
	DeliveryAddress string      `json:"delivery_address"`
	OrderDate      time.Time   `json:"order_date"`
	DeliveredAt    *time.Time  `json:"delivered_at"` // Nullable
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderItemResponse struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
	Product   ProductResponse `json:"product"`
}

type UserAddress struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}

type UserAddressResponse struct {
	ID           int64  `json:"id"`
	UserID       int64  `json:"user_id"`
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	IsDefault    bool   `json:"is_default"`
}
