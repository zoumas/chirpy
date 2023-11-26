package app

import "testing"

func TestValidateChirpLength(t *testing.T) {
	t.Run("small chirp", func(t *testing.T) {
		body := "First post!"
		err := ValidateChirpLength(body)
		assertNoError(t, err)
	})

	t.Run("empty chirp", func(t *testing.T) {
		err := ValidateChirpLength("")
		assertError(t, err, ErrChirpEmpty)
	})

	t.Run("huge chirp", func(t *testing.T) {
		body := `lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.`
		err := ValidateChirpLength(body)
		assertError(t, err, ErrChirpTooLong)
	})
}

func TestCleanChirpBody(t *testing.T) {
	token := struct{}{}
	profane := map[string]struct{}{
		"kerfuffle": token,
		"sharbert":  token,
		"fornax":    token,
	}
	replaceWith := "****"

	cases := []struct {
		Desc        string
		Body        string
		CleanedBody string
	}{
		{
			Desc:        "one profane word",
			Body:        "This is a kerfuffle opinion I need to share with the world",
			CleanedBody: "This is a **** opinion I need to share with the world",
		},
		{
			Desc:        "one profane uppercase word",
			Body:        "This is a Kerfuffle opinion I need to share with the world",
			CleanedBody: "This is a **** opinion I need to share with the world",
		},
		{
			Desc:        "two profane words with uppercase jumbled",
			Body:        "This is a KerFuFFle opinion I need to share with the world and FOrnaX",
			CleanedBody: "This is a **** opinion I need to share with the world and ****",
		},
		{
			Desc:        "words surrounded by puncuation don't count as profane",
			Body:        "Sharbert!",
			CleanedBody: "Sharbert!",
		},
		{
			Desc:        "all the rules combined",
			Body:        "Kerfuffle is an interesting word. Sharbert is misspelled. What is fornax?",
			CleanedBody: "**** is an interesting word. **** is misspelled. What is fornax?",
		},
	}

	for _, cs := range cases {
		t.Run(cs.Desc, func(t *testing.T) {
			if got := CleanChirpBody(cs.Body, profane, replaceWith); got != cs.CleanedBody {
				t.Errorf("\ngot: %q\nwant: %q", got, cs.CleanedBody)
			}
		})
	}
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
