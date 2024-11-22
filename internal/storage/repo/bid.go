package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"tender-backend/model"
	request_model "tender-backend/model/request"
)

type BidService struct {
	db            *gorm.DB
	tenderService *TenderService
	redis         *redis.Client
}

func NewBidService(db *gorm.DB, redisClient *redis.Client) *BidService {
	return &BidService{
		db:            db,
		tenderService: NewTenderService(db, redisClient),
		redis:         redisClient,
	}
}

func (s *BidService) CreateBid(req *request_model.CreateBidReq, tenderID int64, contractorID int64) (*model.Bid, error) {
	if err := s.validateCreateBidRequest(req); err != nil {
		return nil, err
	}

	var tender model.Tender
	if err := s.db.First(&tender, tenderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Tender not found")
		}
		return nil, fmt.Errorf("failed to fetch tender: %s", err.Error())
	}

	if tender.Status != "open" {
		return nil, errors.New("Tender is not open for bids")
	}

	newBid := model.Bid{
		TenderID:     tenderID,
		ContractorID: contractorID,
		Price:        req.Price,
		DeliveryTime: req.DeliveryTime,
		Comments:     req.Comments,
		Status:       "pending",
	}

	if err := s.db.Create(&newBid).Error; err != nil {
		return nil, fmt.Errorf("failed to create bid: %s", err.Error())
	}

	// Clear related cache after creating a bid
	s.clearBidsCache(tenderID)
	return &newBid, nil
}

func (s *BidService) GetBidByID(bidID, tenderID int64) (*model.Bid, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("bid_%d_tender_%d", bidID, tenderID)

	cachedBid, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var bid model.Bid
		if jsonErr := json.Unmarshal([]byte(cachedBid), &bid); jsonErr == nil {
			return &bid, nil
		}
	}

	var bid model.Bid
	if err := s.db.Where("id = ? AND tender_id = ?", bidID, tenderID).First(&bid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bid not found")
		}
		return nil, fmt.Errorf("failed to retrieve bid: %s", err.Error())
	}

	bidJSON, _ := json.Marshal(bid)
	_ = s.redis.Set(ctx, cacheKey, bidJSON, 10*time.Minute).Err()

	return &bid, nil
}

func (s *BidService) GetAllBids(tenderID int64) ([]model.Bid, error) {
	// Ensure the tender exists before fetching bids
	if _, err := s.tenderService.GetTenderById(tenderID); err != nil {
		return nil, err
	}

	ctx := context.Background()
	cacheKey := fmt.Sprintf("bids_tender_%d", tenderID)

	cachedBids, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var bids []model.Bid
		if jsonErr := json.Unmarshal([]byte(cachedBids), &bids); jsonErr == nil {
			return bids, nil
		}
	}

	var bids []model.Bid
	if err := s.db.Where("tender_id = ?", tenderID).Find(&bids).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch bids: %s", err.Error())
	}

	bidsJSON, _ := json.Marshal(bids)
	_ = s.redis.Set(ctx, cacheKey, bidsJSON, 10*time.Minute).Err()

	return bids, nil
}

func (s *BidService) validateCreateBidRequest(req *request_model.CreateBidReq) error {
	if req.DeliveryTime <= 0 {
		return errors.New("Invalid bid data")
	}

	if req.Price <= 0 {
		return errors.New("invalid price")
	}

	return nil
}

func (s *BidService) GetContractorBids(contractorID int64) ([]model.Bid, error) {
	var bids []model.Bid
	if err := s.db.Where("contractor_id = ?", contractorID).Find(&bids).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve bids: %s", err.Error())
	}

	return bids, nil
}

func (s *BidService) DeleteBid(bidID, contractorID int64) error {
	var bid model.Bid
	if err := s.db.Where("id = ? AND contractor_id = ?", bidID, contractorID).First(&bid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("Bid not found or access denied")
		}
		return fmt.Errorf("failed to find bid: %s", err.Error())
	}

	if err := s.db.Delete(&bid).Error; err != nil {
		return fmt.Errorf("failed to delete bid: %s", err.Error())
	}

	// Clear related cache after deleting a bid
	s.clearBidsCache(bid.TenderID)
	return nil
}

func (s *BidService) clearBidsCache(tenderID int64) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("bids_tender_%d", tenderID)
	_ = s.redis.Del(ctx, cacheKey).Err()
}


func (s *BidService) IsBidExists(bidID int64) bool {
	var bid model.Bid
	err := s.db.First(&bid, bidID).Error;
	return err == nil
}

