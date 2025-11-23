package main

import (
	"time"
)

// users table structure
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
	CartID    int       `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`
}

// items that we sell
type Item struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// carts table
type Cart struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Items     []Item    `gorm:"many2many:cart_items;"` // relation for cart items
}

// final orders
type Order struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CartID    int       `json:"cart_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
