package controllers

import (
	"net/http"
	"projetoapi/model"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	var creds model.Users
	var usr model.Users

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}
	services.OpenDatabase()
	services.Db.Find(&usr, "username = ? and password = ?", creds.Username, creds.Password)
	if usr.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Invalid User!"})
		return
	}

	token := services.GenerateTokenJWT(creds)

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "message": "Access denied!"})
		return
	}
	defer services.Db.Close()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success!", "token": token})
}

func RegisterHandler(c *gin.Context) {
	var creds model.Users

	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Bad request!"})
		return
	}
	services.OpenDatabase()
	services.Db.Save(&creds)

	defer services.Db.Close()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Success!", "User ID": creds.ID})
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

	defer services.Db.Close()
	c.JSON(http.StatusOK, gin.H{"status": http.StatusNoContent, "message": "Token atualizado com sucesso!", "token": token})
}
