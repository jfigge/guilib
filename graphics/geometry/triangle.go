/*
 * Copyright (C) 2023 by Jason Figge
 */

package geometry

import (
	"fmt"
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Triangle struct {
	vs       [3]*Vector
	normal   *Vector
	normalI  float64
	incident float64
	color    uint32
}

var origin = &Vector{x: 0, y: 0, z: 1}

func (t *Triangle) NewInstance() *Triangle {
	return &Triangle{
		vs:       [3]*Vector{t.vs[0].NewInstance(), t.vs[1].NewInstance(), t.vs[2].NewInstance()},
		normal:   t.normal.NewInstance(),
		incident: 0,
		normalI:  0,
		color:    t.color,
	}
}

func (t *Triangle) Normal(screen *Vector) *Triangle {
	t.vs[0].z = t.vs[0].w
	t.vs[1].z = t.vs[1].w
	t.vs[2].z = t.vs[2].w
	v1 := t.vs[1].Subtract(t.vs[0])
	v2 := t.vs[2].Subtract(t.vs[0])
	t.normal = v1.CrossProduct(v2).Normalize()
	t.normalI = math.Acos(t.normal.DotProduct(origin) / t.normal.Length())
	t.incident = math.Acos(t.normal.DotProduct(screen) / t.normal.Length())
	return t
}

func (t *Triangle) CompareX(t2 *Triangle) int {
	x1 := (t.vs[0].x + t.vs[1].x + t.vs[2].x) / 3
	x2 := (t2.vs[0].x + t2.vs[1].x + t2.vs[2].x) / 3
	return int(x1 - x2)
}

func (t *Triangle) CompareY(t2 *Triangle) int {
	y1 := (t.vs[0].y + t.vs[1].y + t.vs[2].y) / 3
	y2 := (t2.vs[0].y + t2.vs[1].y + t2.vs[2].y) / 3
	return int(y1 - y2)
}

func (t *Triangle) CompareZ(t2 *Triangle) int {
	z1 := (t.vs[0].z + t.vs[1].z + t.vs[2].z) / 3
	z2 := (t2.vs[0].z + t2.vs[1].z + t2.vs[2].z) / 3
	return int(z1 - z2)
}

func (t *Triangle) FaceColor() sdl.Color {
	return sdl.Color{
		R: uint8(t.color >> 24),
		G: uint8(t.color >> 16),
		B: uint8(t.color >> 8),
		A: uint8(t.color),
	}
}

func (t *Triangle) ShadedColor() sdl.Color {
	if t.normalI > (math.Pi / 2) {
		return sdl.Color{R: 32, G: 32, B: 32, A: 0}
	}
	shade := math.Cos(t.normalI)
	return sdl.Color{
		R: uint8(float64(uint8(t.color>>24)) * shade),
		G: uint8(float64(uint8(t.color>>16)) * shade),
		B: uint8(float64(uint8(t.color>>8)) * shade),
		A: uint8(t.color),
	}
}

func (t *Triangle) Multiply(matrix *Matrix4X4) *Triangle {
	t1 := &Triangle{
		vs: [3]*Vector{
			t.vs[0].MatrixMultiply(matrix),
			t.vs[1].MatrixMultiply(matrix),
			t.vs[2].MatrixMultiply(matrix),
		},
		normal:   t.normal,
		incident: t.incident,
		normalI:  t.normalI,
		color:    t.color,
	}
	return t1
}

func (t *Triangle) Project(projection *Matrix4X4) *Triangle {
	t1 := &Triangle{
		vs: [3]*Vector{
			t.vs[0].MatrixMultiply(projection),
			t.vs[1].MatrixMultiply(projection),
			t.vs[2].MatrixMultiply(projection),
		},
		normal:   t.normal,
		incident: t.incident,
		normalI:  t.normalI,
		color:    t.color,
	}
	t1.vs[0] = t1.vs[0].Divide(t1.vs[0].w)
	t1.vs[1] = t1.vs[1].Divide(t1.vs[1].w)
	t1.vs[2] = t1.vs[2].Divide(t1.vs[2].w)
	return t1
}

func (t *Triangle) Log() *Triangle {
	fmt.Println("\n    X           Y           Z           W")
	fmt.Printf(t.vs[0].String())
	fmt.Printf(t.vs[1].String())
	fmt.Printf(t.vs[2].String())
	return t
}

func (t *Triangle) Vertices() []sdl.Vertex {
	c := t.ShadedColor()
	return []sdl.Vertex{
		{Position: *t.vs[0].Point(), Color: c, TexCoord: sdl.FPoint{}},
		{Position: *t.vs[1].Point(), Color: c, TexCoord: sdl.FPoint{}},
		{Position: *t.vs[2].Point(), Color: c, TexCoord: sdl.FPoint{}},
	}
}

func (t *Triangle) Points() []sdl.FPoint {
	return []sdl.FPoint{
		*t.vs[0].Point(),
		*t.vs[1].Point(),
		*t.vs[2].Point(),
		*t.vs[0].Point(),
	}
}

func (t *Triangle) R() uint8 {
	return uint8(t.color >> 24)
}

func (t *Triangle) G() uint8 {
	return uint8(t.color >> 16)
}

func (t *Triangle) B() uint8 {
	return uint8(t.color >> 8)
}

func (t *Triangle) A() uint8 {
	return uint8(t.color)
}

func (t *Triangle) applyMatrix(m *Matrix4X4) {
	t.vs[0] = t.vs[0].MatrixMultiply(m)
	t.vs[1] = t.vs[1].MatrixMultiply(m)
	t.vs[2] = t.vs[2].MatrixMultiply(m)
}

func (t *Triangle) IsVisible() bool {
	return t.normalI < math.Pi/2
}

func (t *Triangle) TranslateXYZ(x, y, z float64) *Triangle {
	offset := &Vector{x: x, y: y, z: z}
	t1 := &Triangle{
		vs: [3]*Vector{
			t.vs[0].Add(offset),
			t.vs[1].Add(offset),
			t.vs[2].Add(offset),
		},
		normal:   t.normal,
		incident: t.incident,
		normalI:  t.normalI,
		color:    t.color,
	}
	return t1
}
