package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `json:"username" form:"username" binding:"required,min=3"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
	Email    string `json:"email" form:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" form:"phone" binding:"required"`
}
type ResponseUser struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type LoginUser struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}
