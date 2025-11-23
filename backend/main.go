package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// init db connection
func initDB() {
	var err error
	// using sqlite so we dont need installed database server
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{}, &Item{}, &Cart{}, &Order{})
}

// helper to generate random token
func generateToken() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	initDB()

	r := gin.Default()

	// fix cors issue for react
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 1. Create User
	r.POST("/users", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&user)
		c.JSON(http.StatusOK, user)
	})

	// List Users
	r.GET("/users", func(c *gin.Context) {
		var users []User
		db.Find(&users)
		c.JSON(http.StatusOK, users)
	})

	// 2. Login
	r.POST("/users/login", func(c *gin.Context) {
		var loginDetails User
		if err := c.ShouldBindJSON(&loginDetails); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad data"})
			return
		}

		var user User
		result := db.Where("username = ? AND password = ?", loginDetails.Username, loginDetails.Password).First(&user)

		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid creds"})
			return
		}

		// make new token and save
		newToken := generateToken()
		user.Token = newToken
		db.Save(&user)

		c.JSON(http.StatusOK, gin.H{"token": newToken, "user_id": user.ID})
	})

	// Create Item
	r.POST("/items", func(c *gin.Context) {
		var item Item
		c.ShouldBindJSON(&item)
		db.Create(&item)
		c.JSON(http.StatusOK, item)
	})

	// List Items
	r.GET("/items", func(c *gin.Context) {
		var items []Item
		db.Find(&items)
		c.JSON(http.StatusOK, items)
	})

	// 3. Add to Cart
	r.POST("/carts", func(c *gin.Context) {
		token := c.GetHeader("token")
		// find user by token to know who owns cart
		var user User
		if err := db.Where("token = ?", token).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "please login"})
			return
		}

		// getting item id from body
		type CartRequest struct {
			ItemID int `json:"item_id"`
		}
		var req CartRequest
		c.ShouldBindJSON(&req)

		// check if user has cart, if not make one
		var cart Cart
		res := db.Where("user_id = ?", user.ID).First(&cart)
		if res.Error != nil {
			// create new cart
			cart = Cart{UserID: int(user.ID), Name: "My Cart", Status: "active"}
			db.Create(&cart)

			// update user cart_id
			user.CartID = int(cart.ID)
			db.Save(&user)
		}

		// add item to cart
		var item Item
		db.First(&item, req.ItemID)

		// using gorm association to add item
		db.Model(&cart).Association("Items").Append(&item)

		c.JSON(http.StatusOK, gin.H{"message": "item added", "cart_id": cart.ID})
	})

	// List Carts
	r.GET("/carts", func(c *gin.Context) {
		var carts []Cart
		db.Preload("Items").Find(&carts)
		c.JSON(http.StatusOK, carts)
	})

	// 4. Create Order
	r.POST("/orders", func(c *gin.Context) {
		token := c.GetHeader("token")
		var user User
		if err := db.Where("token = ?", token).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		type OrderReq struct {
			CartID int `json:"cart_id"`
		}
		var req OrderReq
		c.ShouldBindJSON(&req)

		// create order
		order := Order{
			CartID: req.CartID,
			UserID: int(user.ID),
		}
		db.Create(&order)

		c.JSON(http.StatusOK, gin.H{"message": "Order successful", "order_id": order.ID})
	})

	// List Orders
	r.GET("/orders", func(c *gin.Context) {
		var orders []Order
		db.Find(&orders)
		c.JSON(http.StatusOK, orders)
	})

	r.Run(":8080")
}
