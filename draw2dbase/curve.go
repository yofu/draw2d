// Copyright 2010 The draw2d Authors. All rights reserved.
// created: 17/05/2011 by Laurent Le Goff

package draw2dbase

import (
	"math"
)

const (
	CurveRecursionLimit = 32
)

// Cubic
//	x1, y1, cpx1, cpy1, cpx2, cpy2, x2, y2 float64

// Subdivide a Bezier cubic curve in 2 equivalents Bezier cubic curves.
// c1 and c2 parameters are the resulting curves
func SubdivideCubic(c, c1, c2 []float64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[6], c2[7] = c[6], c[7]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2

	midX := (c[2] + c[4]) / 2
	midY := (c[3] + c[5]) / 2

	c2[4] = (c[4] + c[6]) / 2
	c2[5] = (c[5] + c[7]) / 2

	c1[4] = (c1[2] + midX) / 2
	c1[5] = (c1[3] + midY) / 2

	c2[2] = (midX + c2[4]) / 2
	c2[3] = (midY + c2[5]) / 2

	c1[6] = (c1[4] + c2[2]) / 2
	c1[7] = (c1[5] + c2[3]) / 2

	// Last Point of c1 is equal to the first point of c2
	c2[0], c2[1] = c1[6], c1[7]
}

// TraceCubic generate lines subdividing the cubic curve using a Flattener
// flattening_threshold helps determines the flattening expectation of the curve
func TraceCubic(t Flattener, cubic []float64, flattening_threshold float64) {
	// Allocation curves
	var curves [CurveRecursionLimit * 8]float64
	copy(curves[0:8], cubic[0:8])
	i := 0

	// current curve
	var c []float64

	var dx, dy, d2, d3 float64

	for i >= 0 {
		c = curves[i*8:]
		dx = c[6] - c[0]
		dy = c[7] - c[1]

		d2 = math.Abs((c[2]-c[6])*dy - (c[3]-c[7])*dx)
		d3 = math.Abs((c[4]-c[6])*dy - (c[5]-c[7])*dx)

		// if it's flat then trace a line
		if (d2+d3)*(d2+d3) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c[6], c[7])
			i--
		} else {
			// second half of bezier go lower onto the stack
			SubdivideCubic(c, curves[(i+1)*8:], curves[i*8:])
			i++
		}
	}
}

// Quad
// x1, y1, cpx1, cpy2, x2, y2 float64

// Subdivide a Bezier quad curve in 2 equivalents Bezier quad curves.
// c1 and c2 parameters are the resulting curves
func SubdivideQuad(c, c1, c2 []float64) {
	// First point of c is the first point of c1
	c1[0], c1[1] = c[0], c[1]
	// Last point of c is the last point of c2
	c2[4], c2[5] = c[4], c[5]

	// Subdivide segment using midpoints
	c1[2] = (c[0] + c[2]) / 2
	c1[3] = (c[1] + c[3]) / 2
	c2[2] = (c[2] + c[4]) / 2
	c2[3] = (c[3] + c[5]) / 2
	c1[4] = (c1[2] + c2[2]) / 2
	c1[5] = (c1[3] + c2[3]) / 2
	c2[0], c2[1] = c1[4], c1[5]
	return
}

// Trace generate lines subdividing the curve using a Flattener
// flattening_threshold helps determines the flattening expectation of the curve
func TraceQuad(t Flattener, quad []float64, flattening_threshold float64) {
	// Allocates curves stack
	var curves [CurveRecursionLimit * 6]float64
	copy(curves[0:6], quad[0:6])
	i := 0
	// current curve
	var c []float64
	var dx, dy, d float64

	for i >= 0 {
		c = curves[i*6:]
		dx = c[4] - c[0]
		dy = c[5] - c[1]

		d = math.Abs(((c[2]-c[4])*dy - (c[3]-c[5])*dx))

		// if it's flat then trace a line
		if (d*d) < flattening_threshold*(dx*dx+dy*dy) || i == len(curves)-1 {
			t.LineTo(c[4], c[5])
			i--
		} else {
			// second half of bezier go lower onto the stack
			SubdivideQuad(c, curves[(i+1)*6:], curves[i*6:])
			i++
		}
	}
}

// TraceArc trace an arc using a Flattener
func TraceArc(t Flattener, x, y, rx, ry, start, angle, scale float64) (lastX, lastY float64) {
	end := start + angle
	clockWise := true
	if angle < 0 {
		clockWise = false
	}
	ra := (math.Abs(rx) + math.Abs(ry)) / 2
	da := math.Acos(ra/(ra+0.125/scale)) * 2
	//normalize
	if !clockWise {
		da = -da
	}
	angle = start + da
	var curX, curY float64
	for {
		if (angle < end-da/4) != clockWise {
			curX = x + math.Cos(end)*rx
			curY = y + math.Sin(end)*ry
			return curX, curY
		}
		curX = x + math.Cos(angle)*rx
		curY = y + math.Sin(angle)*ry

		angle += da
		t.LineTo(curX, curY)
	}
	return curX, curY
}