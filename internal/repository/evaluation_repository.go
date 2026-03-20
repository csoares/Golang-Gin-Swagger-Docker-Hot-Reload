package repository

import (
	"projetoapi/internal/domain"

	"gorm.io/gorm"
)

// EvaluationRepository defines the interface for evaluation data access
type EvaluationRepository interface {
	FindAll() ([]domain.Evaluation, error)
	FindByID(id uint) (*domain.Evaluation, error)
	Create(eval *domain.Evaluation) error
	Update(eval *domain.Evaluation) error
	Delete(id uint) error
	RawQuery(query string, args ...interface{}) (*gorm.DB, error)
	UpdateBatch(increment int) error
	CreateWithAudit(eval *domain.Evaluation) error
}

// evaluationRepository implements EvaluationRepository using GORM
type evaluationRepository struct {
	db *gorm.DB
}

// NewEvaluationRepository creates a new evaluation repository
func NewEvaluationRepository(db *gorm.DB) EvaluationRepository {
	return &evaluationRepository{db: db}
}

func (r *evaluationRepository) FindAll() ([]domain.Evaluation, error) {
	var evaluations []domain.Evaluation
	result := r.db.Find(&evaluations)
	return evaluations, result.Error
}

func (r *evaluationRepository) FindByID(id uint) (*domain.Evaluation, error) {
	var evaluation domain.Evaluation
	result := r.db.First(&evaluation, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &evaluation, nil
}

func (r *evaluationRepository) Create(eval *domain.Evaluation) error {
	return r.db.Create(eval).Error
}

func (r *evaluationRepository) Update(eval *domain.Evaluation) error {
	return r.db.Save(eval).Error
}

func (r *evaluationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Evaluation{}, id).Error
}

func (r *evaluationRepository) RawQuery(query string, args ...interface{}) (*gorm.DB, error) {
	return r.db.Raw(query, args...), nil
}

func (r *evaluationRepository) UpdateBatch(increment int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&domain.Evaluation{}).
			Where("rating + ? <= 5", increment).
			UpdateColumn("rating", gorm.Expr("rating + ?", increment)).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *evaluationRepository) CreateWithAudit(eval *domain.Evaluation) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(eval).Error; err != nil {
			return err
		}
		if err := tx.Exec("INSERT INTO audit_logs (action, resource_id, created_at) VALUES (?, ?, NOW())", "EVALUATION_CREATED", eval.ID).Error; err != nil {
			return err
		}
		return nil
	})
}