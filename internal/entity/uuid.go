package entity

import "github.com/google/uuid"

type UUID struct {
	uuid.UUID
}

var EmptyUUID = UUID{uuid.Nil}

func ParseUUID(str string) UUID {
	parse, _ := uuid.Parse(str)
	return UUID{parse}
}
