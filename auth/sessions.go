package auth

type SessionData struct {
	ActiveCookies []string
	AuthState     map[interface{}]uint
}

func (s *SessionData) Init() *SessionData {
	s.AuthState = make(map[interface{}]uint)
	return s
}
