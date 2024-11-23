package handler

import (
	"encoding/json"
	"net/http"
	"unicode/utf8"

	"github.com/jackgris/twitter-backend/auth/internal/domain/usermodel"
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

	if utf8.RuneCountInString(input.UserName) < 5 {
		http.Error(w, "user name invalid must be longer than 5 characters", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(input.Password) < 4 {
		http.Error(w, "password size should have a minimun of 4 characters", http.StatusBadRequest)
		return
	}

	v := validator.New()
	validator.ValidateEmail(v, input.Email)
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
}
func (u UserHandler) Follow(w http.ResponseWriter, r *http.Request) {
}
func (u UserHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
}
func (u UserHandler) GetUserbyUsername(w http.ResponseWriter, r *http.Request) {
}
func (u UserHandler) GetUserbyID(w http.ResponseWriter, r *http.Request) {
}
func (u UserHandler) Update(w http.ResponseWriter, r *http.Request) {
}
