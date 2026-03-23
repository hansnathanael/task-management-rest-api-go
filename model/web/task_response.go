package web

type TaskResponse struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status string `json:"status"`
}