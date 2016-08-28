package drivers

import "os"

type FileDriver struct {
	file *os.File
	bpp int
	maxValue uint
}

func NewFileDriver(file *os.File, bpp int, maxValue uint) *FileDriver {

	return &FileDriver{ file: file, bpp: bpp, maxValue: maxValue }
}

func (d *FileDriver) Write(p []byte) (n int, err error) {
	return d.file.Write(p)
}

func (d *FileDriver) Bpp() int {
	return d.bpp
}

func (d *FileDriver) MaxValue() uint {
	return d.maxValue
}
