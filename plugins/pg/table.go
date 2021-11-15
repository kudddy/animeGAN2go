package pg

import (
	"time"
)


type UserFileIdStatus struct {
	UserId int
	ChatId int
	Username string
	FileId string
	Date time.Time
	Status string
}


