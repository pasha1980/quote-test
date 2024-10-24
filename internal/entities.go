package internal

type Quote struct {
	Id      string `json:"id"`
	Likes   int64  `json:"likes"`
	Author  string `json:"author"`
	Content string `json:"content"`
}
