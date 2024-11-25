package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
	"github.com/jackgris/twitter-backend/auth/pkg/uuid"
	"github.com/jackgris/twitter-backend/auth/pkg/validator"
)

func (u UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserName string `json:"user_name"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if input.UserName == "" || input.Password == "" {
		http.Error(w, "user_name and password are required", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(input.Password) < 4 {
		http.Error(w, "password size should have a minimun of 4 characters", http.StatusBadRequest)
		return
	}

	v := validator.New()
	validator.ValidateEmail(v, input.Email)
	validator.ValidateName(v, input.UserName)
	if !v.Valid() {
		err := ""
		for key, value := range v.Errors {
			err += key + " " + value + " "
		}

		http.Error(w, err, http.StatusBadRequest)
		return
	}

	user := usermodel.User{
		UserName: input.UserName,
		Password: input.Password,
		Email:    input.Email,
	}
	user, err := u.store.Create(user)
	if err != nil {
		http.Error(w, "Can't save user in database", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(UserToJSON(user))

}

func (u UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	if ok := uuid.IsValid(userID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	err := u.store.Delete(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u UserHandler) Follow(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserID     string `json:"user_id"`
		FollowerID string `json:"follower_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.UserID == "" || input.FollowerID == "" {
		http.Error(w, "user_id and follower_id are required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.UserID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.FollowerID); !ok {
		http.Error(w, "follower id invalid", http.StatusBadRequest)
		return
	}

	follow := usermodel.UserFollowers{
		UserID:     input.UserID,
		FollowerID: input.FollowerID,
	}
	err = u.store.Follow(follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(FollowerToJSON(follow))
}

func (u UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserID     string `json:"user_id"`
		FollowerID string `json:"follower_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.UserID == "" || input.FollowerID == "" {
		http.Error(w, "user_id and follower_id are required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.UserID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(input.FollowerID); !ok {
		http.Error(w, "follower id invalid", http.StatusBadRequest)
		return
	}

	follow := usermodel.UserFollowers{
		UserID:     input.UserID,
		FollowerID: input.FollowerID,
	}
	err = u.store.Unfollow(follow)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Unfollow successful"))
}

func (u UserHandler) GetUserbyID(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(userID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}

	user, err := u.store.GetUserbyID(userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserToJSON(user))
}

func (u UserHandler) GetUserbyUsername(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "name query parameter is required", http.StatusBadRequest)
		return
	}

	user, err := u.store.GetUserbyUsername(name)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to retrieve user: %v", err), http.StatusInternalServerError)
		}
		return
	}

	v := validator.New()
	validator.ValidateName(v, name)
	if !v.Valid() {
		err := ""
		for key, value := range v.Errors {
			err += key + " " + value + " "
		}

		http.Error(w, err, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(UserToJSON(user))
}

func (u UserHandler) Update(w http.ResponseWriter, r *http.Request) {

	var input map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	userID, ok := input["id"].(string)
	if !ok || userID == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	if ok := uuid.IsValid(userID); !ok {
		http.Error(w, "user id invalid", http.StatusBadRequest)
		return
	}
	// TODO
	user := usermodel.User{
		// UserName       string
		// Email          string
		// Password       string
		// FollowerCount  int
		// FollowingCount int
		// Salt           string
		// Token          string
		// DateCreated    time.Time
		// EncodedDate    string
	}
	updatedUser, err := u.store.Update(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(updatedUser)
}
