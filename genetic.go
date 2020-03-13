package gmm

import (
	"errors"
	"fmt"
)

// Candidate ...
type Candidate interface {
	ID(index int) int64
	Points() []interface{}
	Dist(from, to int64) float64
}

// GA ...
type GA struct {
	repeat        int
	numPopulatuin int
	parents       []*gpoints
	childs        [][]*gpoints
	candidates    []Candidate
	saveDists     bool
	mapping       map[interface{}]int64                    // link - ID (point)
	dists         map[interface{}]map[interface{}]*float64 // link - linc - dist (from - to points)
}

// SetRepeat ...
func (ga *GA) SetRepeat(r int) {
	ga.repeat = r
}

// SaveDists ...
func (ga *GA) SaveDists(save bool) {
	ga.saveDists = save
}

// NewGA ...
func NewGA(c []Candidate) (ga *GA, err error) {
	ga = new(GA)
	ga.repeat = 1000
	ga.saveDists = false
	ga.parents = make([]*gpoints, len(c))
	ga.mapping = make(map[interface{}]int64, len(c))
	ga.candidates = c
	count := 0
	for i := range ga.candidates {
		ga.parents[i] = new(gpoints)
		ga.parents[i].pts = make([]interface{}, len(c[i].Points()))
		ga.numPopulatuin = len(ga.parents)
		copy(ga.parents[i].pts, ga.candidates[i].Points())
		for ind, pl := range ga.parents[i].pts {
			ga.mapping[pl] = ga.candidates[i].ID(ind)
		}

		if i == 0 {
			count = len(ga.parents[i].pts)
			if count <= 3 {
				return &GA{}, errors.New("the number of points in the population less 3")
			}

		} else {
			if count != len(ga.parents[i].pts) {
				return &GA{}, errors.New("the number of points in the population does not match")
			}
		}
		ga.parents[i].count = count
		ga.calcDistParent(i)
	}
	if ga.numPopulatuin == 1 {
		newPerant := new(gpoints)
		newPerant.count = ga.parents[0].count
		newPerant.pts = make([]interface{}, len(ga.candidates[0].Points()))
		newCandidate := ga.candidates[0]
		copy(newPerant.pts, ga.parents[0].pts)
		ga.parents = append(ga.parents, newPerant)
		ga.candidates = append(ga.candidates, newCandidate)
		ga.calcDistParent(1)
		ga.numPopulatuin = 2
	}
	ga.childs = make([][]*gpoints, ga.numPopulatuin)
	for i := range ga.parents {
		ga.childs[i] = make([]*gpoints, ga.numPopulatuin-1)
		for j := range ga.childs[i] {
			ga.childs[i][j] = new(gpoints)
		}
	}
	return ga, nil
}

// Run ...
func (ga *GA) Run() []interface{} {

	for i := range ga.parents {
		ga.parents[i].selectionRand()
		ga.calcDistParent(i)
	}
	repeat := ga.repeat
	for {
		for i, parenti := range ga.parents {
			k := 0
			for j, parentj := range ga.parents {
				if i != j {
					child := parenti.crossoverPMX(parentj)
					ga.childs[i][k] = child
					k++
				}
			}
		}
		for i, childi := range ga.childs {
			for _, child := range childi {
				child.mytation()
			}
			ga.calcDistChild(i)
		}
		rep := true
		for i := range ga.parents {
			for j := range ga.childs[i] {
				if ga.parents[i].dist > ga.childs[i][j].dist {
					ga.parents[i].count = ga.childs[i][j].count
					ga.parents[i].dist = ga.childs[i][j].dist
					copy(ga.parents[i].pts, ga.childs[i][j].pts)
					repeat = ga.repeat
					rep = false
				}
			}
		}
		if rep {
			repeat--
		}

		if repeat == 0 {
			break
		}
	}
	rez := ga.parents[0]
	for i := 1; i < ga.numPopulatuin; i++ {
		if ga.parents[i].dist < rez.dist {
			rez = ga.parents[i]
		}
	}
	return rez.pts
}

