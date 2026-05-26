package router

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"linuxmodule/pkg/auth"
	"linuxmodule/pkg/container"
	"linuxmodule/pkg/db"
	"linuxmodule/pkg/email"
	"linuxmodule/pkg/game"
)

// SetupRouter creates the ServeMux and registers all handlers with CORS middleware.
func SetupRouter() http.Handler {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/register", handleRegister)
	mux.HandleFunc("/api/auth/login", handleLogin)
	mux.HandleFunc("/api/auth/verify-email", handleVerifyEmail)
	mux.HandleFunc("/api/auth/resend-code", handleResendCode)

	// Stats routes
	mux.HandleFunc("/api/users/count", handleUsersCount)

	// Game routes
	mux.HandleFunc("/api/game/chapters", handleGetChapters)
	mux.HandleFunc("/api/game/progress", corsMiddleware(authMiddleware(handleGetProgress)))
	mux.HandleFunc("/api/game/verify", corsMiddleware(authMiddleware(handleVerifyTask)))

	// Terminal routes
	mux.HandleFunc("/api/terminal/ws", handleTerminalWS)

	// Wrap in global CORS handler
	return globalCORS(mux)
}

// CORS & Middleware Helpers
func globalCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

type contextKey string

const userClaimsKey contextKey = "userClaims"

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Inject claims into request context by overriding the handler logic or standard middleware flow
		// Go's net/http handles context wrapping. We can do simple type assertion or use r.Header or a custom struct.
		// Since we're in a simple setup, we can write claims to request headers or pass via custom context.
		// Let's pass user claims in request headers or custom request context.
		// In Go, context is standard:
		// ctx := context.WithValue(r.Context(), userClaimsKey, claims)
		// next(w, r.WithContext(ctx))
		// However, to keep it extremely simple and typesafe without context assertion, we can just write the USER_ID and USERNAME directly to request headers!
		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-Username", claims.Username)

		next(w, r)
	}
}

// Request/Response Structs
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
}

type VerifyRequest struct {
	UserInput string `json:"user_input"`
}

type VerifyResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	CurrentChapter   int    `json:"current_chapter"`
	CurrentTask      int    `json:"current_task"`
	ChapterCompleted bool   `json:"chapter_completed"`
	GameCompleted    bool   `json:"game_completed"`
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResendCodeRequest struct {
	Email string `json:"email"`
}

