package GameLoop

import (
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var currentMusic rl.Music
var currentTrack string = ""
var musicLoaded bool = false
var musicPlaying bool = false

// PlayMusic - ładuje i odtwarza wybraną ścieżkę muzyczną
func PlayMusic(track string) {
	// Jeśli już gra ta sama muzyka — nie rób nic
	if musicPlaying && currentTrack == track {
		return
	}

	// Zatrzymaj i zwolnij poprzednią muzykę
	if musicPlaying {
		rl.StopMusicStream(currentMusic)
		rl.UnloadMusicStream(currentMusic)
	}

	// Wybierz plik na podstawie identyfikatora ścieżki
	var path string
	switch track {
	case "menu":
		path = filepath.Join("Assets", "Music", "MenuTheme.ogg")
	case "selection":
		path = filepath.Join("Assets", "Music", "AdventureBegin.ogg")
	case "GameBoard":
		path = filepath.Join("Assets", "Music", "GameBoard.ogg")
	case "Fight":
		path = filepath.Join("Assets", "Music", "Fight.ogg")
	default:
		return
	}

	// Załaduj i odtwórz nową muzykę
	currentMusic = rl.LoadMusicStream(path)
	currentMusic.Looping = true
	rl.PlayMusicStream(currentMusic)

	currentTrack = track
	musicLoaded = true
	musicPlaying = true
}

// UpdateMusic - aktualizuje strumień muzyczny (wywołuj w każdej klatce)
func UpdateMusic() {
	if musicPlaying {
		rl.UpdateMusicStream(currentMusic)
	}
}

// StopMusic - zatrzymuje aktualnie odtwarzaną muzykę (bez zwalniania pamięci)
func StopMusic() {
	if musicPlaying {
		rl.StopMusicStream(currentMusic)
		musicPlaying = false
	}
}

// UnloadMusic - zwalnia wszystkie zasoby związane z muzyką
func UnloadMusic() {
	if musicLoaded {
		rl.UnloadMusicStream(currentMusic)
		rl.CloseAudioDevice()
	}

	currentTrack = ""
	musicLoaded = false
	musicPlaying = false
}
