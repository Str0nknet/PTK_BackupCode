package Pawns

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var NextPawnID int = 1

// Lista wszystkich dostępnych pionków (nie resetowana)
var PawnsList []BasePawn

// Lista pionków obecnie w grze (resetowana na nową grę)
var PawnsInGame []BasePawn

// Pionki obecnie umieszczone na planszy
var PawnsOnBoard []BasePawn

// Tablice pionków dla graczy (jeszcze przed rozstawieniem)
var Player1Pawns []BasePawn
var Player2Pawns []BasePawn

// Pionki dostpne do rzostawienia
var AvailablePawnsP1 []BasePawn // Pionki dostępne do wystawienia dla Gracza 1
var AvailablePawnsP2 []BasePawn // Pionki dostępne do wystawienia dla Gracza 2

var PawnTextures map[string]rl.Texture2D = make(map[string]rl.Texture2D) // Mapa przechowywanych tekstur
var TexturesLoaded bool = false                                          // Flaga, czy tekstury zostały już załadowane

// BasePawn reprezentuje podstawowe dane każdego pionka
type BasePawn struct {
	ID          int
	Type        string       // Typ pionka (Warrior, King itp.)
	X           int32        // Pozycja X na planszy
	Y           int32        // Pozycja Y na planszy
	IsAlive     bool         // Czy pionek jest żywy
	Owner       string       // Właściciel pionka
	Cost        int32        // Koszt w sklepie
	Texture     rl.Texture2D // Tekstura pionka
	AnimTexture rl.Texture2D // Sprite sheet z animacją
	Animation   Animation    // Animacja
}

type ShopPawn struct {
	Type  string
	Cost  int32
	Owner string
}

type PawnVisualConfig struct {
	StaticTexturePath string
	AnimTexturePath   string
	Anim              Animation
}

type Animation struct {
	FrameWidth   int
	FrameHeight  int
	FrameCount   int
	CurrentFrame int
	FrameTime    float32
	ElapsedTime  float32
}

var PawnIcons = make(map[string]rl.Texture2D)

