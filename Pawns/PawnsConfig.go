package Pawns

type Move struct {
	DX int32 // Przesunięcie w poziomie (X)
	DY int32 // Przesunięcie w pionie (Y)
}

type Fight struct {
	Type  string `json:"type"`
	Owner string `json:"owner"`
}

// Definicja możliwych zestawów przeciwników
// Standardowe walki
var AvailableFights = [][]Fight{
	{
		{"Racoon", "Player 2"},
		{"Racoon", "Player 2"},
		{"Racoon", "Player 2"},
		{"Racoon", "Player 2"},
		{"Racoon", "Player 2"},
		{"Racoon", "Player 2"},
	},
	{
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
	},
	{

		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
	},
}

// Walki bossfight – trudniejsze zestawy przeciwników
var BossFights = [][]Fight{
	{
		{"Boss", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Lizard", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
		{"Reptile", "Player 2"},
	},
}

// Resetowanie obu list
func ResetAvailableFights() {
	AvailableFights = [][]Fight{
		{
			{"Racoon", "Player 2"},
			{"Racoon", "Player 2"},
			{"Racoon", "Player 2"},
			{"Racoon", "Player 2"},
			{"Racoon", "Player 2"},
			{"Racoon", "Player 2"},
		},
		{
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
		},
		{

			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
		},
	}
}

func ResetBossFights() {
	BossFights = [][]Fight{
		{
			{"Boss", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Lizard", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
			{"Reptile", "Player 2"},
		},
	}
}

// PawnMoves przechowuje zasady ruchu dla każdego typu pionka
var PawnMoves = map[string][]Move{
	"Warrior": {
		{0, 1}, {0, -1}, // Pionowo góra/dół
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1}, // Po skosie
	},
	"Knight": {
		// Ruch w górę i w dół całą kolumną
		{0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9}, {0, 10}, {0, 11},
		{0, -1}, {0, -2}, {0, -3}, {0, -4}, {0, -5}, {0, -6}, {0, -7}, {0, -8}, {0, -9}, {0, -10}, {0, -11},

		// Ruch w lewo i w prawo całą kolumną
		{1, 0}, {2, 0}, {3, 0}, {4, 0}, {5, 0}, {6, 0}, {7, 0}, {8, 0}, {9, 0}, {10, 0}, {11, 0},
		{-1, 0}, {-2, 0}, {-3, 0}, {-4, 0}, {-5, 0}, {-6, 0}, {-7, 0}, {-8, 0}, {-9, 0}, {-10, 0}, {-11, 0},
	},
	"Monk": {
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	},
	"King": {
		{0, 1}, {0, -1}, // Pionowo
		{1, 0}, {-1, 0}, // Poziomo
	},
	"Master": {
		{1, 1}, {2, 2}, {3, 3}, {4, 4}, {5, 5}, {6, 6}, {7, 7}, {8, 8}, {9, 9}, {10, 10},
		{-1, -1}, {-2, -2}, {-3, -3}, {-4, -4}, {-5, -5}, {-6, -6}, {-7, -7}, {-8, -8}, {-9, -9}, {-10, -10},
		{1, -1}, {2, -2}, {3, -3}, {4, -4}, {5, -5}, {6, -6}, {7, -7}, {8, -8}, {9, -9}, {10, -10},
		{-1, 1}, {-2, 2}, {-3, 3}, {-4, 4}, {-5, 5}, {-6, 6}, {-7, 7}, {-8, 8}, {-9, 9}, {-10, 10},
	},
	"Boss": {
		{0, 1}, {0, -1}, // Pionowo
		{1, 0}, {-1, 0}, // Poziomo
	},
	"Lizard": {
		{0, 1}, {0, -1}, // Pionowo góra/dół
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1}, // Po skosie
	},
	"Reptile": {
		{2, 1}, {2, -1}, {-2, 1}, {-2, -1},
		{1, 2}, {1, -2}, {-1, 2}, {-1, -2},
	},
	"Racoon": {
		{0, 1}, {0, 2}, {0, 3}, {0, -1}, {0, -2}, {0, -3},
		{1, 0}, {2, 0}, {3, 0}, {-1, 0}, {-2, 0}, {-3, 0},
	},
	"LionWarrior": {
		{0, 2}, {0, 3}, {0, 1}, {0, -1}, // Pionowo góra/dół
		{1, 1}, {2, 2}, {3, 3}, {-1, 1}, {-2, 2}, {-3, 3}, // Po skosie
	},
}
