package main

import (
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
	"os"
	"io"
	"fmt"
)

const (
	TopLeft     = 0
	BottomLeft  = 1
	BottomRight = 2
	TopRight    = 3
)



const innerRadius = 1

var traversal = polyhedron.RemapRoute{

TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight,
}

/* 'OpenGL' -> OpenSCAD coordinates x=x, y=z, z=-y */
func WriteOpenScadFilthyHack(out io.Writer, solid polyhedron.Solid) {


	fmt.Fprintf(out, "module eggshell(inner, outer) {\n")
	fmt.Fprintf(out, "polyhedron( points = [ \n")

varname := "inner";
half_points := len(solid.Points)/2;

	for i, p := range solid.Points {
if i == half_points {
varname = "outer";
}
		fmt.Fprintf(out, "[%e * %s, %e * %s, %e * %s]", p.X, varname, p.Z, varname, -p.Y, varname);
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

	fmt.Fprintf(out, " ], convexity = 4);\n")
	fmt.Fprintf(out, "}\n")
}


func reverseIntSlice(a []int) {
	for i,j := 0, len(a)-1 ; i < j ; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func insideOut(face *polyhedron.Face) {

	face.Normal = face.Normal.Mul(-1)
	reverseIntSlice(face.Polygon[1:])

}

func main() {

	solid := polyhedron.RemapSolid(polyhedron.DeltoidalHexecontahedron(), 0, traversal)

	in := solid.Copy()
	out := solid.Copy()

    in.Scale( solid.Faces[0].Center.Magnitude() )
    out.Scale( solid.Faces[0].Center.Magnitude() )

	eggshell := polyhedron.Combine(in, out)

	n := len(in.Points)
	for i,f := range(in.Faces) {
		for j,edge := range(f.Neighbours) {
			sz := len(f.Polygon)
			if edge == -1 {
				var newF polyhedron.Face
				newF.Polygon = []int { f.Polygon[j], f.Polygon[(j+1)%sz], f.Polygon[(j+1)%sz]+n, f.Polygon[j]+n }
				eggshell.Faces = append(eggshell.Faces, newF )
			}
		}
		insideOut(&eggshell.Faces[i])
	}

    north := eggshell.Points[eggshell.Faces[1].Polygon[2]].Normalize()
    center := vector.Vector3{ 0, 0, 0 }
    eye := north.CrossProduct(eggshell.Points[eggshell.Faces[1].Polygon[3]]).Normalize()

    eggshell.Rotate(eye, center, north)

	WriteOpenScadFilthyHack(os.Stdout, eggshell)

}
