package mysql

import (
	"github.com/jmoiron/sqlx"
	"github.com/xiaorui/web_app/models"
	"strings"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post (
	post_id, title, content, author_id, community_id)
	values (?,?,?,?,?)
`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select 
	post_id, title, content, author_id, community_id
	from post where post_id=?
`
	err = db.Get(post, sqlStr, pid)

	return
}

func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time
	from post
	ORDER BY create_time
	DESC 
	limit ?, ?
`
	posts = make([]*models.Post, 0, pageSize)
	err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize)

	return
}

func GetPostListByIDs(pids []string) (posts []*models.Post, err error) {
	sqlStr := `select 
	post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	ORDER BY FIND_IN_SET(post_id, ?)
`
	query, args, err := sqlx.In(sqlStr, pids, strings.Join(pids, ","))
	if err != nil {
		return
	}
	query = db.Rebind(query) // 就是要加上?

	err = db.Select(&posts, query, args...)
	return
}
