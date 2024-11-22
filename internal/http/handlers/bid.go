package handlers

import (
	"net/http"
	"strconv"
	request_model "tender-backend/model/request"

	"github.com/gin-gonic/gin"
)

// CreateBid godoc
// @Summary Create a new bid
// @Description Creates a new bid. Example time: 2024-11-16T15:00:00Z
// @Tags Bid
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param bid body request_model.CreateBidReq true "Bid creation request"
// @Success 201 {object} model.Bid "Bid created successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /api/contractor/tenders/{tender_id}/bid [POST]
func (h *HTTPHandler) CreateBid(c *gin.Context) {
	var req request_model.CreateBidReq
	tenderIdStr := c.Param("tender_id")
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	contractorId := c.GetInt64("user_id")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	createdBid, err := h.BidService.CreateBid(&req, int64(tenderId), contractorId)
	if err != nil {
		if err.Error()=="Tender not found"{
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdBid)
}

// GetBidByID godoc
// @Summary Get Bid by ID
// @Description Retrieves a bid by its ID.
// @Tags Bid
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} model.Bid "Bid retrieved successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Bid not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /api/contractor/tenders/{tender_id}/bid/{bid_id} [GET]
func (h *HTTPHandler) GetBid(c *gin.Context) {
	bidIDStr := c.Param("bid_id")
	bidID, err := strconv.Atoi(bidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}
	tenderIDStr := c.Param("tender_id")
	tenderID, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	bid, err := h.BidService.GetBidByID(int64(bidID), int64(tenderID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bid"})
		return
	}

	c.JSON(http.StatusOK, bid)
}

// GetBids godoc
// @Summary Get all Bids
// @Description Retrieves all Bids for the authenticated user.
// @Tags Bid
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Success 200 {object} []model.Bid "All bids retrieved successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /api/client/contractor/tenders/{tender_id}/bids [get]
func (h *HTTPHandler) GetBids(c *gin.Context) {
	tenderIDStr := c.Param("tender_id")
	tenderID, err := strconv.Atoi(tenderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	bids, err := h.BidService.GetAllBids(int64(tenderID))
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// GetContractorBids godoc
// @Summary Get all Bids for a Contractor
// @Description Retrieves all Bids for the authenticated Contractor.
// @Tags Bid
// @Accept json
// @Produce json
// @Success 200 {object} []model.Bid "All bids retrieved successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /api/contractor/bids [get]
func (h *HTTPHandler) GetContractorBids(c *gin.Context) {
	contractorID := c.GetInt64("user_id")
	bids, err := h.BidService.GetContractorBids(contractorID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Failed to retrieve bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}

// DeleteBid godoc
// @Summary Delete a Bid
// @Description Deletes a Bid by its ID.
// @Tags Bid
// @Accept json
// @Produce json
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} string "Bid deleted successfully"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "Bid not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /api/contractor/bids/{bid_id} [DELETE]
func (h *HTTPHandler) DeleteBid(c *gin.Context) {
	bidIDStr := c.Param("bid_id")
	bidID, err := strconv.Atoi(bidIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	err = h.BidService.DeleteBid(int64(bidID), c.GetInt64("user_id"))
	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bid deleted successfully"})
}

