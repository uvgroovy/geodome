package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"

	"github.com/sbinet/go-gnuplot"
	"github.com/ungerik/go3d/float64/vec3"
	"github.com/uvgroovy/geodome/pkg/dome"
)

func getDataFile(vertices []dome.Vertex, triangles []dome.Triangle) string {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	defer tmpFile.Close()

	for _, tri := range triangles {
		p1, p2, p3 := vec3.T(vertices[tri.P1]), vec3.T(vertices[tri.P2]), vec3.T(vertices[tri.P3])
		/*
			if (p1[2] < 0) || (p2[2] < 0) || (p3[2] < 0) {
				continue
			}
		*/

		fmt.Fprintf(tmpFile, "%0f %0f %0f 0x000000 \n", p1[0], p1[1], p1[2])
		fmt.Fprintf(tmpFile, "%0f %0f %0f %s \n", p2[0], p2[1], p2[2], getColor(p1, p2))
		fmt.Fprintf(tmpFile, "%0f %0f %0f %s \n", p3[0], p3[1], p3[2], getColor(p2, p3))
		fmt.Fprintf(tmpFile, "%0f %0f %0f %s \n", p1[0], p1[1], p1[2], getColor(p3, p1))
		fmt.Fprintf(tmpFile, "\n\n")

	}

	return tmpFile.Name()
}

var colors map[int]string = make(map[int]string)
var numbers map[int]int = make(map[int]int)
var lengths map[int]float64 = make(map[int]float64)
var angles map[int]float64 = make(map[int]float64)
var colorToSet int = 0

var colorPallete []string = []string{"0xff0000", "0x00ff00", "0x0000ff"}

func getNewColor() string {
	if colorToSet == len(colorPallete) {
		return fmt.Sprintf("0x%02x%02x%02x", rand.Intn(0x100), rand.Intn(0x100), rand.Intn(0x100))
	}
	colorToSet += 1
	return colorPallete[colorToSet-1]
}

func getColor(p1, p2 vec3.T) string {

	l := vec3.Distance(&p1, &p2)
	key := int(1000 * l)
	numbers[key] = numbers[key] + 1
	if _, ok := colors[key]; !ok {
		colors[key] = getNewColor()
	}
	if _, ok := angles[key]; !ok {
		lengths[key] = l
		tmp := vec3.Sub(&p1, &p2)
		angles[key] = math.Pi/2 - vec3.Angle(&p1, &tmp)
	}
	return colors[key]

}

func main() {
	geoDome := dome.GenerateGeoSphere(3)
	vertices, triangles := geoDome.Verteices, geoDome.Triangles

	fmt.Printf("got %d verts, %d tris, radius: %0.3f\n", len(vertices), len(triangles), geoDome.SphereRadius)

	dataFile := getDataFile(vertices, triangles)
	defer os.Remove(dataFile)

	fmt.Printf("got %d colors\n", len(colors))
	for k, v := range numbers {
		fmt.Printf("%d of cord factor %0.3f angle: %0.3f°\n", v, lengths[k]/geoDome.SphereRadius, 180.0*angles[k]/math.Pi)

	}

	fname := ""
	persist := false
	debug := true

	p, err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p.CheckedCmd("set terminal qt")
	p.CheckedCmd("set view equal xyz")

	p.CheckedCmd("splot '%s' using 1:2:3:4 with lines linecolor rgb variable", dataFile)
	bufio.NewReader(os.Stdin).ReadString('\n')
	p.CheckedCmd("q")
	return
}
