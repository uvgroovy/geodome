package dome

import (
	"log"
	"math"
	"sort"

	"github.com/ungerik/go3d/float64/vec3"
)

// 	math.Phi is golden ratio
const ϕ = math.Phi

type Vertex vec3.T

/*
type Vertex struct {
	X, Y, Z float64
}
*/

// index to vertex array
type Triangle struct {
	P1, P2, P3 int
}

func NewTriangle(i, j, k int) Triangle {
	var p sort.IntSlice = make([]int, 3)
	p[0] = i
	p[1] = j
	p[2] = k
	sort.Sort(p)

	return Triangle{p[0], p[1], p[2]}
}

/*
func (v Vertex) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}


func (v *Vertex) Scale(factor float64) {
	v.X *= factor
	v.Y *= factor
	v.Z *= factor
}
*/

type GeoDome struct {
	Verteices    []Vertex
	Triangles    []Triangle
	SphereRadius float64
}

func GenerateGeoSphere(n int) GeoDome {
	v, triangles := GenerateIcosahedronVertexes()

	oneV := (vec3.T)(v[0])
	r := oneV.Length()
	// triangles := GenerateIcosahedronTriangles()

	// the isocahendro is on the sphere  with r = ||(0,1,ϕ)||
	if n > 1 {
		s :=
			SubDevider{
				Vertices:  v,
				Triangles: triangles,
				N:         n,
				R:         r,
				cache:     make(map[lineOfTriangle][]int),
			}

		s.SubDevide()
		v, triangles = s.Vertices, s.TrianglesOutput
	}
	return GeoDome{Verteices: v, Triangles: triangles, SphereRadius: r}

}

func GenerateIcosahedronVertexes() ([]Vertex, []Triangle) {
	// generate Icosahedron around the unit sphere
	var ret []Vertex
	var triret []Triangle

	// Recangles (according to wikipedia cyclic permutations of (0, ±1, ±ϕ))
	// initial vertex:
	rectAngle := []Vertex{
		Vertex{0, 1, ϕ},
		Vertex{0, -1, ϕ},
		Vertex{0, 1, -ϕ},
		Vertex{0, -1, -ϕ},
	}
	for i := 0; i < 3; i++ {
		for _, v := range rectAngle {
			ret = append(ret, v)
		}
		for vind := range rectAngle {
			rectAngle[vind] = cycle(rectAngle[vind])
		}
	}

	//	var tri []triangle

	set := make(map[Triangle]bool)
	for ind, vertex := range ret {

		neighbors := getNeighborsInOrder(vertex)
		for i := range neighbors {
			tri := NewTriangle(ind, findIndex(ret, neighbors[i]), findIndex(ret, neighbors[(i+1)%len(neighbors)]))
			set[tri] = true
		}
	}

	for tri := range set {
		triret = append(triret, tri)
	}

	return ret, triret
}

func findIndex(vertices []Vertex, vertex Vertex) int {

	for i, v := range vertices {
		if v == vertex {
			return i
		}
	}

	return -1
}

func getNeighborsInOrder(vertex Vertex) [5]Vertex {
	var res [5]Vertex
	var oneIndex, phiIndex, zeroIndex int
	indexesSet := 0
	for index, val := range vertex {
		switch val {
		case 0:
			zeroIndex = index
			indexesSet |= 1
		case 1:
			fallthrough
		case -1:
			oneIndex = index
			indexesSet |= 2
		default:
			phiIndex = index
			indexesSet |= 4
		}
	}

	if indexesSet != 7 {
		panic("bad vertex!")
	}

	// vertex is center point
	// create all five point of distance 2 and find the verteices in the array

	// first triangle:
	p1 := Vertex{}
	p1[zeroIndex] = ϕ
	p1[oneIndex] = 0
	p1[phiIndex] = math.Copysign(1, vertex[phiIndex])

	p2 := Vertex{}
	p2[zeroIndex] = -ϕ
	p2[oneIndex] = 0
	p2[phiIndex] = math.Copysign(1, vertex[phiIndex])

	p3 := Vertex{}
	p3[zeroIndex] = 1
	p3[oneIndex] = math.Copysign(ϕ, vertex[oneIndex])
	p3[phiIndex] = 0

	p4 := Vertex{}
	p4[zeroIndex] = -1
	p4[oneIndex] = math.Copysign(ϕ, vertex[oneIndex])
	p4[phiIndex] = 0

	p5 := Vertex{}
	p5[zeroIndex] = vertex[zeroIndex]
	p5[oneIndex] = -vertex[oneIndex]
	p5[phiIndex] = vertex[phiIndex]

	if p5[oneIndex] < 0 {
		res[0] = p4
		res[1] = p2
		res[2] = p5
		res[3] = p1
		res[4] = p3
	} else {
		res[0] = p3
		res[1] = p1
		res[2] = p5
		res[3] = p2
		res[4] = p4
	}

	verifyDistance := func(p1, p2 Vertex) {
		p1v, p2v := vec3.T(p1), vec3.T(p2)
		d := vec3.Distance(&p1v, &p2v)
		if d != 2.0 {
			log.Panicf("Bad distance %0.3f", d)
		}

	}

	for i := range res {

		verifyDistance(vertex, res[i])
		verifyDistance(res[i], res[(i+1)%len(res)])
	}

	return res
}

