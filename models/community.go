package models

import "time"

type Community struct {
	ID   int64  `json:"id,string" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetails struct {
	ID           int64     `json:"id,string" db:"community_id"`
	Name         string    `json:"name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"` //omitempty 表示为空时不展示
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
