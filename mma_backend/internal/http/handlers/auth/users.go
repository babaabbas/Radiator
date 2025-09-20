package auth

import (
	"encoding/json"
	"fmt"
	"mma_api/internal/storage"
	"mma_api/internal/storage/postgres"
	"mma_api/internal/utils/response"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

func Register_handler(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}
		_, err := storage.GetUserByEmail(req.Email)
		if err == nil {
			http.Error(w, "email already registered", http.StatusBadRequest)
			return
		}
		password, err := HashPassword(req.Password)
		if err != nil {
			fmt.Print("password didnt get hashed")
		}

		newUser, err := storage.CreateUser(req.Name, req.Role, req.Email, password)
		if err != nil {
			http.Error(w, "failed to create user", http.StatusInternalServerError)
			return
		}

		resp := RegisterResponse{
			Message: "user registered successfully",
			UserID:  newUser.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
	Name    string `json:"name"`
	Role    string `json:"role"`
}

func Login_handler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := storage.GetUserByEmail(req.Email)
		if err != nil {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		if !CheckPasswordHash(req.Password, user.PasswordHash) {
			http.Error(w, "invalid email or password", http.StatusUnauthorized)
			return
		}

		resp := LoginResponse{
			Message: "login successful",
			UserID:  user.ID,
			Name:    user.Name,
			Role:    user.Role,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func GetUsersHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := storage.GetUsers()
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		var respData []UserResponse
		for _, u := range users {
			respData = append(respData, UserResponse{
				ID:    u.ID,
				Name:  u.Name,
				Role:  u.Role,
				Email: u.Email,
			})
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          respData,
		})
	}
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

func GetUserByIDHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 { // ["", "api", "users", "{id}"]
			resp := response.GeneralError(http.ErrMissingFile)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		user, err := storage.GetUserByID(id)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusNotFound, resp)
			return
		}

		respData := UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Role:  user.Role,
			Email: user.Email,
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data":          respData,
		})
	}
}
func DeleteUserByIDHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			resp := response.GeneralError(http.ErrMissingFile)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		err = storage.DeleteUser(id)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"message":       "user deleted successfully",
		})
	}
}

type UpdateUserRequest struct {
	Name     string `json:"name,omitempty"`
	Role     string `json:"role,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func UpdateUserHandler(storage *postgres.Postgres) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			resp := response.GeneralError(http.ErrNotSupported)
			_ = response.WriteJson(w, http.StatusMethodNotAllowed, resp)
			return
		}

		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 4 {
			resp := response.GeneralError(http.ErrMissingFile)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		idStr := pathParts[3]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		var req UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusBadRequest, resp)
			return
		}

		var hashedPassword string
		if req.Password != "" {
			hashedPassword, err = HashPassword(req.Password)
			if err != nil {
				resp := response.GeneralError(err)
				_ = response.WriteJson(w, http.StatusInternalServerError, resp)
				return
			}
		}

		updatedUser, err := storage.UpdateUser(id, req.Name, req.Role, req.Email, hashedPassword)
		if err != nil {
			resp := response.GeneralError(err)
			_ = response.WriteJson(w, http.StatusInternalServerError, resp)
			return
		}

		_ = response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"custom_status": response.Status_Ok,
			"data": map[string]interface{}{
				"id":    updatedUser.ID,
				"name":  updatedUser.Name,
				"role":  updatedUser.Role,
				"email": updatedUser.Email,
			},
		})
	}
}
