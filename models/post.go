package models

import "time"

type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
}

//帖子详情内容的结构体
type ApiPostDetail struct {
	*Post
	*CommunityDetails `json:"community"`
	AuthorName        string `json:"author_name"`
	VoteNumber        int64  `json:"vote_number"`
}
