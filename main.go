package main

import (
	"encoding/json"
	"net/http"
)

// Мок-сервер авторизации
func main() {
	http.HandleFunc("/oauth/authorize", mockAuthorizeHandler)
	http.HandleFunc("/oauth/token", mockTokenHandler)
	http.HandleFunc("/api/userinfo", mockUserInfoHandler)

	http.ListenAndServe(":8089", nil)
}

// Обработчик для эндпоинта авторизации
func mockAuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	// Имитация редиректа с кодом авторизации
	redirectURI := r.URL.Query().Get("redirect_uri")
	state := r.URL.Query().Get("state")
	code := "mock_auth_code_12345" // Мок-код

	http.Redirect(w, r, redirectURI+"?code="+code+"&state="+state, http.StatusFound)
}

// Обработчик для получения токена
func mockTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"access_token":  "mock_access_token_67890",
		"token_type":    "Bearer",
		"expires_in":    3600,
		"refresh_token": "mock_refresh_token_13579",
	}
	json.NewEncoder(w).Encode(response)
}

// Обработчик для получения данных пользователя
func mockUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := map[string]interface{}{
		"sub":         "1234567890",
		"name":        "Иван Иванов",
		"given_name":  "Иван",
		"family_name": "Иванов",
		"birthdate":   "1990-01-01",
		"snils":       "123-456-789 00",
	}
	json.NewEncoder(w).Encode(userInfo)
}
