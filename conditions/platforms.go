package conditions

type Platform string

const (
	android Platform = "android"
	ios     Platform = "ios"
	web     Platform = "web"
)

func AllPlatforms() []Platform {
	return []Platform{android, ios, web}
}