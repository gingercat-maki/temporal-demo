package model

type TeamName int

const (
	TeamReserved TeamName = iota
	TeamMarketing
	TeamSales
	TeamEngineering
)

func (s TeamName) String() string {
	switch s {
	case TeamMarketing:
		return "Marketing"
	case TeamSales:
		return "sales"
	case TeamEngineering:
		return "Engineering"
	}
	return "unknown"
}
