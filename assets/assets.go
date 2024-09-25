package assets

import (
	"bytes"
	"embed"
	"fmt"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const AudioSampleRate = 44100

//go:embed fonts/*.ttf
var _fonts embed.FS

//go:embed sounds/*.mp3
var _sounds embed.FS

// Assets contains all the assets of the game.
type Assets struct {
	fonts  map[string][]byte
	sounds map[string][]byte
}

// Load loads all the assets of the game.
func Load() (*Assets, error) {
	assets := new(Assets)

	// load fonts
	fonts := map[string]string{
		"score": "fonts/score.ttf",
		"stat":  "fonts/stat.ttf",
		"ui":    "fonts/ui.ttf",
	}

	assets.fonts = make(map[string][]byte, len(fonts))

	for key, path := range fonts {
		f, err := _fonts.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read font file %q: %w", key, err)
		}

		assets.fonts[key] = f
	}

	// load sounds
	sounds := map[string]string{
		"paddle": "sounds/bounce-paddle.mp3",
		"wall":   "sounds/bounce-wall.mp3",
		"score":  "sounds/score.mp3",
	}

	assets.sounds = make(map[string][]byte, len(sounds))

	for key, path := range sounds {
		s, err := _sounds.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read sound file %q: %w", key, err)
		}

		assets.sounds[key] = s
	}

	return assets, nil
}

// NewTextFaceSource creates a new text face source from the given font key.
func (a *Assets) NewTextFaceSource(key string) (*text.GoTextFaceSource, error) {
	data, ok := a.fonts[key]
	if !ok {
		return nil, fmt.Errorf("font %q not found", key)
	}

	tt, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create font face source %q: %w", key, err)
	}

	return tt, nil
}

// AllFonts returns all the fonts keys.
func (a *Assets) AllFonts() []string {
	return slices.Collect(maps.Keys(a.fonts))
}

// GetAudioPlayer returns an audio player for the given sound key.
func (a *Assets) GetAudioPlayer(audioContext *audio.Context, key string) (*audio.Player, error) {
	data, ok := a.sounds[key]
	if !ok {
		return nil, fmt.Errorf("sound %q not found", key)
	}

	// decode the MP3 file
	d, err := mp3.DecodeWithSampleRate(AudioSampleRate, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode sound %q: %w", key, err)
	}

	// create an audio player for the decoded data
	player, err := audioContext.NewPlayer(d)
	if err != nil {
		return nil, fmt.Errorf("failed to create audio player for %q: %w", key, err)
	}

	return player, nil
}
