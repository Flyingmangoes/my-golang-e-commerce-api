package services

import (
	"backend/models"
	"database/sql"
	"context"
)

type ProductStoreInterface interface {
	CreateProduct(ctx context.Context, productname, product_desc, storename string, price float64, productpic *[]byte) (*models.Product, error)
	UpdateProduct(ctx context.Context, productid string) (*models.Product, error)
	RemoveProduct(ctx context.Context, productid string) error
	GetProductByID(ctx context.Context, productid string) (*models.Product, error)
	GetProductByName(ctx context.Context, productname string) (*models.Product, error)
}

type ProductStore struct {
	db *sql.DB
}

func NewProductStore(db *sql.DB) *ProductStore{
	return &ProductStore{db: db}
}

func (ps *ProductStore)CreateProduct(ctx context.Context, productname, product_desc, storename string, price float64,productpic *[]byte) (*models.Product, error) {
	product := &models.Product{}
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_products (product_id, product_name, product_desc, storename, price, product_pic)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING product_id, product_name, product_desc, storename, price ,product_pic`,
		productname, product_desc, storename, price, productpic,
	).Scan(&product.ProductName, &product.ProductDesc, &product.StoreName, &product.Price, &product.ProductPic)

	if err != nil {
		return nil, err
	}
	
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return product, nil
}

func (ps *ProductStore)UpdateProduct(ctx context.Context, productid string) (*models.Product, error) {
	return nil, nil
}

func (ps *ProductStore)RemoveProduct(ctx context.Context, productid string) error {
	return nil
}

func (ps *ProductStore)GetProductByID(ctx context.Context, productid string) (*models.Product, error) {
	return nil, nil
}

func (ps *ProductStore)GetProductByName(ctx context.Context, productname string) (*models.Product, error) {
	return nil, nil
}