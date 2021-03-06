package main

import (
	"fmt"
	"gmm"
	"math"
	"math/rand"
	"time"
)

type point struct {
	ID   int64
	X, Y float64
}

func (a *point) distTo(b *point) float64 {
	return math.Sqrt((a.X-b.X)*(a.X-b.X) + (a.Y-b.Y)*(a.Y-b.Y))
}

// Points ...
type Points []*point

// ID ...
func (p Points) ID(index int) int64 {
	return p[index].ID
}

// Points ...
func (p Points) Points() []interface{} {
	out := make([]interface{}, len(p))
	for i := range p {
		out[i] = p[i]
	}
	return out
}

var mapVertex map[int64]*point

// Dist ...
func (p Points) Dist(from, to int64) float64 {
	var froml, tol *point
	if froml = mapVertex[from]; froml == nil {
		fmt.Println("shit from")
		return 0.0
	}
	if tol = mapVertex[to]; tol == nil {
		fmt.Println("shit to")
		return 0.0
	}
	return froml.distTo(tol)
}

func (p Points) print() {
	for _, i := range p {
		fmt.Printf("%d\t%.2f\t%.2f\n", i.ID, i.X, i.Y)
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	pts := make(Points, 10)
	pts[0] = &point{ID: 1, X: 10, Y: 0}
	pts[1] = &point{ID: 2, X: 20, Y: 0}
	pts[2] = &point{ID: 3, X: 30, Y: 0}
	pts[3] = &point{ID: 4, X: 40, Y: 0}
	pts[4] = &point{ID: 5, X: 50, Y: 0}
	pts[5] = &point{ID: 6, X: 60, Y: 0}
	pts[6] = &point{ID: 7, X: 70, Y: 0}
	pts[7] = &point{ID: 8, X: 80, Y: 0}
	pts[8] = &point{ID: 9, X: 90, Y: 0}
	pts[9] = &point{ID: 10, X: 100, Y: 0}

	pts2 := make(Points, 10)
	pts2[0] = &point{ID: 11, X: 0, Y: 10}
	pts2[1] = &point{ID: 12, X: 0, Y: 20}
	pts2[2] = &point{ID: 14, X: 0, Y: 30}
	pts2[3] = &point{ID: 15, X: 0, Y: 40}
	pts2[4] = &point{ID: 16, X: 0, Y: 50}
	pts2[5] = &point{ID: 17, X: 0, Y: 60}
	pts2[6] = &point{ID: 18, X: 0, Y: 70}
	pts2[7] = &point{ID: 19, X: 0, Y: 80}
	pts2[8] = &point{ID: 20, X: 0, Y: 90}
	pts2[9] = &point{ID: 21, X: 0, Y: 100}

	mapVertex = make(map[int64]*point)
	for _, p := range pts {
		mapVertex[p.ID] = p
	}
	for _, p := range pts2 {
		mapVertex[p.ID] = p
	}

	parent := make([]gmm.Candidate, 2)
	parent[0] = pts
	parent[1] = pts2

	ga, err := gmm.NewGA(parent)
	if err != nil {
		fmt.Println(err)
	}
	ga.SetRepeat(10000)
	rez := ga.Run()

	newPts := make(Points, len(rez))
	for i, link := range rez {
		if link == nil {
			fmt.Println("govno: error algoritm")
			continue
		}
		for j, p := range mapVertex {
			if j == link.(*point).ID {
				newPts[i] = p
				break
			}
		}

	}
	newPts.print()
}
