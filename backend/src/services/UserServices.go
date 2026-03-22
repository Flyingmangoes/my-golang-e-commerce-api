package services

import (
	"backend/src/models"
	"backend/src/utils"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type UserStoreInterface interface {
	CreateUser(ctx context.Context, first, last, username, email, hashed, usertype, location string, isagree bool) (*models.User, error)
	UpdateUser(ctx context.Context, id string, first, last, hashedpass, email, username, location *string) (*models.User, error)
	DeleteUser(ctx context.Context, id, email string) (error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByUserType(ctx context.Context, usertype string) (*models.User, error)

	GetPassword(ctx context.Context, id, email, username *string) (*models.User, error)
	ListUsers(ctx context.Context, filter ListUsersFilter) (*utils.Page[*models.User], error)
}

type ListUsersFilter struct {
	UserType *string
	utils.PagFilter
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore)CreateUser(ctx context.Context, first, last, username, email, hashedpass, usertype, location string, isagree bool) (*models.User, error) {
	user := &models.User{}
	tx, err := us.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	} 

	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, 
		`INSERT INTO mkt_users (firstname, lastname, username, email, passwordhashed, user_type, user_location, is_agree)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING user_id, firstname, lastname, username, email, user_type, user_location, is_agree, created_at`,
		first, last, username, email, hashedpass, usertype, location, isagree,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.UserType, &user.UserLocation, &user.IsAgree, &user.CreatedAt)
	
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

	func (us *UserStore)UpdateUser(ctx context.Context, id string, first, last, hashedpass, email, username, location *string) (*models.User, error) {
		updatedUser := &models.User{}
		err := us.db.QueryRowContext(ctx, 
		`UPDATE mkt_users SET 
			firstname		= COALESCE ($1, firstname),
			lastname		= COALESCE ($2, lastname),
			passwordhashed 	= COALESCE ($3, passwordhashed),
			email			= COALESCE ($4, email),
			username		= COALESCE ($5, username),
			user_location  	= COALESCE ($6, user_location),
			updated_at		= NOW()
		WHERE user_id = $7
		RETURNING user_id, firstname, lastname, email`,
		first, last, hashedpass, email, username, location, id,
	).Scan(&updatedUser.UserID, &updatedUser.FirstName, &updatedUser.LastName, &updatedUser.Email)

	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (us *UserStore)DeleteUser(ctx context.Context, id, email string) error {
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
		`SELECT user_id, firstname, lastname, username, email, user_location, user_type FROM mkt_users
		WHERE email = $1`,
		email,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email ,&user.UserLocation, &user.UserType)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, user_location, user_type FROM mkt_users 
		WHERE username = $1`,
		username,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.UserLocation, &user.UserType)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetUserByUserType(ctx context.Context, usertype string) (*models.User, error) {
	user := &models.User{}
	err := us.db.QueryRowContext(ctx,
		`SELECT user_id, firstname, lastname, username, email, user_location, user_type FROM mkt_users 
		WHERE usertype = $1`,
		usertype,
	).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Username, &user.Email, &user.UserLocation, &user.UserType)
	
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) GetPassword(ctx context.Context, id, email, username *string) (*models.User, error) {
	user := &models.User{}

	var err error
	err = us.db.QueryRowContext(ctx,
    	`SELECT passwordhashed FROM mkt_users
    	WHERE 
        	($1::int     IS NULL OR user_id  = $1) AND
        	($2::varchar IS NULL OR email    = $2) AND
        	($3::varchar IS NULL OR username = $3)
    	LIMIT 1`, 
		id, email, username,
	).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserStore) ListUsers(ctx context.Context, filter ListUsersFilter) (*utils.Page[*models.User], error) {
	filter.Normalize()

	createdAt, id := filter.CursorValues()

	rows, err := us.db.QueryContext(ctx, `
		SELECT user_id, firstname, lastname, email, user_type, created_at
        FROM mkt_users
        WHERE
            ($1::varchar    IS NULL OR user_type = $1)
            AND ($2::timestamptz IS NULL OR (created_at, user_id) < ($2, $3))
        ORDER BY created_at DESC, user_id DESC
        LIMIT $4
    `, filter.UserType, createdAt, id, filter.Limit+1)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(
			&u.UserID, &u.FirstName, &u.LastName,
			&u.Email, &u.UserType, &u.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, u)
	}

	return utils.Build(users, filter.Limit, func(u *models.User) (time.Time, string) {
		return u.CreatedAt, u.UserID
	})
}