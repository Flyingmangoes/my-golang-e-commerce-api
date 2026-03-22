package models

import (
    "time"
)

type Product struct {
    ProductID       string      `db:"product_id"`
    SellerID        string      `db:"seller_id"`
    ProductName     string      `db:"product_name"`
    ProductDesc     string      `db:"product_desc"`
    StoreName       string      `db:"store_name"`
    ProductPicUrl   string      `db:"product_pic"`   // URLpath
    Price           float64     `db:"price"`
    Rating          float64     `db:"rating"`
    CreatedAt       time.Time   `db:"created_at"`    // needed for cursor
    UpdatedAt       *time.Time  `db:"updated_at"`
}