package core

import (
	"math/big"
	"reflect"
	"testing"
)

func TestEstimateC(t *testing.T) {
	corpus := []*Instance{
		NewInstance(kDachengObs, kDachengPeriods),
		NewInstance(kGuanObs, kGuanPeriods),
		NewInstance(kYiObs, kYiPeriods)}
	c := EstimateC(corpus)
	if c != kC {
		t.Errorf("Expecting %d, got %d", kC, c)
	}
}

type mockRng struct {
	History []int
}

func (rng *mockRng) Intn(n int) int {
	if len(rng.History) == 0 {
		rng.History = make([]int, 1, 100)
		rng.History[0] = 0
		return 0
	}
	p := rng.History[len(rng.History)-1]
	if p+1 >= n {
		p = 0
	} else {
		p = p + 1
	}
	rng.History = append(rng.History, p)
	return p
}

func TestInit(t *testing.T) {
	corpus := []*Instance{NewInstance(kDachengObs, kDachengPeriods)}
	c := EstimateC(corpus)
	rng := new(mockRng)
	m := Init(kN, c, corpus, rng)

	truth := &Model{
		S1: []*big.Rat{rat(1), rat(0)},
		Σγ: []*big.Rat{rat(5), rat(4)},
		Σξ: [][]*big.Rat{
			{rat(0), rat(5)},
			{rat(4), rat(0)}},
		Σγo: [][]*Multinomial{
			[]*Multinomial{
				&Multinomial{
					Hist: map[string]*big.Rat{
						"founder":   rat(1),
						"president": rat(4),
						"vice":      rat(4)},
					Sum: rat(9)},
				&Multinomial{
					Hist: map[string]*big.Rat{
						"applied":    rat(4),
						"helping":    rat(1),
						"predictive": rat(4)},
					Sum: rat(9)}},
			[]*Multinomial{
				&Multinomial{
					Hist: map[string]*big.Rat{
						"manager":   rat(1),
						"president": rat(4),
						"senior":    rat(1),
						"vice":      rat(4)},
					Sum: rat(10)},
				&Multinomial{
					Hist: map[string]*big.Rat{
						"applied":    rat(4),
						"linkedin":   rat(1),
						"predictive": rat(4)},
					Sum: rat(9)}}}}

	if !reflect.DeepEqual(m, truth) {
		t.Errorf("Expecting %v, got %v", truth, m)
	}
}
