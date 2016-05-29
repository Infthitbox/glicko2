package glicko2

import (
	"math"
	"testing"
)

func TestSystemUpdate(t *testing.T) {
	player1 := &Rating{1500.0, 200.0, 0.06, []Result{}}
	player2 := &Rating{1400.0, 30.0, 0.06, []Result{}}
	player3 := &Rating{1550.0, 100.0, 0.06, []Result{}}
	player4 := &Rating{1700.0, 300.0, 0.06, []Result{}}

	player1.AddResult(Result{player2, 1.0})
	player2.AddResult(Result{player1, 0.0})

	player1.AddResult(Result{player3, 0.0})
	player3.AddResult(Result{player1, 1.0})

	player1.AddResult(Result{player4, 0.0})
	player4.AddResult(Result{player1, 1.0})

	s := &System{Tau: 0.5, Mu: 1500, Phi: 350, Epsilon: 0.000001, Players: []*Rating{player1, player2, player3, player4}}

	s.Update()

	p1 := s.Players[0]
	if round(p1.Mu, 1) != 1464.1 {
		t.Errorf("Rating should be: (%f), rating is (%f)", 1464.1, p1.Mu)
	}
	if round(p1.Phi, 1) != 151.5 {
		t.Errorf("Deviation should be: (%f), deviation is (%f)", 151.5, p1.Phi)
	}
	if round(p1.Sigma, 3) != 0.060 {
		t.Errorf("Volatility should be: (%f), volatility is (%f)", 0.060, p1.Sigma)
	}

}

func round(s float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor((s*shift)+.5) / shift
}
