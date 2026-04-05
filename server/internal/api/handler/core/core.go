package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"

	"chat-application/internal/api/model"
	"chat-application/internal/constants"
	"chat-application/internal/middleware"
	roomRepository "chat-application/internal/repo/room"
	websoc "chat-application/internal/websocket"
	"chat-application/util"

	"github.com/google/uuid"
)

// allowedWebSocketOrigins defines the valid origins for WebSocket connections.
var defaultAllowedWebSocketOrigins = []string{
	"http://localhost:3000",
	"http://localhost:5173",
	"http://localhost:5174",
	"https://yappin.chat",
}

// checkWebSocketOrigin validates the origin header against allowed origins.
func checkWebSocketOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	for _, allowed := range util.GetEnvList("ALLOWED_ORIGINS", defaultAllowedWebSocketOrigins) {
		if origin == allowed {
			return true
		}
	}
	log.Printf("WebSocket connection rejected: invalid origin %s", origin)
	return false
}

// CoreHandler handles HTTP requests for core chat functionality.
type CoreHandler struct {
	core           *websoc.Core
	roomRepository roomRepository.RoomRepositoryInterface
	roomLimit      int
}

// NewCoreHandler creates a new CoreHandler instance.
func NewCoreHandler(c *websoc.Core) *CoreHandler {
	return NewCoreHandlerWithRoomRepository(c, roomRepository.NewRoomRepository(c.GetDB()))
}

func NewCoreHandlerWithRoomRepository(c *websoc.Core, repo roomRepository.RoomRepositoryInterface) *CoreHandler {
	roomLimit := constants.DefaultRoomLimit
	if maxRoomsStr := util.GetEnv("MAX_ROOMS", ""); maxRoomsStr != "" {
		if limit, err := strconv.Atoi(maxRoomsStr); err == nil {
			roomLimit = limit
		}
	}

	return &CoreHandler{
		core:           c,
		roomRepository: repo,
		roomLimit:      roomLimit,
	}
}

func (h *CoreHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req model.CreateRoomReq

	log.Printf("CreateRoom request: %s", r.URL.Path)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	ctx := r.Context()

	var creatorID *uuid.UUID
	if userIDStr, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		log.Printf("User ID from context: %s", userIDStr)
		if uid, err := uuid.Parse(userIDStr); err == nil {
			creatorID = &uid

		} else {
			log.Printf("Failed to parse user ID: %v", err)
		}
	} else {
		log.Printf("User ID not found in context")
	}

	activeRooms, err := h.roomRepository.CountActiveRooms(ctx)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve active rooms")
		return
	}

	log.Printf("Active rooms: %d, limit: %d", activeRooms, h.roomLimit)

	if activeRooms >= h.roomLimit {
		util.WriteErrorResponse(w, http.StatusTooManyRequests, "Room limit reached")
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	if req.ExpiresAt != nil && !req.ExpiresAt.IsZero() {
		expiresAt = *req.ExpiresAt
	}

	room := &roomRepository.Room{
		Name:      req.Name,
		CreatorID: creatorID,
		ExpiresAt: expiresAt,
	}

	room, err = h.roomRepository.CreateRoom(ctx, room)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create room")
		return
	}

	if creatorID != nil {
		if err := h.roomRepository.EnsureRoomMembership(ctx, room.ID, *creatorID); err != nil {
			util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add room owner")
			return
		}

		member, err := h.roomRepository.GetRoomMember(ctx, room.ID, *creatorID)
		if err == nil && member != nil {
			member.Role = "owner"
			member.CanManageRoom = true
			member.CanManageChannels = true
			member.CanModerate = true
			member.CanPost = true
			if err := h.roomRepository.UpdateRoomMember(ctx, *member); err != nil {
				util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update room owner")
				return
			}
		}
	}

	category, err := h.roomRepository.CreateCategory(ctx, &roomRepository.RoomCategory{
		RoomID:   room.ID,
		Name:     "General",
		Position: 0,
	})
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create room category")
		return
	}

	_, err = h.roomRepository.CreateChannel(ctx, &roomRepository.RoomChannel{
		RoomID:      room.ID,
		CategoryID:  &category.ID,
		Name:        "lobby",
		Description: "Default channel",
		Kind:        "text",
		Position:    0,
	})
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create default channel")
		return
	}

	h.core.AddRoom(&websoc.Room{
		ID:               room.ID.String(),
		Name:             room.Name,
		Clients:          make(map[string]*websoc.Client),
		IsPinned:         room.IsPinned,
		TopicTitle:       room.TopicTitle,
		TopicDescription: room.TopicDescription,
		TopicURL:         room.TopicURL,
		TopicSource:      room.TopicSource,
	})

	response := model.CreateRoomReq{
		ID:   room.ID.String(),
		Name: room.Name,
	}

	util.WriteJSONResponse(w, http.StatusOK, response)

}

