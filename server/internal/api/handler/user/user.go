package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-application/internal/api/model"
	"chat-application/internal/middleware"
	service "chat-application/internal/service/user"
	"chat-application/util"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateUser - JSON decode error: %v", err)
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateUser - Request received: username=%s, email=%s", req.Username, req.Email)

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		log.Printf("CreateUser - Service error: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateUser - Success: user created with ID=%s, username=%s", user.ID, user.Username)

	util.SetCookie(w, "jwt", user.AccessToken, 60*60*24) // 24 hours

	util.WriteJSONResponse(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.RequestLoginUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	user, err := h.userService.Login(r.Context(), req)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
		return
	}

	util.SetCookie(w, "jwt", user.AccessToken, 60*60*24) // 24 hours

	util.WriteJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	util.ClearSecureCookie(w, "jwt")
	util.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "logout successful"})
}

func (h *UserHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	req.Username = util.SanitizeString(req.Username)

	if err := util.ValidateUsername(req.Username); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.userService.UpdateUsername(r.Context(), userID, req.Username)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, user)
}

func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), uid)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	response := model.ResponseLoginUser{
		ID:       user.ID.String(),
		Username: user.Username,
	}

	util.WriteJSONResponse(w, http.StatusOK, response)
}
