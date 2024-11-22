package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	request_model "tender-backend/model/request"

	"gorm.io/gorm/utils"

	"github.com/gin-gonic/gin"
)

// CreateTender godoc
// @Security BearerAuth
// @Summary Create a new tender
// @Description Create a new tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param tender body request_model.CreateTenderReq true "Tender information"
// @Success 201 {object} model.Tender
// @Router /tenders [post]
func (h *HTTPHandler) CreateTender(ctx *gin.Context) {
	req := request_model.CreateTenderReq{}

	// Validate JSON input
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input. Please check your request format.",
		})
		return
	}

	// Call the service to create a tender
	res, err := h.TenderService.CreateTender(&req, ctx.GetInt64("user_id"))
	if err != nil {
		// Log the error for internal debugging
		fmt.Printf("Error creating tender: %v\n", err) // Replace with a logger in production

		// Return a user-friendly error message
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid tender data",
		})
		return
	}

	// Respond with the created resource
	ctx.JSON(http.StatusCreated, res)
}

// GetTender godoc
// @Security BearerAuth
// @Summary Get a tender by ID
// @Description Get a tender by ID
// @Tags Tender
// @Produce json
// @Param tender_id path int true "Tender ID"
// @Success 200 {object} model.Tender
// @Router /api/client/tenders/{tender_id} [get]
func (h *HTTPHandler) GetTender(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	res, err := h.TenderService.GetTenderById(int64(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tender not found"})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// GetTenders godoc
// @Security BearerAuth
// @Summary Get all tenders
// @Description Get all tenders
// @Tags Tender
// @Produce json
// @Success 200 {object} []model.Tender
// @Router /api/client/tenders [get]
func (h *HTTPHandler) GetTenders(ctx *gin.Context) {
	res, err := h.TenderService.GetTenders()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// UpdateTender godoc
// @Security BearerAuth
// @Summary Update a tender by ID
// @Description Update a tender by ID
// @Tags Tender
// @Accept json
// @Produce json
// @Param tender_id path int true "Tender ID"
// @Param tender body request_model.UpdateTenderReq true "Tender information"
// @Success 200 {object} model.Tender
// @Router /api/client/tenders/{tender_id} [put]
func (h *HTTPHandler) UpdateTender(ctx *gin.Context) {
	// Get tender ID from the path
	tenderID, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	clientID := ctx.GetInt64("user_id")

	// Bind the request JSON to the UpdateTenderReq struct
	req := request_model.UpdateTenderReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	availableStatus := []string{"open", "closed", "awarded"}
	if !utils.Contains(availableStatus, req.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid tender status"})
		return
	}

	// Call the service method
	_, err = h.TenderService.UpdateTender(int64(tenderID), clientID, &req)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tender not found"})
		return
	}

	// Respond with the updated tender
	ctx.JSON(http.StatusOK, gin.H{"message": "Tender status updated"})
}

// DeleteTender godoc
// @Security BearerAuth
// @Summary Delete a tender by ID
// @Description Delete a tender by ID
// @Tags Tender
// @Param tender_id path int true "Tender ID"
// @Success 204
// @Router /api/client/tenders/{tender_id} [delete]
func (h *HTTPHandler) DeleteTender(ctx *gin.Context) {
	// Get tender ID from the path
	tenderID, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	clientID := ctx.GetInt64("user_id")
	err = h.TenderService.DeleteTender(int64(tenderID), clientID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
		return
	}

	// Respond with no content
	ctx.JSON(http.StatusOK, gin.H{"message": "Tender deleted successfully"})
}

// AwardTender godoc
// @Security BearerAuth
// @Summary Award a tender
// @Description Award a tender
// @Tags Tender
// @Accept json
// @Produce json
// @Param tender_id path int true "Tender ID"
// @Param bid_id path int true "Bid ID"
// @Success 200 {object} model.Tender
// @Router /api/client/tenders/{tender_id}/award/{bid_id} [post]
func (h *HTTPHandler) AwardTender(ctx *gin.Context) {
	tenderID, err := strconv.Atoi(ctx.Param("tender_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tender ID"})
		return
	}

	bidID, err := strconv.Atoi(ctx.Param("bid_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid bid ID"})
		return
	}

	clientID := ctx.GetInt64("user_id")

	if !h.TenderService.IsTenderExists(int64(tenderID)) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Tender not found or access denied"})
		return
	}

	 err=h.BidService.IsBidExists(int64(bidID), int64(tenderID))
	 if err!=nil{
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	 }

	err = h.TenderService.AwardTender(int64(tenderID), clientID, int64(bidID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Tender not found or access denied"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bid awarded successfully"})
}

