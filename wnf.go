package wnf

import (
	"math"

	"github.com/scheibo/calc"
)

const mb = 8.0
const mr = 67.0

const cdaC = 0.325 // 1.80m
const cdaTT = 0.250 // ?? TODO

type course struct {
	d, gr, wkg float64
}

var climb = course{
	d: 5000.0,
	gr: 0.080,
	p: 4.5
}

var tt = course{
	d: 16093.4,
	gr: 0.0,
	p: 4.0
}

// TTT(p), ClimbT(p)
// TTP(t), ClimbP(p)

// or, don't bother with that at all? just take in cda/mr/mb and provide default constants
// Score(rho, cda, crr, va, vg, gr, mt, g, ec, fw float64) float64 {


// TODO(kjs): add functionality to convert from points to slices, keep grade constant
// and ensure tn - t0 = total duration
type Slice struct {
	d, gr, t, db float64
}


func TT(d, gr, t, db, h, rho, vw, dw float64) {
	s := &Slice{d: d, gr: gr, t: t, db: db}
	return SlicedTTM([]Slice{s}, h, rho, vw, db, dw)
}

// TODO(kjs): want both time and power functionality:
// - power functionaly for effort independent model (fixed W/kg),
// - time functionality for figuring out after fact how performance was affected

func SlicedTT(slices []*Slice, h, rho, vw, db, dw float64) {
	for _, s := slices {


	}
}

func time(p, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	return calc.Time(p, d, rho, cda, calc.Crr, vw, dw, db, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func power(t, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	vg := d / t
	return calc.Psimp(rho, cda, calc.Crr, calc.Va(vg, vw, dw, db), vg, gr, mt, calc.G, calc.Ec, calc.Fw)
