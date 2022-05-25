package model

import "time"

type Comment struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id" db:"user_id" `
	Content    string    `json:"content" db:"content" binding:"required"`
	CreateDate time.Time `json:"create_date" db:"create_date"`
	VideoID    int64     `json:"video_id" db:"video_id" binding:"required"`
}
