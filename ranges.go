package main

import (
	"fmt"
)

type Intervalable interface {
	Left() float64
	Right() float64
}

type Interval struct {
	left float64
	right float64
}

func (i Interval) Left() float64 {
	return i.left
}

func (i Interval) Right() float64 {
	return i.right
}

type TaggedInterval struct {
	Interval
	Index int
}

type Set struct {
	QuickIntervals map[int][]TaggedInterval
	Intervals []TaggedInterval
	WindowSize float64
	NextTag int
	WindowCap int
}

func (s *Set) CheckCap() (ok bool) {
	for _, val := range s.QuickIntervals {
		if len(val) > s.WindowCap {
			return false
		}
	}
	return true
}

func NewSet() (s Set) {
	s.QuickIntervals = make(map[int][]TaggedInterval)
	s.WindowSize = 1.0
	s.WindowCap = 10
	return s
}

func WindowOverlapIndices(i Intervalable, window_size float64) (indices []int) {
	winleft := int(i.Left() / window_size)
	winright := int(i.Right() / window_size)
	for pos := winleft; pos <= winright; pos++ {
		indices = append(indices, pos)
	}
	return indices
}

func WindowOverlaps(i Intervalable, window_size float64) (overlaps []Interval) {
	indices := WindowOverlapIndices(i, window_size)
	for _, index := range indices {
		overlaps = append(overlaps, Interval{left: float64(index) * window_size, right: float64(index+1) * window_size})
	}
	return overlaps
}

func (s *Set) addIntervalInternal(i Intervalable) {
	s.Intervals = append(
		s.Intervals,
		TaggedInterval{
			Interval: Interval{
				left: i.Left(),
				right: i.Right(),
			},
			Index: s.NextTag,
		},
	)
	windows := WindowOverlapIndices(i, s.WindowSize)
	tag := s.NextTag
	s.NextTag++
	newInterval := TaggedInterval{
		Interval: Interval{left: i.Left(), right: i.Right()},
		Index: tag,
	}
	for _, window_index := range windows {
		s.QuickIntervals[window_index] = append(s.QuickIntervals[window_index], newInterval)
	}
}

func (s *Set) AddInterval(i Intervalable) {
	s.addIntervalInternal(i)
	s.Recap()
}

func (s *Set) Recap() {
	if ! s.CheckCap() {
		oldSet := s
		newSet := NewSet()
		newSet.WindowSize = oldSet.WindowSize / 2.0
		for _, interval := range oldSet.Intervals {
			newSet.addIntervalInternal(interval)
		}
		s = &newSet
	}
}

func Max(f1 float64, f2 float64) float64 {
	if f1 > f2 {
		return f1
	}
	return f2
}

func Min(f1 float64, f2 float64) float64 {
	if f1 > f2 {
		return f2
	}
	return f1
}

func Intersect(i Intervalable, i2 Intervalable) (intersect Interval, overlap bool) {
	intersect.left = Max(i.Left(), i2.Left())
	intersect.right = Min(i.Right(), i2.Right())
	overlap = intersect.left < intersect.right
	return intersect, overlap
}

func Intersections(i Intervalable, s Set) (intersections []TaggedInterval) {
	hits := make(map[int]struct{})
	window_indices := WindowOverlapIndices(i, s.WindowSize)
	for _, window_index := range window_indices {
		window := s.QuickIntervals[window_index]
		for _, interval := range window {
			if _, does_intersect := Intersect(i, interval) ; does_intersect {
				fmt.Println("intersects!")
				if _,hit := hits[interval.Index]; !hit {
					fmt.Println("intersects!")
					hits[interval.Index] = struct{}{}
					intersections = append(intersections, interval)
				}
			}
		}
	}
	return intersections
}

// func Intersections(i Intervalable, s Set) (intersections []Interval) {
// 	ti := TaggedIntersections(i, s)
// 	for _, tagint := range ti {
// 		intersections = append(intersections, tagint.Interval)
// 	}
// 	return intersections
// }

// func (s *Set) Reduce() Set {
// 	var merged []int
// 	intermediate := NewSet()
// 	newSet := NewSet()
// 	for _, interval := range s.Intervals {
// 		hits := TaggedIntersections(interval, s)
// 		for _, hit := range hits {
// 			merged = append(merged, hit.Index)
// 			newSet.AddInterval(hit)
// 		}
// 	}
// 	return newSet
// }

func main() {
	i := Interval{3.3, 5.5}
	s := NewSet()
	s.AddInterval(Interval{1.2,3.4})
	s.AddInterval(Interval{1.1,2.2})
	s.AddInterval(Interval{5.3,5.6})
	s.AddInterval(Interval{3.4,5.0})
	s.AddInterval(Interval{1.1,7.0})
	fmt.Println(i)
	fmt.Println(s)
	fmt.Println(Intersections(i, s))
}
