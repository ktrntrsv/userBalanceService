package transaction

type Transaction struct {
	Amount        float64
	AccountToId   string
	AccountFromId string
	Status        Status
	ServiceId     *string
}

type Status int

const (
	Pending Status = iota
	Approved
	Canceled
)
