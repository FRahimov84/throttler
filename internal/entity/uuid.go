package entity

import "github.com/google/uuid"

type UUID struct {
	uuid.UUID
}

var EmptyUUID = UUID{uuid.Nil}

// ParseUUID parses from string UUID
func ParseUUID(str string) (UUID, error) {
	parse, err := uuid.Parse(str)
	return UUID{parse}, err
}

// New generate new UUID
func New() UUID {
	u := uuid.New()
	return UUID{u}
}
