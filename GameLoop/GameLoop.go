package GameLoop

import (
	"fmt"
	"math/rand"
	"time"

	"Protect_The_King/Boards"
	"Protect_The_King/Pawns"
	"Protect_The_King/menu"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var selectedPawnMove *Pawns.BasePawn // Wybrany pionek do przesunięcia
var currentTurn string = "Player 1"  // Zaczyna Gracz 1
var board [][]Boards.Tile
var resetSeed int64
var selectedPawn *Pawns.BasePawn // Aktualnie wybrany pionek
var currentPhase int = 1         // 1 = Rozstawienie, 2 = Walka
var placementPhase int = 1       // 1 = Gracz 1 rozstawia pionki, 2 = Gracz 2 rozstawia pionki
var turnCounter int = 1          // Zmienna przechowująca aktualną turę

var pawnSelected bool = false // Czy pionek został wybrany?
var menuActive bool = false   // Czy menu wyboru pionka jest aktywne?

var gameStarted bool = false // zmienna sterująca restartem gry

func RunGame(screenWidth, screenHeight int32) {

	// Załadowanie pionków dla graczy
	Pawns.LoadPawns()
	menu.LoadMenuAssets()
	Boards.LoadGameBoardAssets()

	// Ładowanie tekstur planszy
	tilesetPath := filepath.Join("Assets", "TilesetField.png")
	tileset, sections := Boards.LoadTileset(tilesetPath)
	defer rl.UnloadTexture(tileset) // Zwolnienie pamięci po zakończeniu gry

	// Zmienna sterująca widokami
	currentView := Boards.ViewMainMenu

	menuState := menu.MenuState{} // Przechowujemy stan menu poza pętlą

	// **Obliczenie rozmiaru planszy na podstawie wymiarów ekranu**
	boardWidth, boardHeight, boardX, boardY := Boards.CalculateGameBoardSize(screenWidth, screenHeight)

	// **Generowanie mapy gry**
	Boards.GenerateMap(boardWidth, boardHeight, 9) // Argument określa liczbę węzłów

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite) // Czyści ekran na kolor biały
		UpdateMusic()

		switch currentView {
		case Boards.ViewMainMenu:
			menuState = menu.ShowMenu(screenWidth, screenHeight) // Aktualizacja menu
			PlayMusic("menu")
			if menuState.Exit {
				rl.EndDrawing()
				return
			}
			if menuState.GameRunning {
				StartNewGame(screenWidth, screenHeight) // Rozpoczęcie nowej gry
				currentView = Boards.ViewPawnSelection
				gameStarted = true
				StopMusic()
			}
			if menuState.InOptions {
				currentView = Boards.ViewOptions
			}
			if menuState.LoadGame { // Wczytywanie gry
				LoadGameState(&currentView, screenWidth, screenHeight)
				gameStarted = true
				StopMusic()
			}
		case Boards.ViewPawnSelection:
			PlayMusic("selection")
			Boards.ShowInitialPawnSelectionMenu(screenWidth, screenHeight)
			if Boards.PawnSelectionDone { // Przechodzimy do planszy gry tylko po zatwierdzeniu
				currentView = Boards.ViewGameBoard
				StopMusic()
			}

		case Boards.ViewOptions:
			newResolution, backToMenu := menu.ShowOptions(screenWidth, screenHeight)
			if newResolution != nil {
				screenWidth, screenHeight = newResolution[0], newResolution[1]
				rl.SetWindowSize(int(screenWidth), int(screenHeight)) // Aktualizacja rozdzielczości
			}
			if backToMenu {
				currentView = Boards.ViewMainMenu
			}

		case Boards.ViewGameBoard:
			PlayMusic("GameBoard")
			// Rysowanie interfejsu planszy gry
			nextView := Boards.DrawGameLayout(screenWidth, screenHeight, boardWidth, boardHeight, boardX, boardY, currentView)

			// Obsługa przejść między ekranami
			if nextView == Boards.ViewShopBoard {
				currentView = Boards.ViewShopBoard
			}
			if nextView == Boards.ViewFightBoard {
				StartNewFight(screenWidth, screenHeight) // **Uruchamiamy nową walkę**
				currentView = Boards.ViewFightBoard
				StopMusic()
			}
			if nextView == Boards.ViewMainMenu {
				currentView = Boards.ViewMainMenu
				gameStarted = false
				StopMusic()
			}
			if nextView == Boards.ViewWinScreen {
				currentView = Boards.ViewWinScreen // Przejście do ekranu wygranej
				StopMusic()

			}

		case Boards.ViewFightBoard:
			PlayMusic("Fight")
			// Ustalanie rozmiaru pojedynczej komórki planszy
			cellSize := int32(float32(screenHeight) * 0.8 / 12)

			// Rysowanie planszy walki
			DrawFightBoard(screenWidth, screenHeight, tileset, sections, board, cellSize)
			// Rysowanie logów walki
			Boards.DrawBattleLogs(screenWidth, screenHeight)

			// 1 - FAZA ROZSTAWIANIA
			if currentPhase == 1 {
				HandlePlacementPhase(screenWidth, screenHeight, cellSize, board)
			} else if currentPhase == 2 {
				// Faza walki (z systemem tur)
				HandleFightPhase(board, cellSize, &currentView)
			}

			// Rysowanie pionków na planszy
			Pawns.UpdateAnimations()
			Pawns.DrawPawns(cellSize, board[0][0].PosX, board[0][0].PosY)

			// Dodaj tu wywołanie rysowania przycisku "Pomoc"
			if Boards.DrawHelpButton(screenWidth, screenHeight) {
				// Obsługa kliknięcia przycisku "Pomoc"
				Boards.HelpWindowActive = true
			}
			// Dodanie przycisku "Powrót" do planszy gry
			if Boards.DrawBackButton(screenWidth, screenHeight) {
				ResetFightBoard()
				currentView = Boards.ViewGameBoard
				Boards.CompleteNode(&currentView) // Oznacza etap walki jako zakończony
				SaveGameState()
				StopMusic()

			}

		case Boards.ViewShopBoard:
			HandleShop(screenWidth, screenHeight, &currentView) // Obsługa sklepu

		case Boards.ViewWinScreen:
			currentView = Boards.DrawWinScreen(screenWidth, screenHeight) // Rysowanie ekranu wygranej
			StopMusic()

		case Boards.ViewLoseScreen:
			currentView = Boards.DrawLoseScreen(screenWidth, screenHeight) // Obsługa ekranu przegranej
			StopMusic()

		}

		rl.EndDrawing() // Zakończenie rysowania ramki
	}

	// Zwolnienie zasobów po zakończeniu gry
	if rl.WindowShouldClose() {
		Boards.UnloadShopTextures() // Zwolnienie pamięci tekstur sklepu
		UnloadMusic()

	}
}