func (h *CoreHandler) JoinRoom(w http.ResponseWriter, r *http.Request) {
	log.Printf("JoinRoom called: Method=%s Path=%s", r.Method, r.URL.Path)

	// Get roomId from URL path parameter (matches router: /join-room/{roomId})
	roomID := chi.URLParam(r, "roomId")
	log.Printf("roomID from path param: %s", roomID)

	if roomID == "" {
		log.Printf("Room ID is empty")
		util.WriteErrorResponse(w, http.StatusBadRequest, "Room ID is required")
		return
	}

	roomUUID, err := uuid.Parse(roomID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid Room ID format")
		return
	}

	ctx := r.Context()
	dbRoom, err := h.roomRepository.GetRoomByID(ctx, roomUUID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve room")
		return
	}

	if dbRoom == nil {
		util.WriteErrorResponse(w, http.StatusNotFound, "Room not found")
		return
	}

	var authenticatedUserID string
	if userID, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		authenticatedUserID = userID
		if parsedUserID, err := uuid.Parse(userID); err == nil {
			member, err := h.roomRepository.GetRoomMember(ctx, roomUUID, parsedUserID)
			if err != nil {
				util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to load room membership")
				return
			}
			if member != nil && member.BannedAt != nil {
				util.WriteErrorResponse(w, http.StatusForbidden, "You are banned from this room")
				return
			}
			if member == nil {
				if err := h.roomRepository.EnsureRoomMembership(ctx, roomUUID, parsedUserID); err != nil {
					util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to join room")
					return
				}
			}
		}
	}

	if _, exists := h.core.GetRoom(dbRoom.ID.String()); !exists {
		h.core.AddRoom(&websoc.Room{
			ID:               roomID,
			Name:             dbRoom.Name,
			Clients:          make(map[string]*websoc.Client),
			IsPinned:         dbRoom.IsPinned,
			TopicTitle:       dbRoom.TopicTitle,
			TopicDescription: dbRoom.TopicDescription,
			TopicURL:         dbRoom.TopicURL,
			TopicSource:      dbRoom.TopicSource,
		})
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  constants.WebSocketReadBufferSize,
		WriteBufferSize: constants.WebSocketWriteBufferSize,
		CheckOrigin:     checkWebSocketOrigin,

		EnableCompression: true,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to upgrade connection")
		return
	}

	log.Printf("WebSocket connection upgraded successfully")

	q := r.URL.Query()
	clientID := q.Get("client_id")
	if clientID == "" {
		clientID = uuid.New().String()
	}

	userID := authenticatedUserID
	if userID == "" {
		userID = q.Get("user_id")
	}
	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID format")
			return
		}
	}

	username := q.Get("username")
	if username == "" {
		username = "Anonymous"
	}

	cl := &websoc.Client{
		Conn:     conn,
		Message:  make(chan *websoc.Event, 16),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
		UserID:   userID,
	}

	log.Printf("Registering client: ID=%s Username=%s RoomID=%s", clientID, username, roomID)
	h.core.Register <- cl

	go cl.WriteMessage()
	cl.ReadMessage(h.core)
}

