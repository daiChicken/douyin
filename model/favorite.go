package model

type FavoriteRequest struct {
	UserID     int64  `json:"user_id"` // form:"user_id" gorm:"column:user_id; type:int64"
	Token      string `json:"token" `
	VideoID    int64  `json:"video_id"`
	ActionType int32  `json:"action_type"`
}

type FavoriteListRequest struct {
	UserID int64  `json:"user_id"` // form:"user_id" gorm:"column:user_id; type:int64"
	Token  string `json:"token" `
}
