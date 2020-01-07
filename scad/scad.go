package scad

import (
	"post6.net/goled/polyhedron"
//	"post6.net/goled/vector"
	"io"
	"fmt"
)

/* 'OpenGL' -> OpenSCAD coordinates x=x, y=z, z=-y */
func WriteOpenScad(out io.Writer, solid polyhedron.Solid) {

	fmt.Fprintf(out, "polyhedron( points = [ ")
	for i, p := range solid.Points {
		fmt.Fprintf(out, "[%e, %e, %e]", p.X, p.Z, -p.Y);
		if i != len(solid.Points)-1 {
			fmt.Fprintf(out, ", ");
		}
	}
	fmt.Fprintf(out, " ], faces = [ ")

	for i, f := range solid.Faces {
		fmt.Fprintf(out, "[ ")
		for j:=len(f.Polygon)-1 ; j >=0 ; j-- {
			if j == 0 {
				fmt.Fprintf(out, "%d ] ", f.Polygon[j]);
			} else {
				fmt.Fprintf(out, "%d, ", f.Polygon[j]);
			}
		}
		if i != len(solid.Faces)-1 {
			fmt.Fprintf(out, ", ");
		}
	}

	fmt.Fprintf(out, " ], convexity = 1);\n")
}

