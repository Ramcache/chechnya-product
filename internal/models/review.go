package models

type Review struct {
	ID        int    `db:"id" json:"id"`
	OwnerID   string `db:"owner_id" json:"owner_id"`
	ProductID int    `db:"product_id" json:"product_id"`
	Rating    int    `db:"rating" json:"rating"`
	Comment   string `db:"comment" json:"comment"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type ReviewRequest struct {
	Rating  int    `json:"rating" example:"5"`
	Comment string `json:"comment" example:"Отличный товар!"`
}
