package ani

type Animation interface {
	Next() [][3]byte
}
