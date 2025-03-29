package GameLoop

import (
	"fmt"
	"math/rand"
	"time"

	"Protect_The_King/Boards"
	"Protect_The_King/Pawns"
)

var recentAIMoves []int // Lista ID pionków, którymi AI ruszało w ostatnich ruchach

var bestAttack struct {
	pawn    *Pawns.BasePawn
	moveX   int32
	moveY   int32
	enemyID int
}

// Auto-Placement for AI pawns
func AutoPlaceAIPawns(board [][]Boards.Tile) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Lokalny generator liczb losowych

	// 1️⃣ Znalezienie kolumny, w której stoi Król Gracza 1
	var kingColumn int32 = -1
	for _, pawn := range Pawns.PawnsOnBoard {
		if pawn.Owner == "Player 1" && pawn.Type == "King" {
			kingColumn = pawn.X
			break
		}
	}

	// 2️⃣ Jeśli znaleziono kolumnę Króla, AI umieszcza tam 1 pionek (ale nie Bossa)
	placedInKingColumn := false
	if kingColumn != -1 {
		for _, pawn := range Pawns.AvailablePawnsP2 {
			if pawn.Type == "Boss" {
				continue // Pomijamy Bossa
			}

			for y := 0; y < 3; y++ { // AI może rozstawiać się tylko w pierwszych 3 rzędach
				if board[y][kingColumn].Walkable && !Pawns.IsTileOccupied(kingColumn, int32(y), Pawns.PawnsOnBoard) {
					// Ustawienie pionka na planszy
					pawn.X = kingColumn
					pawn.Y = int32(y)
					Pawns.PawnsOnBoard = append(Pawns.PawnsOnBoard, pawn)

					// Usunięcie pionka z listy dostępnych
					Pawns.RemovePawnFromAvailableList(&Pawns.AvailablePawnsP2, pawn.ID)

					fmt.Printf("AI umieściło pionek %s w kolumnie Króla na (%d, %d)\n", pawn.Type, kingColumn, y)
					placedInKingColumn = true
					break
				}
			}
			if placedInKingColumn {
				break
			}
		}
	}

	// 3️⃣ **Rozstawienie reszty pionków normalnie**
	for len(Pawns.AvailablePawnsP2) > 0 {
		pawn := Pawns.AvailablePawnsP2[0] // Pobranie pierwszego dostępnego pionka

		var placed bool
		for !placed {
			x := rng.Intn(len(board[0])) // Losowa kolumna
			y := rng.Intn(3)             // Ograniczenie do pierwszych 3 rzędów (AI)

			// Sprawdzenie, czy pole jest wolne i przechodnie
			if board[y][x].Walkable && !Pawns.IsTileOccupied(int32(x), int32(y), Pawns.PawnsOnBoard) {
				// Ustawienie pionka na planszy
				pawn.X = int32(x)
				pawn.Y = int32(y)
				Pawns.PawnsOnBoard = append(Pawns.PawnsOnBoard, pawn)

				// Usunięcie pionka z listy dostępnych
				Pawns.RemovePawnFromAvailableList(&Pawns.AvailablePawnsP2, pawn.ID)

				fmt.Printf("AI rozstawiło pionek %s na (%d, %d)\n", pawn.Type, x, y)
				placed = true
			}
		}
	}

	currentPhase = 2
}

