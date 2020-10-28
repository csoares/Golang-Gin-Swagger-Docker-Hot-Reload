package routes

import (
	"projetoapi/controllers"

	"github.com/gin-gonic/gin"
)

// @Summary Echo the data sent on get
// @Description Echo the data sent though the get request.
// @Accept  json
// @Produce  json
// @Success 200 {array} model.Evaluation
// @Router /echo [get]
// @Failure 404 "Not found"
func EchoRepeat(c *gin.Context) {
	controllers.Echo(c)
}

// @Summary Recupera as avaliações
// @Description Exibe a lista, sem todos os campos, de todas as avaliações
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Success 200 {array} model.Evaluation
// @Router /evaluation [get]
// @Failure 404 "Not found"
func GetAllEvaluation(c *gin.Context) {
	controllers.GetAllEvaluations(c)
}

// @Summary Recupera uma avaliação pelo id
// @Description Exibe os detalhes de uma avaliação pelo ID
// @ID get-evaluation-by-int
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Param id path int true "Evaluation ID"
// @Success 200 {object} model.Evaluation
// @Router /evaluation/{id} [get]
// @Failure 404 "Not found"
func GetEvaluationById(c *gin.Context) {
	controllers.GetEvaluationByID(c)
}

// @Summary Atualiza uma avaliação
// @Description Atualiza uma avaliação sobre a utilização da aplicação
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Param evaluation body model.Evaluation true "Udpdate evaluation"
// @Param id path int true "Evaluation ID"
// @Router /evaluation/{id} [put]
// @Success 200 {object} model.Evaluation
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
func UpdateEvaluation(c *gin.Context) {
	controllers.UpdateEvaluation(c)
}

// @Summary Exclui uma avaliação pelo ID
// @Description Exclui uma avaliação realizada
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Param id path int true "Evaluation ID"
// @Router /evaluation/{id} [delete]
// @Success 200 {object} model.Evaluation
// @Failure 404 "Not found"
func DeleteEvaluation(c *gin.Context) {
	controllers.DeleteEvaluation(c)
}

// @Summary Adicionar uma avaliação
// @Description Cria uma avaliação sobre a utilização da aplicação
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @param Authorization header string true "Token"
// @Param evaluation body model.Evaluation true "Add evaluation"
// @Router /evaluation [post]
// @Success 201 {object} model.Evaluation
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
func AddEvaluation(c *gin.Context) {
	controllers.AddEvaluation(c)
}
