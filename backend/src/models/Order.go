package models

import "time"

type Order struct {
    OrderID     string      `db:"order_id"`
    BuyerID     string      `db:"buyer_id"`
    TotalPrice  float64     `db:"price_total"`
    Location    string      `db:"location"`
    Status      string      `db:"status"`
    ETA         time.Time   `db:"ETA"`
    CreatedAt   time.Time   `db:"created_at"`   // needed for cursor
    OrderItems  []OrderItem
}

type OrderItem struct {
    OrderItemID string  `db:"order_item_id"`
    OrderID     string  `db:"order_id"`
    ProductID   string  `db:"product_id"`   // UUID, matches Product
    Quantity    int     `db:"quantity"`
    Price       float64 `db:"price"`        // snapshot of price at time of order
}