func DrawFightBoard(screenWidth, screenHeight int32, tileset rl.Texture2D, sections map[string]rl.Rectangle, board [][]Boards.Tile, cellSize int32) bool {

	if Boards.FightBoardBackgroundLoaded {
		source := rl.Rectangle{X: 0, Y: 0, Width: float32(Boards.FightBoardBackground.Width), Height: float32(Boards.FightBoardBackground.Height)}
		dest := rl.Rectangle{X: 0, Y: 0, Width: float32(screenWidth), Height: float32(screenHeight)}
		rl.DrawTexturePro(Boards.FightBoardBackground, source, dest, rl.Vector2{}, 0, rl.White)
	}

	// Rysowanie pól planszy
	for _, row := range board {
		for _, tile := range row {
			Boards.DrawTileTexture(tileset, tile.SourceRect, tile.PosX, tile.PosY, cellSize)
		}
	}

	DrawVictoryConditions(screenWidth, screenHeight)

	// Pobieranie pozycji myszy
	mouseX := rl.GetMouseX()
	mouseY := rl.GetMouseY()

	// Sprawdzanie, czy mysz znajduje się na planszy i obliczanie współrzędnych pola
	var hoveredX, hoveredY int32 = -1, -1
	for y, row := range board {
		for x, tile := range row {
			if mouseX >= int32(tile.PosX) && mouseX < int32(tile.PosX+cellSize) &&
				mouseY >= int32(tile.PosY) && mouseY < int32(tile.PosY+cellSize) {
				hoveredX = int32(x)
				hoveredY = int32(y)
				break
			}
		}
	}

	// Wyświetlanie współrzędnych w lewym dolnym rogu
	if hoveredX != -1 && hoveredY != -1 {
		positionText := fmt.Sprintf("x: %d, y: %d", hoveredX, hoveredY)
		rl.DrawText(positionText, 10, screenHeight-30, 20, rl.Black)
	}
	//Licznik Tur
	rl.DrawText(fmt.Sprintf("Turn: %d", turnCounter), screenWidth-150, 40, 30, rl.Black)

	return false
}

