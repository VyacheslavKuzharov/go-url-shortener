package softdelete

import uuid "github.com/satori/go.uuid"

type Object struct {
	UserID    uuid.UUID
	ShortKeys []string
}

func GenObjects(currentUserID uuid.UUID, slice []string, batchSize int) []Object {
	var objs []Object

	for i := 0; i < len(slice); i += batchSize {
		end := i + batchSize

		if end > len(slice) {
			end = len(slice)
		}

		newObjs := Object{
			UserID:    currentUserID,
			ShortKeys: slice[i:end],
		}

		objs = append(objs, newObjs)
	}

	return objs
}
