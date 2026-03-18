package controllers

import (
	"backend/src/middlewares"
	"backend/src/services"
	"backend/src/utils"
	"backend/src/validators"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type UserContext struct {
    Users    services.UserStoreInterface
    Products services.ProductStoreInterface
    Orders   services.OrderStoreInterface
}


type CreateUserRequest struct {
    FirstName    string `json:"firstname"    binding:"required"`
    LastName     string `json:"lastname"     binding:"required"`
    Username     string `json:"username"     binding:"required"`
    Email        string `json:"email"        binding:"required,email"`
    Password     string `json:"password"     binding:"required,min=8"`
    UserLocation string `json:"user_location"`
    UserType     string `json:"user_type"    binding:"required"`
    UserAgreed   bool   `json:"user_agreed"  binding:"required"`
}

type UpdateUserRequest struct {
    UserID      	int     `json:"user_id"      binding:"required"`
	Password 		string 	`json:"password"`
    NewFirstname 	*string `json:"newfirstname" binding:"omitempty"`
    NewLastname  	*string `json:"newlastname"  binding:"omitempty"`
    NewPassword  	*string `json:"newpassword"  binding:"omitempty,min=8"`
   	NewEmail     	*string `json:"newemail"     binding:"omitempty,email"`
	NewUsername  	*string `json:"newusername"  binding:"omitempty,min=3"`
    NewLocation  	*string `json:"newlocation"`
}

type RemoveUserRequest struct {
	UserID      	int     `json:"user_id"      binding:"required"`
	Email 			string  `json:"email"        binding:"required,email"`
	Password     	string  `json:"password"     binding:"required"`
}

type GetUserRequest struct {
	UserID		*int 	`json:"user_id"  binding:"omitempty"`
	Username 	*string	`json:"username" binding:"omitempty"`
	Email		*string `json:"email"    binding:"omitempty,email"`
	Password 	 string	`json:"password" binding:"required"`
}

func (uc *UserContext)Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		hashedpass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.Error(middlewares.ErrInternal("Failed to hash password"))
			return
		}

		user, err := uc.Users.CreateUser(
			c.Request.Context(), 
			req.FirstName, req.LastName, req.Username,
			req.Email, string(hashedpass), req.UserType, 
			req.UserLocation, req.UserAgreed,
		)
		if err != nil {
			slog.Error("[DEBUG]", "error", err)
			var PgErr *pq.Error
			if errors.As(err, &PgErr) && PgErr.Code == "23505" {
				c.Error(middlewares.ErrConflict("User already exists"))
				return
			}

			c.Error(middlewares.ErrInternal("Failed to create user"))
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user created": user}) 
	}
}

func (us *UserContext)Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req UpdateUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		//userID := c.GetInt("user_id")
		var newPassword *string = nil
		if req.NewPassword != nil {
			var pw string = *req.NewPassword

			hashedpass, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
			if err != nil {
				c.Error(middlewares.ErrInternal("Failed to hash password"))
				return
			}

			newPassword = utils.Stroptr(string(hashedpass))
		}

		user, err := us.Users.UpdateUser(c.Request.Context(), 
			req.UserID, req.NewFirstname, req.NewLastname,
			newPassword, req.NewEmail, 
			req.NewUsername, req.NewLocation,
		)

		if err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to update user"))
			return
		}

		compares, err := us.Users.GetPassword(c.Request.Context(), req.UserID)
		if err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to compare data"))
			return
		}

		if err := validators.ValidatePassword(compares.PasswordHash, req.Password); err != nil {
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
		}

		c.JSON(http.StatusOK, gin.H{"user updated": user}) 
	}
}

func (us *UserContext)RemoveUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		var req RemoveUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		userID := c.GetInt("user_id")

		hashedpass, err := us.Users.GetPassword(c.Request.Context(), userID)
		if err != nil {
			slog.Error("[DEBUG]", "error", err)
			var PgErr *pq.Error
			if errors.As(err, &PgErr) && PgErr.Code == "20000"{
				slog.Error("[DEBUG]", "error", err)
				c.Error(middlewares.ErrNotFound("Invalid credentials"))
				return
			}

			c.Error(middlewares.ErrInternal("Failed to check password"))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedpass.PasswordHash), []byte(req.Password)); err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
			return
		}

		err = us.Users.DeleteUser(c.Request.Context(), userID, req.Email)
		if err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrInternal("Failed to update user"))
			return
		}
	}
}

func (us *UserContext)Login() gin.HandlerFunc{
	return func(c *gin.Context) {
		var req GetUserRequest
		if err := c.ShouldBindBodyWithJSON(&req); err != nil {
			slog.Error("[DEBUG]", "error", err)
			c.Error(middlewares.ErrBadRequest("Failed to read client request"))
			return
		}

		if req.Email == nil && req.Username == nil {
            c.Error(middlewares.ErrBadRequest("Email or username is required"))
            return
        }

		user, err := us.Users.GetUserByEmail(c.Request.Context(), *req.Email)
        if err != nil {
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }

		if err := validators.ValidatePassword(user.PasswordHash, req.Password); err != nil {
            c.Error(middlewares.ErrUnauthorized("Invalid credentials"))
            return
        }
		
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}