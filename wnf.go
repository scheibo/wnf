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

func Score(h, rho, vw float64) (float64, float64) {
	cc1, cc2 := PowerClimb(climb.wkg*Mr, climb.d, h, rho, vw, climb.gr)
	tt1, tt2 := PowerTT(tt.wkg*Mr, tt.d, h, rho, vw, tt.gr))
	return (cc1 + tt1) / 2, (cc2 + tt2) / 2
}

func PowerClimb(p, d, h, rho, vw, gr float64) (float64, float64) {
	return Power360(p, d, h, rho, CdaClimb, vw, gr, Mt)
}

func PowerClimbLL(p float64, lls []geo.LatLng, d, h, rho, vw, gr float64) (float64, float64) {
	return Power360LL(p, lls, d, h, rho, CdaClimb, vw, gr, Mt)
}

func TimeClimb(t, d, h, rho, vw, gr float64) (float64, float64) {
	return Time360(t, d, h, rho, CdaClimb, vw, gr, Mt)
}

func TimeClimbLL(t float64, lls []geo.LatLng, d, h, rho, vw, gr float64) (float64, float64) {
	return Time360LL(t, lls, d, h, rho, CdaClimb, vw, gr, Mt)
}

func PowerTT(p, d, h, rho, vw, gr float64) (float64, float64) {
	return Power360(p, d, h, rho, CdaTT, vw, gr, Mt)
}

func PowerTTLL(p float64, lls []geo.LatLng, d, h, rho, vw, gr float64) (float64, float64) {
	return Power360LL(p, lls, d, h, rho, CdaTT, vw, gr, Mt)
}

func TimeTT(t, d, h, rho, vw, gr float64) (float64, float64) {
	return Time360(t, d, h, rho, CdaTT, vw, gr, Mt)
}

func TimeTTLL(t float64, lls []geo.LatLng, d, h, rho, vw, gr float64) (float64, float64) {
	return Time360LL(t, lls, d, h, rho, CdaTT, vw, gr, Mt)
}

func Power360(p, d, h, rho, cda, vw, gr, mt float64) (float64, float64) {
	tot1, tot2 := 0.0, 0.0
	for dw := 0; dw < 360; dw++ {
		s1, s2 := Power(p, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
		tot1 += s1
		tot2 += s2
	}
	return tot1 / 360, tot2 / 360
}

func Power360LL(p float64, lls []geo.LatLng, d, h, rho, cda, vw, gr, mt float64) (float64, float64) {
	tot1, tot2 := 0.0, 0.0
	for dw := 0; dw < 360; dw++ {
		s1, s2 := PowerLL(p, lls, d, h, rho, cda, vw, float64(dw), gr, mt)
		tot1 += s1
		tot2 += s2
	}
	return tot1 / 360, tot2 / 360
}

func Time360(t, d, h, rho, cda, vw, gr, mt float64) (float64, float64) {
	tot1, tot2 := 0.0, 0.0
	for dw := 0; dw < 360; dw++ {
		s1, s2 := Time(t, d, h, rho, cda, vw, float64(dw), 0, gr, mt)
		tot1 += s1
		tot2 += s2
	}
	return tot1 / 360, tot2 / 360
}

func Time360LL(t float64, lls []geo.LatLng, d, h, rho, cda, vw, gr, mt float64) (float64, float64) {
	tot1, tot2 := 0.0, 0.0
	for dw := 0; dw < 360; dw++ {
		s1, s2 := TimeLL(t, lls, d, h, rho, cda, vw, float64(dw), gr, mt)
		tot1 += s1
		tot2 += s2
	}
	return tot1 / 360, tot2 / 360
}

// TODO
func Power(p, d, h, rho, cda, vw, dw, db, gr, mt float64) (float64, float64) {
	t := time(p, d, calc.Rho(h, calc.G), cda, 0, dw, db, gr, mt)

	// how much power for same time?
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)

	// how much faster given same power? = wrong, want how much faster given power
	// AND curve

	mr := mt - Mb
	ratio := (p / (mt - Mb)) / (cp(t) / Mr) // fraction of CP curve

	return p2 / p
}


// TODO
func PowerLL(p float64, lls []geo.LatLng, d, h, rho, cda, vw, dw, gr, mt float64) (float64, float64) {
	t := timeLL(p, lls, d, calc.Rho(h, calc.G), cda, 0, dw, gr, mt)
	p2 := powerLL(t, lls, d, rho, cda, vw, dw, gr, mt)
	return p2 / p
}

