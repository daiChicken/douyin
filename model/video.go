package model

/**
 * @author  daijizai
 * @date  2022/5/10 20:22
 * @version  1.0
 * @description
 */

type Video struct {
	Id            int    `db:"id"`
	AuthorId      int    `db:"author_id"`
	PlayUrl       string `db:"play_url"`
	CoverUrl      string `db:"cover_url"`
	FavoriteCount int    `db:"favorite_count"`
	CommentCount  int    `db:"comment_count"`
	//IsFavorite    int    `db:"is_favorite"`
	CreateTime int64 `db:"create_time"`
	IsDeleted  int   `db:"is_deleted"`
}
