package web

type TaskCreateRequest struct {
	Title string `validate:"required,min=1,max=100" json:"title"`
	Description string `validate:"required,min=1,max=200" json:"description"`
	Status string `validate:"required,min=1,max=20" json:"status"`
}