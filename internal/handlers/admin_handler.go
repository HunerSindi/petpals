package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"
	"petapp/internal/utils"

	"database/sql"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type AdminCreateClinicRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ClinicName string `json:"clinic_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	OpenTime   string `json:"open_time"`
	CloseTime  string `json:"close_time"`
	Description string `json:"description"`
	Images     []string `json:"images"`
}

type AdminUpdateClinicRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	ClinicName string `json:"clinic_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	OpenTime   string `json:"open_time"`
	CloseTime  string `json:"close_time"`
	Description string `json:"description"`
}

type AdminCreateProductRequest struct {
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Images      []string `json:"images"`
}

type AdminUpdateProductRequest struct {
	CategoryID  int64   `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type AdminAddProductImagesRequest struct {
	ImgUrl    string `json:"img_url"`
	IsPrimary bool   `json:"is_primary"`
}

type AdminUpdateOrderStatusRequest struct {
	Status string `json:"status"`
}

// Clinic Management
func AdminCreateClinic(c *fiber.Ctx) error {
	req := new(AdminCreateClinicRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not hash password",
			"error":   err.Error(),
		})
	}

	openTime, err := time.Parse("15:04", req.OpenTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid open time format",
			"error":   err.Error(),
		})
	}

	closeTime, err := time.Parse("15:04", req.CloseTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid close time format",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())

	tx, err := database.DB.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not begin transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	clinic, err := qtx.CreateClinic(c.Context(), db.CreateClinicParams{
		FirstName:  sql.NullString{String: req.FirstName, Valid: true},
		LastName:   sql.NullString{String: req.LastName, Valid: true},
		ClinicName: sql.NullString{String: req.ClinicName, Valid: true},
		Email:      sql.NullString{String: req.Email, Valid: true},
		Password:   sql.NullString{String: string(hashedPassword), Valid: true},
		OpenTime:   sql.NullTime{Time: openTime, Valid: true},
		CloseTime:  sql.NullTime{Time: closeTime, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: req.Description != ""},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not register clinic",
			"error":   err.Error(),
		})
	}

	for _, imgURL := range req.Images {
		if imgURL != "" { // Only add non-empty image URLs
			_, err := qtx.CreateClinicImage(c.Context(), db.CreateClinicImageParams{
				ClinicID: sql.NullInt64{Int64: clinic.ID, Valid: true},
				ImgUrl:   imgURL,
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Could not add clinic image",
					"error":   err.Error(),
				})
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not commit transaction",
			"error":   err.Error(),
		})
	}

	clinicResponse := models.ClinicResponse{
		ID:         clinic.ID,
		FirstName:  clinic.FirstName.String,
		LastName:   clinic.LastName.String,
		ClinicName: clinic.ClinicName.String,
		Email:      clinic.Email.String,
		OpenTime:   clinic.OpenTime.Time,
		CloseTime:  clinic.CloseTime.Time,
		Description: clinic.Description.String,
		CreatedAt:  clinic.CreatedAt.Time,
		Images:     req.Images,
	}

	return c.Status(fiber.StatusCreated).JSON(clinicResponse)
}

func AdminGetClinics(c *fiber.Ctx) error {
	queries := db.New(database.DB.DB())
	clinics, err := queries.ListClinics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get clinics",
			"error":   err.Error(),
		})
	}

	clinicResponses := make([]models.ClinicResponse, len(clinics))
	for i, clinic := range clinics {
		images, err := queries.ListClinicImagesByClinicID(c.Context(), sql.NullInt64{Int64: clinic.ID, Valid: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not get clinic images",
				"error":   err.Error(),
			})
		}

		imgUrls := make([]string, 0) // Initialize with capacity 0
		for _, img := range images {
			if img.ImgUrl != "" {
				imgUrls = append(imgUrls, img.ImgUrl)
			}
		}

		clinicResponses[i] = models.ClinicResponse{
			ID:         clinic.ID,
			FirstName:  clinic.FirstName.String,
			LastName:   clinic.LastName.String,
			ClinicName: clinic.ClinicName.String,
			Email:      clinic.Email.String,
			OpenTime:   clinic.OpenTime.Time,
			CloseTime:  clinic.CloseTime.Time,
			Description: clinic.Description.String,
			CreatedAt:  clinic.CreatedAt.Time,
			Images:     imgUrls,
		}
	}

	return c.JSON(clinicResponses)
}

func AdminUpdateClinic(c *fiber.Ctx) error {
	clinicID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid clinic ID",
			"error":   err.Error(),
		})
	}

	req := new(AdminUpdateClinicRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateClinicParams{
		ID: clinicID,
	}

	if req.FirstName != "" {
		params.FirstName = sql.NullString{String: req.FirstName, Valid: true}
	}
	if req.LastName != "" {
		params.LastName = sql.NullString{String: req.LastName, Valid: true}
	}
	if req.ClinicName != "" {
		params.ClinicName = sql.NullString{String: req.ClinicName, Valid: true}
	}
	if req.Email != "" {
		params.Email = sql.NullString{String: req.Email, Valid: true}
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not hash password",
				"error":   err.Error(),
			})
		}
		params.Password = sql.NullString{String: string(hashedPassword), Valid: true}
	}
	if req.OpenTime != "" {
		openTime, err := time.Parse("15:04", req.OpenTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid open time format",
				"error":   err.Error(),
			})
		}
		params.OpenTime = sql.NullTime{Time: openTime, Valid: true}
	}
	if req.CloseTime != "" {
		closeTime, err := time.Parse("15:04", req.CloseTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid close time format",
				"error":   err.Error(),
			})
		}
		params.CloseTime = sql.NullTime{Time: closeTime, Valid: true}
	}
	if req.Description != "" {
		params.Description = sql.NullString{String: req.Description, Valid: true}
	}

	updatedClinic, err := queries.UpdateClinic(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update clinic",
			"error":   err.Error(),
		})
	}

	clinicResponse := models.ClinicResponse{
		ID:         updatedClinic.ID,
		FirstName:  updatedClinic.FirstName.String,
		LastName:   updatedClinic.LastName.String,
		ClinicName: updatedClinic.ClinicName.String,
		Email:      updatedClinic.Email.String,
		OpenTime:   updatedClinic.OpenTime.Time,
		CloseTime:  updatedClinic.CloseTime.Time,
		Description: updatedClinic.Description.String,
		CreatedAt:  updatedClinic.CreatedAt.Time,
	}

	return c.JSON(clinicResponse)
}

func AdminDeleteClinic(c *fiber.Ctx) error {
	clinicID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid clinic ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	err = queries.DeleteClinic(c.Context(), clinicID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete clinic",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Product Management
func AdminCreateProduct(c *fiber.Ctx) error {
	req := new(AdminCreateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())

	tx, err := database.DB.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not begin transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback()

	qtx := queries.WithTx(tx)

	product, err := qtx.CreateProduct(c.Context(), db.CreateProductParams{
		CategoryID:  sql.NullInt64{Int64: req.CategoryID, Valid: true},
		Name:        sql.NullString{String: req.Name, Valid: true},
		Description: sql.NullString{String: req.Description, Valid: true},
		Price:       sql.NullString{String: strconv.FormatFloat(req.Price, 'f', 2, 64), Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not add product",
			"error":   err.Error(),
		})
	}

	for _, imgURL := range req.Images {
		if imgURL != "" { // Only add non-empty image URLs
			_, err := qtx.CreateProductImage(c.Context(), db.CreateProductImageParams{
				ProductID: sql.NullInt64{Int64: product.ID, Valid: true},
				ImgUrl:    sql.NullString{String: imgURL, Valid: true},
				IsPrimary: sql.NullBool{Bool: false, Valid: true},
			})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Could not add product image",
					"error":   err.Error(),
				})
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not commit transaction",
			"error":   err.Error(),
		})
	}

	productResponse := models.ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID.Int64,
		Name:        product.Name.String,
		Description: product.Description.String,
		Price:       utils.ParseFloat(product.Price.String),
		CreatedAt:   product.CreatedAt.Time,
		Images:      req.Images,
	}

	return c.Status(fiber.StatusCreated).JSON(productResponse)
}

func AdminUpdateProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	req := new(AdminUpdateProductRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	params := db.UpdateProductParams{
		ID: productID,
	}

	if req.CategoryID != 0 {
		params.CategoryID = sql.NullInt64{Int64: req.CategoryID, Valid: true}
	}
	if req.Name != "" {
		params.Name = sql.NullString{String: req.Name, Valid: true}
	}
	if req.Description != "" {
		params.Description = sql.NullString{String: req.Description, Valid: true}
	}
	if req.Price != 0 {
		params.Price = sql.NullString{String: strconv.FormatFloat(req.Price, 'f', 2, 64), Valid: true}
	}

	updatedProduct, err := queries.UpdateProduct(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update product",
			"error":   err.Error(),
		})
	}

	productResponse := models.ProductResponse{
		ID:          updatedProduct.ID,
		CategoryID:  updatedProduct.CategoryID.Int64,
		Name:        updatedProduct.Name.String,
		Description: updatedProduct.Description.String,
		Price:       utils.ParseFloat(updatedProduct.Price.String),
		CreatedAt:   updatedProduct.CreatedAt.Time,
	}

	return c.JSON(productResponse)
}

func AdminDeleteProduct(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	err = queries.DeleteProduct(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not delete product",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func AdminAddProductImages(c *fiber.Ctx) error {
	productID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid product ID",
			"error":   err.Error(),
		})
	}

	req := new(AdminAddProductImagesRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	image, err := queries.CreateProductImage(c.Context(), db.CreateProductImageParams{
		ProductID: sql.NullInt64{Int64: productID, Valid: true},
		ImgUrl:    sql.NullString{String: req.ImgUrl, Valid: true},
		IsPrimary: sql.NullBool{Bool: req.IsPrimary, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not add product image",
			"error":   err.Error(),
		})
	}

	imageResponse := models.ProductImageResponse{
		ID:        image.ID,
		ProductID: image.ProductID.Int64,
		ImgUrl:    image.ImgUrl.String,
		IsPrimary: image.IsPrimary.Bool,
	}

	return c.Status(fiber.StatusCreated).JSON(imageResponse)
}

// Order Management
func AdminGetOrders(c *fiber.Ctx) error {
	queries := db.New(database.DB.DB())
	orders, err := queries.ListAllOrders(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get all orders",
			"error":   err.Error(),
		})
	}

	orderResponses := make([]models.OrderResponse, len(orders))
	for i, order := range orders {
		orderResponses[i] = models.OrderResponse{
			ID:              order.ID,
			UserID:          order.UserID.Int64,
			TotalAmount:     utils.ParseFloat(order.TotalAmount.String),
			Status:          models.OrderStatus(order.Status.OrderStatus),
			DeliveryAddress: order.DeliveryAddress.String,
			DeliveredAt:     &order.DeliveredAt.Time,
			OrderDate:       order.OrderDate.Time,
		}
	}

	return c.JSON(orderResponses)
}

func AdminUpdateOrderStatus(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
			"error":   err.Error(),
		})
	}

	req := new(AdminUpdateOrderStatusRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	order, err := queries.UpdateOrderStatus(c.Context(), db.UpdateOrderStatusParams{
		ID:     orderID,
		Status: db.NullOrderStatus{OrderStatus: db.OrderStatus(req.Status), Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update order status",
			"error":   err.Error(),
		})
	}

	orderResponse := models.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID.Int64,
		TotalAmount:     utils.ParseFloat(order.TotalAmount.String),
		Status:          models.OrderStatus(order.Status.OrderStatus),
		DeliveryAddress: order.DeliveryAddress.String,
		DeliveredAt:     &order.DeliveredAt.Time,
		OrderDate:       order.OrderDate.Time,
	}

	return c.JSON(orderResponse)
}