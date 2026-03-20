package service

import (
	"projetoapi/internal/domain"
	"projetoapi/internal/dto"
	"projetoapi/internal/repository"
)

// EvaluationService defines the interface for evaluation business logic
type EvaluationService interface {
	GetAll() ([]dto.EvaluationResponse, error)
	GetByID(id uint) (*dto.EvaluationResponse, error)
	Create(req dto.CreateEvaluationRequest) (*dto.EvaluationResponse, error)
	Update(id uint, req dto.UpdateEvaluationRequest) (*dto.EvaluationResponse, error)
	Delete(id uint) error
	GetByRawQuery(minRating string) ([]dto.EvaluationResponse, error)
	UpdateBatch(increment int) error
	CreateWithAudit(req dto.CreateEvaluationRequest) (*dto.EvaluationResponse, error)
}

// evaluationService implements EvaluationService
type evaluationService struct {
	repo repository.EvaluationRepository
}

// NewEvaluationService creates a new evaluation service
func NewEvaluationService(repo repository.EvaluationRepository) EvaluationService {
	return &evaluationService{repo: repo}
}

func (s *evaluationService) GetAll() ([]dto.EvaluationResponse, error) {
	evaluations, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return dto.ToEvaluationResponseList(evaluations), nil
}

func (s *evaluationService) GetByID(id uint) (*dto.EvaluationResponse, error) {
	eval, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	response := dto.ToEvaluationResponse(*eval)
	return &response, nil
}

func (s *evaluationService) Create(req dto.CreateEvaluationRequest) (*dto.EvaluationResponse, error) {
	eval := &domain.Evaluation{
		Rating: req.Rating,
		Note:   req.Note,
	}
	if err := s.repo.Create(eval); err != nil {
		return nil, err
	}
	response := dto.ToEvaluationResponse(*eval)
	return &response, nil
}

func (s *evaluationService) Update(id uint, req dto.UpdateEvaluationRequest) (*dto.EvaluationResponse, error) {
	eval, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	eval.Rating = req.Rating
	eval.Note = req.Note
	if err := s.repo.Update(eval); err != nil {
		return nil, err
	}
	response := dto.ToEvaluationResponse(*eval)
	return &response, nil
}

func (s *evaluationService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *evaluationService) GetByRawQuery(minRating string) ([]dto.EvaluationResponse, error) {
	if minRating == "" {
		minRating = "3"
	}
	result, err := s.repo.RawQuery("SELECT * FROM evaluations WHERE rating >= ? AND deleted_at IS NULL", minRating)
	if err != nil {
		return nil, err
	}
	var evaluations []domain.Evaluation
	if err := result.Scan(&evaluations).Error; err != nil {
		return nil, err
	}
	return dto.ToEvaluationResponseList(evaluations), nil
}

func (s *evaluationService) UpdateBatch(increment int) error {
	return s.repo.UpdateBatch(increment)
}

func (s *evaluationService) CreateWithAudit(req dto.CreateEvaluationRequest) (*dto.EvaluationResponse, error) {
	eval := &domain.Evaluation{
		Rating: req.Rating,
		Note:   req.Note,
	}
	if err := s.repo.CreateWithAudit(eval); err != nil {
		return nil, err
	}
	response := dto.ToEvaluationResponse(*eval)
	return &response, nil
}
