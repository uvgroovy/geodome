package main

import (
	"fmt"

	"github.com/ungerik/go3d/float64/vec3"
	"github.com/uvgroovy/geodome/pkg/dome"
)

func main() {

	geoDome := dome.GenerateGeoSphere(3)
	vertices, triangles := geoDome.Verteices, geoDome.Triangles

	fmt.Printf("got %d verts, %d tris, radius: %0.3f\n", len(vertices), len(triangles), geoDome.SphereRadius)

	for i, tri := range triangles {
		p1, p2, p3 := vec3.T(vertices[tri.P1]), vec3.T(vertices[tri.P2]), vec3.T(vertices[tri.P3])
		fmt.Printf("Triangle %d: %0.3f %0.3f %0.3f\n", i, p1, p2, p3)
		fmt.Printf("\tLengths: %0.3f %0.3f %0.3f\n", vec3.Distance(&p1, &p2), vec3.Distance(&p1, &p3), vec3.Distance(&p3, &p2))
		//	fmt.Printf("\tindexes: %d %d %d\n", tri.P1, tri.P2, tri.P3)
	}

}
