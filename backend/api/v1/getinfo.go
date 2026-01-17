package v1

import (
	"NFTmarket/internal/database"
	"NFTmarket/internal/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func GetProfile(c *gin.Context) {
	userid, exists := c.Get("user_id")
	if !exists {
		log.Println("missing user id after authorization")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id after authorization"})
		return
	}
	db := database.GetDB()
	var currentUser user.User
	if err := db.Select("id", "username", "email", "phone").First(&currentUser, userid.(uint)).Error; err != nil {
		log.Printf("error getting user%d from db:%v", userid, err)
		c.JSON(http.StatusNotFound, gin.H{"msg": "cannot find user"})
		return
	}
	safeuser := user.ResponseUser{
		UserID:   currentUser.ID,
		Username: currentUser.Username,
		Email:    currentUser.Email,
		Phone:    currentUser.Phone,
	}
	c.JSON(http.StatusOK, safeuser)
}
