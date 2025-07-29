package handlers

import (
	"os"
	"time"

	"petapp/internal/db"
	"petapp/internal/database"
	"petapp/internal/models"

	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
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

	queries := db.New(database.DB.DB())

	user, err := queries.CreateUser(c.Context(), db.CreateUserParams{
		FirstName: sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		LastName:  sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		Phone:     sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Email:     sql.NullString{String: req.Email, Valid: req.Email != ""},
		Password:  sql.NullString{String: string(hashedPassword), Valid: true},
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not create user",
			"error":   err.Error(),
		})
	}

	userResponse := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName.String,
		LastName:  user.LastName.String,
		Phone:     user.Phone.String,
		Email:     user.Email.String,
		CreatedAt: user.CreatedAt.Time,
	}

	return c.Status(fiber.StatusCreated).JSON(userResponse)
}

func Login(c *fiber.Ctx) error {
	req := new(LoginRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())
	user, err := queries.GetUserByEmail(c.Context(), sql.NullString{String: req.Email, Valid: req.Email != ""})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password.String), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid credentials",
		})
	}

	expiresAt := time.Now().Add(time.Hour * 24).Unix()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expiresAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate token",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"token": signedToken})
}

func Logout(c *fiber.Ctx) error {
	// For JWT, logout is typically handled client-side by discarding the token.
	// If using server-side sessions or blacklisting, implement that logic here.
	return c.SendStatus(fiber.StatusOK)
}

func GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	queries := db.New(database.DB.DB())
	profile, err := queries.GetUserByID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not get user profile",
			"error":   err.Error(),
		})
	}

	userResponse := models.UserResponse{
		ID:        profile.ID,
		FirstName: profile.FirstName.String,
		LastName:  profile.LastName.String,
		Phone:     profile.Phone.String,
		Email:     profile.Email.String,
		CreatedAt: profile.CreatedAt.Time,
	}

	return c.JSON(userResponse)
}

func UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := int64(claims["user_id"].(float64))

	req := new(RegisterRequest) // Reusing RegisterRequest for simplicity, consider a dedicated UpdateProfileRequest
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	queries := db.New(database.DB.DB())

	params := db.UpdateUserParams{
		ID: userID,
	}

	if req.FirstName != "" {
		params.FirstName = sql.NullString{String: req.FirstName, Valid: true}
	}
	if req.LastName != "" {
		params.LastName = sql.NullString{String: req.LastName, Valid: true}
	}
	if req.Phone != "" {
		params.Phone = sql.NullString{String: req.Phone, Valid: true}
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

	updatedUser, err := queries.UpdateUser(c.Context(), params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not update user profile",
			"error":   err.Error(),
		})
	}

	userResponse := models.UserResponse{
		ID:        updatedUser.ID,
		FirstName: updatedUser.FirstName.String,
		LastName:  updatedUser.LastName.String,
		Phone:     updatedUser.Phone.String,
		Email:     updatedUser.Email.String,
		CreatedAt: updatedUser.CreatedAt.Time,
	}

	return c.JSON(userResponse)
}