package models

import (
    "time"
)

type Song struct {
    ID          int       `json:"id"`
    GroupName   string    `json:"group" binding:"required"`
    SongName    string    `json:"song" binding:"required"`
    ReleaseDate string    `json:"releaseDate"`
    Text        string    `json:"text"`
    Link        string    `json:"link"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type SongFilter struct {
    GroupName   string `form:"group"`
    SongName    string `form:"song"`
    ReleaseDate string `form:"release_date"`
    Page        int    `form:"page,default=1"`
    PageSize    int    `form:"page_size,default=10"`
}

type SongDetail struct {
    ReleaseDate string `json:"releaseDate"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}

type VersePagination struct {
    Page     int `form:"verse_page,default=1"`
    PageSize int `form:"verse_size,default=4"`
}

type SongWithVerses struct {
    Song
    Verses     []string `json:"verses"`
    TotalVerses int     `json:"total_verses"`
    CurrentPage int     `json:"current_page"`
}