func (h *CoreHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbRooms, err := h.roomRepository.GetAllActiveRooms(ctx)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "failed to fetch rooms")
		return
	}

	rooms := make([]model.RoomRes, 0, len(dbRooms))
	for _, room := range dbRooms {
		var creatorUsername *string
		if room.CreatorID != nil {
			var username string
			err := h.roomRepository.GetDB().QueryRowContext(ctx,
				"SELECT username FROM users WHERE id = $1", room.CreatorID).Scan(&username)
			if err == nil {
				creatorUsername = &username
			}
		}

		participantCount := 0
		if wsRoom, exists := h.core.GetRoom(room.ID.String()); exists {
			participantCount = len(wsRoom.Clients)
		}

		rooms = append(rooms, model.RoomRes{
			ID:               room.ID.String(),
			Name:             room.Name,
			IsPinned:         room.IsPinned,
			CreatedAt:        room.CreatedAt,
			Expires:          room.ExpiresAt,
			TopicTitle:       room.TopicTitle,
			TopicDescription: room.TopicDescription,
			TopicURL:         room.TopicURL,
			TopicSource:      room.TopicSource,
			CreatorUsername:  creatorUsername,
			Participants:     participantCount,
		})

		if _, exists := h.core.GetRoom(room.ID.String()); !exists {
			h.core.AddRoom(&websoc.Room{
				ID:               room.ID.String(),
				Name:             room.Name,
				Clients:          make(map[string]*websoc.Client),
				IsPinned:         room.IsPinned,
				TopicTitle:       room.TopicTitle,
				TopicDescription: room.TopicDescription,
				TopicURL:         room.TopicURL,
				TopicSource:      room.TopicSource,
			})
		}
	}

	util.WriteJSONResponse(w, http.StatusOK, rooms)
}

func (h *CoreHandler) GetRoomDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	room, err := h.roomRepository.GetRoomByID(ctx, roomID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to load room")
		return
	}
	if room == nil {
		util.WriteErrorResponse(w, http.StatusNotFound, "Room not found")
		return
	}

	roomRes, err := h.buildRoomDetailResponse(ctx, room)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to build room detail")
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, roomRes)
}

func (h *CoreHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID, member, ok := h.requireRoomManager(w, r)
	if !ok {
		return
	}
	_ = member

	var req model.CreateCategoryReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Category name is required")
		return
	}

	category, err := h.roomRepository.CreateCategory(ctx, &roomRepository.RoomCategory{
		RoomID:   roomID,
		Name:     req.Name,
		Position: req.Position,
	})
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, model.RoomCategoryRes{
		ID:       category.ID.String(),
		Name:     category.Name,
		Position: category.Position,
		Channels: []model.RoomChannelRes{},
	})
}

func (h *CoreHandler) CreateChannel(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID, _, ok := h.requireRoomManager(w, r)
	if !ok {
		return
	}

	var req model.CreateChannelReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Channel name is required")
		return
	}

	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		parsedID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		categoryID = &parsedID
	}

	channel, err := h.roomRepository.CreateChannel(ctx, &roomRepository.RoomChannel{
		RoomID:      roomID,
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Kind:        defaultString(req.Kind, "text"),
		Position:    req.Position,
		IsPrivate:   req.IsPrivate,
	})
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to create channel")
		return
	}

	response := model.RoomChannelRes{
		ID:          channel.ID.String(),
		Name:        channel.Name,
		Description: channel.Description,
		Kind:        channel.Kind,
		Position:    channel.Position,
		IsPrivate:   channel.IsPrivate,
	}
	if channel.CategoryID != nil {
		response.CategoryID = channel.CategoryID.String()
	}

	util.WriteJSONResponse(w, http.StatusCreated, response)
}

