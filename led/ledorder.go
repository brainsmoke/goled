package led

import (
	"errors"
	"strings"
)

type LedOrder int

const (
	RGB LedOrder = iota
	GRB
	RBG
	BRG
	GBR
	BGR
	RGBW
	GRBW
	RBGW
	BRGW
	GBRW
	BGRW
)

var ledString = []string{
	"RGB",
	"GRB",
	"RBG",
	"BRG",
	"GBR",
	"BGR",
	"RGBW",
	"GRBW",
	"RBGW",
	"BRGW",
	"GBRW",
	"BGRW",
}

var ledMap = [...][]int{

	{0, 1, 2},
	{1, 0, 2},
	{0, 2, 1},
	{2, 0, 1},
	{1, 2, 0},
	{2, 1, 0},
	{0, 1, 2, 3},
	{1, 0, 2, 3},
	{0, 2, 1, 3},
	{2, 0, 1, 3},
	{1, 2, 0, 3},
	{2, 1, 0, 3},
}

func (o *LedOrder) Set(value string) error {
	for i, s := range ledString {
		if strings.ToUpper(value) == s {
			*o = LedOrder(i)
			return nil
		}
	}
	return errors.New("unknown ledorder value")
}

func (o *LedOrder) String() string {
	return ledString[int(*o)]
}

func (o *LedOrder) IndexMap() []int {
	return ledMap[int(*o)]
}

func (o *LedOrder) Valid() bool {
	return RGB <= *o && *o <= BGRW
}