var PawnVisualConfigs = map[string]PawnVisualConfig{
	"Warrior": {
		StaticTexturePath: filepath.Join("Assets", "Player", "Warrior", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Player", "Warrior", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.25},
	},
	"King": {
		StaticTexturePath: filepath.Join("Assets", "Player", "King", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Player", "King", "Idle.png"),
		Anim:              Animation{FrameWidth: 96, FrameHeight: 48, FrameCount: 6, FrameTime: 0.15},
	},
	"Knight": {
		StaticTexturePath: filepath.Join("Assets", "Player", "Knight", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Player", "Knight", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.25},
	},
	"Master": {
		StaticTexturePath: filepath.Join("Assets", "Player", "Master", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Player", "Master", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.25},
	},
	"Monk": {
		StaticTexturePath: filepath.Join("Assets", "Player", "Monk", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Player", "Monk", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.25},
	},
	"Boss": {
		StaticTexturePath: filepath.Join("Assets", "Mobs", "GiantBambooBoss", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Mobs", "GiantBambooBoss", "Idle.png"),
		Anim:              Animation{FrameWidth: 62, FrameHeight: 62, FrameCount: 6, FrameTime: 0.15},
	},
	"Lizard": {
		StaticTexturePath: filepath.Join("Assets", "Mobs", "Lizard", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Mobs", "Lizard", "Lizard.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.2},
	},
	"Reptile": {
		StaticTexturePath: filepath.Join("Assets", "Mobs", "Reptile", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Mobs", "Reptile", "Reptile2.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.2},
	},
	"Racoon": {
		StaticTexturePath: filepath.Join("Assets", "Mobs", "Racoon", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Mobs", "Racoon", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.2},
	},
	"LionWarrior": {
		StaticTexturePath: filepath.Join("Assets", "Mobs", "LionWarrior", "Faceset.png"),
		AnimTexturePath:   filepath.Join("Assets", "Mobs", "LionWarrior", "SpriteSheet.png"),
		Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0.2},
	},
}

/*// Interfejs Pawn dla pionków
type Pawn interface {
	GetPosition() (int32, int32)
	GetTexture() rl.Texture2D
	IsPawnAlive() bool
	GetOwner() string
	GetID() int
}

// Implementacja interfejsu Pawn dla BasePawn
func (p *BasePawn) GetPosition() (int32, int32) { return p.X, p.Y }
func (p *BasePawn) GetTexture() rl.Texture2D    { return p.Texture }
func (p *BasePawn) IsPawnAlive() bool           { return p.IsAlive }
func (p *BasePawn) GetOwner() string            { return p.Owner }
func (p *BasePawn) GetID() int                  { return p.ID }
*/

// ResetPawnsOnBoard czyści listę pionków na planszy
func ResetPawnsOnBoard() {
	PawnsOnBoard = []BasePawn{}
}

// Ładowanie pionków do głównej listy (wywoływane raz na start aplikacji)
func LoadPawns() {

	// Czyszczenie listy pionków
	PawnsList = []BasePawn{}

	// Wczytywanie przeciwników (nie wybieramy ich od razu, tylko wrzucamy do puli)
	for _, fight := range AvailableFights {
		for _, config := range fight {
			newPawn := CreatePawn(config.Type, config.Owner)
			PawnsList = append(PawnsList, newPawn)
		}
	}

}

// Kopiowanie pionków do gry (wywoływane przy rozpoczęciu nowej gry)
func LoadPawnsIntoGame() {

	// Reset list pionków
	PawnsInGame = make([]BasePawn, len(PawnsList))
	copy(PawnsInGame, PawnsList)

	Player2Pawns = []BasePawn{}

	// **Losowanie zestawu przeciwników**
	Player2Pawns = GetRandomFight()

}

// / Tworzenie pionka na podstawie typu
func CreatePawn(pawnType, owner string) BasePawn {
	cfg, exists := PawnVisualConfigs[pawnType]
	if !exists {
		cfg = PawnVisualConfig{
			StaticTexturePath: filepath.Join("Assets", "DefaultPawn.png"),
			AnimTexturePath:   filepath.Join("Assets", "DefaultAnim.png"),
			Anim:              Animation{FrameWidth: 16, FrameHeight: 16, FrameCount: 1, FrameTime: 0},
		}
	}

	staticTex := rl.LoadTexture(cfg.StaticTexturePath)
	animTex := rl.LoadTexture(cfg.AnimTexturePath)

	pawn := BasePawn{
		ID:          NextPawnID,
		Type:        pawnType,
		Owner:       owner,
		X:           -1,
		Y:           -1,
		IsAlive:     true,
		Texture:     staticTex,
		AnimTexture: animTex,
		Animation:   cfg.Anim,
	}

	NextPawnID++
	return pawn
}

func LoadAllPawnIcons() {
	for pawnType, cfg := range PawnVisualConfigs {
		if _, loaded := PawnIcons[pawnType]; !loaded {
			icon := rl.LoadTexture(cfg.StaticTexturePath)
			PawnIcons[pawnType] = icon
		}
	}
}

// **🔹 Zwolnienie pamięci po zamknięciu gry**
func UnloadPawnTextures() {
	if !TexturesLoaded {
		return
	}

	for _, texture := range PawnTextures {
		rl.UnloadTexture(texture)
	}

	PawnTextures = nil
	TexturesLoaded = false
}

func RemovePawnByID(pawnID int) {
	newPawnsOnBoard := []BasePawn{} // Nowa lista pionków (bez zbitego)

	for _, pawn := range PawnsOnBoard {
		if pawn.ID == pawnID {
			fmt.Printf(" Usuwanie pionka %s (ID: %d) z pozycji (%d, %d)\n", pawn.Type, pawn.ID, pawn.X, pawn.Y)
			continue
		}
		newPawnsOnBoard = append(newPawnsOnBoard, pawn)
	}

	PawnsOnBoard = newPawnsOnBoard // Nadpisujemy nową listą
}

func RemovePawnFromAvailableList(playerAvailablePawns *[]BasePawn, pawnID int) {
	for i, pawn := range *playerAvailablePawns {
		if pawn.ID == pawnID {
			*playerAvailablePawns = append((*playerAvailablePawns)[:i], (*playerAvailablePawns)[i+1:]...)
			fmt.Printf("Usunięto pionek ID %d z listy dostępnych do rozstawienia\n", pawnID)
			return
		}
	}
	fmt.Printf("Błąd: Nie znaleziono pionka o ID %d w liście dostępnych\n", pawnID)
}

// Animacje
func DrawPawnPro(pawn *BasePawn, cellSize, boardX, boardY int32) {
	anim := &pawn.Animation

	// Aktualizacja ramki animacji
	anim.ElapsedTime += rl.GetFrameTime()
	//fmt.Printf("Elapsed: %.2f / FrameTime: %.2f\n", anim.ElapsedTime, anim.FrameTime) //DEBUG ON

	if anim.ElapsedTime >= anim.FrameTime {
		anim.ElapsedTime = 0
		anim.CurrentFrame = (anim.CurrentFrame + 1) % anim.FrameCount
	}

	source := rl.Rectangle{
		X:      float32(anim.CurrentFrame * anim.FrameWidth),
		Y:      0,
		Width:  float32(anim.FrameWidth),
		Height: float32(anim.FrameHeight),
	}

	dest := rl.Rectangle{
		X:      float32(boardX + pawn.X*cellSize),
		Y:      float32(boardY + pawn.Y*cellSize),
		Width:  float32(cellSize),
		Height: float32(cellSize),
	}

	rl.DrawTexturePro(pawn.AnimTexture, source, dest, rl.Vector2{}, 0, rl.White)
	//fmt.Printf("Animating %s - frame %d\n", pawn.Type, pawn.Animation.CurrentFrame) //DEBUG ON

}

// DrawPawns rysuje wszystkie pionki umieszczone na planszy
func DrawPawns(cellSize, boardX, boardY int32) {
	for _, pawn := range PawnsOnBoard {
		if pawn.IsAlive {
			DrawPawnPro(&pawn, cellSize, boardX, boardY)
		}
	}
}

func UpdateAnimations() {
	for i := range PawnsOnBoard {
		anim := &PawnsOnBoard[i].Animation
		anim.ElapsedTime += rl.GetFrameTime()

		if anim.ElapsedTime >= anim.FrameTime {
			anim.ElapsedTime = 0
			anim.CurrentFrame = (anim.CurrentFrame + 1) % anim.FrameCount
		}

		fmt.Printf("Animating %s - frame %d\n", PawnsOnBoard[i].Type, anim.CurrentFrame)
	}
}

func IsTileOccupied(x, y int32, pawnsOnBoard []BasePawn) bool {
	for _, pawn := range pawnsOnBoard {
		if pawn.X == x && pawn.Y == y {
			return true // Pole zajęte
		}
	}
	return false // Pole wolne
}

func InitializeAvailablePawns() {
	AvailablePawnsP1 = make([]BasePawn, len(Player1Pawns))
	copy(AvailablePawnsP1, Player1Pawns)

	AvailablePawnsP2 = make([]BasePawn, len(Player2Pawns))
	copy(AvailablePawnsP2, Player2Pawns)

}

// **Losowanie zestawu przeciwników**
func GetRandomFight() []BasePawn {
	// **Sprawdzenie, czy są dostępne walki**
	if len(AvailableFights) == 0 {
		ResetAvailableFights()
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(AvailableFights))

	// **Pobranie wylosowanego zestawu przeciwników**
	selectedFight := AvailableFights[index]

	// **Usunięcie użytej walki, aby uniknąć powtórzeń**
	AvailableFights = append(AvailableFights[:index], AvailableFights[index+1:]...)

	// **Konwersja do listy pionków**
	var enemyPawns []BasePawn
	for _, pawnData := range selectedFight {
		enemyPawns = append(enemyPawns, CreatePawn(pawnData.Type, pawnData.Owner))
	}

	return enemyPawns
}

func GetRandomBossFight() []BasePawn {
	// **Sprawdzenie dostępnych bossfightów**
	if len(BossFights) == 0 {
		ResetBossFights()
	}

	// **Losowanie bossfightu**
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(BossFights))

	// **Pobranie wylosowanego zestawu przeciwników**
	selectedFight := BossFights[index]

	// **Usunięcie użytej walki, aby uniknąć powtórzeń**
	BossFights = append(BossFights[:index], BossFights[index+1:]...)

	// **Konwersja do listy pionków**
	var enemyPawns []BasePawn
	for _, pawnData := range selectedFight {
		enemyPawns = append(enemyPawns, CreatePawn(pawnData.Type, pawnData.Owner))
	}

	return enemyPawns
}

/*func GetPawnValue(pawnType string) int {
	switch pawnType {
	case "King":
		return 100 // Kluczowa jednostka - nie można stracić
	case "Warrior":
		return 30 // Podstawowy wojownik
	case "Boss":
		return 80 // Kluczowy pionek AI
	case "Lizard":
		return 30 // Podstawowa jednostak AI
	default:
		return 10 // Domyślna wartość dla innych jednostek
	}
}*/
