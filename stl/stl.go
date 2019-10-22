package stl

import (
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
	"io"
	"fmt"
)

func writeFacet(out io.Writer, normal, a, b, c vector.Vector3) {

	fmt.Fprintf(out,
"facet normal %e %e %e\n"+
"	outer loop\n"+
"		vertex %e %e %e\n"+
"		vertex %e %e %e\n"+
"		vertex %e %e %e\n"+
"	endloop\n"+
"endfacet\n", normal.X, normal.Y, normal.Z, a.X, a.Y, a.Z, b.X, b.Y, b.Z, c.X, c.Y, c.Z)

}

func WriteStl(out io.Writer, name string, solid polyhedron.Solid) {

	fmt.Fprintf(out, "solid %s\n", name)

	for _, f := range solid.Faces {
		a := solid.Points[f.Polygon[0]]
		for i := 1; i<len(f.Polygon)-1; i++ {
			b, c := solid.Points[f.Polygon[i]], solid.Points[f.Polygon[i+1]]
			writeFacet(out, f.Normal, a, b, c)
		}
	}

	fmt.Fprintf(out, "endsolid %s\n", name)
}

