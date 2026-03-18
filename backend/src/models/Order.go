package models

import "time"

type Order struct {
	OrderID 	string		`db:"order_id"`
	TotalPrice 	float64		`db:"price_total"`
	CreatedAt 	time.Time	`db:"created_at"`
	Location 	string		`db:"location"`
	OrderItem 	[]OrderItem
}

type OrderItem struct {
	OrderID   string	`db:"order_id"`
    ProductID int 		`db:"product_id"`
    Quantity  int 		`db:"quantity"`
}