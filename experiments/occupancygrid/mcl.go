package main

import (
	"gonum.org/v1/gonum/stat/distuv"
)

type Particle struct {
	worldPos WorldCordinate
	worldRot float64
	weight   float64
	err      float64
	use      bool
}

func (p *Particle) Move(deltaX, deltaY, deltaRot float64) {
	p.worldPos.X += deltaX
	p.worldPos.Y += deltaY
	p.worldRot += deltaRot
}

type MonteCarlo struct {
	numParticle int
	width       int
	height      int
	particles   []Particle
	stdDev      float64
}

func NewMonteCarlo(w, h int) MonteCarlo {
	numParticle := 300

	mcl := MonteCarlo{
		numParticle: numParticle,
		width:       w,
		height:      h,
		stdDev:      5,
	}

	particles := mcl.createUniformParticle()
	mcl.particles = particles

	return mcl
}

func (mcl *MonteCarlo) ResetToUniform() {
	mcl.particles = mcl.createUniformParticle()
}

func (mcl *MonteCarlo) EstimatePose() (x, y, rot float64) {
	xSum, ySum, rotSum := 0.0, 0.0, 0.0
	num := 0.0

	skipped := 0
	for _, v := range mcl.particles {
		if !v.use {
			skipped++
			// continue
		}
		// fmt.Println("??11")
		num++
		xSum += v.worldPos.X
		ySum += v.worldPos.Y
		rotSum += v.worldRot
	}
	_x, _y, _rot := xSum/num, ySum/num, rotSum/num
	// fmt.Println(_x, _y, _rot)
	// fmt.Println("SKip", skipped)
	return _x, _y, _rot
}

func (mcl *MonteCarlo) Resample() {
	numParticle := mcl.numParticle
	newParticles := make([]Particle, numParticle)

	lastFilledNewParticle := 0
	totWeight := 0.0
	for _, v := range mcl.particles {
		wNormalRandom := distuv.Normal{Mu: v.worldPos.X, Sigma: mcl.stdDev}
		hNormalRandom := distuv.Normal{Mu: v.worldPos.Y, Sigma: mcl.stdDev}
		rotNormalRandom := distuv.Normal{Mu: v.worldRot, Sigma: mcl.stdDev}

		nextGenNumParticle := int(v.weight * float64(numParticle))
		totWeight += v.weight
		// fmt.Println("next gen", v.weight, nextGenNumParticle)
		for i := 0; i < nextGenNumParticle; i++ {
			newParticles[lastFilledNewParticle] = Particle{
				use:      true,
				worldPos: WorldCordinate{wNormalRandom.Rand(), hNormalRandom.Rand()},
				worldRot: rotNormalRandom.Rand(),
				weight:   v.weight,
			}
			// fmt.Println("Resamplexx")
			lastFilledNewParticle++
		}
	}

	if float64(lastFilledNewParticle) < float64(numParticle)*0.3 {
		// fmt.Println("ha", totWeight, lastFilledNewParticle, len(mcl.particles))
		newParticles = mcl.createUniformParticle()
		// 	// wUniRandom := distuv.Uniform{Min: 0, Max: float64(w)}
		// 	// hUniRandom := distuv.Uniform{Min: 0, Max: float64(h)}
		// 	// rotUniRandom := distuv.Uniform{Min: -179, Max: 180}

		// 	// for i := lastFilledNewParticle; i < numParticle; i++ {
		// 	// 	newParticles[lastFilledNewParticle] = Particle{
		// 	// 		x:        wNormalRandom.Rand(),
		// 	// 		y:        hNormalRandom.Rand(),
		// 	// 		worldRot: rotNormalRandom.Rand(),
		// 	// 	}
		// 	// }
	}

	mcl.particles = newParticles
}

func (mcl *MonteCarlo) Update(deltaX, deltaY, deltaRot float64, errorFunction func(worldPos WorldCordinate, rot float64) float64) {
	biggestErr := 0.0
	sumErr := 0.0

	for i, _ := range mcl.particles {
		mcl.particles[i].Move(deltaX, deltaY, deltaRot)

		p := mcl.particles[i]
		_err := errorFunction(WorldCordinate{p.worldPos.X, p.worldPos.Y}, p.worldRot)
		// fmt.Println("err", _err)
		mcl.particles[i].err = _err
		if _err > biggestErr {
			biggestErr = _err
		}
		sumErr += _err
	}

	totWeight := 0.0
	for i, v := range mcl.particles {
		_w := (biggestErr - v.err)
		mcl.particles[i].weight = _w
		totWeight += _w
	}

	if totWeight == 0 {
		mcl.ResetToUniform()
	}

	for i, v := range mcl.particles {
		mcl.particles[i].weight = v.weight / totWeight
	}
}

// Helper

func (mcl *MonteCarlo) createUniformParticle() []Particle {
	numParticle := mcl.numParticle
	wUniRandom := distuv.Uniform{Min: 0, Max: float64(mcl.width)}
	hUniRandom := distuv.Uniform{Min: 0, Max: float64(mcl.height)}
	rotUniRandom := distuv.Uniform{Min: -179, Max: 180}
	particles := make([]Particle, numParticle)

	weight := 1.0 / float64(numParticle)
	for i := 0; i < numParticle; i++ {
		particles[i] = Particle{
			worldPos: WorldCordinate{wUniRandom.Rand(), hUniRandom.Rand()},
			worldRot: rotUniRandom.Rand(),
			weight:   weight,
			use:      true,
		}
	}
	return particles
}
