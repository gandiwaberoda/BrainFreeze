package main

import "gonum.org/v1/gonum/stat/distuv"

type Particle struct {
	x, y     float64
	worldRot float64
	weight   float64
}

func (p *Particle) Move(deltaX, deltaY float64) {
	p.x += deltaX
	p.y += deltaY
}

func (p *Particle) Rotate(deltaRot float64) {
	p.worldRot += deltaRot
}

type MonteCarlo struct {
	numParticle int
	width       int
	height      int
	particles   []Particle
}

func NewMonteCarlo(w, h int) MonteCarlo {
	numParticle := 50

	wUniRandom := distuv.Uniform{Min: 0, Max: float64(w)}
	hUniRandom := distuv.Uniform{Min: 0, Max: float64(h)}
	rotUniRandom := distuv.Uniform{Min: -179, Max: 180}
	particles := make([]Particle, numParticle)

	weight := 1.0 / float64(numParticle)
	for i := 0; i < numParticle; i++ {
		particles[i] = Particle{
			x:        wUniRandom.Rand(),
			y:        hUniRandom.Rand(),
			worldRot: rotUniRandom.Rand(),
			weight:   weight,
		}
	}

	return MonteCarlo{
		numParticle: numParticle,
		width:       w,
		height:      h,
		particles:   particles,
	}
}

func (mcl *MonteCarlo) Resample() {

}

func (mcl *MonteCarlo) Update(deltaX, deltaY, deltaRot float64, reading map[float64]LidarReading) {
	for i := range mcl.particles {
		mcl.particles[i].Move(deltaX, deltaY)
		mcl.particles[i].Rotate(deltaRot)
	}
}
