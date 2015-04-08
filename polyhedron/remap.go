package polyhedron

type RemapRoute []int
type RemapReorientRoute [][2]int

func RemapFaces(faces []Face, first int, route RemapRoute) []Face {

	newFaces := make([]Face, len(faces))
	mapping := make([]int, len(faces))

	current := first

	for i := range newFaces {

		newFaces[i] = faces[current]
		mapping[current] = i
		current = faces[current].Neighbours[route[i]]
	}

	for i := range newFaces {
		for j, k := range newFaces[i].Neighbours {

			newFaces[i].Neighbours[j] = mapping[k]
		}
	}

	return newFaces
}

