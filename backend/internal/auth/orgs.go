package auth

import "fmt"

type Org int

const (
	MOHW Org = iota
	NAC
	BFLA
	CSO
)

const (
	FirstOrg = MOHW
	LastOrg  = CSO
)

// String converts an Org to a string.
func (o Org) String() string {
	return [...]string{"MOHW", "NAC", "BFLA", "CSO"}[o]
}

// ToOrg converts a string to an Org.
// If there are not matches, an error is returned
func ToOrg(s string) (Org, error) {
	switch s {
	case "MOHW":
		return MOHW, nil
	case "NAC":
		return NAC, nil
	case "BFLA":
		return BFLA, nil
	case "CSO":
		return CSO, nil
	default:
		return -1, fmt.Errorf("invalid org: %s", s) //nolint: goerr113
	}
}
