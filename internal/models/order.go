package models

import "time"

// Order представляет заказ
type Order struct {
	ID         int       `json:"id"`
	UserID     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	OrderDate  time.Time `json:"order_date"`
	Products   []Product `json:"products"`
}
