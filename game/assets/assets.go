package assets

import (
	"bytes"
	"embed"
	"fmt"
	"maps"
	"slices"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed fonts/*.ttf
var _fonts embed.FS

type Assets struct {
	fonts map[string][]byte
}

func Load() (*Assets, error) {
	assets := new(Assets)

	// Load fonts
	fonts := map[string]string{
		"score": "fonts/score.ttf",
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

	return assets, nil
}

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

func (a *Assets) AllFonts() []string {
	return slices.Collect(maps.Keys(a.fonts))
}
