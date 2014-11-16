package drivers

import "os"

type Ws2801Spi struct {
	file *os.File
}

func NewWs2801Spi(filename string, speed int) (spi *Ws2801Spi, err error) {

	spi = new(Ws2801Spi)
	spi.file, err = os.OpenFile(filename, os.O_RDWR, 0)

	if err == nil {
		_, err = SPISetMode(spi.file, 0)
	}
	if err == nil {
		_, err = SPISetSpeed(spi.file, speed)
	}
	if err == nil {
		_, err = SPISetBitsPerWord(spi.file, 8)
	}
	if err != nil {
		spi.file.Close()
		spi = nil
	}

	return
}

func (spi *Ws2801Spi) Write(p []byte) (n int, err error) {

	return SPIWrite(spi.file, p)
}

func (spi *Ws2801Spi) Close() (err error) {
	return spi.file.Close()
}
