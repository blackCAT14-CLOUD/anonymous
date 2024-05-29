package commons

type errorCodes struct {
	InternalError string
	EmptyField    string
	InvalidField  string

	EmailNotFound string
	UsernameNotFound string

	WrongPassword  string
    DuplicateUsername string

	DuplicateEmail string
  DuplicateLabel string
  InvalidToken string
	
}

var Codes = errorCodes{
	InternalError:  "InternalError",
	EmptyField:     "EmptyField",
	InvalidField:   "InvalidField",
	EmailNotFound:  "EmailNotFound",
	WrongPassword:  "WrongPassword",
	DuplicateEmail: "DuplicateEmail",
    DuplicateUsername : "DuplicateUsername",
    UsernameNotFound: "UsernameNotFound",
    InvalidToken: "InvalidToken",
}
