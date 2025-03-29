package menu

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ShowOptions wyświetla ekran opcji (wybór rozdzielczości)
func ShowOptions(screenWidth, screenHeight int32) ([]int32, bool) {
	// Przyciski rozdzielczości
	resolutions := [][]int32{
		{1280, 720},
		{1920, 1080},
		{2560, 1440},
	}

	// Wymiary przycisków
	buttonWidth := int32(200)
	buttonHeight := int32(50)
	margin := int32(20)

	rl.ClearBackground(rl.RayWhite)

	rl.DrawText("OPCJE", screenWidth/2-rl.MeasureText("OPCJE", 20)/2, 50, 20, rl.Black)

	// Rysowanie przycisków rozdzielczości
	newResolution := []int32(nil)
	backToMenu := false

	for i, res := range resolutions {
		resButton := rl.Rectangle{
			X:      float32(screenWidth/2 - buttonWidth/2),
			Y:      float32(200 + int32(i)*(buttonHeight+margin)),
			Width:  float32(buttonWidth),
			Height: float32(buttonHeight),
		}

		if rl.CheckCollisionPointRec(rl.GetMousePosition(), resButton) {
			rl.DrawRectangleRec(resButton, rl.Gray)
			if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
				newResolution = res
			}
		} else {
			rl.DrawRectangleRec(resButton, rl.DarkGray)
		}

		resText := fmt.Sprintf("%dx%d", res[0], res[1])
		rl.DrawText(resText, int32(resButton.X)+30, int32(resButton.Y)+15, 20, rl.White)
	}

	// Przycisk "Back"
	backButton := rl.Rectangle{
		X:      float32(screenWidth/2 - buttonWidth/2),
		Y:      float32(screenHeight - buttonHeight*2),
		Width:  float32(buttonWidth),
		Height: float32(buttonHeight),
	}

	if rl.CheckCollisionPointRec(rl.GetMousePosition(), backButton) {
		rl.DrawRectangleRec(backButton, rl.Gray)
		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			backToMenu = true
		}
	} else {
		rl.DrawRectangleRec(backButton, rl.DarkGray)
	}
	rl.DrawText("Back", int32(backButton.X)+60, int32(backButton.Y)+15, 20, rl.White)

	return newResolution, backToMenu
}
