// wnf provides a CLI for calculating WNF score for arbitrary conditions
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/scheibo/calc"
	"github.com/scheibo/wnf"
)

// COMPASS maps from cardinal direction to degrees
var COMPASS = map[string]float64{
	"N":   0,
	"NNE": 22.5,
	"NE":  45,
	"ENE": 67.5,
	"E":   90,
	"ESE": 122.5,
	"SE":  135,
	"SSE": 157.5,
	"S":   180,
	"SSW": 202.5,
	"SW":  225,
	"WSW": 247.5,
	"W":   270,
	"WNW": 292.5,
	"NW":  315,
	"NNW": 337.5,
}

type DirectionFlag struct {
	Direction float64
}

func (df *DirectionFlag) String() string {
	return strconv.FormatFloat(df.Direction, 'f', -1, 64)
}

func (df *DirectionFlag) Set(v string) error {
	d, ok := COMPASS[strings.ToUpper(v)]
	if ok {
		df.Direction = d
		return nil
	}

	d, err := strconv.ParseFloat(v, 64)
	if err == nil {
		df.Direction = d
		return nil
	}

	return fmt.Errorf("invalid direction '%s'", v)
}

func main() {
	var calc360 bool
	var rho, cda, mr, mb, mt, vw, e, gr, h, t, d, p float64
	var dw, db DirectionFlag
	var dur time.Duration

	flag.BoolVar(&calc360, "360", false, "calculate 360 score")

	flag.Float64Var(&rho, "rho", calc.Rho0, "air density in kg/m*3")
	flag.Float64Var(&cda, "cda", -1, "coefficient of drag area")

	flag.Float64Var(&mr, "mr", wnf.Mr, "total mass of the rider in kg")
	flag.Float64Var(&mb, "mb", wnf.Mb, "total mass of the bicycle in kg")

	flag.Float64Var(&vw, "vw", 0, "the wind speed in m/s")
	flag.Var(&dw, "dw", "the cardinal direction the wind originates from")
	flag.Var(&db, "db", "the cardinal direction the bicycle is travelling")

	flag.Float64Var(&e, "e", 0, "total elevation gained in m")
	flag.Float64Var(&gr, "gr", 0, "average grade")
	flag.Float64Var(&h, "h", 0, "median elevation")

	flag.Float64Var(&d, "d", -1, "distance travelled in m")
	flag.Float64Var(&p, "p", -1, "power in watts")
	flag.DurationVar(&dur, "t", -1, "duration in minutes and seconds ('12m34s')")

	flag.Parse()

	verify("rho", rho)

	verify("mr", mr)
	verify("mb", mb)
	mt = mr + mb

	flag.Parse()

	verify("vw", vw)
	verify("h", h)

	// error correct in case grade was passed in as a %
	if gr > 1 || gr < -1 {
		gr = gr / 100
	}

	if e > 0 {
		// if both are specified, make sure they agree
		if gr > 0 && ((d*gr != e) || (e/d != gr)) {
			exit(fmt.Errorf("specified both e=%f and gr=%f but they do not agree", e, gr))
		}
		gr = e / d
	}

	// Assume Climb/TT mode depending on gr if cda is unspecified
	if cda == -1 {
		if gr == 0 {
			cda = wnf.CdaTT
		} else {
			cda = wnf.CdaClimb
		}
	}

	if d <= 0 {
		output(wnf.Score(h, rho, vw))
	} else {
		if p != -1 {
			verify("p", p)
			if dur != -1 {
				exit(fmt.Errorf("t and p can't both be provided"))
			}

			if calc360 {
				output(wnf.Power360(p, d, h, rho, cda, vw, gr, mt))
			} else {
				output(wnf.Power(p, d, h, rho, cda, vw, dw.Direction, db.Direction, gr, mt))
			}
		} else if dur != -1 {
			verify("t", float64(dur))
			t = float64(dur / time.Second)
			if p != -1 {
				exit(fmt.Errorf("p and t can't both be provided"))
			}

			if calc360 {
				output(wnf.Time360(t, d, h, rho, cda, vw, gr, mt))
			} else {
				output(wnf.Time(t, d, h, rho, cda, vw, dw.Direction, db.Direction, gr, mt))
			}
		} else {
			exit(fmt.Errorf("p or t must be specified"))
		}
	}
}

func output(v float64) {
	fi, _ := os.Stdout.Stat()
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		fmt.Println(strconv.FormatFloat(v, 'f', -1, 64))
	} else {
		fmt.Printf("%.2f%% (%.4f)\n", (1.0-v)*100, v)
	}
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
