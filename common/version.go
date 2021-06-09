package common

var (
	version, build, level string
)

func setVersion(v, b, l string) {
	version = v
	build = b
	level = l
}

func GetVersion() string {
	return version + " build " + build + " level " + level
}
