package Boards

import (
	"Protect_The_King/Pawns"
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Jednostki dostępne w sklepie
var AllShopUnits = []Pawns.ShopPawn{
	{Type: "Warrior", Cost: 20, Owner: "Player 1"},
	{Type: "Knight", Cost: 90, Owner: "Player 1"},
	{Type: "Monk", Cost: 60, Owner: "Player 1"},
	{Type: "Master", Cost: 80, Owner: "Player 1"},
}

// Lista jednostek dostępnych w sklepie (losowe)
var RandomShopUnits []Pawns.ShopPawn
var ShopGenerated bool = false // Flaga, by nie losować za każdym razem

var PlayerGold int32 = 100 // 💰 Złoto gracza
var RollTickets int = 1    // Ilość dostępnych przerzutów sklepu

// Tekstury
var loadedTextures map[string]rl.Texture2D
var texturesLoaded bool = false

// Widok sklepu
func DrawShopBoard(screenWidth int32, screenHeight int32, currentView *GameView) bool {
	// Generujemy jednostki tylko raz na wejście do sklepu
	if !ShopGenerated {
		GenerateRandomShop()
	}

	// Rysowanie interfejsu sklepu
	rl.DrawRectangle(0, 0, screenWidth, screenHeight, rl.Gray)
	rl.DrawText("Sklep", screenWidth/2-50, 20, 30, rl.White)

	// Rysowanie złota gracza
	DrawGold(screenWidth)
	// Rysowanie ticketowResetu
	DrawRollTickets(screenWidth)

	// **Rysujemy trzy okna z losowymi jednostkami**
	windowWidth := screenWidth / 3
	windowHeight := int32(float32(screenHeight) * 0.7)
	windowY := int32(80)

	for i, unit := range RandomShopUnits {
		DrawShopWindow(int32(i)*windowWidth, windowY, windowWidth, windowHeight, unit)
	}

	// Przycisk "Gotowe"
	if DrawReadyButton(screenWidth, screenHeight) {
		ShopGenerated = false // Resetujemy sklep przy zamknięciu
		return true           // Sklep został zamknięty
	}

	// **Przycisk "Roll Shop"**
	if DrawRollButton(screenWidth, screenHeight) {
		GenerateRandomShop() // Nowe losowanie
		RollTickets--        // Zużycie jednego przerzutu
	}

	return false // Sklep nadal otwarty
}

// Rysowanie pojedynczego okna sklepu
func DrawShopWindow(x, y, width, height int32, unit Pawns.ShopPawn) {
	rl.DrawRectangle(x, y, width, height, rl.LightGray)
	rl.DrawRectangleLines(x, y, width, height, rl.DarkGray)

	textWidth := rl.MeasureText(unit.Type, 20)
	rl.DrawText(unit.Type, x+(width-textWidth)/2, y+10, 20, rl.Black)

	// Rysowanie ceny
	priceText := fmt.Sprintf("%d G", unit.Cost)
	rl.DrawText(priceText, x+20, y+40, 20, rl.DarkGray)

	// Pobranie tekstury z mapy zamiast ponownego ładowania
	texture, exists := loadedTextures[unit.Type]
	if exists && texture.ID > 0 {
		rl.DrawTexture(texture, x+20, y+70, rl.White)
	}

	// Przycisk zakupu
	button := rl.Rectangle{
		X:      float32(x + width/2 - 60),
		Y:      float32(y + height - 60),
		Width:  120,
		Height: 40,
	}

	rl.DrawRectangleRec(button, rl.DarkGreen)
	rl.DrawText("Kup", int32(button.X)+40, int32(button.Y)+10, 20, rl.White)

	if rl.CheckCollisionPointRec(rl.GetMousePosition(), button) {
		rl.DrawRectangleRec(button, rl.Green)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			PurchaseUnit(unit.Type)
		}
	}
}

func LoadShopTextures() {
	if texturesLoaded {
		return
	}

	loadedTextures = make(map[string]rl.Texture2D)

	for _, unit := range AllShopUnits {
		texture := rl.LoadTexture(fmt.Sprintf("Assets/Player/%s/Faceset.png", unit.Type))
		loadedTextures[unit.Type] = texture
	}

	texturesLoaded = true
}

func UnloadShopTextures() {
	if !texturesLoaded {
		return
	}

	for _, texture := range loadedTextures {
		rl.UnloadTexture(texture)
	}

	loadedTextures = nil
	texturesLoaded = false
}

// Kupowanie pionka
func PurchaseUnit(unitType string) bool {
	for _, unit := range AllShopUnits {
		if unit.Type == unitType {
			if PlayerGold >= unit.Cost {
				PlayerGold -= unit.Cost

				// Pobranie konfiguracji wizualnej
				cfg, exists := Pawns.PawnVisualConfigs[unitType]
				if !exists {
					fmt.Printf("Nie znaleziono konfiguracji dla %s\n", unitType)
					return false
				}

				// Wczytanie obu tekstur (ikonki + animacja)
				staticTex := rl.LoadTexture(cfg.StaticTexturePath)
				animTex := rl.LoadTexture(cfg.AnimTexturePath)

				newPawn := Pawns.BasePawn{
					ID:          Pawns.NextPawnID,
					Type:        unit.Type,
					X:           -1,
					Y:           -1,
					IsAlive:     true,
					Owner:       "Player 1",
					Cost:        unit.Cost,
					Texture:     staticTex,
					AnimTexture: animTex,
					Animation:   cfg.Anim,
				}

				Pawns.Player1Pawns = append(Pawns.Player1Pawns, newPawn)
				Pawns.NextPawnID++

				fmt.Printf("Kupiono %s za %d G. Pozostałe złoto: %d G\n", unit.Type, unit.Cost, PlayerGold)
				return true
			} else {
				fmt.Println("Za mało złota!")
				return false
			}
		}
	}
	return false
}

