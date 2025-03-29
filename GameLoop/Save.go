package GameLoop

import (
	"Protect_The_King/Boards"
	"Protect_The_King/Pawns"
	"encoding/json"
	"fmt"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// GameState struktura przechowująca dane gry
type GameState struct {
	Player1Pawns      []Pawns.BasePawn `json:"player1_pawns"`
	Nodes             []Boards.Node    `json:"Boards.Nodes"`
	PlayerGold        int32            `json:"player_gold"`
	RollTickets       int              `json:"roll_tickets"`
	CurrentTurn       string           `json:"current_turn"`
	TurnCounter       int              `json:"turn_counter"`
	ResetSeed         int64            `json:"reset_seed"`
	SelectedPawn      *Pawns.BasePawn  `json:"selected_pawn"`
	SelectedPawnMove  *Pawns.BasePawn  `json:"selected_pawn_move"`
	CurrentPhase      int              `json:"current_phase"`
	PlacementPhase    int              `json:"placement_phase"`
	PawnSelected      bool             `json:"pawn_selected"`
	MenuActive        bool             `json:"menu_active"`
	PawnSelectionDone bool             `json:"pawn_selection_done"`
	ShopGenerated     bool             `json:"shop_generated"`
	CurrentNodeID     int              `json:"current_node_id"`
	AvailableFights   [][]Pawns.Fight  `json:"available_fights"`
	BossFights        [][]Pawns.Fight  `json:"boss_fights"`
}

func SaveGameState() {
	gameState := GameState{
		Player1Pawns:      Pawns.Player1Pawns,
		Nodes:             Boards.Nodes,
		CurrentNodeID:     -1,
		PlayerGold:        Boards.PlayerGold,
		RollTickets:       Boards.RollTickets,
		CurrentTurn:       currentTurn,
		TurnCounter:       turnCounter,
		ResetSeed:         resetSeed,
		SelectedPawn:      selectedPawn,
		SelectedPawnMove:  selectedPawnMove,
		CurrentPhase:      currentPhase,
		PlacementPhase:    placementPhase,
		PawnSelected:      pawnSelected,
		MenuActive:        menuActive,
		PawnSelectionDone: Boards.PawnSelectionDone,
		ShopGenerated:     Boards.ShopGenerated,
		AvailableFights:   Pawns.AvailableFights,
		BossFights:        Pawns.BossFights,
	}

	// **Zapisujemy ID aktywnego węzła**
	if Boards.CurrentNode != nil {
		for i := range Boards.Nodes {
			if &Boards.Nodes[i] == Boards.CurrentNode {
				gameState.CurrentNodeID = i
				break
			}
		}
	}

	file, err := os.Create("savegame.json")
	if err != nil {
		rl.TraceLog(rl.LogError, "Error saving game.")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(gameState)
	if err != nil {
		rl.TraceLog(rl.LogError, "Error encoding save data.")
	}
}

func LoadGameState(currentView *Boards.GameView, screenWidth, screenHeight int32) {
	file, err := os.Open("savegame.json")
	if err != nil {
		rl.TraceLog(rl.LogWarning, "No saved game found.")
		return
	}
	defer file.Close()

	var gameState GameState
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gameState)
	if err != nil {
		rl.TraceLog(rl.LogError, "Error loading saved game.")
		return
	}

	// **Przywracanie zmiennych**
	Pawns.Player1Pawns = gameState.Player1Pawns
	Boards.PlayerGold = gameState.PlayerGold
	Boards.RollTickets = gameState.RollTickets
	currentTurn = gameState.CurrentTurn
	turnCounter = gameState.TurnCounter
	resetSeed = gameState.ResetSeed
	selectedPawn = gameState.SelectedPawn
	selectedPawnMove = gameState.SelectedPawnMove
	currentPhase = gameState.CurrentPhase
	placementPhase = gameState.PlacementPhase
	pawnSelected = gameState.PawnSelected
	menuActive = gameState.MenuActive
	Boards.PawnSelectionDone = gameState.PawnSelectionDone
	Boards.ShopGenerated = gameState.ShopGenerated

	// **Przywracanie mapy węzłów**
	Boards.Nodes = gameState.Nodes
	Boards.CurrentNode = nil

	fmt.Printf("🔍 Wczytano %d węzłów\n", len(Boards.Nodes))

	// **Przypisanie aktywnego węzła na podstawie zapisanego ID**
	if gameState.CurrentNodeID >= 0 && gameState.CurrentNodeID < len(Boards.Nodes) {
		Boards.CurrentNode = &Boards.Nodes[gameState.CurrentNodeID]

		if Boards.CurrentNode.Completed && Boards.CurrentNode.Next != nil {
			Boards.CurrentNode.Next.Active = true
		}
	} else {
		rl.TraceLog(rl.LogError, "Error: Invalid active node ID!")
	}

	// **Przywracanie referencji między węzłami**!!!!!!!!!!!!!!!!!!!!
	for i := 0; i < len(Boards.Nodes)-1; i++ {
		Boards.Nodes[i].Next = &Boards.Nodes[i+1]
	}

	// Wczytanie zapisanych walk
	Pawns.AvailableFights = gameState.AvailableFights
	Pawns.BossFights = gameState.BossFights

	// **Odtwarzanie planszy z zapisanego stanu**
	if resetSeed != 0 {
		_, _, boardX, boardY := Boards.CalculateGameBoardSize(screenWidth, screenHeight)
		board = Boards.GenerateBoard(12, int32(float32(screenHeight)*0.8/12), boardX, boardY, nil, resetSeed)
	} else {
		rl.TraceLog(rl.LogWarning, "Warning: resetSeed is 0, generating new board!")
		_, _, boardX, boardY := Boards.CalculateGameBoardSize(screenWidth, screenHeight)
		board = Boards.GenerateBoard(12, int32(float32(screenHeight)*0.8/12), boardX, boardY, nil, time.Now().UnixNano())
	}

	*currentView = Boards.ViewGameBoard
}
