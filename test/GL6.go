package test

func GL6(number int) int {
	if number == 1 {
		return 10
	} else if number == 2 {
		return 20
	} else if number == 3 {
		return 30
	}
	return 0
}

// This function is a false negative in staticcheck, but is supplemented by don't use Yoda conditions
func GL6Alt(number int) int {
	if number == 1 {
		return 10
	} else if number == 2 {
		return 20
	} else if 3 == number {
		return 30
	}
	return 0
}
