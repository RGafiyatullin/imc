package request

type Request interface{}
type request struct{}

func New() Request {
	req := new(request)
	return req
}

func FromBytes(bytes []byte) Request {
	return nil
}
