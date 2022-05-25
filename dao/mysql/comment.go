package mysql

import "BytesDanceProject/model"

func CreateComment(p *model.Comment) (err error) {
	comment := model.Comment{
		VideoID: p.VideoID,
		Content: p.Content,
		UserID:  p.UserID,
	}
	db.Create(&comment)
	return nil
}
