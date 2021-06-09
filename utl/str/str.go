// trivia
package str

func Nvl(str, option string) string {
	if str == "" {
		return option
	}

	return str
}

func Empty(str string) bool {
	if str == "" {
		return true
	}

	return false
}
