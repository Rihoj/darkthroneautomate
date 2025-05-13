package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	minAttackTurns  = 10
	levelDifference = 7
	pageSize        = 100
	defaultUnitType = "soldier_1"
)

var token string
var playerID string
var playerList []Player

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO" // Default log level
	}

	var level slog.Level
	switch logLevel {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	logger.Info("Starting Dark Throne Automate...")
	if len(os.Args) < 5 {
		fmt.Printf("Usage: %s <should_attack> <start_page> <end_page> <gold_threshold>\n", os.Args[0])
		return
	}

	shouldAttack := os.Args[1] == "true"
	startPage, _ := strconv.Atoi(os.Args[2])
	endPage, _ := strconv.Atoi(os.Args[3])
	goldThreshold, _ := strconv.Atoi(os.Args[4])

	logger.Debug("Parsed arguments",
		"should_attack", shouldAttack,
		"start_page", startPage,
		"end_page", endPage,
		"gold_threshold", goldThreshold,
	)

	token = login(logger)
	if token == "" {
		logger.Error("Failed to acquire token. Exiting.")
		return
	}

	currentPlayer := getCurrentPlayer(logger)

	if currentPlayer.AttackTurns < minAttackTurns {
		logger.Error("Insufficient attack turns",
			"current_turns", currentPlayer.AttackTurns,
			"required_turns", minAttackTurns,
		)
		return
	}

	for page := startPage; page <= endPage; page++ {
		logger.Info("Fetching players", "page", page)
		playerList = append(playerList, getPlayers(logger, page)...)
	}

	if len(playerList) == 0 {
		logger.Info("No players found.")
		return
	}

	for _, player := range playerList {
		if player.ID == currentPlayer.ID {
			logger.Info("Skipping current player", "player_name", player.Name)
			continue
		}
		logger.Info("Evaluating player",
			"player_name", player.Name,
			"gold", player.Gold,
			"level", player.Level,
			"army_size", player.ArmySize,
		)

		if isValidTarget(currentPlayer, player, goldThreshold) {
			logger.Warn("Valid target found", "player_name", player.Name)
			if shouldAttack {
				attackPlayer(logger, player.ID, currentPlayer)
				currentPlayer.AttackTurns -= minAttackTurns
			}
		}
	}
	logger.Info("Execution completed.")
}

func getPlayers(logger *slog.Logger, page int) []Player {
	logger.Info("Fetching players from API", "page", page)
	response := makeGetRequest[PlayersListResponse](logger, fmt.Sprintf("players?page=%d&pageSize=%d", page, pageSize))
	return response.Items
}

func attackPlayer(logger *slog.Logger, targetID string, currentPlayer Player) {
	logger.Warn("Attacking player", "target_id", targetID)
	payload := map[string]interface{}{
		"targetID":    targetID,
		"attackTurns": 10,
	}
	response := makePostRequest[AttackResponse](logger, "attack", payload)
	if response.IsAttackerVictor {
		logger.Warn("Attack successful", "target_id", targetID)
	} else {
		logger.Warn("Attack failed", "target_id", targetID)
	}
}

func isValidTarget(self Player, target Player, goldThreshold int) bool {
	isValid := target.Gold > goldThreshold &&
		target.ArmySize < getUnitQuantity(self.Units, defaultUnitType) &&
		abs(self.Level-target.Level) <= levelDifference &&
		target.ID != self.ID
	return isValid
}

func getUnitQuantity(units []Unit, unitType string) int {
	for _, unit := range units {
		if unit.UnitType == unitType {
			return unit.Quantity
		}
	}
	return 0
}

func calculateOffense(logger *slog.Logger, units []Unit) int {
	total := 0
	for _, unit := range units {
		logger.Debug("Calculating offense", "unit_type", unit.UnitType, "quantity", unit.Quantity)
		if unit.UnitType == defaultUnitType {
			total += unit.Quantity
		}
	}
	return total
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
