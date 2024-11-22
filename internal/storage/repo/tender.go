package server

import (
	"context"
	"encoding/json"
	"errors"
	"tender-backend/model"
	request_model "tender-backend/model/request"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type TenderService struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewTenderService initializes a new TenderService with the database connection.
func NewTenderService(db *gorm.DB, redisClient *redis.Client) *TenderService {
	return &TenderService{
		db:    db,
		redis: redisClient,
	}
}

// CreateTender creates a new tender in the database.
func (t *TenderService) CreateTender(req *request_model.CreateTenderReq, clientID int64) (*model.Tender, error) {
	if err := validateCreateTender(req); err != nil {
		return nil, err
	}

	tender := &model.Tender{
		ClientID:    clientID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Budget:      req.Budget,
		Status:      "open",
	}

	// Save the tender to the database.
	if err := t.db.Create(tender).Error; err != nil {
		return nil, err
	}

	// Invalidate the cache after creating a new tender
	t.redis.Del(context.Background(), "tenders_cache")

	return tender, nil
}

// ValidateCreateTender validates the input for creating a tender.
func validateCreateTender(req *request_model.CreateTenderReq) error {
	if req.Title == "" {
		return errors.New("invalid input: title is required")
	}
	if req.Deadline.Before(time.Now()) {
		return errors.New("invalid input: deadline must be in the future")
	}
	if req.Budget <= 0 {
		return errors.New("invalid input: budget must be positive")
	}
	return nil
}

// GetTenderById retrieves a tender by its ID.
func (t *TenderService) GetTenderById(id int64) (*model.Tender, error) {
	var tender model.Tender

	// Try fetching the tender from the database
	if err := t.db.First(&tender, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Tender not found or access denied")
		}
		return nil, err
	}

	return &tender, nil
}

// GetTenders retrieves all tenders from the cache or database.
func (t *TenderService) GetTenders() ([]model.Tender, error) {
	ctx := context.Background()
	cacheKey := "tenders_cache"

	// Try fetching the tenders from Redis
	cachedTenders, err := t.redis.Get(ctx, cacheKey).Result()
	if err == nil && cachedTenders != "" {
		var tenders []model.Tender
		if err := json.Unmarshal([]byte(cachedTenders), &tenders); err == nil {
			return tenders, nil
		}
	}

	// If cache miss or unmarshal error, fetch from the database
	var tenders []model.Tender
	if err := t.db.Find(&tenders).Error; err != nil {
		return nil, err
	}

	// Cache the tenders for 10 minutes
	tendersJSON, err := json.Marshal(tenders)
	if err == nil {
		_ = t.redis.Set(ctx, cacheKey, tendersJSON, 10*time.Minute).Err()
	}

	return tenders, nil
}

// UpdateTender updates the tender with the given ID.
func (t *TenderService) UpdateTender(tenderID, clientID int64, req *request_model.UpdateTenderReq) (*model.Tender, error) {
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tender not found or access denied")
		}
		return nil, err
	}

	if err := ValidateTenderUpdate(tender.Status, req.Status); err != nil {
		return nil, err
	}

	tender.Status = req.Status

	// Save the updated tender to the database
	if err := t.db.Save(&tender).Error; err != nil {
		return nil, err
	}

	// Invalidate the cache after updating the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return &tender, nil
}

func ValidateTenderUpdate(existingStatus, newStatus string) error {
	if existingStatus != "open" {
		return errors.New("updates are only allowed for tenders with 'open' status")
	}
	if newStatus == "awarded" {
		return errors.New("status cannot be updated to 'awarded'")
	}
	return nil
}

// DeleteTender deletes a tender by its ID.
func (t *TenderService) DeleteTender(tenderID, clientID int64) error {
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	if err := t.db.Delete(&model.Tender{}, tenderID).Error; err != nil {
		return err
	}

	// Invalidate the cache after deleting the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return nil
}

// ValidateTenderBelongsToUser ensures that a tender belongs to a specific client.
func (t *TenderService) ValidateTenderBelongsToUser(tenderID, clientID int64) error {
	var tender model.Tender

	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Tender not found or access denied")
		}
		return err
	}

	if tender.ClientID != clientID {
		return errors.New("Tender not found or access denied")
	}

	return nil
}

func (t *TenderService) AwardTender(tenderID, clientID, bidID int64) error {
	if err := t.ValidateTenderBelongsToUser(tenderID, clientID); err != nil {
		return err
	}

	if err := t.ValidateBidBelongsToTender(bidID, tenderID); err != nil {
		return err
	}

	if err := t.db.Model(&model.Tender{}).Where("id = ?", tenderID).Updates(map[string]interface{}{
		"status":                "awarded",
		"awarded_contractor_id": bidID,
	}).Error; err != nil {
		return err
	}

	// Invalidate the cache after awarding the tender
	t.redis.Del(context.Background(), "tenders_cache")

	return nil
}

func (t *TenderService) ValidateBidBelongsToTender(bidID, tenderID int64) error {
	var bid model.Bid

	if err := t.db.First(&bid, bidID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("bid not found or access denied")
		}
		return err
	}

	if bid.TenderID != tenderID {
		return errors.New("bid not found or access denied")
	}

	return nil
}

func (t *TenderService) IsTenderExists(tenderID int64) bool {
	var tender model.Tender
	if err := t.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		return false
	}
	return true
}

