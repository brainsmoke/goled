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
/*
var traversal = polyhedron.RemapRoute{

TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight,
TopRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight, TopRight, TopRight, BottomRight, BottomRight, BottomRight, TopRight, BottomRight, BottomRight, BottomRight, BottomRight,
}
*/
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

	solid := polyhedron.DeltoidalHexecontahedron()
	//solid := polyhedron.RemapSolid(polyhedron.DeltoidalHexecontahedron(), 0, traversal)

    solid.Scale( 1 / solid.Faces[0].Center.Magnitude() )
    north := solid.Points[solid.Faces[1].Polygon[2]].Normalize()
    center := vector.Vector3{ 0, 0, 0 }
    eye := north.CrossProduct(solid.Points[solid.Faces[1].Polygon[3]]).Normalize()

    solid.Rotate(eye, center, north)

	scad.WriteOpenScad(os.Stdout, solid)

}