// AI makes a move
func MakeAIMove() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Lokalny generator losowy

	// 🔹 **Ustal liczbę zapamiętanych ruchów**
	aiPawnCount := 0
	for _, pawn := range Pawns.PawnsOnBoard {
		if pawn.Owner == "Player 2" {
			aiPawnCount++
		}
	}
	maxRememberedMoves := min(3, aiPawnCount) // AI pamięta maksymalnie 3 ostatnie ruchy lub mniej, jeśli ma mniej pionków

	// 🔹 **Krok 1: Znalezienie Bossa AI**
	var boss *Pawns.BasePawn
	for i := range Pawns.PawnsOnBoard {
		if Pawns.PawnsOnBoard[i].Owner == "Player 2" && Pawns.PawnsOnBoard[i].Type == "Boss" {
			boss = &Pawns.PawnsOnBoard[i]
			break
		}
	}

	// 🔹 **Krok 2: Jeśli Boss jest pod szachem, wymuszamy jego ruch**
	if boss != nil && IsPawnUnderThreat(boss) {

		validMoves := GetValidMoves(boss, board, "Player 2")

		// **Jeśli Boss nie ma ruchu → AI przegrywa**
		if len(validMoves) == 0 {
			EndFight("Player 1", nil)
			return
		}

		// 🔥 **Najpierw sprawdzamy, czy Boss może wykonać bicie**
		var bestMove struct{ x, y int32 }
		hasCapture := false

		for _, move := range validMoves {
			enemyID, _ := FindEnemyPawnID(move.x, move.y, "Player 2")
			if enemyID != -1 { // Jeśli na polu jest przeciwnik – wybieramy ten ruch!
				bestMove = move
				hasCapture = true
				fmt.Printf("AI Boss priorytetowo atakuje przeciwnika (ID: %d) na (%d, %d)\n", enemyID, move.x, move.y)

				// **Usuń przeciwnika**
				Pawns.RemovePawnByID(enemyID)

				// **🔥 Teraz przesuwamy Bossa na pole atakowanego pionka!**
				boss.X = bestMove.x
				boss.Y = bestMove.y

				fmt.Printf(" AI Boss przesunięty na (%d, %d) po zbiciu!\n", bestMove.x, bestMove.y)
				break
			}
		}

		// **Jeśli Boss nie znalazł ruchu ataku, wybiera losowy ruch ucieczki**
		if !hasCapture {
			bestMove = validMoves[rng.Intn(len(validMoves))]
			fmt.Printf("AI Boss ucieka na (%d, %d)\n", bestMove.x, bestMove.y)
			boss.X = bestMove.x
			boss.Y = bestMove.y
		}

		currentTurn = swapTurn(currentTurn)
		return
	}

	// 🔹 **Krok 3: AI sprawdza, w której kolumnie znajduje się King Gracza 1**
	var king *Pawns.BasePawn
	for i := range Pawns.PawnsOnBoard {
		if Pawns.PawnsOnBoard[i].Owner == "Player 1" && Pawns.PawnsOnBoard[i].Type == "King" {
			king = &Pawns.PawnsOnBoard[i]
			break
		}
	}

	if king != nil {
		kingColumn := king.X
		pawnInColumn := false

		// **Sprawdzamy, czy AI ma już pionka w kolumnie Kinga**
		for _, pawn := range Pawns.PawnsOnBoard {
			if pawn.Owner == "Player 2" && pawn.X == kingColumn {
				pawnInColumn = true
				break
			}
		}

		// **Jeśli AI nie ma pionka w tej kolumnie, szukamy najlepszego do ruchu**
		if !pawnInColumn {
			fmt.Printf(" AI: W kolumnie %d brakuje pionka! AI przesuwa pionka w tę kolumnę.\n", kingColumn)

			var bestPawn *Pawns.BasePawn
			var bestMove struct{ x, y int32 }
			minDistance := int32(100)

			// **Przeszukanie pionków AI pod kątem możliwości wejścia do tej kolumny**
			for i := range Pawns.PawnsOnBoard {
				pawn := &Pawns.PawnsOnBoard[i]
				if pawn.Owner != "Player 2" {
					continue
				}

				validMoves := []struct{ x, y int32 }{}

				for _, move := range Pawns.PawnMoves[pawn.Type] {
					newX := pawn.X + move.DX
					newY := pawn.Y + move.DY

					if IsValidMove(pawn, newX, newY, board, "Player 2") {
						validMoves = append(validMoves, struct{ x, y int32 }{newX, newY})
					}
				}

				// **Sprawdzamy, czy któryś ruch pozwala wejść do tej kolumny**
				for _, move := range validMoves {
					distance := Abs(move.x - kingColumn)
					if distance < minDistance {
						bestPawn = pawn
						bestMove = move
						minDistance = distance
					}
				}
			}

			// **Jeśli znaleziono pionka, który może wejść do kolumny Kinga**
			if bestPawn != nil {
				fmt.Printf("AI przesuwa %s (ID: %d) na (%d, %d)\n", bestPawn.Type, bestPawn.ID, bestMove.x, bestMove.y)
				bestPawn.X = bestMove.x
				bestPawn.Y = bestMove.y
				currentTurn = swapTurn(currentTurn)
				return
			}
		}

		// **Jeśli AI ma pionka w kolumnie Kinga – nie rusza go**
		for i := range Pawns.PawnsOnBoard {
			pawn := &Pawns.PawnsOnBoard[i]
			if pawn.Owner == "Player 2" && pawn.X == kingColumn {
				fmt.Printf(" AI: Pionek %s (ID: %d) już jest w kolumnie %d i nie rusza się.\n", pawn.Type, pawn.ID, kingColumn)
				goto AttackPhase // Przechodzimy do sprawdzania ataków
			}
		}
	}

