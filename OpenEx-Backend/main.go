package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
)

var db *gorm.DB

// Models
type User struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	Password       string `gorm:"not null"`
	ContactDetails string `gorm:"not null"`
	Role           string `gorm:"default:'user'"`
	HostelID       uint
	Hostel         Hostel `gorm:"foreignKey:HostelID"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Hostel struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
	HostelID    uint      `gorm:"not null"`
	Hostel      Hostel    `gorm:"foreignKey:HostelID"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	Price       float64
	Status      string    `gorm:"default:'pending'"`
	Type        string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TransactionRequest struct {
	ID            uint      `gorm:"primaryKey"`
	BuyerID       uint      `gorm:"not null"`
	Buyer         User      `gorm:"foreignKey:BuyerID"`
	SellerID      uint      `gorm:"not null"`
	Seller        User      `gorm:"foreignKey:SellerID"`
	ItemID        uint      `gorm:"not null"`
	Item          Item      `gorm:"foreignKey:ItemID"`
	OfferedItemID *uint
	OfferedItem   Item      `gorm:"foreignKey:OfferedItemID"`
	Status        string    `gorm:"default:'pending'"`
	Type          string    `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func main() {
	// Initialize database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Hostel{}, &Item{}, &TransactionRequest{})

	// Set up router
	r := gin.Default()

	// Public routes
	r.POST("/signup", signup)
	r.POST("/login", login)

	// Authenticated routes
	auth := r.Group("/")
	auth.Use(AuthMiddleware())
	{
		auth.GET("/hostels", listHostels)
		auth.GET("/hostels/:id/items", listItemsByHostel)
		auth.POST("/items", createItem)
		auth.GET("/items/:id", getItem)
		auth.POST("/requests", createRequest)
		auth.GET("/requests", listRequests)
		auth.PATCH("/requests/:id/approve", approveRequest)
	}

	// Admin routes
	admin := auth.Group("/admin")
	admin.Use(AdminMiddleware())
	{
		admin.GET("/items", listPendingItems)
		admin.PATCH("/items/:id/approve", approveItem)
		admin.PATCH("/items/:id/reject", rejectItem)
	}

	r.Run(":8080")
}

// Middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		var user User
		db.First(&user, uint(claims["user_id"].(float64)))
		c.Set("user", user)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(User)
		if user.Role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			return
		}
		c.Next()
	}
}

// Handlers
func signup(c *gin.Context) {
	type SignupRequest struct {
		Name           string `json:"name" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Password       string `json:"password" binding:"required,min=6"`
		ContactDetails string `json:"contact_details" binding:"required"`
		HostelID       uint   `json:"hostel_id" binding:"required"`
	}

	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var hostel Hostel
	if err := db.First(&hostel, req.HostelID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Hostel not found"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := User{
		Name:           req.Name,
		Email:          req.Email,
		Password:       string(hashedPassword),
		ContactDetails: req.ContactDetails,
		HostelID:       req.HostelID,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func login(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func listHostels(c *gin.Context) {
	var hostels []Hostel
	db.Find(&hostels)
	c.JSON(http.StatusOK, hostels)
}

func listItemsByHostel(c *gin.Context) {
	var items []Item
	db.Where("hostel_id = ? AND status = ?", c.Param("id"), "approved").Find(&items)
	c.JSON(http.StatusOK, items)
}

func createItem(c *gin.Context) {
	user := c.MustGet("user").(User)

	type ItemRequest struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Price       float64 `json:"price"`
		Type        string  `json:"type" binding:"required,oneof=sell exchange"`
	}

	var req ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := Item{
		UserID:      user.ID,
		HostelID:    user.HostelID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		Type:        req.Type,
	}

	db.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func listPendingItems(c *gin.Context) {
	var items []Item
	db.Where("status = 'pending'").Find(&items)
	c.JSON(http.StatusOK, items)
}

func approveItem(c *gin.Context) {
	var item Item
	db.First(&item, c.Param("id"))
	item.Status = "approved"
	db.Save(&item)
	c.JSON(http.StatusOK, item)
}

func rejectItem(c *gin.Context) {
	var item Item
	db.First(&item, c.Param("id"))
	item.Status = "rejected"
	db.Save(&item)
	c.JSON(http.StatusOK, item)
}

func createRequest(c *gin.Context) {
	user := c.MustGet("user").(User)

	type Request struct {
		ItemID        uint  `json:"item_id" binding:"required"`
		OfferedItemID *uint `json:"offered_item_id"`
		Type          string `json:"type" binding:"required,oneof=buy exchange"`
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item Item
	if err := db.Preload("User").First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if item.Status != "approved" || item.UserID == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	tr := TransactionRequest{
		BuyerID:       user.ID,
		SellerID:      item.User.ID,
		ItemID:        item.ID,
		OfferedItemID: req.OfferedItemID,
		Type:          req.Type,
	}

	db.Create(&tr)
	c.JSON(http.StatusCreated, tr)
}

func listRequests(c *gin.Context) {
	user := c.MustGet("user").(User)
	var requests []TransactionRequest
	db.Preload("Item").Preload("OfferedItem").
		Where("buyer_id = ? OR seller_id = ?", user.ID, user.ID).
		Find(&requests)
	c.JSON(http.StatusOK, requests)
}

func approveRequest(c *gin.Context) {
	user := c.MustGet("user").(User)
	var tr TransactionRequest
	db.Preload("Seller").Preload("Buyer").First(&tr, c.Param("id"))

	if tr.SellerID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	tr.Status = "approved"
	db.Save(&tr)

	c.JSON(http.StatusOK, gin.H{
		"request":        tr,
		"seller_contact": tr.Seller.ContactDetails,
		"buyer_contact":  tr.Buyer.ContactDetails,
	})
}

func getItem(c *gin.Context) {
	var item Item
	db.First(&item, c.Param("id"))
	c.JSON(http.StatusOK, item)
}