func (h *CoreHandler) SearchMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid room ID")
		return
	}

	queryText := strings.TrimSpace(r.URL.Query().Get("query"))
	if queryText == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Query is required")
		return
	}

	var channelID *uuid.UUID
	if rawChannelID := strings.TrimSpace(r.URL.Query().Get("channel_id")); rawChannelID != "" {
		parsedID, err := uuid.Parse(rawChannelID)
		if err != nil {
			util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid channel ID")
			return
		}
		channelID = &parsedID
	}

	results, err := h.roomRepository.SearchMessages(ctx, roomID, queryText, channelID, strings.TrimSpace(r.URL.Query().Get("username")), 50)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to search messages")
		return
	}

	response := make([]model.MessageSearchRes, 0, len(results))
	for _, message := range results {
		item := model.MessageSearchRes{
			ID:          message.ID.String(),
			RoomID:      message.RoomID.String(),
			Username:    message.Username,
			Content:     message.Content,
			Highlighted: highlightQuery(message.Content, queryText),
			CreatedAt:   message.CreatedAt,
		}
		if message.ChannelID != nil {
			item.ChannelID = message.ChannelID.String()
		}
		if message.ParentMessageID != nil {
			item.ParentMessageID = message.ParentMessageID.String()
		}
		response = append(response, item)
	}

	util.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *CoreHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	notifications, err := h.roomRepository.GetNotifications(r.Context(), parsedUserID, 50)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to load notifications")
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, h.mapNotifications(notifications))
}

func (h *CoreHandler) MarkNotificationRead(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	notificationID, err := uuid.Parse(chi.URLParam(r, "notificationId"))
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.roomRepository.MarkNotificationRead(r.Context(), notificationID, parsedUserID); err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update notification")
		return
	}
	util.WriteJSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *CoreHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	roomID, actingMember, ok := h.requireRoomManager(w, r)
	if !ok {
		return
	}
	if !actingMember.CanModerate && !actingMember.CanManageRoom {
		util.WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return
	}

	targetUserID, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	targetMember, err := h.roomRepository.GetRoomMember(ctx, roomID, targetUserID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to load member")
		return
	}
	if targetMember == nil {
		util.WriteErrorResponse(w, http.StatusNotFound, "Member not found")
		return
	}
	if targetMember.Role == "owner" && actingMember.Role != "owner" {
		util.WriteErrorResponse(w, http.StatusForbidden, "Only room owners can modify other owners")
		return
	}

	var req model.UpdateMemberRoleReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Role != "" {
		targetMember.Role = req.Role
	}
	if req.CanManageRoom != nil {
		targetMember.CanManageRoom = *req.CanManageRoom
	}
	if req.CanManageChannels != nil {
		targetMember.CanManageChannels = *req.CanManageChannels
	}
	if req.CanModerate != nil {
		targetMember.CanModerate = *req.CanModerate
	}
	if req.CanPost != nil {
		targetMember.CanPost = *req.CanPost
	}
	if req.Ban {
		now := time.Now().UTC()
		targetMember.BannedAt = &now
	}

	if err := h.roomRepository.UpdateRoomMember(ctx, *targetMember); err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update member")
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, map[string]bool{"ok": true})
}

func (h *CoreHandler) GetClients(w http.ResponseWriter, r *http.Request) {
	var clients []model.ClientRes
	roomID := chi.URLParam(r, "room_id")

	room, ok := h.core.GetRoom(roomID)
	if !ok {
		util.WriteErrorResponse(w, http.StatusNotFound, "Room not found")
		return
	}

	for _, c := range room.Clients {
		clients = append(clients, model.ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	util.WriteJSONResponse(w, http.StatusOK, clients)
}

func (h *CoreHandler) requireRoomManager(w http.ResponseWriter, r *http.Request) (uuid.UUID, *roomRepository.RoomMember, bool) {
	ctx := r.Context()
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return uuid.Nil, nil, false
	}
	roomID, err := uuid.Parse(chi.URLParam(r, "roomId"))
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid room ID")
		return uuid.Nil, nil, false
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID")
		return uuid.Nil, nil, false
	}
	member, err := h.roomRepository.GetRoomMember(ctx, roomID, parsedUserID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to load room membership")
		return uuid.Nil, nil, false
	}
	if member == nil || (!member.CanManageRoom && !member.CanManageChannels && member.Role != "owner" && member.Role != "admin") {
		util.WriteErrorResponse(w, http.StatusForbidden, "Insufficient permissions")
		return uuid.Nil, nil, false
	}
	return roomID, member, true
}

