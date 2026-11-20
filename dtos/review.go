package dtos

type CreateReviewRequest struct {
	Name   string `json:"name" binding:"required"`
	Score  int    `json:"score" binding:"required,min=1,max=5"`
	Body   string `json:"body" binding:"required"`
	ShopID uint   `json:"shopID" binding:"required"`
	UserID uint   `json:"userID" binding:"required"`
}
