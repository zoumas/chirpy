package main

import "testing"

func TestValidateChirp(t *testing.T) {
	t.Run("small chirp", func(t *testing.T) {
		body := "First post!"
		err := ValidateChirpLength(body)
		assertNoError(t, err)
	})

	t.Run("empty chirp", func(t *testing.T) {
		err := ValidateChirpLength("")
		assertError(t, err, ErrEmpty)
	})

	t.Run("huge chirp", func(t *testing.T) {
		body := `lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
		err := ValidateChirpLength(body)
		assertError(t, err, ErrTooLong)
	})
}

func assertNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error : %q", err)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Fatalf("got: %s\nwant: %s", got, want)
	}
}
