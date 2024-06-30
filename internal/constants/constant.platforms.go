package constants

type Platform int

const (
	COSTCO Platform = iota + 10
	NOFRILLS
	WALMART
)

func (p Platform) String() string {
	// get the string of constant
	return [...]string{"COSTCO", "NOFRILLS", "WALMART"}[p]
}
