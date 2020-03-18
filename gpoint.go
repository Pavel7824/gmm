package gmm

import (
	"math/rand"
)

type gpoints struct {
	count int
	dist  *float64
	pts   []interface{}
}

func (p *gpoints) selectionRand() {
	c := p.count
	for c != 0 {
		i := rand.Intn(p.count)
		j := rand.Intn(p.count)
		if i != j {
			p.pts[i], p.pts[j] = p.pts[j], p.pts[i]
			c--
		}
	}
}

func (p *gpoints) mytation() {
	for {
		i := rand.Intn(p.count)
		j := rand.Intn(p.count)
		if i != j {
			p.pts[i], p.pts[j] = p.pts[j], p.pts[i]
			break
		}
	}
}

func (p *gpoints) crossoverPMX(parent2 *gpoints) *gpoints {
	child := new(gpoints)
	var start, end int
	for {
		start = rand.Intn(p.count)
		end = start + rand.Intn(p.count-start)
		if start != end {
			break
		}
	}
	child.count = p.count
	child.pts = make([]interface{}, p.count)
	copy(child.pts, p.pts)
	swaps := make(map[interface{}]interface{})
	for i := start; i < end; i++ {
		swaps[parent2.pts[i]] = child.pts[i]
		child.pts[i] = parent2.pts[i]
	}

	for j := 0; j < start; j++ {
		if end := swaps[child.pts[j]]; end != nil {
			for {
				if val := swaps[end]; val != nil {
					end = val
				} else {
					break
				}
			}
			child.pts[j] = end
		}
	}
	for j := end; j < child.count; j++ {
		if end := swaps[child.pts[j]]; end != nil {
			for {
				if val := swaps[end]; val != nil {
					end = val
				} else {
					break
				}
			}
			child.pts[j] = end
		}
	}
	return child
}
