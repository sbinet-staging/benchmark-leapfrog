package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"

	"github.com/gonum/floats"
)

const nParticles = 2

var dx [3]float64

type Vec3 []float64

func newArray() *[nParticles]Vec3 {
	var arr [nParticles]Vec3
	for i := range arr {
		arr[i] = Vec3(make([]float64, 3))
	}
	return &arr
}

func main() {
	{
		f, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}
	var time float64
	timeStep := 0.08
	halfTimeStep := timeStep / 2.
	timeLimit := 365.25 * 1e6

	// Arrays (initialized to zero by default)
	x := newArray()
	v := newArray()
	a := newArray()
	m := &[nParticles]float64{}

	m[0] = 0.08                 // M_SUN
	m[1] = 3.0e-6               // M_SUN
	x[1][0] = 0.0162            // AU
	x[1][1] = 6.57192058353e-15 // AU
	x[1][2] = 5.74968548652e-16 // AU
	v[1][0] = -1.48427302304e-14
	v[1][1] = 0.0399408809121
	v[1][2] = 0.00349437429104

	for time <= timeLimit {
		integrator_leapfrog_part1(x, v, halfTimeStep)
		time += halfTimeStep
		gravity_calculate_acceleration(m, x, a)
		integrator_leapfrog_part2(x, v, a, timeStep, halfTimeStep)
		time += halfTimeStep
	}
	fmt.Println("Positions:", x)
}

func integrator_leapfrog_part1(x *[nParticles]Vec3, v *[nParticles]Vec3, halfTimeStep float64) {
	for i := 0; i < nParticles; i++ {
		floats.AddScaled(x[i], halfTimeStep, v[i])
	}
}

func integrator_leapfrog_part2(x *[nParticles]Vec3, v *[nParticles]Vec3, a *[nParticles]Vec3, timeStep float64, halfTimeStep float64) {
	for i := 0; i < nParticles; i++ {
		floats.AddScaled(v[i], timeStep, a[i])
		floats.AddScaled(x[i], halfTimeStep, v[i])
	}
}

func gravity_calculate_acceleration(m *[nParticles]float64, x *[nParticles]Vec3, a *[nParticles]Vec3) {
	G := 6.6742367e-11 // m^3.kg^-1.s^-2
	dx := dx[:]
	for i := 0; i < nParticles; i++ {
		a[i][0] = 0
		a[i][1] = 0
		a[i][2] = 0
		for j := 0; j < nParticles; j++ {
			if j == i {
				continue
			}
			floats.AddScaledTo(dx, x[i], -1, x[j])
			r := math.Sqrt(floats.Dot(dx, dx))
			prefact := -G / (r * r * r) * m[j]
			floats.AddScaled(a[i], prefact, dx)
		}
	}
}
