package drivers

import "os"

type Teensy16 struct {
	file *os.File
	buf []byte
}

func NewTeensy16(filename string) (t *Teensy16, err error) {

	t = new(Teensy16)
	t.file, err = os.OpenFile(filename, os.O_RDWR, 0)
	SetBaudrate(t.file, 12000000)
	SetBinary(t.file)

	if err != nil {
		t.file.Close()
		t = nil
	}

	return t, err
}

func (t *Teensy16) Write(p []byte) (n int, err error) {
	needed := (len(p)+5) &^ 1
	if len(t.buf) < needed {
		t.buf = make([]byte, needed)

	}
	copy(t.buf, p)
	t.buf[needed-4] = 0xff
	t.buf[needed-3] = 0xff
	t.buf[needed-2] = 0xff
	t.buf[needed-1] = 0xf0
	return t.file.Write(t.buf)
}

func (t *Teensy16) Close() (err error) {
	return t.file.Close()
}

func (t *Teensy16) Bpp() int {
	return 16
}

func (t *Teensy16) MaxValue() uint {
	return 0xff00
}
