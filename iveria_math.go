package main

import (
	"math"
	"math/rand"
)

// Vec2 represents a vector with an x and y integer component.
type Vec2 struct {
	x int
	y int
}

var Vec2Zero = Vec2{0, 0}
var Vec2Up = Vec2{0, 1}
var Vec2Down = Vec2{0, -1}
var Vec2Left = Vec2{-1, 0}
var Vec2Right = Vec2{1, 0}

// Vec2Math contains mathematical functions for vectors
type Vec2Math interface {
	normalise() Vec2
	magnitude() float64
	add(Vec2) Vec2
	invert() Vec2
	random() Vec2
}

// normalise the components of a vector by the magnitude.
func (v Vec2) normalise() Vec2 {
	mag := v.magnitude()
	return Vec2{
		(int)(math.Round((float64)(v.x) / mag)),
		(int)(math.Round((float64)(v.y) / mag))}
}

// magnitude calculates the length of a vector.
func (v Vec2) magnitude() float64 {
	return math.Sqrt((float64)((v.x * v.x) + (v.y * v.y)))
}

// add returns the sum of two vectors.
func (v Vec2) add(v2 Vec2) Vec2 {
	return Vec2{v.x + v2.x, v.y + v2.y}
}

// invert multiples the vectors components by negative one.
func (v Vec2) invert() Vec2 {
	return Vec2{-v.x, -v.y}
}

func (v Vec2) random() Vec2 {
	var index = rand.Intn(3)
	switch index {
	case 0:
		return Vec2Up
	case 1:
		return Vec2Down
	case 2:
		return Vec2Left
	case 3:
		return Vec2Right
	}
	return Vec2Zero
}

// Rect represents a rectangle consisting of a vector for position and size.
type Rect struct {
	pos  Vec2
	size Vec2
}

// RectMath contains mathematical functions for rectangles
type RectMath interface {
	center() Vec2
	topLeft() Vec2
	bottomRight() Vec2
	topRight() Vec2
	bottomLeft() Vec2
	expand(int) Rect
	contains(Vec2) bool
	intersects(Rect) bool
}

// center returns the center of the rectangle.
func (r Rect) center() Vec2 {
	return Vec2{r.pos.x + (r.size.x / 2), r.pos.y + (r.size.y / 2)}
}

// topLeft returns the top left corner of the rectangle.
func (r Rect) topLeft() Vec2 {
	return r.pos
}

// bottomRight returns the bottom right corner of the rectangle.
func (r Rect) bottomRight() Vec2 {
	return Vec2{r.topLeft().x + r.size.x, r.topLeft().y + r.size.y}
}

// topRight returns the top right corner of the rectangle.
func (r Rect) topRight() Vec2 {
	return Vec2{r.bottomRight().x, r.topLeft().y}
}

// bottomLeft returns the bottom left corner of the rectangle.
func (r Rect) bottomLeft() Vec2 {
	return Vec2{r.topLeft().x, r.bottomRight().y}
}

// expand creates an expanded version of the rectangle.
func (r Rect) expand(s int) Rect {
	return Rect{
		pos:  Vec2{r.pos.x - s, r.pos.y - s},
		size: Vec2{r.size.x + (s * 2), r.size.x + (s * 2)},
	}
}

// contains returns true if the point is within the rectangle
func (r Rect) contains(point Vec2) bool {
	if point.x > r.topLeft().x && point.x < r.bottomRight().x &&
		point.y > r.topLeft().y && point.y < r.bottomRight().y {
		return true
	}
	return false
}

// intersects returns true if the rectangle given intersects
func (r Rect) intersects(r2 Rect) bool {
	return r.contains(r2.topLeft()) || r.contains(r2.topRight()) ||
		r.contains(r2.bottomLeft()) || r.contains(r2.bottomRight())
}
