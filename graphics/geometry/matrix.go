/*
 * Copyright (C) 2023 by Jason Figge
 */

package geometry

import (
	"math"
)

type Matrix4X4 [4][4]float64

func (m *Matrix4X4) NewInstance() *Matrix4X4 {
	return &Matrix4X4{
		{m[0][0], m[0][1], m[0][2], m[0][3]},
		{m[1][0], m[1][1], m[1][2], m[1][3]},
		{m[2][0], m[2][1], m[2][2], m[2][3]},
		{m[3][0], m[3][1], m[3][2], m[3][3]},
	}
}

func Identity() *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func InvertXYZ(x, y, z bool) *Matrix4X4 {
	m1 := &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{1, 0, 1, 0},
		{1, 0, 0, 1},
	}
	if x {
		m1[0][0] = -1
	}
	if y {
		m1[1][1] = -1
	}
	if z {
		m1[2][2] = -1
	}
	return m1
}

func TranslateXYZ(x, y, z float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, x},
		{0, 1, 0, y},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func TranslateX(x, y, z float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, x},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func TranslateY(x, y, z float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, y},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func TranslateZ(z float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, z},
		{0, 0, 0, 1},
	}
}

func Projection(width, height, fovDeg, near, far float64) *Matrix4X4 {
	a := height / width
	f := 1 / math.Tan(fovDeg/360*math.Pi)
	return &Matrix4X4{
		{a * f, 0, 0, 0},
		{0, f, 0, 0},
		{0, 0, far / (far - near), (-far * near) / (far - near)},
		{0, 0, 1, 0},
	}
}

func ScaleXYZ(x, y, z float64) *Matrix4X4 {
	return &Matrix4X4{
		{x, 0, 0, 0},
		{0, y, 0, 0},
		{0, 0, z, 0},
		{0, 0, 0, 1},
	}
}

func ScaleX(x float64) *Matrix4X4 {
	return ScaleXYZ(x, 1, 1)
}

func ScaleY(y float64) *Matrix4X4 {
	return ScaleXYZ(1, y, 1)
}

func ScaleZ(z float64) *Matrix4X4 {
	return ScaleXYZ(1, 1, z)
}

func DegToRad(x float64) float64 {
	return x * math.Pi / 180
}

func RotateXYZ(x, y, z float64) *Matrix4X4 {
	cosX := math.Cos(x)
	sinX := math.Sin(x)
	cosY := math.Cos(y)
	sinY := math.Sin(y)
	cosZ := math.Cos(z)
	sinZ := math.Sin(z)
	return &Matrix4X4{
		{cosY * cosZ, sinX*sinY*cosZ - cosX*sinZ, cosX*sinY*cosZ + sinX*sinZ, 0},
		{cosY * sinZ, sinX*sinY*sinZ + cosX*cosZ, cosX*sinY*sinZ - sinX*cosZ, 0},
		{-sinY, sinX * cosY, cosX * cosY, 0},
		{0, 0, 0, 1},
	}
}

func RotateX(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{1, 0, 0, 0},
		{0, math.Cos(angle), -math.Sin(angle), 0},
		{0, math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}
func RotateY(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{math.Cos(angle), 0, math.Sin(angle), 0},
		{0, 1, 0, 0},
		{-math.Sin(angle), 0, math.Cos(angle), 0},
		{0, 0, 0, 1},
	}
}

func RotateZ(angle float64) *Matrix4X4 {
	return &Matrix4X4{
		{math.Cos(angle), -math.Sin(angle), 0, 0},
		{math.Sin(angle), math.Cos(angle), 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func PointAt(pos, target, up *Vector) *Matrix4X4 {
	newForward := target.Subtract(pos).Normalize()
	newUp := up.Subtract(newForward.Multiply(up.DotProduct(newForward))).Normalize()
	newRight := newUp.CrossProduct(newForward)
	return &Matrix4X4{
		{newRight.x, newRight.y, newRight.z, 0},
		{newUp.x, newUp.y, newUp.z, 0},
		{newForward.x, newForward.y, newForward.z, 0},
		{pos.x, pos.y, pos.z, 1},
	}
}

func LookAtMatrix(pos, target, up *Vector) *Matrix4X4 {
	newForward := target.Subtract(pos).Normalize()
	newUp := up.Subtract(newForward.Multiply(up.DotProduct(newForward))).Normalize()
	newRight := newUp.CrossProduct(newForward)
	return &Matrix4X4{
		{newRight.x, newUp.x, newForward.x, 0},
		{newRight.y, newUp.y, newForward.y, 0},
		{newRight.z, newUp.z, newForward.z, 0},
		{
			-(pos.x*newRight.x + pos.y*newRight.y + pos.z*newRight.z),
			-(pos.x*newUp.x + pos.y*newUp.y + pos.z*newUp.z),
			-(pos.x*newForward.x + pos.y*newForward.y + pos.z*newForward.y),
			1,
		},
	}
}

func (m *Matrix4X4) Multiply(m1 *Matrix4X4) *Matrix4X4 {
	mo := &Matrix4X4{}
	for c := 0; c < 4; c++ {
		for r := 0; r < 4; r++ {
			mo[r][c] = m[r][0]*m1[0][c] + m[r][1]*m1[1][c] + m[r][2]*m1[2][c] + m[r][3]*m1[3][c]
		}
	}
	return mo
}
