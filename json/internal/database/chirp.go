package database

// A Chirp is a text-only post, similar to twitter's Tweet.
type Chirp struct {
	ID     int    `json:"id"`
	Body   string `json:"body"`
	UserID int    `json:"author_id"`
}

type CreateChirpParams struct {
	Body   string
	UserID int
}

type DeleteChirpParams struct {
	ID     int
	UserID int
}

type ChirpRepository interface {
	Create(params CreateChirpParams) (Chirp, error)
	GetAll() ([]Chirp, error)
	GetByID(id int) (Chirp, error)
	Delete(params DeleteChirpParams) error
}
