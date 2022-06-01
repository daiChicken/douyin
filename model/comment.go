package model

type Comment struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id" db:"user_id" `
	VideoID    int    `json:"video_id" db:"video_id" binding:"required"`
	Content    string `json:"content" db:"content" binding:"required"`
	CreateDate string `json:"create_date" db:"create_date"`
	IsDeleted  int    `db:"is_deleted"` //0表示显示；1表示删除
}

func (Comment) TableName() string {
	return "comment"
}