func ShowPawnSelectionMenu(screenWidth, screenHeight int32, availablePawns []Pawns.BasePawn) *Pawns.BasePawn {
	menuWidth := int32(float32(screenWidth) * 0.3)   // Szerokość menu
	menuHeight := int32(float32(screenHeight) * 0.4) // Wysokość menu
	menuX := int32(10)
	menuY := int32(10)

	rl.DrawRectangle(menuX, menuY, menuWidth, menuHeight, rl.Gray)
	rl.DrawText("Select a pawn to place:", menuX+10, menuY+10, 20, rl.White)

	columns := 2   // Liczba kolumn
	xOffset := 120 // Odstęp poziomy między kolumnami
	yOffset := 40  // Odstęp pionowy między wierszami
	startX := menuX + 20
	startY := menuY + 50

	for i, pawn := range availablePawns {
		col := i % columns // Określa kolumnę
		row := i / columns // Określa wiersz

		pawnX := startX + int32(col)*int32(xOffset)
		pawnY := startY + int32(row)*int32(yOffset)

		rl.DrawText(pawn.Type, pawnX, pawnY, 20, rl.White)

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			mouseX := rl.GetMouseX()
			mouseY := rl.GetMouseY()

			if mouseX > pawnX && mouseX < pawnX+int32(xOffset) && mouseY > pawnY && mouseY < pawnY+20 {
				fmt.Printf("Selected pawn: %s\n", pawn.Type)
				return &availablePawns[i]
			}
		}
	}

	return nil
}

func DrawReadyButton(screenWidth, screenHeight int32) bool {
	buttonWidth := int32(200)
	buttonHeight := int32(50)
	buttonX := screenWidth - buttonWidth - 20
	buttonY := screenHeight - buttonHeight - 20

	rl.DrawRectangle(buttonX, buttonY, buttonWidth, buttonHeight, rl.DarkGreen)
	rl.DrawText("Ready", buttonX+50, buttonY+15, 20, rl.White)

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mouseX := rl.GetMouseX()
		mouseY := rl.GetMouseY()
		if mouseX > buttonX && mouseX < buttonX+buttonWidth && mouseY > buttonY && mouseY < buttonY+buttonHeight {
			return true
		}
	}
	return false
}

func HandlePlacementPhase(screenWidth, screenHeight int32, cellSize int32, board [][]Boards.Tile) {

	// Jeśli AI powinno się rozstawiać, uruchamiamy AI Placement
	if placementPhase == 2 {
		AutoPlaceAIPawns(board)
		return
	}

	// Jeśli menu wyboru pionków jest aktywne
	if menuActive {
		if placementPhase == 1 {
			selectedPawn = ShowPawnSelectionMenu(screenWidth, screenHeight, Pawns.AvailablePawnsP1)
		} else {
			selectedPawn = ShowPawnSelectionMenu(screenWidth, screenHeight, Pawns.AvailablePawnsP2)
		}

		// Jeśli wybrano pionek, zamknij menu i przejdź dalej
		if selectedPawn != nil {
			logMessage := fmt.Sprintf("Selected pawn: %s", selectedPawn.Type)
			Boards.AddBattleLog(logMessage)
			pawnSelected = true
			menuActive = false
		}
		return
	}

	// Pierwsze kliknięcie - aktywacja menu wyboru pionka
	if !pawnSelected && !menuActive && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		menuActive = true
		return
	}

	// Drugie kliknięcie - ustawienie pionka na planszy
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) && pawnSelected && selectedPawn != nil {
		mouseX := rl.GetMouseX()
		mouseY := rl.GetMouseY()

		var targetX, targetY int32 = -1, -1
		for y, row := range board {
			for x, tile := range row {
				if mouseX >= int32(tile.PosX) && mouseX < int32(tile.PosX+cellSize) &&
					mouseY >= int32(tile.PosY) && mouseY < int32(tile.PosY+cellSize) {
					targetX = int32(x)
					targetY = int32(y)
					break
				}
			}
		}

		// **SPRAWDZENIE OGRANICZEŃ STREF STARTOWYCH**
		if placementPhase == 1 && targetY <= 8 { // Gracz 1 może umieszczać tylko w rzędach 9, 10, 11
			Boards.AddBattleLog("Error: Pawns can only be placed in the bottom 3 rows for Player 1!")
			return
		}
		if placementPhase == 2 && targetY >= 3 { // Gracz 2 może umieszczać tylko w rzędach 0, 1, 2
			Boards.AddBattleLog("Error: Pawns can only be placed in the top 3 rows for Player 2!")
			return
		}

		// **SPRAWDZENIE CZY POLE JEST PRZECHODNIE**
		if !board[targetY][targetX].Walkable {
			Boards.AddBattleLog("Error: This field is not passable!")
			return
		}

		// **SPRAWDZENIE CZY POLE JEST ZAJĘTE**
		if Pawns.IsTileOccupied(targetX, targetY, Pawns.PawnsOnBoard) {
			Boards.AddBattleLog("Error: The field is already occupied!")
			return
		}

		// **UMIESZCZANIE PIONKA**
		selectedPawn.X = targetX
		selectedPawn.Y = targetY
		Pawns.PawnsOnBoard = append(Pawns.PawnsOnBoard, *selectedPawn)

		logMessage := fmt.Sprintf("Pawn %s placed at (%d, %d)", selectedPawn.Type, targetX, targetY)
		Boards.AddBattleLog(logMessage)

		// **USUWANIE PIONKA Z LISTY DOSTĘPNYCH**
		if placementPhase == 1 {
			Pawns.RemovePawnFromAvailableList(&Pawns.AvailablePawnsP1, selectedPawn.ID)
			if len(Pawns.AvailablePawnsP1) == 0 {
				logMessage = "Player 1 has finished setting up. Now it's AI turn!"
				Boards.AddBattleLog(logMessage)
				placementPhase = 2
			}
		} else {
			Pawns.RemovePawnFromAvailableList(&Pawns.AvailablePawnsP2, selectedPawn.ID)
			if len(Pawns.AvailablePawnsP2) == 0 {
				logMessage = "AI has finished setting up. Moving on to the battle phase!"
				Boards.AddBattleLog(logMessage)
				currentPhase = 2
			}
		}

		// Reset wyboru pionka
		selectedPawn = nil
		pawnSelected = false
	}
}

