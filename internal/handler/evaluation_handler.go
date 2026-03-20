package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"projetoapi/internal/dto"
	"projetoapi/internal/service"
)

// EvaluationHandler handles HTTP requests for evaluations
type EvaluationHandler struct {
	service service.EvaluationService
}

// NewEvaluationHandler creates a new evaluation handler
func NewEvaluationHandler(service service.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{service: service}
}

// Echo godoc
// @Summary Echo the data sent on get
// @Description Echo the data sent through the get request.
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /echo [get]
// @Param name query string false "string valid" minlength(1) maxlength(10)
// @Failure 404 "Not found"
func (h *EvaluationHandler) Echo(c *gin.Context) {
	echo := c.Query("name")
	c.JSON(http.StatusOK, gin.H{
		"echo": echo,
	})
}

// GetAll godoc
// @Summary Recupera as avaliações
// @Description Exibe a lista de todas as avaliações
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Success 200 {array} domain.Evaluation
// @Router /evaluation [get]
// @Failure 404 "Not found"
func (h *EvaluationHandler) GetAll(c *gin.Context) {
	evaluations, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to fetch evaluations"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": evaluations})
}

// GetByID godoc
// @Summary Recupera uma avaliação pelo id
// @Description Exibe os detalhes de uma avaliação pelo ID
// @ID get-evaluation-by-int
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param id path int true "Evaluation ID"
// @Success 200 {object} domain.Evaluation
// @Router /evaluation/{id} [get]
// @Failure 404 "Not found"
func (h *EvaluationHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid ID"})
		return
	}

	evaluation, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Evaluation not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": evaluation})
}

// Create godoc
// @Summary Adicionar uma avaliação
// @Description Cria uma avaliação sobre a utilização da aplicação
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param evaluation body dto.CreateEvaluationRequest true "Add evaluation"
// @Router /evaluation [post]
// @Success 201 {object} dto.EvaluationResponse
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
func (h *EvaluationHandler) Create(c *gin.Context) {
	var req dto.CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check syntax!"})
		return
	}

	evaluation, err := h.service.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Failed to create evaluation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Create successful!", "resourceId": evaluation.ID})
}

// Update godoc
// @Summary Atualiza uma avaliação
// @Description Atualiza uma avaliação sobre a utilização da aplicação
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param evaluation body dto.UpdateEvaluationRequest true "Update evaluation"
// @Param id path int true "Evaluation ID"
// @Router /evaluation/{id} [put]
// @Success 200 {object} dto.EvaluationResponse
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
func (h *EvaluationHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid ID"})
		return
	}

	var req dto.UpdateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check request!"})
		return
	}

	evaluation, err := h.service.Update(uint(id), req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Evaluation not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Update succeeded!", "data": evaluation})
}

// Delete godoc
// @Summary Exclui uma avaliação pelo ID
// @Description Exclui uma avaliação realizada
// @ID delete-evaluation-by-int
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param id path int true "Evaluation ID"
// @Router /evaluation/{id} [delete]
// @Success 200 {object} map[string]string
// @Failure 404 "Not found"
func (h *EvaluationHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Evaluation not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Delete succeeded!"})
}

// GetByRawQuery godoc
// @Summary Get evaluations using raw SQL query
// @Description Demonstrates GORM Raw SQL query with parameter binding
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param min_rating query string false "Minimum rating filter"
// @Success 200 {array} domain.Evaluation
// @Router /evaluation/raw [get]
// @Failure 500 "Internal server error"
func (h *EvaluationHandler) GetByRawQuery(c *gin.Context) {
	minRating := c.Query("min_rating")

	evaluations, err := h.service.GetByRawQuery(minRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Query failed!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": evaluations})
}

// UpdateBatch godoc
// @Summary Batch update ratings with transaction
// @Description Demonstrates GORM Transaction handling - batch update
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param increment body dto.BatchUpdateRequest true "Increment value"
// @Success 200 {object} map[string]string
// @Router /evaluation/batch [put]
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
func (h *EvaluationHandler) UpdateBatch(c *gin.Context) {
	var req dto.BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check request!"})
		return
	}

	if err := h.service.UpdateBatch(req.Increment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Transaction failed!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Batch update successful!"})
}

// CreateWithAudit godoc
// @Summary Create evaluation with audit log using transaction
// @Description Demonstrates manual transaction with rollback/commit control
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Token"
// @Param evaluation body dto.CreateEvaluationRequest true "Add evaluation"
// @Success 201 {object} dto.EvaluationResponse
// @Router /evaluation/audit [post]
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
func (h *EvaluationHandler) CreateWithAudit(c *gin.Context) {
	var req dto.CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "message": "Check Syntax!"})
		return
	}

	evaluation, err := h.service.CreateWithAudit(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "message": "Transaction failed!"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Create successful with audit!", "resourceId": evaluation.ID})
}
