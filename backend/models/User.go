package models

import "time"

type User struct {
	UserID 			int 		`db:"user_id"`		 
	FirstName 		string 		`db:"firstname"` 
	LastName 		string		`db:"lastname"`
	Email 			string 		`db:"email"`
	PasswordHash 	string 		`db:"passwordhashed"`
	UserType 		string		`db:"user_type"`
	IsVerified		bool 		`db:"is_verified"`	
	IsAgree 		bool		`db:"is_agree"`
	CreatedAt 		time.Time 	`db:"created_at"`
	Updatedat  		*time.Time	`db:"updated_at"`
}	