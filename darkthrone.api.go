package main

import (
	"log/slog"
	"os"
)

const (
	baseURL              = "https://api.darkthronereborn.com" // Ensure this is correct
	loginEndpoint        = "auth/login"
	currentUserEndpoint  = "auth/current-user"
	playersListEndpoint  = "auth/current-user/players"
	assumePlayerEndpoint = "auth/assume-player"
)

func login(logger *slog.Logger) string {
	logger.Info("Logging in...")

	email := os.Getenv("DARK_THRONE_EMAIL")
	if email == "" {
		logger.Error("Environment variable DARK_THRONE_EMAIL is not set")
		panic("DARK_THRONE_EMAIL not set")
	}

	password := os.Getenv("DARK_THRONE_PASSWORD")
	if password == "" {
		logger.Error("Environment variable DARK_THRONE_PASSWORD is not set")
		panic("DARK_THRONE_PASSWORD not set")
	}

	payload := map[string]string{
		"email":    email,
		"password": password,
	}
	response := makePostRequest[LoginResponse](logger, loginEndpoint, payload)
	token := response.Token
	logger.Info("Login successful. Token acquired.")
	return token
}

// func currentUser(logger *slog.Logger) (string, error) {
// 	logger.Debug("Fetching current user...")
// 	if token == "" {
// 		logger.Error("Token is not set. Please ensure login() is called before making requests.")
// 	}
// 	response := makeGetRequest[CurrentUserResponse](logger, currentUserEndpoint)
// 	if response.Player.ID == "" {
// 		err := fmt.Errorf("failed to retrieve current user ID")
// 		return "", err
// 	}
// 	log.Printf("Current user ID: %s\n", response.Player.ID)
// 	return response.Player.ID, nil
// }

// getCurrentPlayer retrieves the current player's information.
func getCurrentPlayer(logger *slog.Logger) Player {
	logger.Debug("Fetching current player...")
	if playerID != "" {
		response := makeGetRequest[CurrentUserResponse](logger, currentUserEndpoint)
		return response.Player
	}

	// Fetch players if playerID is not set
	response := makeGetRequest[UserPlayersListResponse](logger, playersListEndpoint)
	if len(response) == 0 { // Directly check the length of the slice
		logger.Error("No players found in the response from auth/current-user/players")
	}

	// Debug log to ensure players list is populated
	logger.Debug("Players list response: %+v\n", response)

	// Ensure playerID is set to the first player's ID
	playerID = response[0].ID // Access the first player directly
	if playerID == "" {
		logger.Error("Failed to set playerID from the players list")
	}

	payload := map[string]string{"playerID": playerID}
	assumeResponse := makePostRequest[CurrentUserResponse](logger, assumePlayerEndpoint, payload)
	return assumeResponse.Player
}
