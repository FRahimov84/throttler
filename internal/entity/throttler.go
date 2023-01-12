package entity

type Request struct {
	ID       UUID   `json:"id"`
	Status   string `json:"status"`
	Response string `json:"response"`
}

type Filter struct {
	ID     UUID
	Status string
}
