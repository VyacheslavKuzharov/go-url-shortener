package entity

import uuid "github.com/satori/go.uuid"

// CurrentUserID - Describes Context Value key
var CurrentUserID uuid.UUID

type User struct {
	UUID uuid.UUID `json:"uuid"`
}
