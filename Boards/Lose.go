package Boards

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var loseTimer float32 = 0

func DrawLoseScreen(screenWidth, screenHeight int32) GameView {

	rl.ClearBackground(rl.Black)

	// Wyśrodkowany napis "YOU LOSE!"
	text := "YOU LOSE!"
	textSize := 50
	textWidth := rl.MeasureText(text, int32(textSize))
	rl.DrawText(text, screenWidth/2-textWidth/2, screenHeight/3, int32(textSize), rl.Red)

	// Informacja o powrocie do menu
	rl.DrawText("Returning to menu...", screenWidth/2-100, screenHeight/2, 20, rl.White)

	// **Automatyczny powrót do menu po 3 sekundach**
	loseTimer += rl.GetFrameTime()
	if loseTimer > 3.0 {
		loseTimer = 0
		return ViewMainMenu
	}

	return ViewLoseScreen
}
