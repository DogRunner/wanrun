package dto

import "time"

type CheckinsRes struct {
	DogID       int64     `json:"dog_id"`
	DogrunID    int64     `json:"dogrun_id"`
	CheckinAt   time.Time `json:"checkin_at"`
	ReCheckinAt time.Time `json:"re_checkin_at"`
}
