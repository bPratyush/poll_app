package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"poll_app/ent"
	"poll_app/ent/notification"
	"poll_app/ent/poll"
	"poll_app/ent/polloption"
	"poll_app/ent/user"
	"poll_app/ent/vote"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your-secret-key-change-in-production")

// SetJWTSecret allows setting the JWT secret from environment
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Handler struct {
	client *ent.Client
}

func NewHandler(client *ent.Client) *Handler {
	return &Handler{client: client}
}

type contextKey string

const userContextKey contextKey = "user"

// Response helpers
func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

// Auth handlers
type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		errorResponse(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	u, err := h.client.User.Create().
		SetUsername(req.Username).
		SetEmail(req.Email).
		SetPassword(string(hashedPassword)).
		Save(context.Background())
	if err != nil {
		if ent.IsConstraintError(err) {
			errorResponse(w, http.StatusConflict, "Username or email already exists")
			return
		}
		errorResponse(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT
	token, err := generateToken(u.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	jsonResponse(w, http.StatusCreated, AuthResponse{
		Token: token,
		User: UserDTO{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		},
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Find user by email
	u, err := h.client.User.Query().Where(user.Email(req.Email)).Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		errorResponse(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Generate JWT
	token, err := generateToken(u.ID)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	jsonResponse(w, http.StatusOK, AuthResponse{
		Token: token,
		User: UserDTO{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		},
	})
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)
	jsonResponse(w, http.StatusOK, UserDTO{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	})
}

func generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (h *Handler) AuthMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errorResponse(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			errorResponse(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			errorResponse(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		userID := int(claims["user_id"].(float64))
		u, err := h.client.User.Get(context.Background(), userID)
		if err != nil {
			errorResponse(w, http.StatusUnauthorized, "User not found")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, u)
		next(w, r.WithContext(ctx), ps)
	}
}

// Poll handlers
type CreatePollRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
}

type UpdatePollRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Options     []OptionUpdate `json:"options"`
}

type OptionUpdate struct {
	ID   int    `json:"id,omitempty"`
	Text string `json:"text"`
}

type PollDTO struct {
	ID                  int         `json:"id"`
	Title               string      `json:"title"`
	Description         string      `json:"description"`
	Creator             UserDTO     `json:"creator"`
	Options             []OptionDTO `json:"options"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	UserVotedOptionID   *int        `json:"user_voted_option_id,omitempty"`
	PollEditedAfterVote bool        `json:"poll_edited_after_vote"`
}

type OptionDTO struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	VoteCount int    `json:"vote_count"`
}

func (h *Handler) CreatePoll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	var req CreatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Title == "" || len(req.Options) < 2 {
		errorResponse(w, http.StatusBadRequest, "Title and at least 2 options are required")
		return
	}

	// Create poll with options in a transaction
	tx, err := h.client.Tx(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	p, err := tx.Poll.Create().
		SetTitle(req.Title).
		SetDescription(req.Description).
		SetCreator(u).
		Save(context.Background())
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "Failed to create poll")
		return
	}

	for _, optText := range req.Options {
		_, err := tx.PollOption.Create().
			SetText(optText).
			SetPoll(p).
			Save(context.Background())
		if err != nil {
			tx.Rollback()
			errorResponse(w, http.StatusInternalServerError, "Failed to create option")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Fetch the created poll with all relations
	p, _ = h.client.Poll.Query().
		Where(poll.ID(p.ID)).
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		Only(context.Background())

	jsonResponse(w, http.StatusCreated, pollToDTO(p, nil, nil))
}

func (h *Handler) ListPolls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	polls, err := h.client.Poll.Query().
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		Order(ent.Desc(poll.FieldCreatedAt)).
		All(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch polls")
		return
	}

	// Get user's votes with timestamps
	userVotes, _ := h.client.Vote.Query().
		Where(vote.HasUserWith(user.ID(u.ID))).
		WithOption().
		All(context.Background())

	userVoteMap := make(map[int]int)           // pollID -> optionID
	userVoteTimeMap := make(map[int]time.Time) // pollID -> vote time
	for _, v := range userVotes {
		opt := v.Edges.Option
		if opt != nil {
			pollID, _ := h.client.Poll.Query().Where(poll.HasOptionsWith(polloption.ID(opt.ID))).OnlyID(context.Background())
			userVoteMap[pollID] = opt.ID
			userVoteTimeMap[pollID] = v.CreatedAt
		}
	}

	dtos := make([]PollDTO, len(polls))
	for i, p := range polls {
		var votedOptionID *int
		var voteTime *time.Time
		if optID, ok := userVoteMap[p.ID]; ok {
			votedOptionID = &optID
			t := userVoteTimeMap[p.ID]
			voteTime = &t
		}
		dtos[i] = pollToDTO(p, votedOptionID, voteTime)
	}

	jsonResponse(w, http.StatusOK, dtos)
}

func (h *Handler) GetPoll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	p, err := h.client.Poll.Query().
		Where(poll.ID(id)).
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Poll not found")
		return
	}

	// Check if user has voted and get vote time
	var votedOptionID *int
	var userVoteTime *time.Time
	for _, opt := range p.Edges.Options {
		for _, v := range opt.Edges.Votes {
			voter, _ := h.client.Vote.QueryUser(v).Only(context.Background())
			if voter != nil && voter.ID == u.ID {
				votedOptionID = &opt.ID
				userVoteTime = &v.CreatedAt
				break
			}
		}
	}

	jsonResponse(w, http.StatusOK, pollToDTO(p, votedOptionID, userVoteTime))
}

func (h *Handler) UpdatePoll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	p, err := h.client.Poll.Query().
		Where(poll.ID(id)).
		WithCreator().
		WithOptions().
		Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Poll not found")
		return
	}

	// Check ownership
	if p.Edges.Creator.ID != u.ID {
		errorResponse(w, http.StatusForbidden, "You can only edit your own polls")
		return
	}

	var req UpdatePollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update poll
	tx, err := h.client.Tx(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	_, err = tx.Poll.UpdateOneID(id).
		SetTitle(req.Title).
		SetDescription(req.Description).
		Save(context.Background())
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "Failed to update poll")
		return
	}

	// Handle options update
	existingOptionIDs := make(map[int]bool)
	for _, opt := range p.Edges.Options {
		existingOptionIDs[opt.ID] = true
	}

	newOptionIDs := make(map[int]bool)
	for _, opt := range req.Options {
		if opt.ID > 0 {
			// Update existing option
			_, err := tx.PollOption.UpdateOneID(opt.ID).SetText(opt.Text).Save(context.Background())
			if err != nil {
				tx.Rollback()
				errorResponse(w, http.StatusInternalServerError, "Failed to update option")
				return
			}
			newOptionIDs[opt.ID] = true
		} else {
			// Create new option
			_, err := tx.PollOption.Create().SetText(opt.Text).SetPollID(id).Save(context.Background())
			if err != nil {
				tx.Rollback()
				errorResponse(w, http.StatusInternalServerError, "Failed to create option")
				return
			}
		}
	}

	// Delete removed options
	for optID := range existingOptionIDs {
		if !newOptionIDs[optID] {
			// Delete votes for this option first
			_, _ = tx.Vote.Delete().Where(vote.HasOptionWith(polloption.ID(optID))).Exec(context.Background())
			// Delete the option
			_ = tx.PollOption.DeleteOneID(optID).Exec(context.Background())
		}
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Fetch updated poll
	p, _ = h.client.Poll.Query().
		Where(poll.ID(id)).
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		Only(context.Background())

	jsonResponse(w, http.StatusOK, pollToDTO(p, nil, nil))
}

func (h *Handler) DeletePoll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	p, err := h.client.Poll.Query().
		Where(poll.ID(id)).
		WithCreator().
		WithOptions().
		Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Poll not found")
		return
	}

	// Check ownership
	if p.Edges.Creator.ID != u.ID {
		errorResponse(w, http.StatusForbidden, "You can only delete your own polls")
		return
	}

	// Delete in transaction
	tx, err := h.client.Tx(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Delete votes for all options
	for _, opt := range p.Edges.Options {
		_, _ = tx.Vote.Delete().Where(vote.HasOptionWith(polloption.ID(opt.ID))).Exec(context.Background())
	}

	// Delete options
	_, _ = tx.PollOption.Delete().Where(polloption.HasPollWith(poll.ID(id))).Exec(context.Background())

	// Delete poll
	err = tx.Poll.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "Failed to delete poll")
		return
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Vote handlers
type VoteRequest struct {
	OptionID int `json:"option_id"`
}

func (h *Handler) Vote(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	pollID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid poll ID")
		return
	}

	var req VoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Verify option belongs to poll
	opt, err := h.client.PollOption.Query().
		Where(polloption.ID(req.OptionID), polloption.HasPollWith(poll.ID(pollID))).
		Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid option for this poll")
		return
	}

	// Check if user already voted on this poll
	p, _ := h.client.Poll.Query().
		Where(poll.ID(pollID)).
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes(func(vq *ent.VoteQuery) {
				vq.Where(vote.HasUserWith(user.ID(u.ID)))
			})
		}).
		Only(context.Background())

	// Track if this is a vote change
	isVoteChange := false
	var previousOptionText string
	for _, o := range p.Edges.Options {
		if len(o.Edges.Votes) > 0 {
			isVoteChange = true
			previousOptionText = o.Text
			break
		}
	}

	tx, err := h.client.Tx(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}

	// Remove existing vote if any
	for _, o := range p.Edges.Options {
		for _, v := range o.Edges.Votes {
			_ = tx.Vote.DeleteOneID(v.ID).Exec(context.Background())
		}
	}

	// Create new vote
	_, err = tx.Vote.Create().
		SetUser(u).
		SetOption(opt).
		Save(context.Background())
	if err != nil {
		tx.Rollback()
		errorResponse(w, http.StatusInternalServerError, "Failed to create vote")
		return
	}

	// Create notification for poll creator if vote was changed (not for creator's own votes)
	if isVoteChange && p.Edges.Creator.ID != u.ID {
		message := fmt.Sprintf("%s changed their vote on \"%s\" from \"%s\" to \"%s\"",
			u.Username, p.Title, previousOptionText, opt.Text)
		_, err = tx.Notification.Create().
			SetMessage(message).
			SetType("vote_changed").
			SetPollID(pollID).
			SetUserID(p.Edges.Creator.ID).
			Save(context.Background())
		if err != nil {
			// Log but don't fail the vote
			_ = err
		}
	}

	if err := tx.Commit(); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	// Fetch updated poll
	p, _ = h.client.Poll.Query().
		Where(poll.ID(pollID)).
		WithCreator().
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		Only(context.Background())

	votedOptionID := req.OptionID
	now := time.Now()
	jsonResponse(w, http.StatusOK, pollToDTO(p, &votedOptionID, &now))
}

func (h *Handler) GetVoters(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	optionID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid option ID")
		return
	}

	votes, err := h.client.Vote.Query().
		Where(vote.HasOptionWith(polloption.ID(optionID))).
		WithUser().
		All(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch voters")
		return
	}

	voters := make([]UserDTO, len(votes))
	for i, v := range votes {
		voters[i] = UserDTO{
			ID:       v.Edges.User.ID,
			Username: v.Edges.User.Username,
			Email:    v.Edges.User.Email,
		}
	}

	jsonResponse(w, http.StatusOK, voters)
}

func pollToDTO(p *ent.Poll, votedOptionID *int, userVoteTime *time.Time) PollDTO {
	options := make([]OptionDTO, len(p.Edges.Options))
	for i, opt := range p.Edges.Options {
		options[i] = OptionDTO{
			ID:        opt.ID,
			Text:      opt.Text,
			VoteCount: len(opt.Edges.Votes),
		}
	}

	// Check if poll was edited after user voted
	pollEditedAfterVote := false
	if userVoteTime != nil && p.UpdatedAt.After(*userVoteTime) {
		pollEditedAfterVote = true
	}

	return PollDTO{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Creator: UserDTO{
			ID:       p.Edges.Creator.ID,
			Username: p.Edges.Creator.Username,
			Email:    p.Edges.Creator.Email,
		},
		Options:             options,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
		UserVotedOptionID:   votedOptionID,
		PollEditedAfterVote: pollEditedAfterVote,
	}
}

// Notification DTOs
type NotificationDTO struct {
	ID        int       `json:"id"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	PollID    int       `json:"poll_id,omitempty"`
	Read      bool      `json:"read"`
	CreatedAt time.Time `json:"created_at"`
}

// GetNotifications returns all notifications for the current user
func (h *Handler) GetNotifications(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	notifications, err := h.client.Notification.Query().
		Where(notification.HasUserWith(user.ID(u.ID))).
		Order(ent.Desc(notification.FieldCreatedAt)).
		Limit(50).
		All(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to fetch notifications")
		return
	}

	dtos := make([]NotificationDTO, len(notifications))
	for i, n := range notifications {
		dtos[i] = NotificationDTO{
			ID:        n.ID,
			Message:   n.Message,
			Type:      n.Type,
			PollID:    n.PollID,
			Read:      n.Read,
			CreatedAt: n.CreatedAt,
		}
	}

	jsonResponse(w, http.StatusOK, dtos)
}

// GetUnreadCount returns the count of unread notifications
func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	count, err := h.client.Notification.Query().
		Where(
			notification.HasUserWith(user.ID(u.ID)),
			notification.Read(false),
		).
		Count(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to count notifications")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]int{"count": count})
}

// MarkNotificationRead marks a notification as read
func (h *Handler) MarkNotificationRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	notifID, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	// Verify notification belongs to user
	n, err := h.client.Notification.Query().
		Where(
			notification.ID(notifID),
			notification.HasUserWith(user.ID(u.ID)),
		).
		Only(context.Background())
	if err != nil {
		errorResponse(w, http.StatusNotFound, "Notification not found")
		return
	}

	_, err = h.client.Notification.UpdateOne(n).
		SetRead(true).
		Save(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update notification")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Notification marked as read"})
}

// MarkAllNotificationsRead marks all notifications as read for the current user
func (h *Handler) MarkAllNotificationsRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := r.Context().Value(userContextKey).(*ent.User)

	_, err := h.client.Notification.Update().
		Where(
			notification.HasUserWith(user.ID(u.ID)),
			notification.Read(false),
		).
		SetRead(true).
		Save(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to update notifications")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}