func cycle(v Vertex) Vertex {
	return Vertex{v[2], v[0], v[1]}
}

type SubDevider struct {
	Vertices  []Vertex
	Triangles []Triangle
	N         int
	R         float64

	TrianglesOutput []Triangle

	cache map[lineOfTriangle][]int
}

func (s *SubDevider) SubDevide() {
	for _, triangle := range s.Triangles {
		newTriangles := s.getDevisionVertexes(triangle)
		s.TrianglesOutput = append(s.TrianglesOutput, newTriangles...)
		// delete	s.TrianglesOutput = append(s.TrianglesOutput, triangle)
		// delete	return
	}
}

func (s *SubDevider) getDevisionVertexes(tri Triangle) []Triangle {
	var res []Triangle
	// for each line in the tri angle, devide to N
	// p1p2 and p1p3

	divisions1 := s.getSubDevisionsForLine(tri.P1, tri.P2)
	divisions2 := s.getSubDevisionsForLine(tri.P1, tri.P3)

	divisions1 = append(divisions1, tri.P2)
	divisions2 = append(divisions2, tri.P3)

	// add inside verteices and triangles
	previousLine := []int{tri.P1}
	for i := 0; i < s.N; i++ {
		subdivisions := s.getSubDevisionsForLineN(divisions1[i], divisions2[i], i+1)
		curLine := []int{divisions1[i]}
		curLine = append(curLine, subdivisions...)
		curLine = append(curLine, divisions2[i])

		// combine and create tri angels from cur and prev line
		for triind := range previousLine {
			newtriangle := NewTriangle(previousLine[triind], curLine[triind], curLine[triind+1])
			res = append(res, newtriangle)
		}
		for triind := range previousLine[:len(previousLine)-1] {
			newtriangle := NewTriangle(previousLine[triind], previousLine[triind+1], curLine[triind+1])
			res = append(res, newtriangle)
		}

		previousLine = curLine
	}

	return res
}

func (s *SubDevider) getSubDevisionsForLine(i, j int) []int {
	return s.getSubDevisionsForLineN(i, j, s.N)
}

func (s *SubDevider) getSubDevisionsForLineN(i, j int, n int) []int {

	if n == 1 {
		return []int{}
	}

	reverse := false
	if i > j {
		i, j = j, i
		reverse = true
	}

	lofT := lineOfTriangle{i, j}
	if _, ok := s.cache[lofT]; !ok {
		// add to verteices
		p1 := (vec3.T)(s.Vertices[i])
		p2 := (vec3.T)(s.Vertices[j])
		var cache []int
		for index := 1; index < n; index++ {
			// create new vertex
			//append it
			newVertex := (Vertex)(vec3.Interpolate(&p1, &p2, float64(index)/float64(n)))
			newVertex = s.project(newVertex)
			s.Vertices = append(s.Vertices, newVertex)
			cache = append(cache, len(s.Vertices)-1)
		}

		s.cache[lofT] = cache
	}

	if reverse {
		return reverseArray(s.cache[lofT])
	}
	return s.cache[lofT]

}

func (s *SubDevider) project(input Vertex) Vertex {
	// project vertex to unit spehre
	inp := (vec3.T)(input)

	l := inp.Length()
	res := (Vertex)(*inp.Scale(s.R / l))
	return res
}

func reverseArray(input []int) []int {
	res := make([]int, len(input))
	for i := 0; i < len(input); i++ {
		res[len(input)-1-i] = input[i]
	}
	return res
}

type lineOfTriangle struct {
	firstIndex, secondIndex int
}
