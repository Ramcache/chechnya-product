package models

import "time"

type Verification struct {
	Phone     string
	Code      string
	CreatedAt time.Time
	Confirmed bool
}
