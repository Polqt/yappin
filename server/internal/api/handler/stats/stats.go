package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-application/internal/middleware"
	statsService "chat-application/internal/service/stats"
	"chat-application/util"

	"github.com/go-chi/chi/v5"
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

func (h *StatsHandler) CheckIn(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userIDString, ok := ctx.Value(middleware.UserIDKey).(string)
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

func (h *StatsHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	targetUserIDString := chi.URLParam(r, "userID")
	targetUserID, err := uuid.Parse(targetUserIDString)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	viewerUserID := uuid.Nil
	if viewerIDString, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		if viewerID, err := uuid.Parse(viewerIDString); err == nil {
			viewerUserID = viewerID
		}
	}

	profile, err := h.statsService.GetUserProfile(ctx, targetUserID, viewerUserID)
	if err != nil {
		log.Printf("Error retrieving user profile for user %s: %v", targetUserID, err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to retrieve user profile")
	}

	util.WriteJSONResponse(w, http.StatusOK, profile)
}

func (h *StatsHandler) GivenUpvote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fromUserIDString, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	fromUserID, err := uuid.Parse(fromUserIDString)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var req struct {
		ToUserID string `json:"to_user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid request body")
		return
	}

	toUserID, err := uuid.Parse(req.ToUserID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid to_user_id format")
		return
	}

	log.Printf("User %s is giving an upvote to user %s", fromUserIDString, req.ToUserID)

	err = h.statsService.GivenUpvote(ctx, fromUserID, toUserID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to give upvote")
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]string{"status": "upvote given successfully"})
}
