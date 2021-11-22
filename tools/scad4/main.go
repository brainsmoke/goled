package main

import (
	"post6.net/goled/polyhedron"
	"post6.net/goled/vector"
	"os"
	"fmt"
	"flag"
)

var name string

func init() {
	flag.StringVar(&name, "model", "oD", "conway name")
}

func main() {

	flag.Parse()

	model, ok := polyhedron.ConwayModels[name]

	if !ok {
		panic("meh.")
	}

	solid := model()

    north := solid.Faces[0].Normal
    center := vector.Vector3{ 0, 0, 0 }
    eye := north.CrossProduct(solid.Points[solid.Faces[0].Polygon[0]]).Normalize()

    solid.Rotate(eye, center, north)

fmt.Fprintf(os.Stdout, "points = [\n")
	for i, ix := range solid.Faces[0].Polygon {
		p := solid.Points[ix]
		fmt.Fprintf(os.Stdout, "[%e,%e]", p.Z, p.X)
		if i < len(solid.Faces[0].Polygon)-1 {
			fmt.Fprintf(os.Stdout, ",")
		}
	}
	fmt.Fprintf(os.Stdout, "];\n")

fmt.Fprintf(os.Stdout, "vertex_holes = [\n\t");

for i := range solid.Faces[0].Polygon {
	fmt.Fprintf(os.Stdout, ".0", )
	if i < len(solid.Faces[0].Polygon)-1 {
		fmt.Fprintf(os.Stdout, ",")
	}
}

fmt.Fprintf(os.Stdout, "];\n");


fmt.Fprintf(os.Stdout, "edge_holes = [\n\t");

for i := range solid.Faces[0].Polygon {
	fmt.Fprintf(os.Stdout, ".0", )
	if i < len(solid.Faces[0].Polygon)-1 {
		fmt.Fprintf(os.Stdout, ",")
	}
}

fmt.Fprintf(os.Stdout, "];\n");

fmt.Fprintf(os.Stdout, `
midpoints = [ for (i=[0:len(points)-1]) (points[i]+points[(i+1)%%len(points)])/2 ];

module shape()
{

	difference()
	{
		polygon( points );
		union()
		{
			for (i=[0:len(points)-1])
				translate(points[i])
					circle(r=vertex_holes[i], $fn=50);
			for (i=[0:len(midpoints)-1])
				translate(midpoints[i])
					circle(r=edge_holes[i], $fn=50);
		}
	}

}

module facet(s,t)
{
	translate([0,0,s])
		linear_extrude(height=t)
			scale([s,s])
				shape();
}`);


fmt.Fprintf(os.Stdout, "module ball(r,t) {\n")

	for _,f := range(solid.Faces) {

		top := solid.Points[f.Polygon[0]]

		newZ := f.Normal
		vY := top.Sub(f.Center).Normalize()
		newX := vY.CrossProduct(newZ).Normalize()
		newY := newZ.CrossProduct(newX).Normalize()

		fmt.Fprintf(os.Stdout, "multmatrix (m=[[%e,%e,%e,%e],[%e,%e,%e,%e],[%e,%e,%e,%e],[0, 0, 0, 1]]) {facet(r,t);};\n",
			newX.X, newY.X, newZ.X, 0.,
			newX.Z, newY.Z, newZ.Z, 0.,
			-newX.Y, -newY.Y, -newZ.Y, 0.,
			)
	}
fmt.Fprintf(os.Stdout, "}\n")


fmt.Fprintf(os.Stdout, "ball(160, .8);\n")
}
