// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// point
type Point struct{ x, y float64 }

// returns x value of point
func (p Point) X() float64{
	return p.x
}

//returns y value of point
func (p Point) Y() float64{
	return p.y
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.x-p.x, q.y-p.y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Round(math.Hypot(q.x-p.x, q.y-p.y)*100) / 100
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point

func onSegment(p, q, r Point) bool{
	if(q.X() <= math.Max(p.X(), r.X())) && (q.X() >= math.Min(p.X(), r.X())) && (q.Y() <= math.Max(q.Y(),r.Y())) && (q.Y() >= math.Min(p.Y(), r.Y())){
		return true
	}
	return false
}

func orientation (p, q, r Point) int{
	val := (q.Y() - p.Y()) * (r.X() - q.X()) - (q.X() - p.X()) * (r.Y() - q.Y())
	if (val == 0) {
		return 0
	} else if val > 0{
		return 1
	} 
	return 2
}

func doIntersect (p1, q1, p2, q2 Point) bool{
	o1 := orientation(p1,q1, p2)
	o2 := orientation(p1,q1, p2)
	o3 := orientation(p2,q2, p1)
	o4 := orientation(p2,q2, p1)

	if(o1 != o2) && (o3 != o4){
		return true
	} else if (o1 == 0) && onSegment(p1, p2, q1) {
		return true
	} else if (o2 == 0) && onSegment(p1, q2, q1) {
		return true
	} else if (o3 == 0) && onSegment(p2, p1, q2) {
		return true
	} else if (o4 == 0) && onSegment(p2, q1, q2) {
		return true
	}
	return false
}

func (path Path) pointsIntersect() bool{
	if len(path) <= 3 {
		return false
	}
	intersect := false
	for i := 0; i < len(path)-3; i++ {
		for j := 0; j < i+1; j++ {
			intersect = doIntersect(path[j], path[j+1], path[len(path)-2], path[len(path)-1])
		}
	}
	return intersect
}

func createPoint() float64{
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	return (random.Float64() * 200) - 100
}


// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func main() {

	if len(os.Args) <= 1{
		fmt.Println("Not enough arguments")
		return
	}
	sides, _ := strconv.Atoi(os.Args[1])
	if sides < 3 {
		fmt.Println("Must assign 3 or more points")
		return
	}

	figure := Path{}
	var perimeter = 0.0
	fmt.Printf("- Generating a [%d] sides figure\n", sides)
	fmt.Println("- Figure's vertices")
	for i := 0; i < sides; i++{
		figure = append(figure, Point{createPoint(),createPoint()})
		for figure.pointsIntersect(){
			figure[i] = Point{createPoint(), createPoint()}
		}
		fmt.Printf("\t- ( %.2f, %.2f)\n", figure[i].X(), figure[i].Y())
	}
	fmt.Println("- Figure's Perimeter")
	fmt.Printf("- ")
	for i := 0; i < sides; i++ {
		perimeter = perimeter + figure[i].Distance(figure[(i+1)%sides])
		if i == sides-1 {
			fmt.Printf("%.2f = ", figure[i].Distance(figure[(i+1)%sides]))
		} else {
			fmt.Printf("%.2f + ", figure[i].Distance(figure[(i+1)%sides]))
		}
	}
	fmt.Printf("%.2f\n", perimeter)

}

//!-path
