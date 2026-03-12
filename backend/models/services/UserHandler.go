package services

import (
	"backend/models"
	"backend/utils"
	"context"
	"database/sql"
	"fmt"
)

// table name
// mkt_users
// mkt_orders
// mkt_products
// mkt_orders_item

type UserStoreInterface interface {
	CreateUser(ctx context.Context, firstName, lastName, email, password string, isagree bool) (*models.User, error)
	UpdateUser(ctx context.Context, id int, newFirstName, newLastName, newPasswordHashed, newEmail *string) (*models.User, error)
	DeleteUser(ctx context.Context, id int, email string) (error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	ValidateUser(ctx, id string) (*models.User, error)
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore)CreateUser(ctx context.Context, firstName, lastName, email, passwordhashed string, isagree bool) (*models.User, error) {
	user := &models.User{}
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	} 

	defer tx.Rollback()

	var id string 

	for {
		id = utils.GenerateId()
		var exists bool
		tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM mkt_users WHERE user_id = $1)`, id).Scan(&exists)
		if !exists {
			break
		}
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO mkt_users (user_id, firstname, lastname, email, passwordhashed, is_agree)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, firstname, lastname, email`,
		id, firstName, lastName, email, passwordhashed, isagree,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email)
	
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore)UpdateUser(ctx context.Context, id int, newFirstName, newLastName, newPasswordHashed, newEmail *string) (*models.User, error) {
	updatedUser := &models.User{}
	err := us.db.QueryRowContext(ctx, 
		`UPDATE mkt_users SET 
			firstname		= COALASCE ($1, firstname),
			lastname		= COALASCE ($2, lastname),
			passwordhashed 	= COALASCE ($3, passwordhashed),
			email			= COALASCE ($4, email)
		WHERE id = $5
		RETURNING user_id, firstname, lastname, email`,
		newFirstName, newLastName, newPasswordHashed, newEmail, id,
	).Scan(&updatedUser.UserID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Email)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (us *UserStore)DeleteUser(ctx context.Context, id int, email string) error {
	result, err := us.db.ExecContext(ctx, 
		`DELETE FROM mkt_users
		WHERE user_id = $1 AND email = $2`,
		id, email,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("user not found or credentials do not match")
	}

	return nil
}

func (us *UserStore)GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, email FROM mkt_users 
		WHERE email = $1`,
		email,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) ValidateUser(ctx, id string) (*models.User, error) {

}
