package services

import (
	"backend/src/models"
	"context"
	"database/sql"
	"fmt"
)

type ProductStoreInterface interface {
	CreateProduct(ctx context.Context, sellerID, name, desc, store, pic string, price float64) (*models.Product, error)
	UpdateProduct(ctx context.Context, prodid string, new_name, new_desc, new_storename, new_pic *string, new_price *float64) (*models.Product, error)
	RemoveProduct(ctx context.Context, ProdID, sellerID string) error
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	GetProductByName(ctx context.Context, name string) (*models.Product, error)
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore{
	return &ProductStore{db: db}
}

func (ps *ProductStore) CreateProduct(ctx context.Context, sellerID, name, desc, store, pic string, price float64) (*models.Product, error) {
	newproduct := &models.Product{}
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_products (seller_id, product_name, product_desc, storename, price, product_pic)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING product_id, seller_id, product_name, product_desc, storename, price ,product_pic, created_at`,
		sellerID, name, desc, store, price, pic,
	).Scan(&newproduct.ProductID, &newproduct.SellerID, &newproduct.ProductName, &newproduct.ProductDesc, &newproduct.StoreName, &newproduct.Price, &newproduct.ProductPicUrl, &newproduct.CreatedAt)

	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return newproduct, nil
}

func (ps *ProductStore) UpdateProduct(ctx context.Context, id string, name, desc, store, pic *string, price *float64)(*models.Product, error) {
	updatedprod := &models.Product{}

	err := ps.db.QueryRowContext(ctx,
		`UPDATE mkt_products SET 
			product_name	= COALESCE ($1, product_name),
			product_desc	= COALESCE ($2, product_desc),
			store_name		= COALESCE ($3, store_name),
			product_pic		= COALESCE ($4, product_pic),
			price			= COALESCE ($5, price),
			updated_at		= NOW()
		WHERE product_id = $6
		RETURNING product_name, product_desc, store_name, product_pic, price, updated_at`,
		name, desc, store, pic, price, id,
	).Scan(&updatedprod.ProductName, &updatedprod.ProductDesc, &updatedprod.StoreName, &updatedprod.ProductPicUrl, &updatedprod.Price, &updatedprod.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return updatedprod, nil
}

func (ps *ProductStore)RemoveProduct(ctx context.Context, ProdID, sellerID string) error {
	results, err := ps.db.ExecContext(ctx,
		`DELETE FROM mkt_products
		WHERE product_id = $1 AND seller_id = $2`,
		ProdID, sellerID,
	)

	if err != nil {
		return err
	}

	rows, err := results.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("product not found or product id not matched")
	}

	return nil
}

func (ps *ProductStore)GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	product := &models.Product{}
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_name, product_desc, seller_id, store_name, product_pic, price, rating, created_at, updated_at FROM mkt_products
		WHERE product_id = $1`,
		id,
	).Scan(&product.ProductName, &product.ProductDesc, &product.SellerID, &product.StoreName, &product.ProductPicUrl, &product.Price, &product.Rating, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductStore)GetProductByName(ctx context.Context, name string) (*models.Product, error) {
	product := &models.Product{}
	err := ps.db.QueryRowContext(ctx,
		`SELECT product_name, product_desc, seller_id, store_name, product_pic, price, rating, created_at, updated_at FROM mkt_products
		WHERE product_name = $1`,
		name,
	).Scan(&product.ProductName, &product.ProductDesc, &product.SellerID, &product.StoreName, &product.ProductPicUrl, &product.Price, &product.Rating, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return product, nil 
}