func HandleFightPhase(board [][]Boards.Tile, cellSize int32, currentView *Boards.GameView) {
	var targetX, targetY int32 = -1, -1

	if currentTurn == "Player 2" {
		MakeAIMove() // Komputer wykonuje swój ruch automatycznie
		return
	}

	// **🔹 Sprawdzanie, czy King (P1) lub Boss (P2) są pod szachem**
	king, _ := FindKingAndBoss(currentTurn)
	var threatenedPawn *Pawns.BasePawn

	if currentTurn == "Player 1" && king != nil && IsPawnUnderThreat(king) {
		threatenedPawn = king
	}

	// **Jeśli King (P1) lub Boss (P2) są pod szachem, sprawdzamy ich ruchy**
	if threatenedPawn != nil {
		validMoves := GetValidMoves(threatenedPawn, board, currentTurn)

		// **Jeśli brak ruchów → szach-mat i koniec gry**
		if len(validMoves) == 0 {
			fmt.Printf("%s of Player %s is in check and cannot move! Game over!\n", threatenedPawn.Type, currentTurn)
			EndFight(swapTurn(currentTurn), currentView)
			return
		}

		// **🔹 Wymuszenie ruchu Kinga (P1) lub Bossa (P2) – inne pionki nie mogą się ruszać**
		fmt.Printf("%s of Player %s is in check! You must move it.\n", threatenedPawn.Type, currentTurn)
		Boards.AddBattleLog(fmt.Sprintf("%s of Player %s is in check! You must move it.", threatenedPawn.Type, currentTurn))

		// Ustawiamy, że tylko King/Boss może się ruszać
		selectedPawnMove = threatenedPawn
	}

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mouseX := rl.GetMouseX()
		mouseY := rl.GetMouseY()

		// Znalezienie klikniętej pozycji na planszy
		for y, row := range board {
			for x, tile := range row {
				if mouseX >= int32(tile.PosX) && mouseX < int32(tile.PosX+cellSize) &&
					mouseY >= int32(tile.PosY) && mouseY < int32(tile.PosY+cellSize) {
					targetX = int32(x)
					targetY = int32(y)
					break
				}
			}
		}

		if targetX != -1 && targetY != -1 {
			// **Wybór pionka do ruchu, jeśli King/Boss nie jest pod szachem**
			if selectedPawnMove == nil {
				for i := range Pawns.PawnsOnBoard {
					if Pawns.PawnsOnBoard[i].X == targetX && Pawns.PawnsOnBoard[i].Y == targetY &&
						Pawns.PawnsOnBoard[i].Owner == currentTurn {

						// **Jeśli King/Boss jest pod szachem, inne pionki nie mogą się ruszać**
						if threatenedPawn != nil && &Pawns.PawnsOnBoard[i] != threatenedPawn {
							fmt.Printf("%s is in check! You must move it.\n", threatenedPawn.Type)
							Boards.AddBattleLog(fmt.Sprintf("%s is in check! You must move it.", threatenedPawn.Type))
							return
						}

						selectedPawnMove = &Pawns.PawnsOnBoard[i]
						logMessage := fmt.Sprintf("Selected pawn to move: %s (ID: %d)", selectedPawnMove.Type, selectedPawnMove.ID)
						Boards.AddBattleLog(logMessage)
						break
					}
				}
			} else {
				// Pobranie ID atakującego pionka
				attackingPawnID := selectedPawnMove.ID

				// Sprawdzenie, czy na polu jest przeciwnik
				enemyID, enemyIndex := FindEnemyPawnID(targetX, targetY, currentTurn)

				if IsValidMove(selectedPawnMove, targetX, targetY, board, currentTurn) {
					if enemyID != -1 && enemyIndex != -1 {
						logMessage := fmt.Sprintf("Pawn %s (ID: %d) attacked the enemy (ID: %d)!", selectedPawnMove.Type, attackingPawnID, enemyID)
						Boards.AddBattleLog(logMessage)

						// Usunięcie pionka przeciwnika
						Pawns.RemovePawnByID(enemyID)

						// Aktualizacja pionka atakującego
						movingPawnIndex := FindPawnIndexByID(attackingPawnID)
						if movingPawnIndex != -1 {
							selectedPawnMove = &Pawns.PawnsOnBoard[movingPawnIndex]
						}
					}

					// Przeniesienie pionka
					selectedPawnMove.X = targetX
					selectedPawnMove.Y = targetY
					logMessage := fmt.Sprintf("Pawn %s moved to (%d, %d)", selectedPawnMove.Type, targetX, targetY)
					Boards.AddBattleLog(logMessage)

					// **🔹 Sprawdzenie, czy King/Boss jest teraz w szachu bez ruchu**
					if CheckWinCondition(currentView) {
						return
					}

					// Zmiana tury
					selectedPawnMove = nil
					currentTurn = swapTurn(currentTurn)

					logMessage = fmt.Sprintf("Now it's %s's turn", currentTurn)
					Boards.AddBattleLog(logMessage)
				} else {
					logMessage := "Invalid move!"
					Boards.AddBattleLog(logMessage)
				}
			}
		}
	}
}

