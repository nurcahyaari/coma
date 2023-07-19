package dto

import "fmt"

type RequestAuthValidate struct {
	Token string
}

type Method int

const (
	Apikey Method = iota
	Oauth
)

func (m Method) String() (string, error) {
	var (
		str string
		err error
	)
	switch m {
	case 0:
		str = "Apikey"
	case 1:
		str = "Oauth2"
	default:
		return str, fmt.Errorf("method is not implemented yet")
	}

	return str, err
}

type RequestValidateToken struct {
	Method Method
	Token  string
}