func (ga *GA) calcDistParent(ind int) {
	dist := 0.0
	for i := 1; i < ga.parents[ind].count; i++ {
		if ga.saveDists {
			d := 0.0
			if l1 := ga.dists[ga.parents[ind].pts[i-1]]; l1 != nil {
				if l2 := l1[ga.parents[ind].pts[i]]; l2 != nil {
					d = *l2
				} else {
					l1 = make(map[interface{}]*float64)
					d = ga.candidates[ind].Dist(ga.mapping[ga.parents[ind].pts[i-1]], ga.mapping[ga.parents[ind].pts[i]])
					ga.dists[ga.parents[ind].pts[i-1]][ga.parents[ind].pts[i]] = &d
				}
			} else {
				ga.dists = make(map[interface{}]map[interface{}]*float64)
				ga.dists[ga.parents[ind].pts[i-1]] = make(map[interface{}]*float64)
				d = ga.candidates[ind].Dist(ga.mapping[ga.parents[ind].pts[i-1]], ga.mapping[ga.parents[ind].pts[i]])
				ga.dists[ga.parents[ind].pts[i-1]][ga.parents[ind].pts[i]] = &d
			}
			dist += d
		} else {
			dist += ga.candidates[ind].Dist(ga.mapping[ga.parents[ind].pts[i-1]], ga.mapping[ga.parents[ind].pts[i]])
		}

	}
	ga.parents[ind].dist = dist
}
func (ga *GA) calcDistChild(ind int) {
	for j := range ga.childs[ind] {
		dist := 0.0
		for i := 1; i < ga.childs[ind][j].count; i++ {
			if ga.saveDists {
				d := 0.0
				if l1 := ga.dists[ga.childs[ind][j].pts[i-1]]; l1 != nil {
					if l2 := l1[ga.childs[ind][j].pts[i]]; l2 != nil {
						d = *l2
					} else {
						l1 = make(map[interface{}]*float64)
						d = ga.candidates[ind].Dist(ga.mapping[ga.childs[ind][j].pts[i-1]], ga.mapping[ga.childs[ind][j].pts[i]])
						ga.dists[ga.childs[ind][j].pts[i-1]][ga.childs[ind][j].pts[i]] = &d
					}
				} else {
					ga.dists = make(map[interface{}]map[interface{}]*float64)
					ga.dists[ga.childs[ind][j].pts[i-1]] = make(map[interface{}]*float64)
					d = ga.candidates[ind].Dist(ga.mapping[ga.childs[ind][j].pts[i-1]], ga.mapping[ga.childs[ind][j].pts[i]])
					ga.dists[ga.childs[ind][j].pts[i-1]][ga.childs[ind][j].pts[i]] = &d
				}
				dist += d
			} else {
				dist += ga.candidates[ind].Dist(ga.mapping[ga.childs[ind][j].pts[i-1]], ga.mapping[ga.childs[ind][j].pts[i]])
			}
		}
		ga.childs[ind][j].dist = dist
	}
}

func (ga *GA) printParent(ind int) {
	fmt.Printf("parent %d  |", ind+1)
	for i := 0; i < ga.parents[ind].count; i++ {
		fmt.Printf("%d ", ga.mapping[ga.parents[ind].pts[i]])
	}
	fmt.Printf("|  dist = %.02f\n", ga.parents[ind].dist)
}

func (ga *GA) printChiild(ind int) {
	for j := range ga.childs[ind] {
		fmt.Printf("parent %d child %d  |", ind+1, j+1)
		for i := 0; i < ga.childs[ind][j].count; i++ {
			fmt.Printf("%d ", ga.mapping[ga.childs[ind][j].pts[i]])
		}
		fmt.Printf("|  dist = %.02f\n", ga.childs[ind][j].dist)
	}
}
