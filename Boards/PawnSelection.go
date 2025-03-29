package Boards

import (
	"Protect_The_King/Pawns"
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var PawnSelectionDone bool = false
var selectedPawns []Pawns.BasePawn  // Lista wybranych pionków
var availablePawns []Pawns.BasePawn // Lista 3 losowych pionków
var maxSelections int = 5
var maxRemovalsPerPhase int = 1
var removalCount int = 0 // Licznik usuniętych pionków

var commonPawns = []string{"Warrior", "Monk"} // 75% szansy
var rarePawns = []string{"Master", "Knight"}  // 25% szansy

func ShowInitialPawnSelectionMenu(screenWidth, screenHeight int32) {
	menuWidth := int32(float32(screenWidth) * 0.5)
	menuHeight := int32(float32(screenHeight) * 0.6)
	menuX := (screenWidth - menuWidth) / 2
	menuY := (screenHeight - menuHeight) / 2

	source := rl.Rectangle{X: 0, Y: 0, Width: float32(PawnSelectionBackground.Width), Height: float32(PawnSelectionBackground.Height)}
	dest := rl.Rectangle{X: 0, Y: 0, Width: float32(screenWidth), Height: float32(screenHeight)}
	rl.DrawTexturePro(PawnSelectionBackground, source, dest, rl.Vector2{X: 0, Y: 0}, 0, rl.White)

	rl.DrawRectangle(menuX, menuY, menuWidth, menuHeight, rl.DarkGray)
	rl.DrawRectangleLines(menuX, menuY, menuWidth, menuHeight, rl.White)
	rl.DrawText("Choose your starting units:", menuX+20, menuY+20, 20, rl.White)

	// Dodajemy Kinga automatycznie
	if len(selectedPawns) == 0 {
		king := Pawns.CreatePawn("King", "Player 1")
		//test := Pawns.CreatePawn("LionWarrior", "Player 1")
		selectedPawns = append(selectedPawns, king)
		//selectedPawns = append(selectedPawns, test)
		GenerateRandomPawnSelection()
	}

	// Rysujemy dostępne pionki do wyboru
	DrawPawnOptions(menuX+20, menuY+60, availablePawns)

	buttonWidth := int32(150)
	buttonHeight := int32(40)
	buttonX := menuX + (menuWidth-buttonWidth)/2
	buttonY := menuY + menuHeight - 60

	rl.DrawRectangle(buttonX, buttonY, buttonWidth, buttonHeight, rl.Green)
	rl.DrawText("Confirm", buttonX+40, buttonY+10, 20, rl.White)

	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		mouseX := rl.GetMouseX()
		mouseY := rl.GetMouseY()

		if mouseX > buttonX && mouseX < buttonX+buttonWidth &&
			mouseY > buttonY && mouseY < buttonY+buttonHeight {

			if len(selectedPawns) > 0 {

				Pawns.Player1Pawns = make([]Pawns.BasePawn, len(selectedPawns))
				copy(Pawns.Player1Pawns, selectedPawns)

				PawnSelectionDone = true
				ResetRemovalCountForPhase()
				selectedPawns = []Pawns.BasePawn{}
			}
		}
	}

	DrawSelectedPawnsPanel(screenWidth, screenHeight)

}

// Losowanie 3 pionków do wyboru
func GenerateRandomPawnSelection() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	availablePawns = []Pawns.BasePawn{}

	for i := 0; i < 3; i++ {
		roll := rng.Intn(100) + 1

		var chosenPawnType string
		if roll <= 75 {
			// 75% szansy na podstawowy pionek
			chosenPawnType = commonPawns[rng.Intn(len(commonPawns))]
		} else {
			// 25% szansy na rzadki pionek
			chosenPawnType = rarePawns[rng.Intn(len(rarePawns))]
		}

		newPawn := Pawns.CreatePawn(chosenPawnType, "Player 1")
		availablePawns = append(availablePawns, newPawn)
	}
}

// **🔹 Rysowanie dostępnych pionków do wyboru**
func DrawPawnOptions(x, y int32, options []Pawns.BasePawn) {
	pawnSize := int32(80)
	spacing := int32(20)

	for i, pawn := range options {
		posX := x + int32(i)*(pawnSize+spacing)
		posY := y

		rl.DrawRectangle(posX, posY, pawnSize, pawnSize, rl.LightGray)
		rl.DrawRectangleLines(posX, posY, pawnSize, pawnSize, rl.DarkGray)

		cfg := Pawns.PawnVisualConfigs[pawn.Type]
		staticTex := rl.LoadTexture(cfg.StaticTexturePath)
		texture := staticTex
		rl.DrawTextureEx(texture, rl.Vector2{X: float32(posX + 5), Y: float32(posY + 5)}, 0, 0.8, rl.White)

		rl.DrawText(pawn.Type, posX+10, posY+60, 18, rl.Black)

		if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
			mouseX := rl.GetMouseX()
			mouseY := rl.GetMouseY()

			if mouseX > posX && mouseX < posX+pawnSize &&
				mouseY > posY && mouseY < posY+pawnSize {

				if len(selectedPawns) < maxSelections {
					selectedPawns = append(selectedPawns, pawn)
					GenerateRandomPawnSelection()
					break
				} else {
				}
			}
		}
	}
}

func DrawSelectedPawnsPanel(screenWidth, screenHeight int32) {
	width := int32(220)
	height := int32(300)
	x := screenWidth - width - 20
	y := int32(50)

	rl.DrawRectangle(x, y, width, height, rl.LightGray)
	rl.DrawRectangleLines(x, y, width, height, rl.Black)
	rl.DrawText("Selected Pawns:", x+10, y+10, 22, rl.Black)
	pawnY := y + 40
	lineSpacing := 35

	for i, pawn := range selectedPawns {
		text := fmt.Sprintf("%d. %s", i+1, pawn.Type)
		rl.DrawText(text, x+10, pawnY, 20, rl.Black)

		if pawn.Type != "King" && removalCount < maxRemovalsPerPhase {
			removeButtonX := x + width - 30
			removeButtonY := pawnY - 5
			rl.DrawRectangle(removeButtonX, removeButtonY, 25, 25, rl.Red)
			rl.DrawText("X", removeButtonX+7, removeButtonY+4, 20, rl.White)

			// Obsług  usunięcia pionka po kliknięciu X
			if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
				mouseX := rl.GetMouseX()
				mouseY := rl.GetMouseY()

				if mouseX > removeButtonX && mouseX < removeButtonX+25 &&
					mouseY > removeButtonY && mouseY < removeButtonY+25 {

					selectedPawns = append(selectedPawns[:i], selectedPawns[i+1:]...)

					removalCount++

					// Po usunięciu losujemy nowe pionki
					availablePawns = []Pawns.BasePawn{}
					GenerateRandomPawnSelection()
					break
				}
			}
		}
		pawnY += int32(lineSpacing)
	}
}

func ResetRemovalCountForPhase() {
	removalCount = 0
}
