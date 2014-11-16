package drivers

import (
	"flag"
	"io"
	"os"
)

var spiDevice string
var spiSpeed int
var useStdout bool

func init() {

	flag.StringVar(&spiDevice, "spidev", "/dev/spidev0.0", "used SPI device")
	flag.IntVar(&spiSpeed, "spispeed", 2000000, "used speed for SPI transfer")
	flag.BoolVar(&useStdout, "stdout", false, "use stdout instead of spi")

}

func LedDriver() io.Writer {

	if useStdout {
		return os.Stdout
	} else {
		out, err := NewWs2801Spi(spiDevice, spiSpeed)

		if err != nil {
			panic("could not open SPI device")
		}

		return out
	}
}