// TODO
func Time(t, d, h, rho, cda, vw, dw, db, gr, mt float64) (float64, float64) {
	p1 := power(t, d, calc.Rho(h, calc.G), cda, 0, dw, db, gr, mt)
	p2 := power(t, d, rho, cda, vw, dw, db, gr, mt)
	return p2 / p1
}

// TODO
func TimeLL(t float64, lls []geo.LatLng, d, h, rho, cda, vw, dw, gr, mt float64) (float64, float64) {
	p1 := powerLL(t, lls, d, calc.Rho(h, calc.G), cda, 0, dw, gr, mt)
	p2 := powerLL(t, lls, d, rho, cda, vw, dw, gr, mt)
	return p2 / p1
}

// TODO
func Power2(p, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, db, gr, mt float64) (float64, float64) {
	t := time(p, d, rho1, cda, vw1, dw1, db, gr, mt)
	p2 := power(t, d, rho2, cda, vw2, dw2, db, gr, mt)
	return p2 / p
}

// TODO
func Power2LL(p float64, lls []geo.LatLng, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, gr, mt float64) (float64, float64) {
	t := timeLL(p, lls, d, rho1, cda, vw1, dw1, gr, mt)
	p2 := powerLL(t, lls, d, rho2, cda, vw2, dw2, gr, mt)
	return p2 / p
}

// TODO
func Time2(t, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, db, gr, mt float64) (float64, float64) {
	p1 := power(t, d, rho1, cda, vw1, dw1, db, gr, mt)
	p2 := power(t, d, rho2, cda, vw2, dw2, db, gr, mt)
	return p2 / p1
}

// TODO
func Time2LL(t float64, lls []geo.LatLng, d, rho1, rho2, cda, vw1, vw2, dw1, dw2, gr, mt float64) (float64, float64) {
	p1 := powerLL(t, lls, d, rho1, cda, vw1, dw1, gr, mt)
	p2 := powerLL(t, lls, d, rho2, cda, vw2, dw2, gr, mt)
	return p2 / p1
}

func time(p, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	return calc.Time(p, d, rho, cda, calc.Crr, vw, dw, db, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func timeLL(p float64, lls []geo.LatLng, d, rho, cda, vw, dw, gr, mt float64) float64 {
	t := 0.0
	if len(lls) <= 1 {
		return t
	}

	// We need to determine the computed distance D and figure out how it compares
	// to the provided distance d so that we can adjust each individual computed
	// distance segment proportionally.
	td := totalDistance(lls, gr)

	ll := lls[0]
	for i := 1; i < len(lls); i++ {
		t += time(p, distance(ll, lls[i], gr)*d/td, rho, cda, vw, dw, geo.Bearing(ll, lls[i]), gr, mt)
		ll = lls[i]
	}

	return t
}

func power(t, d, rho, cda, vw, dw, db, gr, mt float64) float64 {
	vg := d / t
	return calc.Psimp(rho, cda, calc.Crr, calc.Va(vg, vw, dw, db), vg, gr, mt, calc.G, calc.Ec, calc.Fw)
}

func powerLL(tt float64, lls []geo.LatLng, d, rho, cda, vw, dw, gr, mt float64) float64 {
	p := 0.0
	if len(lls) <= 1 {
		return p
	}

	td := 0.0
	var dis []float64
	var dbs []float64

	ll := lls[0]
	for i := 1; i < len(lls); i++ {
		di := distance(ll, lls[i], gr)
		td += di
		dis = append(dis, di)
		dbs = append(dbs, geo.Bearing(ll, lls[i]))
		ll = lls[i]
	}

	for i := 0; i < len(dis); i++ {
		// Adjust the computed segment distance di to be proportional to the
		// total provided distance d.
		dadj := dis[i] * (d / td)
		// The time t to complete this small segment compared to the total time tt
		// is proportional to the adjusted segment distance dadj compared to the
		// total provided distance d.
		t := tt * (dadj / d)
		// We need to weight this segments contribution to the total average power
		// by the fraction of the total time that was spent at this power.
		p += power(t, dadj, rho, cda, vw, dw, dbs[i], gr, mt) * (t / tt)
	}

	return p
}

func cp(t float64) float64 {
	return 422.58 + 23296.7801287949/t
}

func distance(p1, p2 geo.LatLng, gr float64) float64 {
	run := geo.Distance(p1, p2)
	// NOTE: Assuming even average gradient for the entire track.
	rise := gr * run
	return math.Sqrt(run*run + rise*rise)
}

func totalDistance(lls []geo.LatLng, gr float64) float64 {
	d := 0.0
	if len(lls) <= 1 {
		return d
	}

	ll := lls[0]
	for i := 1; i < len(lls); i++ {
		d += distance(ll, lls[i], gr)
		ll = lls[i]
	}
	return d
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
