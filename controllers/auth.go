package controllers

import (
	"net/http"
	"projetoapi/model"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func LoginHandler(c *gin.Context) {
	var creds model.Users
	var usr model.Users

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}
	services.OpenDatabase()
	services.Db.Find(&usr, "username = ?", creds.Username)
	if usr.Username == "" || !CheckPasswordHash(creds.Password, usr.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid User!"})
		return
	}

	token := services.GenerateTokenJWT(creds)

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Access denied!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success!", "token": token})
}

func RegisterHandler(c *gin.Context) {
	var creds model.Users

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}
	services.OpenDatabase()
	hash, _ := HashPassword(creds.Password)

	creds.Password = hash
	result := services.Db.Save(&creds)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusCreated, gin.H{"status": http.StatusOK, "message": "Success!", "User ID": creds.ID})
		return
	}
	c.JSON(http.StatusNotAcceptable, gin.H{"status": http.StatusNotAcceptable, "message": "Cannot be created!"})
}

func RefreshHandler(c *gin.Context) {

	user := model.Users{
		Username: c.GetHeader("username"),
	}

	token := services.GenerateTokenJWT(user)

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Acesso não autorizado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusNoContent, "message": "Token atualizado com sucesso!", "token": token})
}
