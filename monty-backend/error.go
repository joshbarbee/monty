package main

type ErrInvalidPassword struct{}

type ErrUserExists struct {
	Username string
}

func (e ErrInvalidPassword) Error() string {
	return "Invalid password"
}

func (e ErrUserExists) Error() string {
	return "User already exists: " + e.Username
}
