package persistent

type RestoreReq interface{}

type RestoreMsg interface{}

type RestoreComplete struct{}

func NewRestoreComplete() *RestoreComplete {
	rc := new(RestoreComplete)
	return rc
}
