package font

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/gandarez/pong-multiplayer-go/assets"
)

// Font is a font manager.
type Font struct {
	assets *assets.Assets
	fonts  map[string]*text.GoTextFace
}

// New creates a new Font.
func New(assets *assets.Assets) *Font {
	return &Font{
		assets: assets,
		fonts:  make(map[string]*text.GoTextFace),
	}
}

// Face returns a text face with the given key and size.
func (f *Font) Face(key string, size float64) (*text.GoTextFace, error) {
	if f, ok := f.fonts[key]; ok {
		new := *f
		new.Size = size

		return &new, nil
	}

	source, err := f.assets.NewTextFaceSource(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create text face source: %w", err)
	}

	face := &text.GoTextFace{
		Source: source,
		Size:   size,
	}

	f.fonts[key] = face

	return face, nil
}
