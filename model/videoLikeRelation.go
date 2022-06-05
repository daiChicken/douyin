package model

type VideoLikeRelation struct {
	Id      int64 `json:"id"`
	VideoId int64 `json:"video_id"`
	UserId  int64 `json:"user_id"`
	Status  int32 `json:"status"`
}
