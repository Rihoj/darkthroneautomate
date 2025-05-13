package main

type Player struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Gold        int    `json:"gold"`
	Level       int    `json:"level"`
	ArmySize    int    `json:"armySize"`
	Units       []Unit `json:"units"`
	AttackTurns int    `json:"attackTurns"`
}

type Unit struct {
	UnitType string `json:"unitType"`
	Quantity int    `json:"quantity"`
}

type CurrentUserResponse struct {
	Player Player `json:"player"`
}

type UserPlayersListResponse []Player // Update to represent a top-level array of Player objects
type PlayersListResponse struct {
	Items []Player `json:"items"` // Match the "items" field in the API response
	// Add other fields if needed, e.g., meta information
}

type AttackResponse struct {
	IsAttackerVictor bool `json:"is_attacker_victor"`
}

type LoginResponse struct {
	Session struct {
		ID                string  `json:"id"`
		Email             string  `json:"email"`
		PlayerID          *string `json:"playerID"` // Nullable field
		HasConfirmedEmail bool    `json:"hasConfirmedEmail"`
		ServerTime        string  `json:"serverTime"`
	} `json:"session"`
	Token string `json:"token"`
}
