package glicko2

import (
	"math"
)

type Result struct {
	Opponent *Rating
	Score    float64
}

type Rating struct {
	Mu      float64
	Phi     float64
	Sigma   float64
	Results []Result
}

type System struct {
	Tau     float64
	Epsilon float64
	Mu      float64
	Phi     float64
	Players []*Rating
}

func (p *Rating) AddResult(r Result) {
	s := p.Results
	s = append(s, r)
	p.Results = s
}

func (p *Rating) estimateFunctionG() (g float64) {
	denom := math.Sqrt(1 + ((3 * (math.Pow(p.Phi, 2))) / math.Pow(math.Pi, 2)))
	return 1 / denom
}

func (p *Rating) estimateFunctionE(r *Rating) (e float64) {
	denom := 1 + math.Exp(-r.estimateFunctionG()*(p.Mu-r.Mu))
	return 1 / denom
}

func (p *Rating) estimateVariance() (v float64) {
	denom := 0.0
	for _, res := range p.Results {
		e := p.estimateFunctionE(res.Opponent)
		denom += (math.Pow(res.Opponent.estimateFunctionG(), 2)) * e * (1 - e)
	}
	return 1 / denom
}

func (p *Rating) estimateDelta(v float64) (delta float64) {
	sum := 0.0
	for _, res := range p.Results {
		sum += res.Opponent.estimateFunctionG() * (res.Score - p.estimateFunctionE(res.Opponent))
	}
	return v * sum
}

func (s *System) ScaleDown(p *Rating) {
	mu := (p.Mu - s.Mu) / 173.7178
	phi := p.Phi / 173.7178
	p.Mu = mu
	p.Phi = phi
}

func (s *System) ScaleUp(p *Rating) {
	mu := p.Mu*173.7178 + s.Mu
	phi := p.Phi * 173.7178
	p.Mu = mu
	p.Phi = phi
}

func convergenceFunction(x float64, delta float64, phi float64, v float64, tau float64, a float64) (f float64) {
	denomleft := 2 * math.Pow(math.Pow(phi, 2)+v+math.Exp(x), 2)
	numleft := math.Exp(x) * (math.Pow(delta, 2) - math.Pow(phi, 2) - v - math.Exp(x))
	denomright := math.Pow(tau, 2)
	numright := x - a
	return (numleft / denomleft) - (numright / denomright)
}

func (s *System) updateVolatility(p *Rating, v float64, delta float64) {
	a := math.Log(math.Pow(p.Sigma, 2))
	bigA := a
	bigB := 0.0
	if math.Pow(delta, 2) > (math.Pow(p.Phi, 2) + v) {
		bigB = math.Log(math.Pow(delta, 2) - math.Pow(p.Phi, 2) - v)
	} else {
		k := 1.0
		for convergenceFunction(a-(k*s.Tau), delta, p.Phi, v, s.Tau, a) < 0 {
			k += 1.0
		}
		bigB = a - (k * s.Tau)
	}
	fA := convergenceFunction(bigA, delta, p.Phi, v, s.Tau, a)
	fB := convergenceFunction(bigB, delta, p.Phi, v, s.Tau, a)
	for math.Abs(bigB-bigA) > s.Epsilon {
		bigC := bigA + (bigA-bigB)*fA/(fB-fA)
		fC := convergenceFunction(bigC, delta, p.Phi, v, s.Tau, a)
		if fC*fB < 0 {
			bigA = bigB
			fA = fB
		} else {
			fA = fA / 2
		}
		bigB = bigC
		fB = fC
	}
	p.Sigma = math.Exp(bigA / 2)
}

func (p *Rating) determinePhiStar() (phistar float64) {
	return math.Sqrt(math.Pow(p.Phi, 2) + math.Pow(p.Sigma, 2))
}

func (p *Rating) updatePhi(phistar float64, v float64) {
	denom := math.Sqrt((1 / math.Pow(phistar, 2)) + (1 / v))
	phiprime := 1 / denom
	p.Phi = phiprime
}

func (p *Rating) updateMu() {
	sum := 0.0
	for _, res := range p.Results {
		sum += res.Opponent.estimateFunctionG() * (res.Score - p.estimateFunctionE(res.Opponent))
	}
	p.Mu = p.Mu + (math.Pow(p.Phi, 2) * sum)
}

func (s *System) ratePlayer(p *Rating) (r *Rating) {
	r = &Rating{Mu: p.Mu, Phi: p.Phi, Sigma: p.Sigma, Results: p.Results}
	variance := r.estimateVariance()
	delta := r.estimateDelta(variance)
	s.updateVolatility(r, variance, delta)
	phistar := r.determinePhiStar()
	r.updatePhi(phistar, variance)
	r.updateMu()
	return r
}

func (s *System) Update() {
	updatedplayers := make([]*Rating, 0)
	for _, player := range s.Players {
		s.ScaleDown(player)
	}
	for _, player := range s.Players {
		updatedplayers = append(updatedplayers, s.ratePlayer(player))
	}
	s.Players = updatedplayers
	for _, player := range s.Players {
		s.ScaleUp(player)
	}
}
