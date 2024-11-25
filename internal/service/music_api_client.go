package service

import (
    "encoding/json"
    "fmt"
    "music-library/internal/models"
    "net/http"
    "net/url"
)

type MusicAPIClient struct {
    baseURL string
    client  *http.Client
}

func NewMusicAPIClient(baseURL string) *MusicAPIClient {
    return &MusicAPIClient{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (c *MusicAPIClient) GetSongInfo(group, song string) (*models.SongDetail, error) {
    params := url.Values{}
    params.Add("group", group)
    params.Add("song", song)

    url := fmt.Sprintf("%s/info?%s", c.baseURL, params.Encode())
    
    resp, err := c.client.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
    }

    var songDetail models.SongDetail
    if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }

    return &songDetail, nil
}
