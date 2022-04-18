package schema

// StatusText Define status text
type StatusText string

func (t StatusText) String() string {
	return string(t)
}

// NextServer
type NextServer struct {
	Host string
	Port string
}
