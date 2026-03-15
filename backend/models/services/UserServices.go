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
	CreateUser(ctx context.Context, firstName, lastName, username, email, passwordhashed, usertype, userlocation string, isagree bool) (*models.User, error)
	UpdateUser(ctx context.Context, id int, newFirstName, newLastName, newPasswordHashed, newEmail, newUsername, newLocation *string) (*models.User, error)
	DeleteUser(ctx context.Context, id int, email string) (error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ValidateUser(ctx context.Context, id int) (*models.User, error)
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore)CreateUser(ctx context.Context, firstName, lastName, username, email, passwordhashed, usertype, userlocation string, isagree bool) (*models.User, error) {
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
		`INSERT INTO mkt_users (user_id, firstname, lastname, username, email, passwordhashed, user_type, user_location, is_agree)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING user_id, firstname, lastname, username, email, user_type, user_location`,
		id, firstName, lastName, username, email, passwordhashed, usertype, userlocation, isagree,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.UserType)
	
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

	func (us *UserStore)UpdateUser(ctx context.Context, id int, newFirstName, newLastName, newPasswordHashed, newEmail, newUsername, newLocation *string) (*models.User, error) {
		updatedUser := &models.User{}
	err := us.db.QueryRowContext(ctx, 
		`UPDATE mkt_users SET 
			firstname		= COALESCE ($1, firstname),
			lastname		= COALESCE ($2, lastname),
			passwordhashed 	= COALESCE ($3, passwordhashed),
			email			= COALESCE ($4, email),
			username		= COALESCE ($5, username),
			user_location  	= COALESCE ($6, user_location)
		WHERE user_id = $5
		RETURNING user_id, firstname, lastname, email`,
		newFirstName, newLastName, newPasswordHashed, newEmail, newUsername,newLocation, id,
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
		`SELECT user_id, firstname, lastname, email, user_location, user_type FROM mkt_users
		WHERE email = $1`,
		email,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.UserLocation, &user.UserType)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, email, user_location, user_type FROM mkt_users 
		WHERE username = $1`,
		username,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.UserLocation, &user.UserType)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) ValidateUser(ctx context.Context, id int) (*models.User, error) {
	return nil, nil
}
