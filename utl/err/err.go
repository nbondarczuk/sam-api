// trivia
package err

func Ok(e error) bool {
	if e == nil {
		return true
	}

	return false
}