func generateVerificationCode() string {
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	b := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, b, 6)
	if err != nil || n != 6 {
		return "384912"
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

// Handlers
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Validate username format (only alphanumeric, underscores and dashes allowed for safe Docker naming)
	for _, r := range req.Username {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-') {
			http.Error(w, "Username can only contain alphanumeric characters, underscores, and hyphens", http.StatusBadRequest)
			return
		}
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	code := generateVerificationCode()

	// Insert user
	var userID string
	err = db.DB.QueryRow(
		"INSERT INTO users (username, email, password_hash, is_verified, verification_code) VALUES ($1, $2, $3, FALSE, $4) RETURNING id",
		req.Username, req.Email, passwordHash, code,
	).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			http.Error(w, "Username or email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create user progress entry
	_, err = db.DB.Exec(
		"INSERT INTO user_progress (user_id, current_chapter, current_task, completed_chapters) VALUES ($1, 1, 1, '[]'::jsonb)",
		userID,
	)
	if err != nil {
		log.Printf("Failed to create progress row for user %s: %v", userID, err)
	}

	// Dispatch email asynchronously
	go func() {
		err := email.SendVerificationEmail(req.Email, req.Username, code)
		if err != nil {
			log.Printf("Failed to dispatch verification email to %s: %v", req.Email, err)
		}
	}()

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"message": "Registration successful. Verification email sent."}`))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.UsernameOrEmail = strings.TrimSpace(req.UsernameOrEmail)

	var userID, username, emailAddress, passwordHash string
	var isVerified bool
	err := db.DB.QueryRow(
		"SELECT id, username, email, password_hash, is_verified FROM users WHERE username = $1 OR email = $1",
		req.UsernameOrEmail,
	).Scan(&userID, &username, &emailAddress, &passwordHash, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if !auth.CheckPasswordHash(req.Password, passwordHash) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !isVerified {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "Email not verified",
			"email": emailAddress,
		})
		return
	}

	// Generate JWT
	token, err := auth.GenerateJWT(userID, username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(LoginResponse{
		Token:    token,
		Username: username,
	})
}

func handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Code = strings.TrimSpace(req.Code)

	if req.Email == "" || req.Code == "" {
		http.Error(w, "Email and code are required", http.StatusBadRequest)
		return
	}

	var userID, correctCode string
	var isVerified bool
	err := db.DB.QueryRow(
		"SELECT id, verification_code, is_verified FROM users WHERE email = $1 OR username = $1",
		req.Email,
	).Scan(&userID, &correctCode, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Wizard profile not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if isVerified {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message": "Wand already verified!"}`))
		return
	}

	if correctCode != req.Code {
		http.Error(w, "Invalid verification code", http.StatusForbidden)
		return
	}

	// Update verification status
	_, err = db.DB.Exec(
		"UPDATE users SET is_verified = TRUE, verification_code = NULL WHERE id = $1",
		userID,
	)
	if err != nil {
		http.Error(w, "Failed to update verification status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "Wand successfully verified! You may now login."}`))
}

func handleResendCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ResendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	var username, emailAddress, correctCode string
	var isVerified bool
	err := db.DB.QueryRow(
		"SELECT username, email, verification_code, is_verified FROM users WHERE email = $1 OR username = $1",
		req.Email,
	).Scan(&username, &emailAddress, &correctCode, &isVerified)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Wizard profile not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if isVerified {
		http.Error(w, "Wand already verified!", http.StatusBadRequest)
		return
	}

	code := generateVerificationCode()

	_, err = db.DB.Exec(
		"UPDATE users SET verification_code = $1 WHERE email = $2 OR username = $2",
		code, req.Email,
	)
	if err != nil {
		http.Error(w, "Failed to update verification code", http.StatusInternalServerError)
		return
	}

	go func() {
		err := email.SendVerificationEmail(emailAddress, username, code)
		if err != nil {
			log.Printf("Failed to dispatch resent verification email to %s: %v", emailAddress, err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"message": "A new verification rune has been dispatched to your owl email."}`))
}

func handleGetChapters(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(game.Chapters)
}

func handleGetProgress(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	var currentChapter, currentTask int
	var completedJSON []byte
	err := db.DB.QueryRow(
		"SELECT current_chapter, current_task, completed_chapters FROM user_progress WHERE user_id = $1",
		userID,
	).Scan(&currentChapter, &currentTask, &completedJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			// If progress row doesn't exist, create it
			_, err = db.DB.Exec(
				"INSERT INTO user_progress (user_id, current_chapter, current_task, completed_chapters) VALUES ($1, 1, 1, '[]'::jsonb)",
				userID,
			)
			if err != nil {
				http.Error(w, "Failed to create user progress: "+err.Error(), http.StatusInternalServerError)
				return
			}
			currentChapter, currentTask, completedJSON = 1, 1, []byte("[]")
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf(`{"current_chapter": %d, "current_task": %d, "completed_chapters": %s}`,
		currentChapter, currentTask, string(completedJSON))))
}

func handleVerifyTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Header.Get("X-User-ID")

	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.UserInput = strings.TrimSpace(req.UserInput)

	// 1. Fetch user's current progress
	var currentChapter, currentTask int
	var completedJSON []byte
	err := db.DB.QueryRow(
		"SELECT current_chapter, current_task, completed_chapters FROM user_progress WHERE user_id = $1",
		userID,
	).Scan(&currentChapter, &currentTask, &completedJSON)

	if err != nil {
		http.Error(w, "Could not load progress", http.StatusInternalServerError)
		return
	}

	// Find the current chapter definition
	var targetChapter *game.Chapter
	for _, c := range game.Chapters {
		if c.ID == currentChapter {
			targetChapter = &c
			break
		}
	}

	if targetChapter == nil {
		http.Error(w, "Active chapter not found", http.StatusBadRequest)
		return
	}

	// Find the current task definition
	var targetTask *game.Task
	for _, t := range targetChapter.Tasks {
		if t.ID == currentTask {
			targetTask = &t
			break
		}
	}

	if targetTask == nil {
		http.Error(w, "Active task not found", http.StatusBadRequest)
		return
	}

	// 2. Spin or fetch the user's container
	username := r.Header.Get("X-Username")
	containerID, err := container.GetOrCreateContainer(userID, username)
	if err != nil {
		log.Printf("Error spinning/fetching container for user %s (%s): %v", username, userID, err)
		http.Error(w, "Sandbox container error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Verify the task in the container
	verified, err := game.VerifyTask(containerID, *targetTask, req.UserInput)
	if err != nil {
		log.Printf("Verification error for user %s, task %d: %v", userID, targetTask.ID, err)
		http.Error(w, "Verification engine error", http.StatusInternalServerError)
		return
	}

	if !verified {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(VerifyResponse{
			Success:        false,
			Message:        "Verification failed. Double check your commands and try again!",
			CurrentChapter: currentChapter,
			CurrentTask:    currentTask,
		})
		return
	}

	// 4. If verified, advance progress
	nextChapter := currentChapter
	nextTask := currentTask + 1
	chapterCompleted := false
	gameCompleted := false

	// Check if we finished tasks in current chapter
	if nextTask > len(targetChapter.Tasks) {
		chapterCompleted = true
		// Save to completed chapters
		var completedList []int
		_ = json.Unmarshal(completedJSON, &completedList)

		// Append if not exists
		exists := false
		for _, v := range completedList {
			if v == currentChapter {
				exists = true
				break
			}
		}
		if !exists {
			completedList = append(completedList, currentChapter)
		}
		completedJSON, _ = json.Marshal(completedList)

		// Advance to next chapter
		nextChapter = currentChapter + 1
		nextTask = 1

		// Check if that was the last chapter
		if nextChapter > len(game.Chapters) {
			gameCompleted = true
			nextChapter = currentChapter // keep cap
			nextTask = len(targetChapter.Tasks)
		}
	}

	// Save back to DB
	_, err = db.DB.Exec(
		"UPDATE user_progress SET current_chapter = $1, current_task = $2, completed_chapters = $3, updated_at = CURRENT_TIMESTAMP WHERE user_id = $4",
		nextChapter, nextTask, completedJSON, userID,
	)
	if err != nil {
		log.Printf("Failed to update progress: %v", err)
		http.Error(w, "Database save error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(VerifyResponse{
		Success:          true,
		Message:          "Objective completed successfully!",
		CurrentChapter:   nextChapter,
		CurrentTask:      nextTask,
		ChapterCompleted: chapterCompleted,
		GameCompleted:    gameCompleted,
	})
}

func handleTerminalWS(w http.ResponseWriter, r *http.Request) {
	// For Websockets, JWT is passed in query params because xterm.js Websocket addon or browser WS API doesn't easily support headers.
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized: missing token query parameter", http.StatusUnauthorized)
		return
	}

	claims, err := auth.ValidateJWT(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Get or spin container for this user
	containerID, err := container.GetOrCreateContainer(claims.UserID, claims.Username)
	if err != nil {
		log.Printf("Failed to get/create container for WS: %v", err)
		http.Error(w, "Container spinner failure", http.StatusInternalServerError)
		return
	}

	// Pipe the socket to the container shell with connection tracking
	container.HandleTerminalWS(w, r, claims.UserID, claims.Username, containerID)
}

func handleUsersCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Printf("Failed to count users: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(fmt.Sprintf(`{"count": %d}`, count)))
}
