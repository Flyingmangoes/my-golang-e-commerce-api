package models

import "time"

type UserRefreshToken struct {
	Id				string
	UserID 			string 				 
	Username		string		
	UserType 		string		
	CreatedAt 		time.Time 	
	ExpiredAt		*time.Time	
	Revoked 		bool
}	