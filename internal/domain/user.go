package domain

import "time"

type User struct {
	ID        int32
	Email     string
	Password  string
	CreatedAt time.Time
}