// Obsługuje widok sklepu i przejście z powrotem na planszę gry
func HandleShop(screenWidth, screenHeight int32, currentView *Boards.GameView) {
	if Boards.DrawShopBoard(screenWidth, screenHeight, currentView) {
		Boards.CompleteNode(currentView) // Oznacza węzeł jako ukończony i wraca do GameBoard
		SaveGameState()
	}
}

func CheckWinCondition(currentView *Boards.GameView) bool {
	p1Alive := false
	p2Alive := false

	var king, boss *Pawns.BasePawn

	// Sprawdzamy obecność Kinga i Bossa dla obu graczy
	for i := range Pawns.PawnsOnBoard {
		pawn := &Pawns.PawnsOnBoard[i]
		if pawn.Owner == "Player 1" {
			p1Alive = true
			if pawn.Type == "King" {
				king = pawn
			}
		} else if pawn.Owner == "Player 2" {
			p2Alive = true
			if pawn.Type == "Boss" {
				boss = pawn
			}
		}
	}

	// **Nowy warunek: jeśli minęło 50 tur, Gracz 1 automatycznie przegrywa**
	if turnCounter >= 40 {
		EndFight("Player 2", currentView)
		return true
	}

	// **Sprawdzenie, czy King Gracza 1 dotarł do górnej krawędzi**
	if king != nil && king.Y == 0 {
		EndFight("Player 1", currentView)
		Boards.CompleteNode(currentView)
		SaveGameState()
		return true
	}

	// **Jeśli jeden z graczy stracił wszystkie pionki, drugi wygrywa**
	if !p1Alive {
		EndFight("Player 2", currentView)
		return true
	}
	if !p2Alive {
		EndFight("Player 1", currentView)
		Boards.CompleteNode(currentView)
		SaveGameState()
		return true
	}

	// **🔹 Sprawdzenie, czy King lub Boss są szachowani i nie mają ruchu**
	if king != nil && IsPawnUnderThreat(king) && len(GetValidMoves(king, board, king.Owner)) == 0 {
		EndFight("Player 2", currentView)
		return true
	}
	if boss != nil && IsPawnUnderThreat(boss) && len(GetValidMoves(boss, board, boss.Owner)) == 0 {
		EndFight("Player 1", currentView)
		Boards.CompleteNode(currentView)
		SaveGameState()
		return true
	}

	return false
}

func EndFight(winner string, currentView *Boards.GameView) {
	fmt.Printf("%s wins the battle!\n", winner)

	// **Resetowanie pola walki**
	ResetFightBoard()

	// **Ustawienie odpowiedniego ekranu po wygranej**
	if winner == "Player 1" {
		reward := rand.Int31n(21) + 20 // Losowa wartość od 20 do 40
		Boards.PlayerGold += reward

		logMessage := fmt.Sprintf(" Reward received: %d G for winning the battle!", reward)
		Boards.AddBattleLog(logMessage)

		fmt.Printf(" Player 1 earned %d G! Current gold: %d G\n", reward, Boards.PlayerGold)

		*currentView = Boards.ViewGameBoard // Powrót do ekranu gry
		StopMusic()

	} else {
		*currentView = Boards.ViewLoseScreen // Przejście do ekranu przegranej
		StopMusic()

	}
}

