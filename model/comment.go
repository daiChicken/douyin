package model

import (
	"encoding"
	"encoding/json"
	"time"
)

type Comment struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id" db:"user_id" `
	UserName   string    `json:"user_name" db:"user_name" binding:"required"`
	VideoID    int       `json:"video_id" db:"video_id" binding:"required"`
	Content    string    `json:"content" db:"content" binding:"required"`
	CreateDate time.Time `json:"create_date" db:"create_date"`
	IsDeleted  int       `db:"is_deleted"` //0表示显示；1表示删除
	UpdateDate time.Time `json:"update_date" db:"update_date"`
}

func (Comment) TableName() string {
	return "comment"
}

var _ encoding.BinaryMarshaler = new(Comment)
var _ encoding.BinaryUnmarshaler = new(Comment)

func (m *Comment) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *Comment) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)

}
