package api

import (
	"NFTmarket/internal/auth"
	"NFTmarket/internal/database"
	"NFTmarket/internal/user"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type AuthResponse struct {
	Token string            `json:"token"`
	User  user.ResponseUser `json:"user"`
}

func Register(c *gin.Context) {
	db := database.GetDB()
	var uuser user.User
	if err := c.ShouldBindJSON(&uuser); err != nil {
		log.Println("format bind error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"format bind error": err.Error()})
		return
	}
	if db.Where("username = ?", uuser.Username).First(&user.User{}).Error == nil {
		log.Println("username already exists")
		c.JSON(http.StatusBadRequest, gin.H{"msg": "username exist"})
		return
	}
	db.Create(&uuser)
	token, err := auth.GenerateToken(uuser.ID, uuser.Username)
	if err != nil {
		log.Println("token generate error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	safeUser := user.ResponseUser{UserID: uuser.ID, Username: uuser.Username, Email: uuser.Email}
	c.JSON(http.StatusOK, AuthResponse{Token: token, User: safeUser})
}
func Login(c *gin.Context) {
	db := database.GetDB()
	var reuser user.LoginUser
	if err := c.ShouldBindJSON(&reuser); err != nil {
		log.Println("format bind error:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"format bind error": err.Error()})
	}
	var uuser user.User
	result := db.Where("username = ?", reuser.Username).First(&uuser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			log.Println("user not found")
			c.JSON(http.StatusNotFound, gin.H{
				"error": "该用户未注册",
			})
		} else {
			log.Printf("查询用户失败: %s, 错误: %v", uuser.Username, result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "无法查询用户信息",
			})
		}
		return
	}
	if uuser.Password != reuser.Password {
		log.Println("password not match")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "password error",
		})
	}
	token, err := auth.GenerateToken(uuser.ID, uuser.Username)
	if err != nil {
		log.Println("token generate error:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "token generate error",
		})
		return
	}
	safeUser := user.ResponseUser{UserID: uuser.ID, Username: uuser.Username, Email: uuser.Email}
	c.JSON(http.StatusOK, AuthResponse{Token: token, User: safeUser})
}
