package mysql

import (
	"BytesDanceProject/model"
	"gorm.io/gorm"
)

func CreateFavoriteAction(vlr *model.VideoLikeRelation) (*gorm.DB, error) {
	dbWithTransaction := db.Begin()
	if err := dbWithTransaction.Create(&vlr).Error; err != nil {
		return nil, err
	}
	return dbWithTransaction, nil
}
