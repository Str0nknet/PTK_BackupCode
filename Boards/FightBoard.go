package Boards

import (
	"fmt"
	"math/rand"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Przechowuje logi z przebiegu walki
var BattleLogs []string

const maxLogs = 5 // Maksymalna liczba wyświetlanych logów

// Tile reprezentuje jedno pole na planszy
type Tile struct {
	CoordX     int          // Współrzędna X (liczba)
	CoordY     string       // Współrzędna Y (litera)
	PosX       int32        // Pozycja X w pikselach na ekranie
	PosY       int32        // Pozycja Y w pikselach na ekranie
	Walkable   bool         // Czy pole jest dostępne do ruchu
	SourceRect rl.Rectangle // Fragment tilesetu przypisany do pola
}

// GetX zwraca współrzędną X pola
func (t Tile) GetX() int32 {
	return t.PosX
}

// GetY zwraca współrzędną Y pola
func (t Tile) GetY() int32 {
	return t.PosY
}

// IsWalkable sprawdza, czy pole jest dostępne do ruchu
func (t Tile) IsWalkable() bool {
	return t.Walkable
}

// LoadTileset ładuje tileset i zwraca teksturę oraz mapę sekcji
func LoadTileset(filePath string) (rl.Texture2D, map[string]rl.Rectangle) {
	tileset := rl.LoadTexture(filePath)

	// Każdy kafelek ma wymiary 32x32 piksele
	tileWidth := float32(28)
	tileHeight := float32(28)

	// Mapowanie tekstur na podstawie współrzędnych
	sections := map[string]rl.Rectangle{
		// Pierwszy rząd (pomarańczowe kafelki)
		"orange_grass": {X: 50, Y: 0, Width: tileWidth, Height: tileHeight},
		// Drugi rząd (zielone kafelki)
		"green_grass": {X: 50, Y: 48, Width: tileWidth, Height: tileHeight},
		// Trzeci rząd (ciemnozielone kafelki)
		"dark_green_grass": {X: 50, Y: 96, Width: tileWidth, Height: tileHeight},
		// Czwarty rząd (różowe kafelki)
		"pink_grass": {X: 50, Y: 144, Width: tileWidth, Height: tileHeight},
		// Piąty rząd (białe kafelki)
		"white_grass": {X: 50, Y: 192, Width: tileWidth, Height: tileHeight},
	}

	return tileset, sections
}

// generateRandomPositions generuje losowe pozycje white_grass
func generateRandomPositions(boardSize int32, count int) [][2]int32 {
	positions := make([][2]int32, 0, count)
	usedPositions := make(map[[2]int32]bool)

	for len(positions) < count {
		row := rand.Int31n(boardSize)
		col := rand.Int31n(boardSize)
		pos := [2]int32{row, col}

		// **Pomijamy 2 górne i 2 dolne rzędy**
		if row < 2 || row >= boardSize-2 {
			continue
		}

		// **Sprawdzamy, czy pozycja nie była wcześniej użyta**
		if !usedPositions[pos] {
			positions = append(positions, pos)
			usedPositions[pos] = true
		}
	}

	return positions
}

// drawTileTexture rysuje fragment tilesetu na polu
func DrawTileTexture(texture rl.Texture2D, sourceRect rl.Rectangle, posX, posY, cellSize int32) {
	// Docelowy prostokąt na ekranie (dopasowany do pola)
	destRect := rl.Rectangle{
		X:      float32(posX),
		Y:      float32(posY),
		Width:  float32(cellSize),
		Height: float32(cellSize),
	}

	// Punkt obrotu (środek tekstury, tutaj brak obrotu)
	origin := rl.Vector2{X: 0, Y: 0}

	// Rysowanie fragmentu tekstury
	rl.DrawTexturePro(texture, sourceRect, destRect, origin, 0, rl.White)
}

// displayCoordinates wyświetla koordynaty najechanego pola w lewym dolnym rogu
func DisplayCoordinates(board [][]Tile, screenHeight, cellSize int32) {
	mousePos := rl.GetMousePosition()
	mouseX := int32(mousePos.X)
	mouseY := int32(mousePos.Y)

	for _, row := range board {
		for _, tile := range row {
			// Sprawdzamy, czy mysz znajduje się nad tym polem
			if mouseX >= tile.PosX && mouseX < tile.PosX+cellSize &&
				mouseY >= tile.PosY && mouseY < tile.PosY+cellSize {
				// Koordynaty najechanego pola
				coordText := fmt.Sprintf("Koordynaty: %d%s", tile.CoordX, tile.CoordY)

				// Wyświetlanie koordynatów w lewym dolnym rogu
				rl.DrawText(coordText, 10, screenHeight-30, 20, rl.Black)
				return
			}
		}
	}
}

// drawBackButton rysuje przycisk "Powrót" w dolnej części ekranu
func DrawBackButton(screenWidth, screenHeight int32) bool {
	// Wymiary przycisku
	buttonWidth := int32(150)
	buttonHeight := int32(40)
	buttonX := screenWidth/2 - buttonWidth/2
	buttonY := screenHeight - buttonHeight - 20 // 20px margines od dolnej krawędzi

	// Prostokąt przycisku
	button := rl.Rectangle{
		X:      float32(buttonX),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}

	// Rysowanie przycisku
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), button) {
		rl.DrawRectangleRec(button, rl.Gray)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			return true // Kliknięcie przycisku
		}
	} else {
		rl.DrawRectangleRec(button, rl.DarkGray)
	}
	rl.DrawText("Powrót", buttonX+30, buttonY+10, 20, rl.White)

	return false
}

