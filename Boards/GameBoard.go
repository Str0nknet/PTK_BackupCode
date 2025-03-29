package Boards

import (
	"Protect_The_King/Pawns"
	"fmt"
	"path/filepath"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var scrollOffsetX int32 = 0       // Przesunięcie w poziomie
var scrollOffsetY int32 = 0       // Przesunięcie w pionie
var HelpWindowActive bool = false // Zmienna sterująca wyświetlaniem okna pomocy

// GameView określa aktualny widok gry
type GameView int

const (
	ViewMainMenu      GameView = iota // Widok menu głównego
	ViewGameBoard                     // Widok głównego panelu gry
	ViewFightBoard                    // Widok planszy walki
	ViewOptions                       // Widok opcji
	ViewShopBoard                     // Widok Sklepu
	ViewWinScreen                     //WinScrean
	ViewLoseScreen                    //LoseScrean
	ViewPawnSelection                 // Pawn Selection
)

// Struktura węzła ścieżki
type Node struct {
	X, Y      int32
	Next      *Node
	Active    bool
	Completed bool // Nowy znacznik zakończenia węzła
}

var Nodes []Node      // Tablica węzłów na mapie
var CurrentNode *Node // Aktualny aktywny węzeł

var Word1Background rl.Texture2D
var Word1BackgroundLoaded bool = false
var FightBoardBackground rl.Texture2D
var FightBoardBackgroundLoaded bool = false
var PawnSelectionBackground rl.Texture2D
var PawnSelectionBackgroundLoaded bool = false

func LoadGameBoardAssets() {
	if !Word1BackgroundLoaded {
		path := filepath.Join("Assets", "Backgrounds", "Word_1_Back.png")
		Word1Background = rl.LoadTexture(path)
		Word1BackgroundLoaded = true
	}
	if !PawnSelectionBackgroundLoaded {
		path := filepath.Join("Assets", "Backgrounds", "Pawn_Selection.png")
		PawnSelectionBackground = rl.LoadTexture(path)
		PawnSelectionBackgroundLoaded = true
	}
	if !FightBoardBackgroundLoaded {
		path := filepath.Join("Assets", "Backgrounds", "Fight_1_Back.png")
		FightBoardBackground = rl.LoadTexture(path)
		FightBoardBackgroundLoaded = true
	}
}

// **Rysowanie układu ekranu**
func DrawGameLayout(screenWidth, screenHeight, boardWidth, boardHeight, boardX, boardY int32, currentView GameView) GameView {

	leftPanelWidth := float32(screenWidth) * 0.225 // 22.5% szerokości ekranu
	rightPanelWidth := leftPanelWidth / 2          // Prawy panel jest połową lewego

	// **Nowy spójny margines, aby wyrównać wszystko**
	marginX := int32(0)

	// **Tworzenie sekcji GUI**
	leftPanel := rl.Rectangle{X: 0, Y: 0, Width: float32(boardX), Height: float32(screenHeight)}
	rightPanel := rl.Rectangle{X: float32(screenWidth) - rightPanelWidth, Y: 0, Width: rightPanelWidth, Height: float32(screenHeight)}
	topPanel := rl.Rectangle{X: float32(boardX - marginX), Y: 0, Width: float32(boardWidth + marginX), Height: float32(screenHeight) * 0.05}
	bottomPanel := rl.Rectangle{X: float32(boardX - marginX), Y: float32(screenHeight) - (float32(screenHeight) * 0.05), Width: float32(boardWidth + marginX), Height: float32(screenHeight) * 0.05}
	gameBoard := rl.Rectangle{X: float32(boardX), Y: float32(boardY), Width: float32(boardWidth), Height: float32(boardHeight)}

	// **Rysowanie interfejsu**
	rl.DrawRectangleRec(leftPanel, rl.DarkGray) // ✔️ Pokrywa całą lewą stronę
	rl.DrawRectangleRec(rightPanel, rl.Gray)    // ✔️ Pokrywa całą prawą stronę
	rl.DrawRectangleRec(topPanel, rl.LightGray) // ✔️ Zaczyna się od `boardX - marginX`
	rl.DrawRectangleRec(bottomPanel, rl.LightGray)

	rl.DrawRectangleLines(int32(gameBoard.X), int32(gameBoard.Y), int32(gameBoard.Width), int32(gameBoard.Height), rl.Red)

	//Pokrywa całą planszę gry tłem Word1
	source := rl.Rectangle{X: 0, Y: 0, Width: float32(Word1Background.Width), Height: float32(Word1Background.Height)}
	dest := gameBoard
	rl.DrawTexturePro(Word1Background, source, dest, rl.Vector2{}, 0, rl.White)

	// **Rysowanie informacji o graczu w lewym panelu**
	DrawLeftPanel(int32(leftPanel.Width), int32(leftPanel.Height))

	// **Rysowanie mapy i obsługa przewijania, jeśli jesteśmy w widoku gry**
	if currentView == ViewGameBoard {
		HandleScrolling(int32(gameBoard.Width), int32(gameBoard.Height))
		DrawMap(int32(gameBoard.X), int32(gameBoard.Y), &currentView)
	}

	// **Przycisk menu**
	if drawMenuAndHelpButtons(bottomPanel) {
		return ViewMainMenu
	}

	return currentView
}

func DrawLeftPanel(panelWidth, panelHeight int32) {
	// **Czyści tło panelu, aby uniknąć nakładania się elementów**
	rl.DrawRectangle(0, 0, panelWidth, panelHeight, rl.DarkGray)

	// 📌 Dostosowanie wartości do pełnego pokrycia `leftPanel`
	marginX := int32(15)     // Margines od lewej krawędzi
	marginY := int32(20)     // Margines od góry panelu
	lineSpacing := int32(35) // Odstęp między liniami tekstu
	textSize := int32(22)    // Rozmiar czcionki

	// **Wyświetlanie złota gracza**
	rl.DrawText(fmt.Sprintf("Gold: %d G", PlayerGold), marginX, marginY, textSize, rl.Yellow)

	// **Wyświetlanie ilości przerzutów**
	rl.DrawText(fmt.Sprintf("Shop Rolls: %d", RollTickets), marginX, marginY+lineSpacing, textSize, rl.White)

	// **Separator**
	rl.DrawText("-----------------", marginX, marginY+(2*lineSpacing), textSize, rl.White)

	// **Lista pionków gracza**
	rl.DrawText("Your Units:", marginX, marginY+(3*lineSpacing), textSize, rl.White)

	// **Dynamiczne rozmieszczanie pionków**
	unitY := marginY + (4 * lineSpacing)
	maxLines := (panelHeight - unitY) / lineSpacing // Ilość linii, jakie mieszczą się w panelu

	for i, pawn := range Pawns.Player1Pawns {
		if int32(i) >= maxLines { // Zapobiega wychodzeniu tekstu poza panel
			break
		}
		rl.DrawText(pawn.Type, marginX, unitY, textSize, rl.White)
		unitY += lineSpacing
	}
}

// **Generowanie mapy z wycentrowaną ścieżką**
func GenerateMap(boardWidth, boardHeight int32, length int32) {
	Nodes = make([]Node, length)

	spacing := int32(100) // Odstęp między węzłami

	// Wyśrodkowanie mapy w `GameBoard`
	centerX := boardWidth / 2
	startY := boardHeight - 50 // Pierwszy węzeł blisko dolnej krawędzi

	for i := int32(0); i < length; i++ {
		Nodes[i] = Node{
			X:         centerX,
			Y:         startY - i*spacing, // Pozycjonowanie w pionie
			Active:    i == 0,             // Tylko pierwszy węzeł aktywny
			Completed: false,              // Na start żaden węzeł nie jest ukończony
		}
		if i > 0 {
			Nodes[i-1].Next = &Nodes[i] // Łączenie węzłów w ścieżkę
		}
	}

	// Pierwszy węzeł jako startowy
	CurrentNode = &Nodes[0]
}

// **Rysowanie mapy**
func DrawMap(boardX, boardY int32, currentView *GameView) {
	rl.BeginScissorMode(boardX, boardY, 800, 600)

	// **Rysowanie węzłów i ścieżki**
	for i := range Nodes {
		node := &Nodes[i]
		nodeX := boardX + node.X + scrollOffsetX
		nodeY := boardY + node.Y + scrollOffsetY

		if node.Next != nil {
			nextX := boardX + node.Next.X + scrollOffsetX
			nextY := boardY + node.Next.Y + scrollOffsetY
			rl.DrawLine(nodeX, nodeY, nextX, nextY, rl.DarkGray)
		}

		color := rl.DarkGray
		if node.Completed {
			color = rl.Green
		} else if node.Active {
			color = rl.Blue
		}

		rl.DrawCircle(nodeX, nodeY, 10, color)
	}

	// **Rysowanie przycisków dla aktywnych węzłów**
	DrawNodeButtons(boardX, boardY, currentView)

	rl.EndScissorMode()
}

// **Rysowanie przycisków „Walka” i „Shop” na aktywnych węzłach**
func DrawNodeButtons(boardX, boardY int32, currentView *GameView) {
	for i := range Nodes {
		node := &Nodes[i]

		if !node.Active {
			continue
		}

		nodeX := boardX + node.X + scrollOffsetX
		nodeY := boardY + node.Y + scrollOffsetY

		buttonWidth := int32(100)
		buttonHeight := int32(40)
		buttonX := nodeX - buttonWidth/2
		buttonY := nodeY - buttonHeight/2 // Przycisk dokładnie na węźle

		button := rl.Rectangle{
			X:      float32(buttonX),
			Y:      float32(buttonY),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		}

		if i != 0 && i%3 == 0 {
			// Co trzeci węzeł (oprócz pierwszego) dostaje przycisk "Shop"
			rl.DrawRectangleRec(button, rl.DarkBlue)
			rl.DrawText("Shop", buttonX+30, buttonY+10, 20, rl.White)

			if rl.CheckCollisionPointRec(rl.GetMousePosition(), button) {
				rl.DrawRectangleRec(button, rl.SkyBlue) // Podświetlenie
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					*currentView = ViewShopBoard
				}
			}
		} else {
			// Na pozostałych węzłach przycisk "Walka"
			rl.DrawRectangleRec(button, rl.Red)
			rl.DrawText("Fight", buttonX+20, buttonY+10, 20, rl.White)

			if rl.CheckCollisionPointRec(rl.GetMousePosition(), button) {
				rl.DrawRectangleRec(button, rl.Green) // Podświetlenie
				if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
					*currentView = ViewFightBoard
					CurrentNode = node
				}
			}
		}
	}
}

