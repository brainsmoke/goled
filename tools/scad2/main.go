package main

import (
	"post6.net/goled/polyhedron"
	"post6.net/goled/scad"
	"post6.net/goled/vector"
	"os"
//	"fmt"
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

    in.Scale( 9 / solid.Faces[0].Center.Magnitude() )
    out.Scale( 10.3 / solid.Faces[0].Center.Magnitude() )

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

	scad.WriteOpenScad(os.Stdout, eggshell)

}
