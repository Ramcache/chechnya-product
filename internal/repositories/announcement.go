package repositories

import (
	"chechnya-product/internal/handlers"
	"github.com/jmoiron/sqlx"
)

type AnnouncementRepository interface {
	Create(title, content string) (*handlers.Announcement, error)
	GetAll() ([]handlers.Announcement, error)
	GetByID(id int) (*handlers.Announcement, error)
	Update(id int, title, content string) error
	Delete(id int) error
}

type AnnouncementRepo struct {
	db *sqlx.DB
}

func NewAnnouncementRepo(db *sqlx.DB) *AnnouncementRepo {
	return &AnnouncementRepo{db: db}
}

func (r *AnnouncementRepo) Create(title, content string) (*handlers.Announcement, error) {
	var ann handlers.Announcement
	err := r.db.Get(&ann, `
		INSERT INTO announcements (title, content)
		VALUES ($1, $2) RETURNING id, title, content
	`, title, content)
	return &ann, err
}

func (r *AnnouncementRepo) GetAll() ([]handlers.Announcement, error) {
	var anns []handlers.Announcement
	err := r.db.Select(&anns, `SELECT * FROM announcements ORDER BY id DESC`)
	return anns, err
}

func (r *AnnouncementRepo) GetByID(id int) (*handlers.Announcement, error) {
	var ann handlers.Announcement
	err := r.db.Get(&ann, `SELECT * FROM announcements WHERE id = $1`, id)
	return &ann, err
}

func (r *AnnouncementRepo) Update(id int, title, content string) error {
	_, err := r.db.Exec(`UPDATE announcements SET title=$1, content=$2 WHERE id=$3`, title, content, id)
	return err
}

func (r *AnnouncementRepo) Delete(id int) error {
	_, err := r.db.Exec(`DELETE FROM announcements WHERE id=$1`, id)
	return err
}
