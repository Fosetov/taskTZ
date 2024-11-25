package service

import (
    "fmt"
    "go.uber.org/zap"
    "music-library/internal/models"
    "music-library/internal/repository"
    "strings"
)

type SongService struct {
    repo      *repository.SongRepository
    apiClient *MusicAPIClient
    logger    *zap.Logger
}

func NewSongService(repo *repository.SongRepository, apiClient *MusicAPIClient, logger *zap.Logger) *SongService {
    return &SongService{
        repo:      repo,
        apiClient: apiClient,
        logger:    logger,
    }
}

func (s *SongService) CreateSong(song *models.Song) error {
    s.logger.Info("Creating new song",
        zap.String("group", song.GroupName),
        zap.String("song", song.SongName))

    // Get additional info from external API
    songDetail, err := s.apiClient.GetSongInfo(song.GroupName, song.SongName)
    if err != nil {
        s.logger.Error("Failed to get song info from API",
            zap.Error(err),
            zap.String("group", song.GroupName),
            zap.String("song", song.SongName))
        return fmt.Errorf("failed to get song info: %w", err)
    }

    // Enrich song with API data
    song.ReleaseDate = songDetail.ReleaseDate
    song.Text = songDetail.Text
    song.Link = songDetail.Link

    if err := s.repo.Create(song); err != nil {
        s.logger.Error("Failed to create song in database",
            zap.Error(err),
            zap.String("group", song.GroupName),
            zap.String("song", song.SongName))
        return fmt.Errorf("failed to create song: %w", err)
    }

    s.logger.Info("Successfully created song",
        zap.Int("id", song.ID),
        zap.String("group", song.GroupName),
        zap.String("song", song.SongName))

    return nil
}

func (s *SongService) UpdateSong(song *models.Song) error {
    s.logger.Info("Updating song",
        zap.Int("id", song.ID),
        zap.String("group", song.GroupName),
        zap.String("song", song.SongName))

    if err := s.repo.Update(song); err != nil {
        s.logger.Error("Failed to update song",
            zap.Error(err),
            zap.Int("id", song.ID))
        return fmt.Errorf("failed to update song: %w", err)
    }

    s.logger.Info("Successfully updated song", zap.Int("id", song.ID))
    return nil
}

func (s *SongService) DeleteSong(id int) error {
    s.logger.Info("Deleting song", zap.Int("id", id))

    if err := s.repo.Delete(id); err != nil {
        s.logger.Error("Failed to delete song",
            zap.Error(err),
            zap.Int("id", id))
        return fmt.Errorf("failed to delete song: %w", err)
    }

    s.logger.Info("Successfully deleted song", zap.Int("id", id))
    return nil
}

func (s *SongService) GetSong(id int) (*models.Song, error) {
    s.logger.Debug("Getting song by ID", zap.Int("id", id))

    song, err := s.repo.GetByID(id)
    if err != nil {
        s.logger.Error("Failed to get song",
            zap.Error(err),
            zap.Int("id", id))
        return nil, fmt.Errorf("failed to get song: %w", err)
    }

    return song, nil
}

func (s *SongService) GetSongWithVerses(id int, pagination *models.VersePagination) (*models.SongWithVerses, error) {
    s.logger.Debug("Getting song with verses",
        zap.Int("id", id),
        zap.Int("page", pagination.Page),
        zap.Int("page_size", pagination.PageSize))

    song, err := s.repo.GetByID(id)
    if err != nil {
        s.logger.Error("Failed to get song",
            zap.Error(err),
            zap.Int("id", id))
        return nil, fmt.Errorf("failed to get song: %w", err)
    }

    // Split text into verses (assuming verses are separated by empty lines)
    verses := strings.Split(strings.TrimSpace(song.Text), "\n\n")
    totalVerses := len(verses)

    // Calculate pagination
    startIdx := (pagination.Page - 1) * pagination.PageSize
    endIdx := startIdx + pagination.PageSize
    if endIdx > totalVerses {
        endIdx = totalVerses
    }

    // Check if page is valid
    if startIdx >= totalVerses {
        return nil, fmt.Errorf("page number exceeds total verses")
    }

    result := &models.SongWithVerses{
        Song:        *song,
        Verses:      verses[startIdx:endIdx],
        TotalVerses: totalVerses,
        CurrentPage: pagination.Page,
    }

    return result, nil
}

func (s *SongService) ListSongs(filter *models.SongFilter) ([]models.Song, error) {
    s.logger.Debug("Listing songs with filter",
        zap.Any("filter", filter))

    songs, err := s.repo.List(filter)
    if err != nil {
        s.logger.Error("Failed to list songs",
            zap.Error(err),
            zap.Any("filter", filter))
        return nil, fmt.Errorf("failed to list songs: %w", err)
    }

    return songs, nil
}
