package omdb

import "errors"

type tokenFactory struct {
	availableTokens []string
	nextTokenIdx int
}

func newTokenFactory(availableTokens []string) (*tokenFactory, error) {
	if len(availableTokens) == 0 {
		return nil, errors.New("cannot initialize omdb token tFactory with empty list of tokens")
	}
	return &tokenFactory{availableTokens: availableTokens, nextTokenIdx: 0}, nil
}

// round-robin use of omdb tokens
func (tf *tokenFactory) next() string {
	token := tf.availableTokens[tf.nextTokenIdx]
	noOfTokens := len(tf.availableTokens)
	tf.nextTokenIdx = (tf.nextTokenIdx + 1) % noOfTokens
	return token
}
