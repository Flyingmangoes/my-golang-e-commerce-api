package services

import (
	"backend/models"
	"database/sql"
	"context"
)

type ProductStoreInterface interface {
	CreateProduct(ctx context.Context, productname, product_desc, storename string, productpic *[]byte) (*models.Product, error)
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

func (ps *ProductStore)CreateProduct(ctx context.Context, productname, product_desc, storename string, productpic *[]byte) (*models.Product, error) {
	return nil, nil
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