// **Obsługa przewijania mapy**
func HandleScrolling(boardWidth, boardHeight int32) {
	scrollSpeed := int32(3)

	// **Zmiana przesunięcia na podstawie klawiszy strzałek**
	if rl.IsKeyDown(rl.KeyUp) {
		scrollOffsetY += scrollSpeed
	}
	if rl.IsKeyDown(rl.KeyDown) {
		scrollOffsetY -= scrollSpeed
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		scrollOffsetX += scrollSpeed
	}
	if rl.IsKeyDown(rl.KeyRight) {
		scrollOffsetX -= scrollSpeed
	}

	// **Ograniczenia przesuwania mapy**
	maxScrollX := boardWidth / 2
	minScrollX := -maxScrollX
	maxScrollY := boardHeight / 2
	minScrollY := -maxScrollY

	if scrollOffsetX > maxScrollX {
		scrollOffsetX = maxScrollX
	}
	if scrollOffsetX < minScrollX {
		scrollOffsetX = minScrollX
	}
	if scrollOffsetY > maxScrollY {
		scrollOffsetY = maxScrollY
	}
	if scrollOffsetY < minScrollY {
		scrollOffsetY = minScrollY
	}
}

// **Rysowanie przycisku menu oraz przycisku help**
func drawMenuAndHelpButtons(bottomPanel rl.Rectangle) bool {
	buttonWidth := int32(120)
	buttonHeight := int32(30)
	buttonSpacing := int32(10) // Odstęp między przyciskami

	// 📌 Pozycja przycisku "Menu"
	menuButtonX := int32(bottomPanel.X) + 10
	menuButtonY := int32(bottomPanel.Y) + int32(bottomPanel.Height)/2 - buttonHeight/2
	menuButton := rl.Rectangle{
		X:      float32(menuButtonX),
		Y:      float32(menuButtonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}

	// 📌 Pozycja przycisku "Help" (obok "Menu")
	helpButtonX := menuButtonX + buttonWidth + buttonSpacing
	helpButtonY := menuButtonY
	helpButton := rl.Rectangle{
		X:      float32(helpButtonX),
		Y:      float32(helpButtonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}

	// **Obsługa kliknięcia w "Menu"**
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), menuButton) {
		rl.DrawRectangleRec(menuButton, rl.Gray)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return true // Powrót do menu
		}
	} else {
		rl.DrawRectangleRec(menuButton, rl.DarkGray)
	}
	rl.DrawText("Menu", menuButtonX+30, menuButtonY+10, 20, rl.White)

	// **Obsługa kliknięcia w "Help"**
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), helpButton) {
		rl.DrawRectangleRec(helpButton, rl.Gray)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			HelpWindowActive = !HelpWindowActive // Włącz/wyłącz okno pomocy
		}
	} else {
		rl.DrawRectangleRec(helpButton, rl.DarkGray)
	}
	rl.DrawText("Help", helpButtonX+30, helpButtonY+10, 20, rl.White)

	// **Wyświetlanie okna pomocy, jeśli aktywne**
	if HelpWindowActive {
		drawHelpWindow()
	}

	return false
}

