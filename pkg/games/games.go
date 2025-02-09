package games

import (
	"time"
)

type AnalyzerOption func(*Analyzer)

func WithLastPlayedBefore(lastPlayedBefore time.Time) AnalyzerOption {
	return func(a *Analyzer) {
		a.lastPlayedBefore = lastPlayedBefore
	}
}

func WithMaxPlaytime(maxPlaytime time.Duration) AnalyzerOption {
	return func(a *Analyzer) {
		a.maxPlaytime = maxPlaytime
	}
}

type Analyzer struct {
	lastPlayedBefore time.Time
	maxPlaytime      time.Duration
}

func NewAnalyzer(options ...AnalyzerOption) *Analyzer {
	analyzer := &Analyzer{
		lastPlayedBefore: time.Now().AddDate(-1, 0, 0),
		maxPlaytime:      time.Hour * 20,
	}
	for _, option := range options {
		option(analyzer)
	}
	return analyzer
}

func (a *Analyzer) Analyze() ([]SteamGame, error) {
	gms, err := getSteamGames()
	if err != nil {
		return nil, err
	}
	var result []SteamGame
	for _, g := range gms {
		if g.LastPlayed.Before(a.lastPlayedBefore) && g.Playtime < a.maxPlaytime {
			result = append(result, g)
		}
	}
	return result, nil
}