func DrawVictoryConditions(screenWidth, screenHeight int32) {
	panelHeight := int32(float32(screenHeight) * 0.05) // 5% wysokości ekranu na panel
	panelWidth := screenWidth
	panelX := int32(0)
	panelY := int32(0) // Umiejscowienie na górze ekranu

	// Rysowanie panelu
	panelRect := rl.Rectangle{
		X:      float32(panelX),
		Y:      float32(panelY),
		Width:  float32(panelWidth),
		Height: float32(panelHeight),
	}
	rl.DrawRectangleRec(panelRect, rl.LightGray)

	// Tekst warunków zwycięstwa
	text := "Victory conditions: Capture all opponent's pawns before turn 40 Move the 'King' to Y=0"
	textSize := rl.MeasureText(text, 20)
	textX := (screenWidth - textSize) / 2
	textY := panelY + 10

	rl.DrawText(text, textX, textY, 20, rl.Black)
}

// Reset Pola Walki
func ResetFightBoard() {
	Pawns.PawnsOnBoard = []Pawns.BasePawn{} // Czyszczenie listy pionków na planszy
	Pawns.InitializeAvailablePawns()        // Ponowne załadowanie dostępnych pionków
	currentPhase = 1                        // Resetowanie do fazy rozstawienia
	turnCounter = 1                         // Resetowanie licznika tur
	currentTurn = "Player 1"                // Resetowanie tury do Gracza 1
	selectedPawnMove = nil                  // Resetowanie wybranego pionka do ruchu
}

