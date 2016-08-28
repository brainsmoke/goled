package drivers

import (
	"flag"
	"io"
	"os"
	"strings"
)

type LedDriver interface {
	io.Writer
	Bpp() int
	MaxValue() uint
}

var driver string
var serialDevice string
var spiDevice string
var spiSpeed int

var driverList []string = []string{
	"teensy16",
	"stdout",
	"stdout16",
	"spi",
}

func init() {

	flag.StringVar(&driver, "driver", "teensy16", "used driver ["+strings.Join(driverList, ", ")+"]")
	flag.StringVar(&serialDevice, "serialdev", "/dev/ttyACM0", "used serial device (for teensy16)")
	flag.StringVar(&spiDevice, "spidev", "/dev/spidev0.0", "used SPI device")
	flag.IntVar(&spiSpeed, "spispeed", 2000000, "used speed for SPI transfer")
}

func GetLedDriver() LedDriver {

	switch driver {
		case "teensy16":
			d, err := NewTeensy16(serialDevice)
			if err != nil {
				panic("could not open serial device")
			}
			return d

		case "stdout":
			return NewFileDriver(os.Stdout, 8, 0xff)

		case "stdout16":
			return NewFileDriver(os.Stdout, 16, 0xff00)

		case "spi":
			spidev, err := NewWs2801Spi(spiDevice, spiSpeed)

			if err != nil {
				panic("could not open SPI device")
			}

			return spidev

		default:
			panic("unknown driver")
	}
}

