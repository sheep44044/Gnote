package models

type PersonalPage struct {
	ID        uint        `json:"id"`
	Username  string      `json:"username"`
	Avatar    string      `json:"avatar,omitempty"`
	Bio       string      `json:"bio,omitempty"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
	Documents []NoteBrief `json:"documents"` // 只返回笔记摘要，不暴露内容
}

type NoteBrief struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	IsPrivate bool   `json:"is_private"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