func DrawGold(screenWidth int32) {
	// Ustawienie pozycji tekstu na ekranie
	text := fmt.Sprintf("Gold: %d G", PlayerGold)
	textX := screenWidth - 200 // Odstęp od prawej krawędzi
	textY := 20                // Odstęp od góry

	// Rysowanie tła dla lepszej czytelności
	background := rl.Rectangle{
		X:      float32(textX - 10),
		Y:      float32(textY - 5),
		Width:  140,
		Height: 30,
	}
	rl.DrawRectangleRec(background, rl.DarkGray)

	// Rysowanie tekstu informującego o ilości złota
	rl.DrawText(text, textX, int32(textY), 25, rl.Yellow)
}

func DrawRollTickets(screenWidth int32) {
	// Wymiary i pozycjonowanie względem prawej krawędzi ekranu
	paddingLeft := int32(20) // Odstęp od prawej krawędzi
	textSize := 25           // Rozmiar czcionki
	text := fmt.Sprintf("Rolls: %d", RollTickets)

	textWidth := rl.MeasureText(text, int32(textSize))
	textX := paddingLeft // Pozycja X (na lewo)
	textY := 20          // Pozycja Y na ekranie

	// Rysowanie tła dla lepszej widoczności
	background := rl.Rectangle{
		X:      float32(textX - 10),
		Y:      float32(textY - 5),
		Width:  float32(textWidth + 20),
		Height: 30,
	}
	rl.DrawRectangleRec(background, rl.DarkGray)

	// Rysowanie tekstu
	rl.DrawText(text, textX, int32(textY), int32(textSize), rl.Yellow)
}

// Funkcja losująca 3 jednostki do sklepu
func GenerateRandomShop() {
	if !texturesLoaded {
		LoadShopTextures() // Załaduj tekstury tylko raz
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	RandomShopUnits = nil // Resetujemy poprzednie jednostki

	// Wybieramy 3 losowe jednostki z listy dostępnych
	for i := 0; i < 3; i++ {
		randomIndex := rng.Intn(len(AllShopUnits))
		RandomShopUnits = append(RandomShopUnits, AllShopUnits[randomIndex])
	}

	ShopGenerated = true
}

func DrawRollButton(screenWidth, screenHeight int32) bool {
	buttonWidth := int32(200)
	buttonHeight := int32(50)
	buttonSpacing := int32(20) // Odstęp między przyciskami

	// Pozycja przycisku "Ready"
	readyButtonX := (screenWidth - buttonWidth) / 2
	readyButtonY := screenHeight - buttonHeight - 20

	// Pozycja przycisku "Roll Shop" (na lewo od "Ready")
	rollButtonX := readyButtonX - buttonWidth - buttonSpacing
	rollButtonY := readyButtonY // Wyrównanie wysokości z "Ready"

	// Rysowanie przycisku "Roll Shop"
	rl.DrawRectangle(rollButtonX, rollButtonY, buttonWidth, buttonHeight, rl.DarkBlue)
	rl.DrawText("Roll Shop", rollButtonX+40, rollButtonY+15, 20, rl.White)

	// Podświetlenie przycisku, jeśli myszka nad nim
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), rl.Rectangle{
		X: float32(rollButtonX), Y: float32(rollButtonY),
		Width: float32(buttonWidth), Height: float32(buttonHeight),
	}) {
		rl.DrawRectangleRec(rl.Rectangle{
			X: float32(rollButtonX), Y: float32(rollButtonY),
			Width: float32(buttonWidth), Height: float32(buttonHeight),
		}, rl.Blue)
		if rl.IsMouseButtonReleased(rl.MouseLeftButton) && RollTickets > 0 {
			return true
		}
	}

	return false
}

// Przycisk "Ready"
func DrawReadyButton(screenWidth, screenHeight int32) bool {
	buttonWidth := int32(200)
	buttonHeight := int32(50)
	buttonX := (screenWidth - buttonWidth) / 2
	buttonY := screenHeight - buttonHeight - 20

	rl.DrawRectangle(buttonX, buttonY, buttonWidth, buttonHeight, rl.DarkGreen)
	rl.DrawText("Ready", buttonX+50, buttonY+15, 20, rl.White)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		mouseX := rl.GetMouseX()
		mouseY := rl.GetMouseY()
		if mouseX > buttonX && mouseX < buttonX+buttonWidth && mouseY > buttonY && mouseY < buttonY+buttonHeight {
			return true
		}
	}

	return false
}
