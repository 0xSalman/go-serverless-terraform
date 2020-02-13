package errors

type Request struct {
	StatusCode int

	Err          error
	UserFriendly error
}

func (r Request) Error() string {
	return r.UserFriendly.Error()
}
