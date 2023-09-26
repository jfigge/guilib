/*
 * Copyright (C) 2023 by Jason Figge
 */

package components

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jfigge/guilib/graphics/matrix"
)

type Object struct {
	ts     []*Triangle
	offset *vector.Vector
}

const invertY = true

func (o *Object) newInstance() *Object {
	ts := make([]*Triangle, len(o.ts))
	for i, t := range o.ts {
		ts[i] = t.NewInstane()
	}
	return &Object{
		ts:     ts,
		offset: o.offset,
	}
}

func (o *Object) Rotate(xDeg, yDeg, zDeg float64) *Object {
	xRad := xDeg * math.Pi / 180
	yRad := yDeg * math.Pi / 180
	zRad := zDeg * math.Pi / 180
	return o.newInstance().applyMatrix(matrix.RotateXYZ(xRad, yRad, zRad))
}

func (o *Object) Scale(x, y, z float64) *Object {
	scale := matrix.ScaleXYZ(x, y, z)
	return o.newInstance().applyMatrix(scale)
}

func (o *Object) Translate(x, y, z float64) *Object {
	o1 := o.newInstance()
	o1.offset = &vector.Vector{x: x, y: y, z: z, w: 0}
	return o1
}

func (o *Object) Offset() *vector.Vector {
	if o.offset == nil {
		o.offset = &vector.Vector{x: 0, y: 0, z: 0, w: 0}
	}
	return o.offset
}

func (o *Object) OffsetXYZ() (float64, float64, float64) {
	o1 := o.Offset()
	return o1.x, o1.y, o1.z
}

func (o *Object) Triangles() []*Triangle {
	return o.ts
}

func (o *Object) applyMatrix(m *matrix.Matrix4X4) *Object {
	for _, t := range o.ts {
		t.applyMatrix(m)
	}
	return o
}

