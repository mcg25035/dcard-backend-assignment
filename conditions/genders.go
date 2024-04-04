package conditions

type Gender string

const (
	M Gender = "M"
	F Gender = "F"
)

func AllGenders() []Gender {
	return []Gender{M, F}
}