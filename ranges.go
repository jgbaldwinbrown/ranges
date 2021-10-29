package ranges

import (
	"fmt"
)

type Intervalable interface {
	Left() float64
	Right() float64
}

type Interval struct {
	LeftEdge float64
	RightEdge float64
}

func (i Interval) Left() float64 {
	return i.LeftEdge
}

func (i Interval) Right() float64 {
	return i.RightEdge
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
		overlaps = append(overlaps, Interval{LeftEdge: float64(index) * window_size, RightEdge: float64(index+1) * window_size})
	}
	return overlaps
}

func (s *Set) addIntervalInternal(i Intervalable) {
	s.Intervals = append(
		s.Intervals,
		TaggedInterval{
			Interval: Interval{
				LeftEdge: i.Left(),
				RightEdge: i.Right(),
			},
			Index: s.NextTag,
		},
	)
	windows := WindowOverlapIndices(i, s.WindowSize)
	tag := s.NextTag
	s.NextTag++
	newInterval := TaggedInterval{
		Interval: Interval{LeftEdge: i.Left(), RightEdge: i.Right()},
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
	intersect.LeftEdge = Max(i.Left(), i2.Left())
	intersect.RightEdge = Min(i.Right(), i2.Right())
	overlap = intersect.LeftEdge < intersect.RightEdge
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

func IntersectSets(s1 Set, s2 Set) (all_intersections [][]TaggedInterval) {
	for _, interval := range s1.Intervals {
		all_intersections = append(all_intersections, Intersections(interval, s2))
	}
	return all_intersections
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