func DrawHelpButton(screenWidth, screenHeight int32) bool {
	// Wymiary i pozycja przycisku w prawym dolnym rogu
	buttonWidth := int32(150)
	buttonHeight := int32(40)
	buttonX := screenWidth - buttonWidth - 20
	buttonY := screenHeight - buttonHeight - 20

	// Prostokąt przycisku
	button := rl.Rectangle{
		X:      float32(buttonX),
		Y:      float32(buttonY),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}

	// Rysowanie przycisku
	// **Obsługa kliknięcia w "Help"**
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), button) {
		rl.DrawRectangleRec(button, rl.Gray)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			HelpWindowActive = !HelpWindowActive // Włącz/wyłącz okno pomocy
		}
	} else {
		rl.DrawRectangleRec(button, rl.DarkGray)
	}
	rl.DrawText("Help", buttonX+30, buttonY+10, 20, rl.White)

	// **Wyświetlanie okna pomocy, jeśli aktywne**
	if HelpWindowActive {
		drawHelpWindow()
	}
	return false
}

func GenerateBoard(boardSize, cellSize, boardX, boardY int32, sections map[string]rl.Rectangle, seed int64) [][]Tile {

	board := make([][]Tile, boardSize)
	whiteGrassPositions := generateRandomPositions(boardSize, 6)

	for row := int32(0); row < boardSize; row++ {
		board[row] = make([]Tile, boardSize)

		for col := int32(0); col < boardSize; col++ {
			terrainType := "green_grass"
			isWalkable := true

			for _, pos := range whiteGrassPositions {
				if pos[0] == row && pos[1] == col {
					terrainType = "white_grass"
					isWalkable = false
					break
				}
			}

			posX := boardX + col*cellSize
			posY := boardY + row*cellSize

			board[row][col] = Tile{
				CoordX:     int(col + 1),
				CoordY:     string('A' + rune(row)),
				PosX:       posX,
				PosY:       posY,
				Walkable:   isWalkable,
				SourceRect: sections[terrainType],
			}
		}
	}

	return board
}

// Dodaje nową wiadomość do logów walki
func AddBattleLog(message string) {
	if len(BattleLogs) >= maxLogs {
		BattleLogs = BattleLogs[1:] // Usuwa najstarszy log, jeśli przekroczono limit
	}
	BattleLogs = append(BattleLogs, message)
}

// Rysuje okno z logami walki w dolnym lewym rogu
func DrawBattleLogs(screenWidth, screenHeight int32) {
	logBoxWidth := int32(320)                   // Szerokość okna logów
	logBoxHeight := int32(150)                  // Wysokość okna logów
	logBoxX := screenWidth - logBoxWidth - 10   // Pozycja X
	logBoxY := screenHeight - logBoxHeight - 80 // Pozycja Y (px od dołu)

	// Tło okna logów
	rl.DrawRectangle(logBoxX, logBoxY, logBoxWidth, logBoxHeight, rl.Fade(rl.Black, 0.7))
	rl.DrawRectangleLines(logBoxX, logBoxY, logBoxWidth, logBoxHeight, rl.White)

	// Maksymalna szerokość tekstu (aby nie wychodził poza ramkę)
	textMargin := int32(10)
	maxTextWidth := logBoxWidth - 2*textMargin
	maxLines := logBoxHeight / 20 // Liczba linii, które mieszczą się w oknie

	// Przetwarzanie logów do wyświetlenia
	var wrappedLogs []string
	fontSize := int32(18) // Teraz `int32`, żeby pasowało do `MeasureText`
	for _, log := range BattleLogs {
		wrappedLogs = append(wrappedLogs, WrapText(log, maxTextWidth, fontSize)...) // Przetwarzamy logi
	}

	// Ograniczenie liczby linii do maksymalnej wysokości okna
	if len(wrappedLogs) > int(maxLines) {
		wrappedLogs = wrappedLogs[len(wrappedLogs)-int(maxLines):] // Pokazuje tylko ostatnie linie
	}

	// Wyświetlanie logów
	textY := logBoxY + textMargin
	for _, log := range wrappedLogs {
		rl.DrawText(log, logBoxX+textMargin, textY, fontSize, rl.Red)
		textY += 20
	}
}

// WrapText - Zawija tekst do nowej linii, jeśli przekracza określoną szerokość
func WrapText(text string, maxWidth int32, fontSize int32) []string {
	var wrappedLines []string
	words := strings.Fields(text) // Zamienia tekst na listę słów
	currentLine := ""
	firstLine := true // Flaga dla pierwszej linii

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		textWidth := rl.MeasureText(testLine, fontSize) // Sprawdza szerokość

		if textWidth > maxWidth && currentLine != "" {
			if firstLine {
				wrappedLines = append(wrappedLines, "# "+currentLine) // `#` tylko w pierwszej linii
				firstLine = false
			} else {
				wrappedLines = append(wrappedLines, currentLine)
			}
			currentLine = word // Nowa linia zaczyna się od bieżącego słowa
		} else {
			currentLine = testLine
		}
	}

	if currentLine != "" {
		if firstLine {
			wrappedLines = append(wrappedLines, "# "+currentLine) // `#` w pierwszej linii loga
		} else {
			wrappedLines = append(wrappedLines, currentLine)
		}
	}

	return wrappedLines
}
