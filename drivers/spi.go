package drivers

/*
#define _BSD_SOURCE
#include <sys/ioctl.h>

#include <linux/spi/spidev.h>

#include <unistd.h>
#include <stdint.h>

int spi_transfer(int fd, void *in, void *out, size_t len)
{
	struct spi_ioc_transfer transfer =
	{
		.tx_buf        = (unsigned long)out,
		.rx_buf        = (unsigned long)in,
		.len           = len,
		.delay_usecs   = 0,
	};

	return ioctl(fd, SPI_IOC_MESSAGE(1), &transfer);
}

int spi_set_mode(int fd, int mode)
{
	return ioctl(fd, SPI_IOC_WR_MODE, &mode);
}

int spi_set_speed(int fd, int speed)
{
	return ioctl(fd, SPI_IOC_WR_MAX_SPEED_HZ, &speed);
}

int spi_set_bits_per_word(int fd, int bits_per_word)
{
	return ioctl(fd, SPI_IOC_WR_BITS_PER_WORD, &bits_per_word);
}
*/
import "C"

import (
	"os"
	"unsafe"
)

func SPISetMode(file *os.File, mode int) (ok int, err error) {

	var cint_ok C.int
	cint_ok, err = C.spi_set_mode(C.int(file.Fd()), C.int(mode))
	ok = int(cint_ok)
	return
}

func SPISetSpeed(file *os.File, speed int) (ok int, err error) {

	var cint_ok C.int
	cint_ok, err = C.spi_set_speed(C.int(file.Fd()), C.int(speed))
	ok = int(cint_ok)
	return
}

func SPISetBitsPerWord(file *os.File, bitsPerWord int) (ok int, err error) {

	var cint_ok C.int
	cint_ok, err = C.spi_set_bits_per_word(C.int(file.Fd()), C.int(bitsPerWord))
	ok = int(cint_ok)
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func SPITransfer(file *os.File, in, out []byte) (ok int, err error) {

	var cint_ok C.int
	cint_ok, err = C.spi_transfer(C.int(file.Fd()), unsafe.Pointer(&in[0]), unsafe.Pointer(&out[0]), C.size_t(min(len(in), len(out))))
	ok = int(cint_ok)
	if err != nil {
		ok = min(len(in), len(out))
	}
	return
}

func SPIWrite(file *os.File, out []byte) (ok int, err error) {

	var cint_ok C.int
	cint_ok, err = C.spi_transfer(C.int(file.Fd()), nil, unsafe.Pointer(&out[0]), C.size_t(len(out)))
	ok = int(cint_ok)
	if err != nil {
		ok = len(out)
	}
	return
}
