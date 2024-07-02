package constants

type Platform int

const (
	COSTCO Platform = iota + 10
	FOODBASICS
	NOFRILLS
	WALMART
)

func (p Platform) String() string {
	// get the string of constant
	return [...]string{"COSTCO", "FOODBASICS", "NOFRILLS", "WALMART"}[p]
}
