package menu

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gandarez/pong-multiplayer-go/internal/network"
	"github.com/gandarez/pong-multiplayer-go/internal/ui"
	"github.com/gandarez/pong-multiplayer-go/pkg/geometry"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// sessionInfo represents the information of a session.
type sessionInfo struct {
	ID      string `json:"id"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2"`
}

// spectateSessionsState represents the state where the player can spectate active sessions.
type spectateSessionsState struct {
	menu          *Menu
	sessions      []sessionInfo
	selectedIndex int
	fetched       bool
	errorMessage  string
}

// newSpectateSessionsState creates a new spectateSessionsState.
func newSpectateSessionsState(menu *Menu) *spectateSessionsState {
	return &spectateSessionsState{
		menu: menu,
	}
}

// Update updates the state.
func (s *spectateSessionsState) Update() {
	if !s.fetched {
		s.fetchSessions()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if s.selectedIndex < len(s.sessions)-1 {
			s.selectedIndex++
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if s.selectedIndex > 0 {
			s.selectedIndex--
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(s.sessions) > 0 {
		selectedSession := s.sessions[s.selectedIndex]
		s.menu.SessionID = selectedSession.ID
		s.menu.ChangeState(newSpectatorConnectingState(s.menu))
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		s.menu.ChangeState(newMainMenuState(s.menu))
	}
}

// Draw draws the state.
func (s *spectateSessionsState) Draw(screen *ebiten.Image) {
	if s.errorMessage != "" {
		s.drawErrorMessage(screen)
		return
	}

	if !s.fetched {
		s.drawFetchingMessage(screen)
		return
	}

	if len(s.sessions) == 0 {
		s.drawNoSessionsMessage(screen)
		return
	}

	s.drawSessionsList(screen)
}

// String returns the string representation of the state.
func (*spectateSessionsState) String() string {
	return "spectateSessionsState"
}

func (s *spectateSessionsState) fetchSessions() {
	sessions, err := fetchSessions()
	if err != nil {
		s.errorMessage = "Failed to fetch sessions"
		return
	}

	s.sessions = sessions
	s.fetched = true
}

func (s *spectateSessionsState) drawFetchingMessage(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("Failed to create text face", slog.Any("error", err))
		return
	}

	message := "Fetching sessions..."
	width, _ := text.Measure(message, textFace, 1)
	uiText := ui.Text{
		Value:    message,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - width) / 2,
			Y: 250.0,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)
}

func (s *spectateSessionsState) drawNoSessionsMessage(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))
		return
	}

	message := "No active sessions"
	width, _ := text.Measure(message, textFace, 1)
	uiText := ui.Text{
		Value:    message,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - width) / 2,
			Y: 250.0,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)
}

func (s *spectateSessionsState) drawErrorMessage(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))
		return
	}

	message := s.errorMessage
	width, _ := text.Measure(message, textFace, 1)
	uiText := ui.Text{
		Value:    message,
		FontFace: textFace,
		Position: geometry.Vector{
			X: (float64(s.menu.screenWidth) - width) / 2,
			Y: 250.0,
		},
		Color: ui.DefaultColor,
	}
	uiText.Draw(screen)
}

func (s *spectateSessionsState) drawSessionsList(screen *ebiten.Image) {
	textFace, err := s.menu.font.Face("ui", 20)
	if err != nil {
		slog.Error("failed to create text face", slog.Any("error", err))
		return
	}

	y := 200.0

	for i, session := range s.sessions {
		sessionTitle := fmt.Sprintf("%s X %s", session.Player1, session.Player2)
		width, _ := text.Measure(sessionTitle, textFace, 1)

		color := ui.DefaultColor
		if i == s.selectedIndex {
			color = ui.HighlightColor
		}

		uiText := ui.Text{
			Value:    sessionTitle,
			FontFace: textFace,
			Position: geometry.Vector{
				X: (float64(s.menu.screenWidth) - width) / 2,
				Y: y,
			},
			Color: color,
		}
		uiText.Draw(screen)

		y += 40.0
	}
}

func fetchSessions() ([]sessionInfo, error) {
	resp, err := http.Get("https://" + network.BaseURL + "/sessions")
	if err != nil {
		slog.Error("failed to fetch sessions", slog.Any("error", err))
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", slog.Any("error", err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to fetch sessions", slog.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	var sessions []sessionInfo
	if err := json.NewDecoder(resp.Body).Decode(&sessions); err != nil {
		slog.Error("failed to parse sessions", slog.Any("error", err))
		return nil, err
	}

	return sessions, nil
}
