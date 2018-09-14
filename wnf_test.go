package wnf

import (
	"math"
	"testing"
)

func TestPower(t *testing.T) {
	tests := []struct {
		p, d, h, rho, cda, vw, dw, db, gr, mt, expected float64
	}{
		{300, 4808.9, 302.8, 1.225, CdaClimb, 2.8, 192.63, 192.63, 0.08118, Mt, 1.1}, // Headwind
		{300, 4808.9, 302.8, 1.3, CdaClimb, 1.0, 272.38, 192.63, 0.08118, Mt, 1.01},  // 90* + density
		{300, 4808.9, 302.8, 1.1, CdaClimb, 1.0, 272.38, 192.63, 0.08118, Mt, 1.00},  // 90* - density
		{300, 4808.9, 302.8, 1.225, CdaClimb, 2.8, 192.63, 12.63, 0.08118, Mt, 0.95}, // Tailwind
	}
	for _, tt := range tests {
		actual := Power(tt.p, tt.d, tt.h, tt.rho, tt.cda, tt.vw, tt.dw, tt.db, tt.gr, tt.mt)
		if !Eqf(actual, tt.expected) {
			t.Errorf("Power(%.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f, %.3f): got: %.3f, want: %.3f",
				tt.p, tt.d, tt.h, tt.rho, tt.cda, tt.vw, tt.dw, tt.db, tt.gr, tt.mt, actual, tt.expected)
		}
	}
}

// Eqf returns true when floats a and b are equal to within some small epsilon eps.
func Eqf(a, b float64, eps ...float64) bool {
	e := 1e-3
	if len(eps) > 0 {
		e = eps[0]
	}
	// min is the smallest normal value possible
	const min = float64(2.2250738585072014E-308) // 1 / 2**(1022)

	absA := math.Abs(a)
	absB := math.Abs(b)
	diff := math.Abs(a - b)

	if a == b {
		return true
	} else if a == b || b == 0 || diff < min {
		// a or b is zero or both are extremely close to it relative error is less meaningful here
		return diff < (e * min)
	} else {
		// use relative error
		return diff/(absA+absB) < e
	}
}
