package api

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "music-library/internal/models"
    "music-library/internal/service"
    "net/http"
    "strconv"
)

type Handler struct {
    songService *service.SongService
    logger      *zap.Logger
}

func NewHandler(songService *service.SongService, logger *zap.Logger) *Handler {
    return &Handler{
        songService: songService,
        logger:      logger,
    }
}

// @Summary Create a new song
// @Description Create a new song with the provided information
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song object"
// @Success 201 {object} models.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs [post]
func (h *Handler) CreateSong(c *gin.Context) {
    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        h.logger.Error("Failed to bind JSON", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
        return
    }

    if err := h.songService.CreateSong(&song); err != nil {
        h.logger.Error("Failed to create song", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusCreated, song)
}

// @Summary Update a song
// @Description Update an existing song's information
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song object"
// @Success 200 {object} models.Song
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [put]
func (h *Handler) UpdateSong(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        h.logger.Error("Invalid song ID", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
        return
    }

    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        h.logger.Error("Failed to bind JSON", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
        return
    }

    song.ID = id
    if err := h.songService.UpdateSong(&song); err != nil {
        h.logger.Error("Failed to update song", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, song)
}

// @Summary Delete a song
// @Description Delete a song by its ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        h.logger.Error("Invalid song ID", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
        return
    }

    if err := h.songService.DeleteSong(id); err != nil {
        h.logger.Error("Failed to delete song", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }

    c.Status(http.StatusNoContent)
}

// @Summary Get a song
// @Description Get a song by its ID
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id} [get]
func (h *Handler) GetSong(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        h.logger.Error("Invalid song ID", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
        return
    }

    song, err := h.songService.GetSong(id)
    if err != nil {
        h.logger.Error("Failed to get song", zap.Error(err))
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, song)
}

// @Summary Get a song with verses
// @Description Get a song by its ID with paginated verses
// @Tags songs
// @Produce json
// @Param id path int true "Song ID"
// @Param verse_page query int false "Verse page number (default: 1)"
// @Param verse_size query int false "Verses per page (default: 4)"
// @Success 200 {object} models.SongWithVerses
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs/{id}/verses [get]
func (h *Handler) GetSongVerses(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        h.logger.Error("Invalid song ID", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid song ID"})
        return
    }

    var pagination models.VersePagination
    if err := c.ShouldBindQuery(&pagination); err != nil {
        h.logger.Error("Failed to bind query parameters", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters"})
        return
    }

    song, err := h.songService.GetSongWithVerses(id, &pagination)
    if err != nil {
        h.logger.Error("Failed to get song verses",
            zap.Error(err),
            zap.Int("id", id))
        c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, song)
}

// @Summary List songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Produce json
// @Param group query string false "Filter by group name"
// @Param song query string false "Filter by song name"
// @Param release_date query string false "Filter by release date"
// @Param page query int false "Page number (default: 1)"
// @Param page_size query int false "Page size (default: 10)"
// @Success 200 {array} models.Song
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /songs [get]
func (h *Handler) ListSongs(c *gin.Context) {
    var filter models.SongFilter
    if err := c.ShouldBindQuery(&filter); err != nil {
        h.logger.Error("Failed to bind query parameters", zap.Error(err))
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid query parameters"})
        return
    }

    songs, err := h.songService.ListSongs(&filter)
    if err != nil {
        h.logger.Error("Failed to list songs", zap.Error(err))
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, songs)
}

type ErrorResponse struct {
    Error string `json:"error"`
}
