package util

func TrigTable(size int, min, max float64, f func(float64) float64) []float64 {

	table := make([]float64, size)

	for i := range table {
		table[i] = f(min + (max-min)*float64(i)/float64(size))
	}

	return table
}