// **Rysowanie okna pomocy**
func drawHelpWindow() {
	windowWidth := int32(600)  // Zwiększenie szerokości okna
	windowHeight := int32(600) // Zwiększenie wysokości okna
	windowX := int32(rl.GetScreenWidth())/2 - windowWidth/2
	windowY := int32(rl.GetScreenHeight())/2 - windowHeight/2

	rl.DrawRectangle(windowX, windowY, windowWidth, windowHeight, rl.DarkGray)
	rl.DrawRectangleLines(windowX, windowY, windowWidth, windowHeight, rl.White)

	textMargin := int32(20)
	iconSize := int32(40)
	textOffsetX := iconSize + 15
	startY := windowY + textMargin + 10
	lineSpacing := iconSize + 15
	textWidth := windowWidth - textMargin*2 - textOffsetX

	// Lista jednostek i ich ruchów
	unitMoves := []struct {
		UnitType string
		MoveDesc string
		Icon     rl.Texture2D
	}{
		{"King", "Moves 1 space in any direction", Pawns.PawnIcons["King"]},
		{"Warrior", "Moves 1 space forward or diagonally", Pawns.PawnIcons["Warrior"]},
		{"Knight", "Moves any number of spaces vertically or horizontally", Pawns.PawnIcons["Knight"]},
		{"Monk", "Moves like a knight in chess (L-shape)", Pawns.PawnIcons["Monk"]},
		{"Racoon", "Moves 1, 2, or 3 spaces horizontally or vertically", Pawns.PawnIcons["Racoon"]},
		{"Master", "Moves diagonally any number of spaces", Pawns.PawnIcons["Master"]},
		{"Boss", "Moves 1 space in any direction", Pawns.PawnIcons["Boss"]},
		{"Lizard", "Moves forward, backward, or diagonally", Pawns.PawnIcons["Lizard"]},
		{"Reptile", "Moves in an L-shape, similar to a knight", Pawns.PawnIcons["Reptile"]},
		{"LionWarrior", "Moves 3 spaces forward or diagonally and 1 space backward", Pawns.PawnIcons["LionWarrior"]},
	}

	for i, unit := range unitMoves {
		iconX := windowX + textMargin
		iconY := startY + int32(i)*lineSpacing

		// Rysowanie ikony pionka
		destRect := rl.Rectangle{
			X:      float32(iconX),
			Y:      float32(iconY),
			Width:  float32(iconSize),
			Height: float32(iconSize),
		}
		sourceRect := rl.Rectangle{
			X: 0, Y: 0, Width: float32(unit.Icon.Width), Height: float32(unit.Icon.Height),
		}
		rl.DrawTexturePro(unit.Icon, sourceRect, destRect, rl.Vector2{}, 0, rl.White)

		// Zawijanie opisu ruchów
		wrappedText := WrapHelpText(unit.MoveDesc, textWidth)
		lineOffsetY := int32(0)

		for _, line := range wrappedText {
			textX := iconX + textOffsetX
			rl.DrawText(line, textX, iconY+10+lineOffsetY, 20, rl.White)
			lineOffsetY += lineSpacing / 2
		}
	}

	// Przycisk zamykania okna pomocy (X w rogu)
	closeButtonSize := int32(30)
	closeButtonX := windowX + windowWidth - closeButtonSize - 10
	closeButtonY := windowY + 10
	closeButton := rl.Rectangle{
		X:      float32(closeButtonX),
		Y:      float32(closeButtonY),
		Width:  float32(closeButtonSize),
		Height: float32(closeButtonSize),
	}

	rl.DrawRectangleRec(closeButton, rl.Red)
	rl.DrawText("X", closeButtonX+10, closeButtonY+5, 20, rl.White)

	if rl.CheckCollisionPointRec(rl.GetMousePosition(), closeButton) && rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		HelpWindowActive = false
	}
}

