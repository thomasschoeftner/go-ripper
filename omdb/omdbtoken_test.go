package omdb

import "testing"

func TestRoundRobinTokenUsage(t *testing.T) {
	tokens := []string {"sepp", "hat", "gelbe", "eier"}
	tf, err := newTokenFactory(tokens)
	if err != nil {
		t.Errorf("omdb token tFactory failed unexpectedly due to %v", err)
	}
	validateToken(t, tokens[0], tf.next())
	validateToken(t, tokens[1], tf.next())
	validateToken(t, tokens[2], tf.next())
	validateToken(t, tokens[3], tf.next())
	validateToken(t, tokens[0], tf.next())
}

func validateToken(t *testing.T, expected string, got string) {
	if expected != got {
		t.Errorf("expected token \"%s\", but got \"%s\"", expected, got)
	}
}

func TestEmptyTokens(t *testing.T) {
	tf, err := newTokenFactory([]string{})
	if err == nil || tf != nil {
		t.Errorf("expected omdb token tFactory to fail - did not happen")
	}
}
