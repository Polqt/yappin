package handler

import (
	"log"
	"net/http"

	statsService "chat-application/internal/service/stats"
	"chat-application/util"

	"github.com/google/uuid"
)

type StatsHandler struct {
	statsService *statsService.StatsService
}

func NewStatsHandler(statsService *statsService.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

func(h *StatsHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDString, ok := ctx.Value("userID").(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
		return 
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	log.Printf("Processing daily check-in for user: %s", userIDString)

	result, err := h.statsService.ProcessDailyCheckin(ctx, userID)
	if err != nil {
		log.Printf("Error processing daily check-in for user %s: %v", userID, err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to process daily check-in")
		return 
	}

	util.WriteJSONResponse(w, http.StatusOK, result)
}

func(h *StatsHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {

}