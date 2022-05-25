package service

import (
	"BytesDanceProject/dao/mysql"
	"BytesDanceProject/model"
)

func CreateComment(p *model.Comment) error {

	return mysql.CreateComment(p)

}
