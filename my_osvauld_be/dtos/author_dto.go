package dtos

type CreateAuthor struct {
	Name string `json:"name" binding:"required"`
	Bio  string `json:"bio" binding:"required"`
}
