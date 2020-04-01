package util

func EqualAny(str string, any []string) bool {
	for _, s := range any {
		if s == str {
			return true
		}
	}
	return false
}

func HasAny(str string, any []string) (string, bool) {
	for _, s := range any {
		if s == str {
			return s, true
		}
	}
	return "", false
}
