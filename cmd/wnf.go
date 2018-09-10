// wnf provides a CLI for calculating WNF score for arbitrary conditions
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/scheibo/wnf"
)

func main() {
	var rho, cda, mr, mb, mt, vw, dw, db, e, gr, h, mt, mr, mb, t, d, p float64
	var dwS, dbS string
	var err error
	var dur time.Duration

	flag.Float64Var(&rho, "rho", calc.Rho0, "air density in kg/m*3")
	flag.Float64Var(&cda, "cda", 0.325, "coefficient of drag area")

	flag.Float64Var(&mr, "mr", 67.0, "total mass of the rider in kg")
	flag.Float64Var(&mb, "mb", 8.0, "total mass of the bicycle in kg")

	flag.Float64Var(&vw, "vw", 0, "the wind speed in m/s")
	// TODO flag arbitrary direction or number..
	flag.StringVar(&dwS, "dw", "N", "the cardinal direction the wind originates from")
	flag.StringVar(&dbS, "db", "N", "the cardinal direction the bicycle is travelling")

	flag.Float64Var(&e, "e", 0, "total elevation gained in m")
	flag.Float64Var(&gr, "gr", 0, "average grade")
	flag.Float64Var(&h, "h", 0, "median elevation")

	flag.Float64Var(&d, "d", -1, "distance travelled in m")
	flag.Float64Var(&p, "p", -1, "power in watts")
	flag.DurationVar(&dur, "t", -1, "duration in minutes and seconds ('12m34s')")

	flag.Parse()

	verify("rho", rho)
	verify("cda", cda)

	verify("mr", mr)
	verify("mb", mb)
	mt = mr + mb


}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%d:%02d", m, s)
}

func verify(s string, x float64) {
	if x < 0 {
		exit(fmt.Errorf("%s must be non negative but was %f", s, x))
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "%s\n", err)
	flag.PrintDefaults()
	os.Exit(1)
}
