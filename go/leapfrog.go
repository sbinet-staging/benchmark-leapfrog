package main

import (
	"fmt"
	"math"
)

const (
	nParticles = 2
	G          = 6.6742367e-11 // m^3.kg^-1.s^-2
)

func main() {
	var (
		dt   = 0.08
		hdt  = dt / 2.
		tmax = 365.25 * 1e6
	)
	// Arrays (initialized to zero by default)
	x := [nParticles][3]float64{
		{0, 0, 0},
		{0.0162, 6.57192058353e-15, 5.74968548652e-16}, // AU
	}
	v := [nParticles][3]float64{
		{0, 0, 0},
		{-1.48427302304e-14, 0.0399408809121, 0.00349437429104},
	}
	a := [nParticles][3]float64{}
	m := [nParticles]float64{0.08, 3.0e-6} // M_SUN

	for t := 0.0; t <= tmax; t += dt {
		// first step of leapfrog
		for i := 0; i < nParticles; i++ {
			x[i][0] += hdt * v[i][0]
			x[i][1] += hdt * v[i][1]
			x[i][2] += hdt * v[i][2]
		}

		// compute forces
		dx := x[0][0] - x[1][0]
		dy := x[0][1] - x[1][1]
		dz := x[0][2] - x[1][2]
		d2 := dx*dx + dy*dy + dz*dz
		dr := math.Sqrt(d2)
		prefact := -G / (dr * d2)

		pm0 := prefact * m[0]
		pm1 := prefact * m[1]

		a[0][0] = pm1 * dx
		a[0][1] = pm1 * dy
		a[0][2] = pm1 * dz

		a[1][0] = pm0 * -dx
		a[1][1] = pm0 * -dy
		a[1][2] = pm0 * -dz

		// second step of leapfrog
		for i := 0; i < nParticles; i++ {
			v[i][0] += dt * a[i][0]
			v[i][1] += dt * a[i][1]
			v[i][2] += dt * a[i][2]
			x[i][0] += hdt * v[i][0]
			x[i][1] += hdt * v[i][1]
			x[i][2] += hdt * v[i][2]
		}
	}
	fmt.Println("Positions:", x)
}