func (h *CoreHandler) buildRoomDetailResponse(ctx context.Context, room *roomRepository.Room) (*model.RoomDetailRes, error) {
	categories, err := h.roomRepository.GetRoomCategories(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	channels, err := h.roomRepository.GetRoomChannels(ctx, room.ID)
	if err != nil {
		return nil, err
	}
	defaultChannel, err := h.roomRepository.GetDefaultChannel(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	var defaultChannelID *uuid.UUID
	if defaultChannel != nil {
		defaultChannelID = &defaultChannel.ID
	}
	messages, err := h.roomRepository.GetRoomMessagesByChannel(ctx, room.ID, defaultChannelID, 100, 0)
	if err != nil {
		return nil, err
	}
	members, err := h.roomRepository.GetRoomMembers(ctx, room.ID)
	if err != nil {
		return nil, err
	}

	var currentUser *model.RoomPermissionRes
	var notifications []model.NotificationRes
	if userID, ok := ctx.Value(middleware.UserIDKey).(string); ok {
		parsedUserID, err := uuid.Parse(userID)
		if err == nil {
			member, err := h.roomRepository.GetRoomMember(ctx, room.ID, parsedUserID)
			if err == nil && member != nil {
				currentUser = &model.RoomPermissionRes{
					Role:              member.Role,
					CanManageRoom:     member.CanManageRoom,
					CanManageChannels: member.CanManageChannels,
					CanModerate:       member.CanModerate,
					CanPost:           member.CanPost,
					IsMuted:           member.MutedUntil != nil && member.MutedUntil.After(time.Now()),
					IsBanned:          member.BannedAt != nil,
				}
			}
			rawNotifications, err := h.roomRepository.GetNotifications(ctx, parsedUserID, 20)
			if err == nil {
				filtered := make([]roomRepository.Notification, 0, len(rawNotifications))
				for _, notification := range rawNotifications {
					if notification.RoomID != nil && *notification.RoomID == room.ID {
						filtered = append(filtered, notification)
					}
				}
				notifications = h.mapNotifications(filtered)
			}
		}
	}

	categoryMap := make(map[uuid.UUID]*model.RoomCategoryRes, len(categories))
	responseCategories := make([]model.RoomCategoryRes, 0, len(categories))
	for _, category := range categories {
		item := model.RoomCategoryRes{
			ID:       category.ID.String(),
			Name:     category.Name,
			Position: category.Position,
			Channels: []model.RoomChannelRes{},
		}
		responseCategories = append(responseCategories, item)
		categoryMap[category.ID] = &responseCategories[len(responseCategories)-1]
	}

	for _, channel := range channels {
		channelRes := model.RoomChannelRes{
			ID:          channel.ID.String(),
			Name:        channel.Name,
			Description: channel.Description,
			Kind:        channel.Kind,
			Position:    channel.Position,
			IsPrivate:   channel.IsPrivate,
		}
		if channel.CategoryID != nil {
			channelRes.CategoryID = channel.CategoryID.String()
			if category := categoryMap[*channel.CategoryID]; category != nil {
				category.Channels = append(category.Channels, channelRes)
				continue
			}
		}
		if len(responseCategories) == 0 {
			responseCategories = append(responseCategories, model.RoomCategoryRes{
				ID:       "",
				Name:     "Uncategorized",
				Position: 0,
				Channels: []model.RoomChannelRes{},
			})
		}
		responseCategories[0].Channels = append(responseCategories[0].Channels, channelRes)
	}

	messageResponses := make([]model.MessageRes, 0, len(messages))
	threadCount := 0
	for _, message := range messages {
		item := model.MessageRes{
			ID:        message.ID.String(),
			Content:   message.Content,
			RoomID:    message.RoomID.String(),
			Username:  message.Username,
			System:    message.IsSystem,
			CreatedAt: message.CreatedAt,
		}
		if message.ChannelID != nil {
			item.ChannelID = message.ChannelID.String()
		}
		if message.ParentMessageID != nil {
			item.ParentMessageID = message.ParentMessageID.String()
			threadCount++
		}
		if message.UserID != nil {
			item.UserID = message.UserID.String()
		}
		if len(message.Metadata) > 0 {
			var metadata map[string]any
			if err := json.Unmarshal(message.Metadata, &metadata); err == nil {
				item.Metadata = metadata
			}
		}
		reactions, err := h.roomRepository.GetReactions(ctx, message.ID.String())
		if err == nil {
			item.Reactions = reactions
		}
		messageResponses = append(messageResponses, item)
	}

	memberResponses := make([]model.RoomMemberRes, 0, len(members))
	for _, member := range members {
		memberResponses = append(memberResponses, model.RoomMemberRes{
			UserID:    member.UserID.String(),
			Username:  member.Username,
			Role:      member.Role,
			CreatedAt: member.CreatedAt,
		})
	}

	participantCount := 0
	if wsRoom, exists := h.core.GetRoom(room.ID.String()); exists {
		participantCount = len(wsRoom.Clients)
	}

	roomRes := model.RoomRes{
		ID:               room.ID.String(),
		Name:             room.Name,
		IsPinned:         room.IsPinned,
		CreatedAt:        room.CreatedAt,
		Expires:          room.ExpiresAt,
		TopicTitle:       room.TopicTitle,
		TopicDescription: room.TopicDescription,
		TopicURL:         room.TopicURL,
		TopicSource:      room.TopicSource,
		Participants:     participantCount,
	}

	res := &model.RoomDetailRes{
		Room:               roomRes,
		Categories:         responseCategories,
		Members:            memberResponses,
		CurrentUser:        currentUser,
		Messages:           messageResponses,
		Notifications:      notifications,
		NotificationCount:  len(notifications),
		OnlineMemberCount:  participantCount,
		ThreadedReplyCount: threadCount,
	}
	if defaultChannel != nil {
		res.DefaultChannelID = defaultChannel.ID.String()
	}
	return res, nil
}

func (h *CoreHandler) mapNotifications(items []roomRepository.Notification) []model.NotificationRes {
	response := make([]model.NotificationRes, 0, len(items))
	for _, item := range items {
		notification := model.NotificationRes{
			ID:        item.ID.String(),
			Kind:      item.Kind,
			Title:     item.Title,
			Body:      item.Body,
			IsRead:    item.IsRead,
			CreatedAt: item.CreatedAt,
		}
		if item.RoomID != nil {
			notification.RoomID = item.RoomID.String()
		}
		if item.MessageID != nil {
			notification.MessageID = item.MessageID.String()
		}
		if len(item.Payload) > 0 {
			var payload map[string]any
			if err := json.Unmarshal(item.Payload, &payload); err == nil {
				notification.Payload = payload
			}
		}
		response = append(response, notification)
	}
	return response
}

func highlightQuery(content, query string) string {
	if query == "" {
		return content
	}
	return strings.ReplaceAll(content, query, "<mark>"+query+"</mark>")
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func (h *CoreHandler) AddReaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		util.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req model.RequestAddReaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !constants.IsValidReactionEmoji(req.Emoji) {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Invalid emoji")
		return
	}

	reaction := &model.MessageReaction{
		MessageID: req.MessageID,
		UserID:    userID,
		Emoji:     req.Emoji,
	}

	if err := h.roomRepository.AddReaction(r.Context(), reaction); err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to add reaction")
		return
	}

	util.WriteJSONResponse(w, http.StatusCreated, reaction)
}

func (h *CoreHandler) GetReactions(w http.ResponseWriter, r *http.Request) {
	messageID := chi.URLParam(r, "messageID")
	if messageID == "" {
		util.WriteErrorResponse(w, http.StatusBadRequest, "Message ID is required")
		return
	}

	reactions, err := h.roomRepository.GetReactions(r.Context(), messageID)
	if err != nil {
		util.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to fetch reactions")
		return
	}

	util.WriteJSONResponse(w, http.StatusOK, reactions)
}
