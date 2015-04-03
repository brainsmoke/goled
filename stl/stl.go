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

func WriteStl(out io.Writer, name string, faces []polyhedron.Face) {

	fmt.Fprintf(out, "solid %s\n", name)

	for _, f := range faces {
		a := f.Polygon[0]
		for i := 1; i<len(f.Polygon)-1; i++ {
			b, c := f.Polygon[i], f.Polygon[i+1]
			writeFacet(out, f.Normal, a, b, c)
		}
	}

	fmt.Fprintf(out, "endsolid %s\n", name)
}

