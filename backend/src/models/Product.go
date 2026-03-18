package models

type Product struct {
	ProductID 		string	`db:"product_id"`
	Rating 		    float64	`db:"rating"`
	Price 			float64 `db:"price"`	
	ProductName 	string	`db:"product_name"`
	ProductDesc 	string	`db:"product_desc"`
	StoreName 		string	`db:"store_name"`
	ProductPic 		[]byte	`db:"product_pic"`
}