/*
 * Copyright (C) 2023 by Jason Figge
 */

package components

import (
	"fmt"
	"math"

	"github.com/jfigge/guilib/graphics/matrix"
	"github.com/veandco/go-sdl2/sdl"
)

type Vector struct {
	x float64
	y float64
	z float64
	w float64
}

func New(x, y, z float64) *Vector {
	return &Vector{x: x, y: y, z: z, w: 1}
}

func (v *Vector) NewInstance() *Vector {
	return &Vector{x: v.x, y: v.y, z: v.z, w: v.w}
}

func (v *Vector) DotProduct(v1 *Vector) float64 {
	return v.x*v1.x + v.y*v1.y + v.z*v1.z
}

func (v *Vector) Length() float64 {
	return math.Sqrt(v.DotProduct(v))
}

func (v *Vector) Normalize() *Vector {
	l := v.Length()
	return &Vector{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
		w: v.w,
	}
}

func (v *Vector) CrossProduct(v1 *Vector) *Vector {
	return &Vector{
		x: v.y*v1.z - v.z*v1.y,
		y: v.z*v1.x - v.x*v1.z,
		z: v.x*v1.y - v.y*v1.x,
		w: v.w,
	}
}

func (v *Vector) Add(v1 *Vector) *Vector {
	return &Vector{
		x: v.x + v1.x,
		y: v.y + v1.y,
		z: v.z + v1.z,
		w: v.w,
	}
}

func (v *Vector) Subtract(v1 *Vector) *Vector {
	return &Vector{
		x: v.x - v1.x,
		y: v.y - v1.y,
		z: v.z - v1.z,
		w: v.w,
	}
}

func (v *Vector) Multiply(l float64) *Vector {
	return &Vector{
		x: v.x * l,
		y: v.y * l,
		z: v.z * l,
		w: v.w,
	}
}

func (v *Vector) Divide(l float64) *Vector {
	if l == 0 {
		l = 1
	}
	return &Vector{
		x: v.x / l,
		y: v.y / l,
		z: v.z / l,
		w: l,
	}
}

func (v *Vector) MatrixMultiply(matrix *matrix.Matrix4X4) *Vector {
	return &Vector{
		x: v.x*matrix[0][0] + v.y*matrix[0][1] + v.z*matrix[0][2] + v.w*matrix[0][3],
		y: v.x*matrix[1][0] + v.y*matrix[1][1] + v.z*matrix[1][2] + v.w*matrix[1][3],
		z: v.x*matrix[2][0] + v.y*matrix[2][1] + v.z*matrix[2][2] + v.w*matrix[2][3],
		w: v.x*matrix[3][0] + v.y*matrix[3][1] + v.z*matrix[3][2] + v.w*matrix[3][3],
	}
}

func (v *Vector) X() float64 {
	return v.x
}

func (v *Vector) Y() float64 {
	return v.y
}

func (v *Vector) Z() float64 {
	return v.z
}

func (v *Vector) W() float64 {
	return v.w
}

func (v *Vector) SetZ(z float64) float64 {
	v.z = z
	return v.z
}

func (v *Vector) Point() *sdl.FPoint {
	return &sdl.FPoint{
		X: float32(v.x),
		Y: float32(v.y),
	}
}

func (v *Vector) String() string {
	return fmt.Sprintf("[ %+8.3f ]  [ %+8.3f ]  [ %+8.3f ]  [ %+8.3f ]\n", v.z, v.y, v.z, v.w)
}
