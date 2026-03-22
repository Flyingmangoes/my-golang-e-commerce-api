package models

import "time"

type User struct {
	UserID 			string 		`db:"user_id"`		 
	FirstName 		string 		`db:"firstname"` 
	LastName 		string		`db:"lastname"`
	Username		string		`db:"username"`
	Email 			string 		`db:"email"`
	PasswordHash 	string 		`db:"passwordhashed"`
	UserLocation 	string 		`db:"user_location"`
	UserType 		string		`db:"user_type"`
	IsVerified		bool 		`db:"is_verified"`	
	IsAgree 		bool		`db:"is_agree"`
	CreatedAt 		time.Time 	`db:"created_at"`
	Updatedat  		*time.Time	`db:"updated_at"`
}	