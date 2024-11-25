package repository

import (
    "database/sql"
    "fmt"
    "music-library/internal/models"
)

type SongRepository struct {
    db *sql.DB
}

func NewSongRepository(db *sql.DB) *SongRepository {
    return &SongRepository{db: db}
}

func (r *SongRepository) Create(song *models.Song) error {
    query := `
        INSERT INTO songs (group_name, song_name, release_date, text, link)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at`

    return r.db.QueryRow(
        query,
        song.GroupName,
        song.SongName,
        song.ReleaseDate,
        song.Text,
        song.Link,
    ).Scan(&song.ID, &song.CreatedAt, &song.UpdatedAt)
}

func (r *SongRepository) Update(song *models.Song) error {
    query := `
        UPDATE songs 
        SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6
        RETURNING updated_at`

    return r.db.QueryRow(
        query,
        song.GroupName,
        song.SongName,
        song.ReleaseDate,
        song.Text,
        song.Link,
        song.ID,
    ).Scan(&song.UpdatedAt)
}

func (r *SongRepository) Delete(id int) error {
    query := "DELETE FROM songs WHERE id = $1"
    result, err := r.db.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("song with id %d not found", id)
    }

    return nil
}

func (r *SongRepository) GetByID(id int) (*models.Song, error) {
    song := &models.Song{}
    query := `
        SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
        FROM songs
        WHERE id = $1`

    err := r.db.QueryRow(query, id).Scan(
        &song.ID,
        &song.GroupName,
        &song.SongName,
        &song.ReleaseDate,
        &song.Text,
        &song.Link,
        &song.CreatedAt,
        &song.UpdatedAt,
    )
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("song with id %d not found", id)
    }
    return song, err
}

func (r *SongRepository) List(filter *models.SongFilter) ([]models.Song, error) {
    query := `
        SELECT id, group_name, song_name, release_date, text, link, created_at, updated_at
        FROM songs
        WHERE ($1 = '' OR group_name ILIKE '%' || $1 || '%')
        AND ($2 = '' OR song_name ILIKE '%' || $2 || '%')
        AND ($3 = '' OR release_date::text LIKE $3)
        ORDER BY id
        LIMIT $4 OFFSET $5`

    offset := (filter.Page - 1) * filter.PageSize

    rows, err := r.db.Query(
        query,
        filter.GroupName,
        filter.SongName,
        filter.ReleaseDate,
        filter.PageSize,
        offset,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var songs []models.Song
    for rows.Next() {
        var song models.Song
        err := rows.Scan(
            &song.ID,
            &song.GroupName,
            &song.SongName,
            &song.ReleaseDate,
            &song.Text,
            &song.Link,
            &song.CreatedAt,
            &song.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }
        songs = append(songs, song)
    }

    return songs, nil
}
