package test

func GL13() bool {
	if 1 == 1 {
		return true
	}
	return false
}

func GL13_no_positive() string {
	if 1 == 1 {
		return "true"
	}
	return ""
}

func GL13_no_positive_2() bool {
	if 1 == 1 {
		return true
	}

	GL13() // some side effect
	return false
}
