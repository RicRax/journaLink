package main

type sessionData struct {
	activeCookies []string
	authState     map[interface{}]uint
}

func (s *sessionData) init() *sessionData {
	s.authState = make(map[interface{}]uint)
	return s
}