func IsValidMove(pawn *Pawns.BasePawn, targetX, targetY int32, board [][]Boards.Tile, player string) bool {
	if pawn == nil {
		return false
	}

	// Sprawdzenie, czy celowe pole istnieje na planszy
	if targetX < 0 || targetY < 0 || int(targetX) >= len(board) || int(targetY) >= len(board[0]) {
		return false
	}

	// Sprawdzenie, czy pole jest dostępne do ruchu
	if !board[targetY][targetX].Walkable {
		fmt.Printf("Cannot move to (%d, %d) - the field is impassable!\n", targetX, targetY)
		return false
	}

	// Sprawdzenie, czy na polu znajduje się własny pionek
	for _, otherPawn := range Pawns.PawnsOnBoard {
		if otherPawn.X == targetX && otherPawn.Y == targetY && otherPawn.Owner == player {
			return false
		}
	}

	// **🔹 Sprawdzenie dla King/Boss - nie mogą wejść na zagrożone pole**
	if pawn.Type == "King" || pawn.Type == "Boss" {
		if IsTileThreatened(targetX, targetY, player, board) {
			fmt.Printf("%s cannot move to (%d, %d) - it is threatened by an attack!\n", pawn.Type, targetX, targetY)
			return false
		}
	}

	// **🔹 Specjalnych Pionkow**
	if pawn.Type == "Knight" {
		// Może poruszać się tylko w tej samej kolumnie lub wierszu
		if pawn.X != targetX && pawn.Y != targetY {
			return false
		}

		// Sprawdzenie czy droga jest wolna (nie może przechodzić przez inne pionki ani pola nie-walkable)
		if pawn.X == targetX { // Ruch w pionie
			step := int32(1)
			if targetY < pawn.Y {
				step = -1
			}
			for y := pawn.Y + step; y != targetY; y += step {
				if Pawns.IsTileOccupied(pawn.X, y, Pawns.PawnsOnBoard) || !board[y][pawn.X].Walkable {
					fmt.Printf(" Move blocked by a pawn or terrain at (%d, %d)\n", pawn.X, y)
					return false
				}
			}
		} else if pawn.Y == targetY { // Ruch w poziomie
			step := int32(1)
			if targetX < pawn.X {
				step = -1
			}
			for x := pawn.X + step; x != targetX; x += step {
				if Pawns.IsTileOccupied(x, pawn.Y, Pawns.PawnsOnBoard) || !board[pawn.Y][x].Walkable {
					fmt.Printf(" Move blocked by a pawn or terrain at (%d, %d)\n", x, pawn.Y)
					return false
				}
			}
		}
		return true
	}
	if pawn.Type == "Racoon" { // Sprawdzenie, czy ruch jest o 1, 2 lub 3 pola w pionie lub poziomie
		if (Abs(targetX-pawn.X) >= 1 && Abs(targetX-pawn.X) <= 3 && targetY == pawn.Y) ||
			(Abs(targetY-pawn.Y) >= 1 && Abs(targetY-pawn.Y) <= 3 && targetX == pawn.X) {

			stepX := int32(0)
			stepY := int32(0)

			// Ustalanie kierunku ruchu
			if targetX > pawn.X {
				stepX = 1
			} else if targetX < pawn.X {
				stepX = -1
			}

			if targetY > pawn.Y {
				stepY = 1
			} else if targetY < pawn.Y {
				stepY = -1
			}

			// Sprawdzenie czy droga jest wolna
			for i := int32(1); i <= Abs(targetX-pawn.X+targetY-pawn.Y); i++ {
				checkX := pawn.X + (stepX * i)
				checkY := pawn.Y + (stepY * i)

				if Pawns.IsTileOccupied(checkX, checkY, Pawns.PawnsOnBoard) || !board[checkY][checkX].Walkable {
					fmt.Printf(" Move blocked by a pawn or terrain at (%d, %d)\n", checkX, checkY)
					return false
				}
			}

			return true
		}
		return false

	}
	if pawn.Type == "LionWarrior" {
		allowedMoves, exists := Pawns.PawnMoves["LionWarrior"]
		if !exists {
			fmt.Printf("No movement config found for LionWarrior!\n")
			return false
		}

		deltaX := targetX - pawn.X
		deltaY := targetY - pawn.Y

		isAllowed := false
		for _, move := range allowedMoves {
			if move.DX == deltaX && move.DY == deltaY {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return false
		}

		stepX, stepY := int32(0), int32(0)

		if deltaX != 0 {
			if deltaX > 0 {
				stepX = 1
			} else {
				stepX = -1
			}
		}
		if deltaY != 0 {
			if deltaY > 0 {
				stepY = 1
			} else {
				stepY = -1
			}
		}

		moveDistance := deltaX
		if deltaY > deltaX {
			moveDistance = deltaY
		}

		for i := int32(1); i < moveDistance; i++ {
			checkX := pawn.X + (stepX * i)
			checkY := pawn.Y + (stepY * i)

			if Pawns.IsTileOccupied(checkX, checkY, Pawns.PawnsOnBoard) || !board[checkY][checkX].Walkable {
				fmt.Printf("LionWarrior move blocked by a pawn at (%d, %d)\n", checkX, checkY)
				return false
			}
		}

		return true

	}
	if pawn.Type == "Master" {
		// Sprawdzenie, czy ruch jest po skosie
		if Abs(targetX-pawn.X) != Abs(targetY-pawn.Y) {
			return false
		}

		// Sprawdzenie czy droga jest wolna
		stepX := int32(1)
		stepY := int32(1)

		if targetX < pawn.X {
			stepX = -1
		}

		if targetY < pawn.Y {
			stepY = -1
		}

		moveDistance := Abs(targetX - pawn.X) // Długość ruchu
		for i := int32(1); i <= moveDistance; i++ {
			checkX := pawn.X + (stepX * i)
			checkY := pawn.Y + (stepY * i)

			if Pawns.IsTileOccupied(checkX, checkY, Pawns.PawnsOnBoard) || !board[checkY][checkX].Walkable {
				fmt.Printf(" Move blocked by a pawn or terrain at (%d, %d)\n", checkX, checkY)
				return false
			}
		}
		return true
	}

	// **Sprawdzenie poprawnego ruchu dla każdego innego pionka**
	deltaX := targetX - pawn.X
	deltaY := targetY - pawn.Y

	return CanPawnAttack(pawn.Type, deltaX, deltaY)
}

func FindKingAndBoss(player string) (*Pawns.BasePawn, *Pawns.BasePawn) {
	var king, boss *Pawns.BasePawn

	for i := range Pawns.PawnsOnBoard {
		if Pawns.PawnsOnBoard[i].Owner == player {
			if Pawns.PawnsOnBoard[i].Type == "King" {
				king = &Pawns.PawnsOnBoard[i]
			}
			if Pawns.PawnsOnBoard[i].Type == "Boss" {
				boss = &Pawns.PawnsOnBoard[i]
			}
		}
	}

	return king, boss
}

func IsTileThreatened(x, y int32, player string, board [][]Boards.Tile) bool {
	for _, pawn := range Pawns.PawnsOnBoard {
		if pawn.Owner != player {
			// **Nie wywołujemy `IsValidMove()` dla sprawdzania zagrożenia!**
			deltaX := x - pawn.X
			deltaY := y - pawn.Y

			// **Sprawdzamy mapę PawnMoves zamiast duplikować kod**
			if CanPawnAttack(pawn.Type, deltaX, deltaY) {
				return true
			}
		}
	}
	return false
}

func IsPawnUnderThreat(pawn *Pawns.BasePawn) bool {
	return IsTileThreatened(pawn.X, pawn.Y, pawn.Owner, board)
}
func GetValidMoves(pawn *Pawns.BasePawn, board [][]Boards.Tile, player string) []struct{ x, y int32 } {
	validMoves := []struct{ x, y int32 }{}

	// Sprawdzenie wszystkich sąsiednich pól
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			newX := pawn.X + int32(dx)
			newY := pawn.Y + int32(dy)
			if IsValidMove(pawn, newX, newY, board, player) {
				validMoves = append(validMoves, struct{ x, y int32 }{newX, newY})
			}
		}
	}

	return validMoves
}

