package controllers

import (
	"net/http"
	"projetoapi/model"
	"projetoapi/services"

	"github.com/gin-gonic/gin"
)

func Echo(c *gin.Context) {
	echo := c.Param("echo")

	c.JSON(http.StatusOK, gin.H{
		"echo": echo,
	})
}

func GetAllEvaluations(c *gin.Context) {
	var evaluations []model.Evaluation

	services.Db.Find(&evaluations)

	if len(evaluations) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "None found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": evaluations})
}

func GetEvaluationByID(c *gin.Context) {
	var evaluation model.Evaluation
	id := c.Param("id")

	services.Db.First(&evaluation, id)
	if evaluation.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Evaluation not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": evaluation})
}

func UpdateEvaluation(c *gin.Context) {
	var evaluation model.Evaluation

	id := c.Param("id")
	services.Db.First(&evaluation, id)

	if evaluation.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Evaluation not found!"})
		return
	}

	if err := c.ShouldBindJSON(&evaluation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check request!"})
		return
	}

	services.Db.Save(evaluation)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Update succeeded!"})
}

func AddEvaluation(c *gin.Context) {
	var evaluation model.Evaluation

	if err := c.ShouldBindJSON(&evaluation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check syntax!"})
		return
	}
	services.Db.Save(&evaluation)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Create successful!", "resourceId": evaluation.ID})
}

func DeleteEvaluation(c *gin.Context) {
	var evaluation model.Evaluation

	id := c.Param("id")
	services.Db.First(&evaluation, id)

	if evaluation.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "None found!"})
		return
	}

	services.Db.Delete(&evaluation)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Delete succeeded!"})
}
