package domain

type Error struct {
	Description string
}

func (e Error) Error() string {
	return e.Description
}
