package wnf

import (
	"github.com/scheibo/calc"
	"github.com/scheibo/geo"
	"math"
)

type Wpt struct {
	geo.LatLng
	Ele  float64
	Time float64
}

const Mb = 8.0
const Mr = 67.0
const Mt = Mr + Mb

const CdaClimb = 0.325 // 1.80m
const CdaTT = 0.250    // ??

func Score(h, rho, vw float64) float64 {
	return (PowerClimb(climb.wkg*Mr, climb.d, h, rho, vw, climb.gr) +
		PowerTT(tt.wkg*Mr, tt.d, h, rho, vw, tt.gr)) / 2
}

func PowerClimb(p, d, h, rho, vw, gr float64) float64 {
	return Power360(p, d, h, rho, CdaClimb, vw, gr, Mt)
}

func PowerClimbLL(p float64, lls []geo.LatLng, h, rho, vw, gr float64) float64 {
	return Power360LL(p, lls, h, rho, CdaClimb, vw, gr, Mt)
}

func TimeClimb(t, d, h, rho, vw, gr float64) float64 {
	return Time360(t, d, h, rho, CdaClimb, vw, gr, Mt)
}

func TimeClimbLL(t float64, lls []geo.LatLng, h, rho, vw, gr float64) float64 {
	return Time360LL(t, lls, h, rho, CdaClimb, vw, gr, Mt)
}

func PowerTT(p, d, h, rho, vw, gr float64) float64 {
	return Power360(p, d, h, rho, CdaTT, vw, gr, Mt)
}

func PowerTTLL(p float64, lls []geo.LatLng, h, rho, vw, gr float64) float64 {
	return Power360LL(p, lls, h, rho, CdaTT, vw, gr, Mt)
}

func TimeTT(t, d, h, rho, vw, gr float64) float64 {
	return Time360(t, d, h, rho, CdaTT, vw, gr, Mt)
}

func TimeTTLL(t float64, lls []geo.LatLng, h, rho, vw, gr float64) float64 {
	return Time360LL(t, lls, h, rho, CdaTT, vw, gr, Mt)
}

func Power360(p, d, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += Power(p, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
	}
	return s / 360
}

func Power360LL(p float64, lls []geo.LatLng, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += PowerLL(p, lls, h, rho, cda, vw, float64(dw), gr, mt)
	}
	return s / 360
}

func Time360(t, d, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += Time(t, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
	}
	return s / 360
}

func Time360LL(t float64, lls []geo.LatLng, h, rho, cda, vw, gr, mt float64) float64 {
	s := 0.0
	for dw := 0; dw < 360; dw++ {
		s += TimeLL(t, lls, h, rho, cda, vw, float64(dw), gr, mt)
	}
	return s / 360
}

func Power(p, d, h, rho, cda, vw, dw, db, gr, mt float64) float64 {
	t := time(p, d, calc.Rho(h, calc.G), cda, 0, dw, db, gr, mt)
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)
	return p2 / p
}

func PowerLL(p float64, lls []geo.LatLng, h, rho, cda, vw, dw, gr, mt float64) float64 {
	t := timeLL(p, lls, calc.Rho(h, calc.G), cda, 0, dw, gr, mt)
	p2 := powerLL(t, lls, rho, cda, vw, dw, gr, mt)
	return p2 / p
}

func Time(t, d, h, rho, cda, vw, dw, db, gr, mt float64) float64 {
	p1 := power(t, d, calc.Rho(h, calc.G), cda, 0, dw, db, gr, mt)
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)
	return p2 / p1
}

func TimeLL(t float64, lls []geo.LatLng, h, rho, cda, vw, dw, gr, mt float64) float64 {
	p1 := powerLL(t, lls, calc.Rho(h, calc.G), cda, 0, dw, gr, mt)
	p2 := powerLL(t, lls, rho, cda, vw, dw, gr, mt)
	return p2 / p1
}

func Power2(p, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, db, gr, mt float64) float64 {
	t := time(p, d, rho1, cda, vw1, dw1, db, gr, mt)
	p2 := power(t, d, rho2, cda, vw2, dw2, db, gr, mt)
	return p2 / p
}

func Power2LL(p float64, lls []geo.LatLng, rho1, rho2, cda, vw1, vw2, dw1, dw2, gr, mt float64) float64 {
	t := timeLL(p, lls, rho1, cda, vw1, dw1, gr, mt)
	p2 := powerLL(t, lls, rho2, cda, vw2, dw2, gr, mt)
	return p2 / p
}

func Time2(t, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, db, gr, mt float64) float64 {
	p1 := power(t, d, rho1, cda, vw1, dw1, db, gr, mt)
	p2 := power(t, d, rho2, cda, vw2, dw2, db, gr, mt)
	return p2 / p1
}

func Time2LL(t float64, lls []geo.LatLng, rho1, rho2, cda, vw1, vw2, dw1, dw2, gr, mt float64) float64 {
	p1 := powerLL(t, lls, rho1, cda, vw1, dw1, gr, mt)
	p2 := powerLL(t, lls, rho2, cda, vw2, dw2, gr, mt)
	return p2 / p1
}

func time(p, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	return calc.Time(p, d, rho, cda, calc.Crr, vw, dw, db, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func timeLL(p float64, lls []geo.LatLng, rho, cda, vw, dw, gr, mt float64) float64 {
	t := 0.0
	if len(lls) <= 1 {
		return t
	}

	ll := lls[0]
	for i := 1; i < len(lls); i++ {
		t += time(p, distance(ll, lls[i], gr), rho, cda, vw, dw, geo.Bearing(ll, lls[i]), gr, mt)
	}

	return t
}

// TODO(kjs): confirm sum of segment distances (with hypotenuse) = strava reported distance!

func power(t, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	vg := d / t
	return calc.Psimp(rho, cda, calc.Crr, calc.Va(vg, vw, dw, db), vg, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func powerLL(tt float64, lls []geo.LatLng, rho, cda, vw, dw, gr, mt float64) float64 {
	p := 0.0
	if len(lls) <= 1 {
		return p
	}

	td := 0.0
	var ds []float64
	var dbs []float64

	ll := lls[0]
	for i := 1; i < len(lls); i++ {
		d := distance(ll, lls[i], gr)
		td += d
		ds = append(ds, d)
		dbs = append(dbs, geo.Bearing(ll, lls[i]))
	}

	for i := 0; i < len(ds); i++ {
		d := ds[i]
		// The time t to complete this small segment compared to the total time tt is
		// proportiona to the segment distance d compared to the total distance td.
		t := tt * (d / td)
		// We need to weight this segments contribution to the total average power by
		// the fraction of the total time that was spent at this power.
		p += power(t, d, rho, cda, vw, dw, dbs[i], gr, mt) * (t / tt)
	}

	return p
}

func distance(p1, p2 geo.LatLng, gr float64) float64 {
	run := geo.Distance(p1, p2)
	// NOTE: Assuming even average gradient for the entire track.
	rise := gr * run
	return math.Sqrt(run*run + rise*rise)
}

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