func LoadObject(filename string) (*Object, error) {
	if _, err := os.Stat(filename); err != nil {
		name := filename
		if !strings.HasSuffix(name, ".obj") {
			name = fmt.Sprintf("%s.obj", name)
		}
		if _, err2 := os.Stat(name); err2 != nil {
			dir, _ := os.Getwd()
			fmt.Printf("%s\n", dir)
			name = filepath.Join(dir, "resources", "objects", name)
			if _, err2 = os.Stat(name); err2 != nil {
				return nil, fmt.Errorf("%s not found: %v", name, err)
			}
		}
		filename = name
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%s can not be opened: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var pts []*vector.Vector
	var ts []*Triangle

	lineCnt := 0
	var pt *vector.Vector
	var t *Triangle

	// https://en.wikipedia.org/wiki/Wavefront_.obj_file#:~:text=The%20OBJ%20file%20format%20is,of%20vertices%2C%20and%20texture%20vertices.
	for scanner.Scan() {
		lineCnt++
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		switch line[:2] {
		case "# ":
			// Ignore comments

		case "v ":
			// List of geometric vertices, with (x, y, z, [w]) coordinates, w is optional and defaults to 1.0.
			pt, err = parseVector(line[1:], lineCnt)
			pts = append(pts, pt)

		case "vt":
			// List of texture coordinates, in (u, [v, w]) coordinates, these will vary between 0 and 1. v, w are optional and default to 0.

		case "vn":
			// List of vertex normals in (x,y,z) form; normals might not be unit vectors.

		case "vp":
			// Parameter space vertices in (u, [v, w]) form; free form geometry statement

		case "f ":
			// Polygonal face element x3 = triangle, x4 = polygon
			// v1 v2 v3 => vertex, no text, no normal
			// v1/t1 v2/t2 v3/t3  => vertex, text, no normal
			// v1/t1/n1 v2/t2/n2 v3/t3/n3  => vertex, text, normal
			t, err = parseFace(pts, line[1:], lineCnt)
			ts = append(ts, t)

		case "l ":
			// Line element

		case "xc":
			// color extension
			err = parseFaceColor(ts, line[2:], lineCnt)

		default:
			log.Printf("ignoring linr %d: %s", lineCnt, line)
		}
		if err != nil {
			return nil, err
		}
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &Object{
		ts: ts,
	}, nil
}

func parseVector(line string, lineCnt int) (*vector.Vector, error) {
	xyzw := strings.Split(strings.TrimSpace(line), " ")
	if len(xyzw) < 3 {
		return nil, fmt.Errorf("too few values in line %d: %s", lineCnt, line)
	} else if len(xyzw) > 4 {
		return nil, fmt.Errorf("too many values in line %d: %s", lineCnt, line)
	}
	var x, y, z, w float64
	var err error

	x, err = strconv.ParseFloat(xyzw[0], 32)
	if err != nil {
		return nil, fmt.Errorf("bad x in line %d: %s", lineCnt, line)
	}

	y, err = strconv.ParseFloat(xyzw[1], 32)
	if err != nil {
		return nil, fmt.Errorf("bad Y in line %d: %s", lineCnt, line)
	}

	z, err = strconv.ParseFloat(xyzw[2], 32)
	if err != nil {
		return nil, fmt.Errorf("bad Z in line %d: %s", lineCnt, line)
	}

	if len(xyzw) == 4 {
		w, err = strconv.ParseFloat(xyzw[3], 32)
		if err != nil {
			return nil, fmt.Errorf("bad W in line %d: %s", lineCnt, line)
		}
	} else {
		w = 1
	}
	return &vector.Vector{x: x, y: y, z: z, w: w}, nil
}

func parseFace(pts []*vector.Vector, line string, lineCnt int) (*Triangle, error) {
	xyz := strings.Split(strings.TrimSpace(line), " ")
	if len(xyz) != 3 {
		return nil, fmt.Errorf("bad face in line %d: %s", lineCnt, line)
	}
	var i1, i2, i3 int64
	var err error

	i1, err = strconv.ParseInt(xyz[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("bad x in line %d: %s", lineCnt, line)
	}

	i2, err = strconv.ParseInt(xyz[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("bad Y in line %d: %s", lineCnt, line)
	}

	i3, err = strconv.ParseInt(xyz[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("bad Z in line %d: %s", lineCnt, line)
	}

	return &Triangle{
		vs:     [3]*vector.Vector{pts[i1-1], pts[i2-1], pts[i3-1]},
		normal: &vector.Vector{0, 0, 0, 1},
		color:  uint32(0xFFFFFFFF),
	}, nil
}

func parseFaceColor(ts []*Triangle, line string, lineCnt int) error {
	irgba := strings.Split(strings.TrimSpace(line), " ")
	if len(irgba) < 1 {
		return fmt.Errorf("too few color entries in line %d: %s", lineCnt, line)
	}
	var i, r, g, b, a uint64
	var err error
	l := len(irgba)

	i, err = strconv.ParseUint(irgba[0], 10, 8)
	if err != nil {
		return fmt.Errorf("bad i in line %d: %s", lineCnt, line)
	}
	if i < 1 || i > uint64(len(ts)) {
		return fmt.Errorf("bad i (%d) in line %d. Must be between 1 & %d: %s", i, lineCnt, len(ts), line)
	}

	if l > 1 {
		r, err = strconv.ParseUint(irgba[1], 10, 8)
		if err != nil {
			return fmt.Errorf("bad r in line %d: %s", lineCnt, line)
		}
	}

	if l > 2 {
		g, err = strconv.ParseUint(irgba[2], 10, 8)
		if err != nil {
			return fmt.Errorf("bad g in line %d: %s", lineCnt, line)
		}
	}

	if l > 3 {
		b, err = strconv.ParseUint(irgba[3], 10, 8)
		if err != nil {
			return fmt.Errorf("bad b in line %d: %s", lineCnt, line)
		}
	}

	if l > 4 {
		a, err = strconv.ParseUint(irgba[4], 10, 8)
		if err != nil {
			return fmt.Errorf("bad a in line %d: %s", lineCnt, line)
		}
	} else {
		a = 255
	}

	ts[i-1].color = uint32(uint8(r))<<24 + uint32(uint8(g))<<16 + uint32(uint8(b))<<8 + uint32(uint8(a))
	return nil
}
