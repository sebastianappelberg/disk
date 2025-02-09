package torrents

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Torrent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Username string `json:"username"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	Category string `json:"category"`
	IMDb     string `json:"imdb"`
}

func Search(query string) ([]Torrent, error) {
	searchURL := fmt.Sprintf("https://apibay.org/q.php?q=%s", url.QueryEscape(query))
	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("response not ok: %d, %s", resp.StatusCode, string(body))
	}

	var torrents []Torrent
	err = json.NewDecoder(resp.Body).Decode(&torrents)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return torrents, nil
}
