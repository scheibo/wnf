package wnf

import (
	"math"

	"github.com/scheibo/calc"
)

const Mb = 8.0
const Mr = 67.0

const CdaClimb = 0.325 // 1.80m
const CdaTT = 0.250    // ?? TODO

type course struct {
	d, gr, wkg float64
}

var climb = course{
	d:   5000.0,
	gr:  0.080,
	wkg: 4.5,
}

var tt = course{
	d:   16093.4,
	gr:  0.0,
	wkg: 4.0,
}

func Score(h, rho, vw float64) float64 {
	return (PowerClimb(climb.wkg*Mr, climb.d, h, rho, vw, climb.gr) +
		PowerTT(tt.wkg*Mr, tt.d, h, rho, vw, tt.gr)) / 2
}

func TimeClimb(p, d, h, rho, vw, gr float64) float64 {
	return Time360(p, d, h, rho, CdaClimb, vw, gr, Mr+Mb)
}

func PowerClimb(p, d, h, rho, vw, gr float64) float64 {
	return Power360(p, d, h, rho, CdaClimb, vw, gr, Mr+Mb)
}

func TimeTT(p, d, h, rho, vw, gr float64) float64 {
	return Time360(p, d, h, rho, CdaTT, vw, gr, Mr+Mb)
}

func PowerTT(p, d, h, rho, vw, gr float64) float64 {
	return Power360(t, d, h, rho, CdaTT, vw, gr, Mr+Mb)
}

func Time360(p, d, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += Time(p, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
	}
	return s / 360
}

func Power360(t, d, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += Power(t, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
	}
	return s / 360
}

func Time(p, d, h, rho, cda, vw, dw, db, gr, mt float64) float64 {
	t := time(p, d, calc.Rho(h, calc.G), cda, vw, dw, db, gr, mt)
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)
	return p / p2
}

func Power(t, d, h, rho, cda, vw, dw, db, gr, mt float64) float64 {
	p1 := power(t, d, calc.Rho(h, calc.G), cda, vw, dw, db, gr, mt)
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)
	return p1 / p2
}

func time(p, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	return calc.Time(p, d, rho, cda, calc.Crr, vw, dw, db, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func power(t, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	vg := d / t
	return calc.Psimp(rho, cda, calc.Crr, calc.Va(vg, vw, dw, db), vg, gr, mt, calc.G, calc.Ec, calc.Fw)
}
