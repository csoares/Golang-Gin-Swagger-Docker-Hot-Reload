package dto

import "projetoapi/internal/domain"

// CreateEvaluationRequest represents the request body for creating an evaluation
type CreateEvaluationRequest struct {
	Rating int    `json:"Rating" binding:"required,oneof=0 1 2 3 4 5"`
	Note   string `json:"Note"`
}

// UpdateEvaluationRequest represents the request body for updating an evaluation
type UpdateEvaluationRequest struct {
	Rating int    `json:"Rating" binding:"required,oneof=0 1 2 3 4 5"`
	Note   string `json:"Note"`
}

// EvaluationResponse represents the response for an evaluation
type EvaluationResponse struct {
	ID     uint   `json:"id"`
	Rating int    `json:"rating"`
	Note   string `json:"note"`
}

// BatchUpdateRequest represents the request body for batch updating ratings
type BatchUpdateRequest struct {
	Increment int `json:"increment" binding:"required"`
}

// ToEvaluationResponse converts a domain.Evaluation to EvaluationResponse
func ToEvaluationResponse(e domain.Evaluation) EvaluationResponse {
	return EvaluationResponse{
		ID:     e.ID,
		Rating: e.Rating,
		Note:   e.Note,
	}
}

// ToEvaluationResponseList converts a slice of domain.Evaluation to a slice of EvaluationResponse
func ToEvaluationResponseList(evaluations []domain.Evaluation) []EvaluationResponse {
	responses := make([]EvaluationResponse, len(evaluations))
	for i, e := range evaluations {
		responses[i] = ToEvaluationResponse(e)
	}
	return responses
}