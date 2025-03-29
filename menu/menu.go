package menu

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// MenuState przechowuje stan menu (czy gra działa, czy wyjść)
type MenuState struct {
	GameRunning bool
	Exit        bool
	InOptions   bool
	LoadGame    bool
}

// Globalna zmienna dla obrazu tła menu
var (
	menuBackground rl.Texture2D
	buttonNormal   rl.Texture2D
	buttonHover    rl.Texture2D
	buttonPressed  rl.Texture2D
)

// LoadMenuAssets wczytuje zasoby menu
func LoadMenuAssets() {

	path := filepath.Join("Assets", "MenuScrean.jpg")

	// Wczytanie tekstury
	menuBackground = rl.LoadTexture(path)

	// Wczytanie tekstur przycisków
	buttonNormal = rl.LoadTexture(filepath.Join("Assets", "Buttons", "ButtonNormal.png"))
	buttonHover = rl.LoadTexture(filepath.Join("Assets", "Buttons", "ButtonHover.png"))
	buttonPressed = rl.LoadTexture(filepath.Join("Assets", "Buttons", "ButtonPressed.png"))
}

// UnloadMenuAssets zwalnia pamięć po zamknięciu gry
func UnloadMenuAssets() {
	rl.UnloadTexture(menuBackground)
	rl.UnloadTexture(buttonNormal)
	rl.UnloadTexture(buttonHover)
	rl.UnloadTexture(buttonPressed)
}

// ShowMenu wyświetla menu główne
func ShowMenu(screenWidth, screenHeight int32) MenuState {
	state := MenuState{}

	// Rysowanie tła menu
	rl.DrawTexture(menuBackground, 0, 0, rl.White)

	// Margines od lewej krawędzi
	marginLeft := float32(100)

	// Rysowanie przycisków
	startButton := rl.Rectangle{X: marginLeft, Y: float32(screenHeight/2 - 150), Width: 200, Height: 50}
	loadGameButton := rl.Rectangle{X: marginLeft, Y: float32(screenHeight/2 - 50), Width: 200, Height: 50} // Nowy przycisk
	optionsButton := rl.Rectangle{X: marginLeft, Y: float32(screenHeight/2 + 50), Width: 200, Height: 50}
	exitButton := rl.Rectangle{X: marginLeft, Y: float32(screenHeight/2 + 150), Width: 200, Height: 50}

	DrawButton(startButton, "Start", &state.GameRunning)
	DrawButton(loadGameButton, "Load Game", &state.LoadGame) // Nowy przycisk
	DrawButton(optionsButton, "Options", &state.InOptions)
	DrawButton(exitButton, "Exit", &state.Exit)

	return state
}

func DrawButton(rect rl.Rectangle, text string, state *bool) {
	var texture rl.Texture2D

	// Sprawdzenie kolizji myszy
	mouseOver := rl.CheckCollisionPointRec(rl.GetMousePosition(), rect)
	if mouseOver && rl.IsMouseButtonDown(rl.MouseLeftButton) {
		texture = buttonPressed
	} else if mouseOver {
		texture = buttonHover
	} else {
		texture = buttonNormal
	}

	// Definiujemy oryginalny rozmiar tekstury (zachowujemy całą grafikę)
	source := rl.Rectangle{X: 0, Y: 0, Width: float32(texture.Width), Height: float32(texture.Height)}

	// Docelowy rozmiar (200x50 px, niezależnie od rzeczywistego rozmiaru obrazka)
	dest := rl.Rectangle{X: rect.X, Y: rect.Y, Width: 200, Height: 50}

	// Rysowanie przeskalowanej tekstury
	rl.DrawTexturePro(texture, source, dest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

	// Automatyczne centrowanie tekstu na przycisku
	textWidth := rl.MeasureText(text, 20)
	textX := int32(rect.X) + (200-textWidth)/2
	textY := int32(rect.Y) + 15

	rl.DrawText(text, textX, textY, 20, rl.Black)

	// Obsługa kliknięcia
	if mouseOver && rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		*state = true
	}
}