// **Funkcja do zawijania tekstu**
func WrapHelpText(text string, maxWidth int32) []string {
	lines := []string{}
	words := strings.Fields(text) // Dzielimy tekst na słowa
	line := ""

	for _, word := range words {
		testLine := line
		if line != "" {
			testLine += " " + word
		} else {
			testLine = word
		}

		if rl.MeasureText(testLine, 20) > int32(maxWidth) {
			if line != "" {
				lines = append(lines, line)
			}
			line = word
		} else {
			line = testLine
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
}

// CalculateGameBoardSize oblicza rozmiar i pozycję GameBoard na ekranie.
func CalculateGameBoardSize(screenWidth, screenHeight int32) (int32, int32, int32, int32) {
	topHeight := int32(float32(screenHeight) * 0.05)    // 5% wysokości ekranu dla górnego panelu
	bottomHeight := int32(float32(screenHeight) * 0.05) // 5% wysokości dla dolnego panelu
	verticalSpace := screenHeight - topHeight - bottomHeight

	leftWidth := int32(float32(screenWidth) * 0.18) // Zwiększamy lewy panel (22.5% szerokości ekranu)
	rightWidth := int32(float32(leftWidth) / 2)     // Zmniejszamy prawy panel (połowa lewego)

	boardWidth := screenWidth - leftWidth - rightWidth // Szerokość GameBoard
	boardHeight := verticalSpace                       // Wysokość GameBoard - powinno dokładnie wypełnić przestrzeń między panelami

	// 📌 Spójny margines dla GameBoard i FightingBoard
	marginX := int32(50) // Można dostosować wartość w zależności od potrzeb
	boardX := leftWidth + marginX
	boardY := topHeight

	return boardWidth, boardHeight, boardX, boardY
}

// **Obsługa zakończenia etapu (bitwa / sklep)**
func CompleteNode(currentView *GameView) {
	if CurrentNode != nil {
		CurrentNode.Completed = true // Oznaczamy obecny węzeł jako ukończony
		CurrentNode.Active = false   // Dezaktywujemy obecny węzeł

		// **Jeśli to ostatni węzeł, przechodzimy do ekranu wygranej**
		if CurrentNode.Next == nil {
			*currentView = ViewWinScreen
			return
		}

		// **Aktywujemy następny węzeł**
		if CurrentNode.Next != nil {
			CurrentNode.Next.Active = true
			CurrentNode = CurrentNode.Next // Przechodzimy do następnego węzła
		}
	}

	// **Powrót do GameBoard**
	*currentView = ViewGameBoard
}
