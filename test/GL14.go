package test

func GL14() (int32, uint8) {
	var thisRune int32 = 'O'
	var thisByte uint8 = 'K'
	return thisRune, thisByte
}

func GL14_Nopositive() (int32, uint8) {
	var bigInt int32 = 2147483647
	var smallerInt uint8 = 127

	return bigInt, smallerInt
}
