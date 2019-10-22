package icosidode

import (
	"post6.net/goled/model"
	"post6.net/goled/model/poly"
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
)

var ledPositions = []poly.FacePosition { // in mm, polygon point 0 defined as up, Center == (0, 0)

	{   0.000,  -4.029, true }, // 0 -> 1
	{  -2.368,   3.260, true }, // 1 -> 0
	{   2.368,   3.260, true }, // 0 -> 2
	{  -3.832,  -1.245, true }, // 2 -> 0
	{   3.832,  -1.245, true }, // 1 -> 2
	{   0.000,   0.000, true }, // 2 -> 1

}

func polyhedronePositions() [][]poly.FacePosition {

	facesList := make([][]poly.FacePosition, 6)

	for i := range facesList {
		facesList[i] = ledPositions;
	}

	return facesList
}

var ledball *model.Model3D

func rotateFace(orig polyhedron.Face, first int) polyhedron.Face {

	return polyhedron.Face {
		Normal : orig.Normal,
		Center : orig.Center,
		Polygon : append(append([]vector.Vector3{}, orig.Polygon[first:]...), orig.Polygon[:first]...),
		Neighbours : append(append([]int{}, orig.Neighbours[first:]...), orig.Neighbours[:first]...),
		Angles : append(append([]float64{}, orig.Angles[first:]...), orig.Angles[:first]...),
	}
}

func cacheLedball() {

	top := 0
	faces := polyhedron.IcosidodecahedronFaces()
	for len(faces[top].Polygon) != 5 {
		top++
	}

	factor := 8 / ( faces[top].Polygon[0].Distance(faces[top].Polygon[1]) )
	faces = polyhedron.Scale(faces, factor)

	north := faces[top].Center.Normalize()
	center := vector.Vector3{ 0, 0, 0 }
	eye := north.CrossProduct(faces[top].Polygon[0]).Normalize()

	faces = polyhedron.Rotate(faces, eye, center, north)

	newFaces := make([]polyhedron.Face, 6)
	newFaces[0] = rotateFace( faces[top], 3)

	faceId := 0
	for i:=0 ; i<5 ; i++ {

		triangleId := faces[top].Neighbours[i]
		for j:=0; j<3; j++ {
			if faces[triangleId].Neighbours[j] == top {
				faceId = faces[triangleId].Neighbours[(j+1)%3]
				break
			}
		}

		for k:=0; k<5; k++ {
			if faces[faceId].Neighbours[k] == triangleId {
				newFaces[1+i] = rotateFace( faces[faceId], (k+1)%5)
				break;
			}
		}
	}

	ledball = new(model.Model3D)
	ledball.Leds = poly.PopulateLeds(newFaces, polyhedronePositions())
}

func Ledball() *model.Model3D {

	return ledball.Copy()
}

func init() {

	cacheLedball()
}

