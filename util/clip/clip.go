package clip

func FloatToByte(f float64) byte {

	if f > 255. {
		return 255
	} else if f < 1 {
		return 0
	} else {
		return byte(f)
	}
}

func IntToByte(i int) byte {

	if i > 255 {
		return 255
	} else if i < 1 {
		return 0
	} else {
		return byte(i)
	}
}

func AddBytes(a, b byte) byte {
	if int(a)+int(b) > 255 {
		return 255
	} else {
		return a + b
	}
}
