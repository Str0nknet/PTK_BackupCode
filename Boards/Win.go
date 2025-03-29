package Boards

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var winTimer float32 = 0 // Timer do automatycznego powrotu do menu

func DrawWinScreen(screenWidth, screenHeight int32) GameView {

	rl.ClearBackground(rl.Black) // Tło ekranu wygranej

	// Wyśrodkowany napis "YOU WIN!"
	text := "YOU WIN!"
	textSize := 50
	textWidth := rl.MeasureText(text, int32(textSize))
	rl.DrawText(text, screenWidth/2-textWidth/2, screenHeight/3, int32(textSize), rl.Gold)

	// Informacja o powrocie do menu
	rl.DrawText("Returning to menu...", screenWidth/2-100, screenHeight/2, 20, rl.White)

	// **Automatyczny powrót do menu po 3 sekundach**
	winTimer += rl.GetFrameTime()
	if winTimer > 8.0 {
		winTimer = 0 // Reset timera
		return ViewMainMenu
	}

	return ViewWinScreen
}
