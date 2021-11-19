package xpression

type tokenStack struct {
	data []*Token
}

func (s *tokenStack) push(tok *Token) *Token {
	if tok != nil {
		s.data = append(s.data, tok)
	}
	return tok
}

func (s *tokenStack) pop() *Token {
	l := len(s.data)
	if l == 0 {
		return nil
	}
	value := s.data[l-1]
	s.data = s.data[:l-1]
	return value
}

func (s *tokenStack) pushDouble(empty *Token, tok *Token) *Token {
	if tok != nil {
		s.data = append(s.data, empty)
		s.data = append(s.data, tok)
	}
	return tok
}

func (s *tokenStack) popDouble() (*Token, *Token) {
	l := len(s.data)
	if l == 0 {
		return nil, nil
	}
	value := s.data[l-1]
	s.data = s.data[:l-1]
	return &Token{}, value
}

func (s *tokenStack) peek() *Token {
	l := len(s.data)
	if l == 0 {
		return nil
	}
	return s.data[l-1]
}

func (s *tokenStack) get() []*Token {
	return s.data
}

func reverse(s []*Token) []*Token {
	l := len(s)
	for i := 0; i < l/2; i++ {
		s[i], s[l-i-1] = s[l-i-1], s[i]
	}
	return s
}
