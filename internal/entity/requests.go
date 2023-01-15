package entity

type Request struct {
	ID       UUID   `json:"id"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

var requestStatusMap = map[string]struct{}{
	"new":       {},
	"processed": {},
}

type Filter struct {
	ID     UUID
	Status string
}

func (r Request) Validate() error {
	_, ok := requestStatusMap[r.Status]
	if ok {
		return nil
	}
	return RequestStatusErr
}
