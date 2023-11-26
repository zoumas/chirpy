package database

// A Chirp is a text-only post, similar to twitter's Tweet.
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type ChirpRepository interface {
	Create(body string) (Chirp, error)
	GetAll() ([]Chirp, error)
	GetByID(id int) (Chirp, error)
}
