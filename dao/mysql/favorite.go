package mysql

import "BytesDanceProject/model"

func CreateFavoriteAction(vlr *model.VideoLikeRelation) error {
	if err := db.Create(&vlr).Error; err != nil {
		return err
	}
	return nil
}