AttackPhase:
	// 🔹 **Krok 4: AI sprawdza, czy może zbić przeciwnika**

	hasAttack := false

	for i := range Pawns.PawnsOnBoard {
		pawn := &Pawns.PawnsOnBoard[i]

		if pawn.Owner != "Player 2" {
			continue
		}

		for _, move := range Pawns.PawnMoves[pawn.Type] {
			newX := pawn.X + move.DX
			newY := pawn.Y + move.DY

			if IsValidMove(pawn, newX, newY, board, "Player 2") {
				enemyID, _ := FindEnemyPawnID(newX, newY, "Player 2")

				if enemyID != -1 { // AI znalazło bicie!
					hasAttack = true

					// **Zapisanie najlepszego ruchu ataku**
					bestAttack.pawn = pawn
					bestAttack.moveX = newX
					bestAttack.moveY = newY
					bestAttack.enemyID = enemyID

					fmt.Printf("⚔ AI znalazło bicie: %s (ID: %d) może zaatakować na (%d, %d)\n",
						pawn.Type, pawn.ID, newX, newY)
				}
			}
		}
	}

	// **Jeśli AI znalazło ruch ataku – wykonuje go!**
	if hasAttack {
		fmt.Printf("✅ AI atakuje pionka gracza na (%d, %d)\n", bestAttack.moveX, bestAttack.moveY)

		attackingPawnID := bestAttack.pawn.ID
		Pawns.RemovePawnByID(bestAttack.enemyID)

		// Przesuwamy pionek AI na miejsce zbitego pionka
		for j := range Pawns.PawnsOnBoard {
			if Pawns.PawnsOnBoard[j].ID == attackingPawnID {
				Pawns.PawnsOnBoard[j].X = bestAttack.moveX
				Pawns.PawnsOnBoard[j].Y = bestAttack.moveY
				break
			}
		}

		currentTurn = swapTurn(currentTurn)
		return
	}

	// 🔹 **Krok 5: AI wykonuje losowy ruch, jeśli nie ma innej opcji**
	availablePawns := []*Pawns.BasePawn{} // Lista dostępnych pionków AI

	for i := range Pawns.PawnsOnBoard {
		pawn := &Pawns.PawnsOnBoard[i]

		if pawn.Owner != "Player 2" {
			continue
		}

		// 🔹 Sprawdzamy, czy AI już ruszało tym pionkiem w ostatnich X turach
		if !isRecentMove(pawn.ID, maxRememberedMoves) {
			availablePawns = append(availablePawns, pawn)
		}
	}

	// Jeśli nie ma żadnych dostępnych pionków, AI **resetuje** listę recentAIMoves, aby kontynuować grę
	if len(availablePawns) == 0 {
		recentAIMoves = []int{}                                         // Reset pamięci AI
		availablePawns = append(availablePawns, &Pawns.PawnsOnBoard[0]) // Awaryjne dodanie pionka
	}

	// **AI losuje pionek tylko spośród dostępnych**
	if len(availablePawns) > 0 {
		for i := range Pawns.PawnsOnBoard {
			pawn := &Pawns.PawnsOnBoard[i]

			if pawn.Owner != "Player 2" {
				continue
			}
			selectedPawnIndex := rng.Intn(len(availablePawns)) // Zapamiętaj indeks pionka
			selectedPawn := availablePawns[selectedPawnIndex]  // Pobierz wybrany pionek
			validMoves := []struct{ x, y int32 }{}

			for _, move := range Pawns.PawnMoves[selectedPawn.Type] {
				newX := pawn.X + move.DX
				newY := pawn.Y + move.DY

				if IsValidMove(selectedPawn, newX, newY, board, "Player 2") {
					validMoves = append(validMoves, struct{ x, y int32 }{newX, newY})
				}
			}

			if len(validMoves) > 0 {
				move := validMoves[rng.Intn(len(validMoves))] // Wybierz losowy ruch z dostępnych
				fmt.Printf(" AI przesuwa %s (ID: %d) z (%d, %d) na (%d, %d)\n",
					selectedPawn.Type, selectedPawn.ID, selectedPawn.X, selectedPawn.Y, move.x, move.y)

				// 🔹 **WAŻNE! Aktualizujemy właściwego pionka na planszy**
				for i := range Pawns.PawnsOnBoard {
					if Pawns.PawnsOnBoard[i].ID == selectedPawn.ID {
						Pawns.PawnsOnBoard[i].X = move.x
						Pawns.PawnsOnBoard[i].Y = move.y
						break
					}
				}
				// **Dodajemy pionek do pamięci recentAIMoves**
				updateRecentMoves(selectedPawn.ID, maxRememberedMoves)

				// **Zmieniamy turę**
				currentTurn = swapTurn(currentTurn)
				return
			}
		}
	}

	// 🔹 **Krok 6: AI kończy turę, jeśli nie ma ruchów**
	currentTurn = swapTurn(currentTurn)
}

// Aktualizuje historię ostatnich ruchów AI
func updateRecentMoves(pawnID int, maxRememberedMoves int) {
	recentAIMoves = append(recentAIMoves, pawnID)

	// Usuwanie najstarszego ruchu, jeśli przekroczyliśmy limit
	if len(recentAIMoves) > maxRememberedMoves {
		recentAIMoves = recentAIMoves[len(recentAIMoves)-maxRememberedMoves:] // Zatrzymujemy tylko ostatnie ruchy
	}
}

// Sprawdza, czy pionek był ruszany w ostatnich X ruchach
func isRecentMove(pawnID int, maxRememberedMoves int) bool {
	if len(recentAIMoves) == 0 {
		return false
	}

	// Sprawdzamy tylko ostatnie maxRememberedMoves ruchów
	for i := len(recentAIMoves) - 1; i >= max(0, len(recentAIMoves)-maxRememberedMoves); i-- {
		if recentAIMoves[i] == pawnID {
			return true
		}
	}
	return false
}

// **Pomocnicza funkcja do obliczania wartości bezwzględnej**
func Abs(value int32) int32 {
	if value < 0 {
		return -value
	}
	return value
}