func CanPawnAttack(pawnType string, deltaX, deltaY int32) bool {
	moves, exists := Pawns.PawnMoves[pawnType]
	if !exists {
		return false // Jeśli typ pionka nie istnieje w mapie
	}

	for _, move := range moves {
		if move.DX == deltaX && move.DY == deltaY {
			return true
		}
	}
	return false
}

func FindEnemyPawnID(x, y int32, currentPlayer string) (int, int) {
	for i, pawn := range Pawns.PawnsOnBoard {
		if pawn.X == x && pawn.Y == y && pawn.Owner != currentPlayer {
			return pawn.ID, i // Zwracamy ID pionka przeciwnika
		}
	}
	return -1, -1 // Nie znaleziono przeciwnika
}

func FindPawnIndexByID(pawnID int) int {
	for i, pawn := range Pawns.PawnsOnBoard {
		if pawn.ID == pawnID {
			return i
		}
	}
	return -1 // Nie znaleziono pionka
}

func swapTurn(current string) string {
	if current == "Player 1" {
		return "Player 2"
	} else {
		turnCounter++ // Zwiększamy licznik tur
		return "Player 1"
	}

}

// **StartNewFight - Rozpoczyna nową losową walkę**
func StartNewFight(screenWidth, screenHeight int32) {

	// **Sprawdzenie, czy to bossfight** (czyli czy aktualny węzeł jest ostatni)
	isBossFight := (Boards.CurrentNode != nil && Boards.CurrentNode.Next == nil)

	// **Resetowanie zmiennych walki** (ale NIE resetujemy całej gry)
	resetSeed = time.Now().UnixNano()
	currentPhase = 1
	placementPhase = 1
	currentTurn = "Player 1"
	selectedPawn = nil
	selectedPawnMove = nil
	turnCounter = 1 // Reset liczby tur

	// **Reset pionków na planszy**
	Pawns.PawnsOnBoard = []Pawns.BasePawn{}
	Pawns.InitializeAvailablePawns()

	// **Losowanie przeciwników zależnie od tego, czy to bossfight**
	if isBossFight {
		Pawns.Player2Pawns = Pawns.GetRandomBossFight()
	} else {
		Pawns.Player2Pawns = Pawns.GetRandomFight()
	}

	// **Przygotowanie dostępnych pionków do rozstawienia**
	Pawns.InitializeAvailablePawns()

	// **Ładowanie tilesetu i sekcji planszy**
	tileset, sections := Boards.LoadTileset("Assets/TilesetField.png")
	defer rl.UnloadTexture(tileset) // Zwolnienie tekstury po zakończeniu gry

	// **Generowanie nowej mapy walki**
	_, _, boardX, boardY := Boards.CalculateGameBoardSize(screenWidth, screenHeight)
	board = Boards.GenerateBoard(12, int32(float32(screenHeight)*0.8/12), boardX, boardY, sections, resetSeed)

}

// ** Funkcja StartNewGame - Resetuje całą grę do początkowego stanu**
func StartNewGame(screenWidth, screenHeight int32) {

	// **Resetowanie zmiennych ekonomicznych**
	Boards.PlayerGold = 100      // Reset złota gracza
	Boards.RollTickets = 1       // Reset przerzutów sklepu
	Boards.ShopGenerated = false // Reset flagi sklepu (losowanie nowego zestawu)

	// Resetowanie zmiennych gry
	selectedPawnMove = nil
	currentTurn = "Player 1"
	turnCounter = 1 // Resetowanie licznika tur
	resetSeed = time.Now().UnixNano()
	selectedPawn = nil
	currentPhase = 1
	placementPhase = 1
	pawnSelected = false
	menuActive = false
	Boards.PawnSelectionDone = false

	// **Ładowanie pionków do gry z konfiguracji**
	Pawns.LoadPawnsIntoGame()

	// **Inicjalizacja pionków do rozstawienia**
	Pawns.InitializeAvailablePawns()

	// Generowanie nowej mapy
	boardWidth, boardHeight, _, _ := Boards.CalculateGameBoardSize(screenWidth, screenHeight)
	Boards.GenerateMap(boardWidth, boardHeight, 9) // Reset mapy

}
