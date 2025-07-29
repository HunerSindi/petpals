
package handlers

import (
	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"
	"petapp/internal/utils"

	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type CreateOrderRequest struct {
	TotalAmount    float64 `json:"total_amount"`
	DeliveryAddress string  `json:"delivery_address"`
	Items          []struct {
		ProductID int64   `json:"product_id"`
		Quantity  int32   `json:"quantity"`
		Price     float64 `json:"price"`
	} `json:"items"`
}

func CreateOrder(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	req := new(CreateOrderRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())

	// Start a transaction
	tx, err := database.DB.DB().BeginTx(c.Context(), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not begin transaction",
			"error":   err.Error(),
		})
	}
	defer tx.Rollback() // Rollback on error, commit on success

	qtx := queries.WithTx(tx)

	order, err := qtx.CreateOrder(c.Context(), db.CreateOrderParams{
		UserID:          sql.NullInt64{Int64: userID, Valid: true},
		TotalAmount:     sql.NullString{String: strconv.FormatFloat(req.TotalAmount, 'f', 2, 64), Valid: true},
		Status:          db.NullOrderStatus{OrderStatus: db.OrderStatusPending, Valid: true},
		DeliveryAddress: sql.NullString{String: req.DeliveryAddress, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create order",
			"error":   err.Error(),
		})
	}

	for _, item := range req.Items {
		_, err := qtx.CreateOrderItem(c.Context(), db.CreateOrderItemParams{
			OrderID:   sql.NullInt64{Int64: order.ID, Valid: true},
			ProductID: sql.NullInt64{Int64: item.ProductID, Valid: true},
			Quantity:  sql.NullInt32{Int32: item.Quantity, Valid: true},
			Price:     sql.NullString{String: strconv.FormatFloat(item.Price, 'f', 2, 64), Valid: true},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not create order item",
				"error":   err.Error(),
			})
		}
	}

	if err := tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not commit transaction",
			"error":   err.Error(),
		})
	}

	orderResponse := models.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID.Int64,
		TotalAmount:     utils.ParseFloat(order.TotalAmount.String),
		Status:          models.OrderStatus(order.Status.OrderStatus),
		DeliveryAddress: order.DeliveryAddress.String,
		OrderDate:       order.OrderDate.Time,
		DeliveredAt:     &order.DeliveredAt.Time,
	}

	return c.Status(fiber.StatusCreated).JSON(orderResponse)
}

func GetOrders(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	queries := db.New(database.DB.DB())
	orders, err := queries.ListOrdersByUserID(c.Context(), sql.NullInt64{Int64: userID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get user orders",
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
			OrderDate:       order.OrderDate.Time,
			DeliveredAt:     &order.DeliveredAt.Time,
		}
	}

	return c.JSON(orderResponses)
}

func GetOrderDetails(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	order, err := queries.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get order details",
			"error":   err.Error(),
		})
	}

	items, err := queries.ListOrderItemsByOrderID(c.Context(), sql.NullInt64{Int64: orderID, Valid: true})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get order items",
			"error":   err.Error(),
		})
	}

	orderResponse := models.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID.Int64,
		TotalAmount:     utils.ParseFloat(order.TotalAmount.String),
		Status:          models.OrderStatus(order.Status.OrderStatus),
		DeliveryAddress: order.DeliveryAddress.String,
		OrderDate:       order.OrderDate.Time,
		DeliveredAt:     &order.DeliveredAt.Time,
	}

	itemResponses := make([]models.OrderItemResponse, len(items))
	for i, item := range items {
		product, err := queries.GetProduct(c.Context(), item.ProductID.Int64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not get product details for order item",
				"error":   err.Error(),
			})
		}

		productImages, err := queries.ListProductImages(c.Context(), sql.NullInt64{Int64: product.ID, Valid: true})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Could not get product images for order item",
				"error":   err.Error(),
			})
		}

		imgUrls := make([]string, 0)
		for _, img := range productImages {
			if img.ImgUrl.Valid && img.ImgUrl.String != "" {
				imgUrls = append(imgUrls, img.ImgUrl.String)
			}
		}

		itemResponses[i] = models.OrderItemResponse{
			ID:        item.ID,
			OrderID:   item.OrderID.Int64,
			ProductID: item.ProductID.Int64,
			Quantity:  item.Quantity.Int32,
			Price:     utils.ParseFloat(item.Price.String),
			Product: models.ProductResponse{
				ID:          product.ID,
				CategoryID:  product.CategoryID.Int64,
				Name:        product.Name.String,
				Description: product.Description.String,
				Price:       utils.ParseFloat(product.Price.String),
				CreatedAt:   product.CreatedAt.Time,
				Images:      imgUrls,
			},
		}
	}

	return c.JSON(fiber.Map{
		"order": orderResponse,
		"items": itemResponses,
	})
}

func CancelOrder(c *fiber.Ctx) error {
	orderID, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid order ID",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	order, err := queries.GetOrderByID(c.Context(), orderID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not retrieve order",
			"error":   err.Error(),
		})
	}

	if order.Status.OrderStatus != db.OrderStatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Order can only be cancelled if status is pending",
			"current_status": order.Status.OrderStatus,
		})
	}

	updatedOrder, err := queries.UpdateOrderStatus(c.Context(), db.UpdateOrderStatusParams{
		ID:     orderID,
		Status: db.NullOrderStatus{OrderStatus: db.OrderStatusCancelled, Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not cancel order",
			"error":   err.Error(),
		})
	}

	orderResponse := models.OrderResponse{
		ID:              updatedOrder.ID,
		UserID:          updatedOrder.UserID.Int64,
		TotalAmount:     utils.ParseFloat(updatedOrder.TotalAmount.String),
		Status:          models.OrderStatus(updatedOrder.Status.OrderStatus),
		DeliveryAddress: updatedOrder.DeliveryAddress.String,
		OrderDate:       updatedOrder.OrderDate.Time,
		DeliveredAt:     &updatedOrder.DeliveredAt.Time,
	}

	return c.JSON(orderResponse)
}
