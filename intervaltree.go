package intervaltree

import (
	"fmt"
	"github.com/ag0st/binarytree"
	"github.com/ag0st/bst"
	"log"
	"sort"
)

// -----------------------------------------------------
// 				INTERVAL TREE
// -----------------------------------------------------

// IntervalTree struct used to represent an interval tree
// An IntervalTree is a simple BinaryTree with specific values as data. Here data are of type elt
type IntervalTree struct {
	tree *binarytree.BinaryTree
	bst  *bst.BST
}

// NewIntervalTree creates a new interval tree with the intervals given in parameters
func NewIntervalTree(intervals []*Interval) *IntervalTree {
	return &IntervalTree{fromIntervals(intervals[:]), buildBST(intervals[:])}
}

// fromIntervals create a binary tree containing elt struct as data
// Build complexity: O(n), n = len(intervals) cause of searching the median point
func fromIntervals(intervals []*Interval) *binarytree.BinaryTree {
	tree := &binarytree.BinaryTree{}
	length := len(intervals)
	if length == 0 {
		return tree
	}
	// Get the xMid by creating array and sort it
	allPoints := make([]int, length*2)
	for i, in := range intervals {
		allPoints[i] = in.Start
		allPoints[length+i] = in.End
	}

	sort.Ints(allPoints)
	xMid := allPoints[length]

	// divide left and right part
	var left []*Interval
	var right []*Interval
	var mid []*Interval
	for _, in := range intervals {
		if in.End < xMid {
			left = append(left, in)
		} else if in.Start > xMid {
			right = append(right, in)
		} else {
			mid = append(mid, in)
		}
	}

	if len(mid)+len(right)+len(left) != len(intervals) {
		log.Fatalln("MID + LEFT + RIGHT == INTERVALS")
	}
	itr := tree.Root()
	itr.Insert(newElt(mid[:], xMid))
	err := itr.Left().Paste(fromIntervals(left[:]))
	if err != nil {
		log.Fatalln("Cannot paste left tree")
	}
	err = itr.Right().Paste(fromIntervals(right[:]))
	if err != nil {
		log.Fatalln("Cannot paste right tree")
	}
	return tree
}

// intersecting returns all intervals intersecting the value x int he IntervalTree
// Output sensitive: Complexity of O(ln n + k), n = len(intervals in struct) and k = returned intervals
func intersecting(itr *binarytree.Iterator, x int) []*Interval {
	var res []*Interval

	if itr.IsBottom() {
		return res
	}
	e := itr.Consult().(*elt) // must be of this type or panic
	res = append(res, e.intersecting(x)...)
	if x > e.xMid {
		res = append(res, intersecting(itr.Right(), x)...)
	} else if x < e.xMid {
		res = append(res, intersecting(itr.Left(), x)...)
	}
	return res
}

// Containing returns all intervals containing the value x int he IntervalTree
// Output sensitive: Complexity of O(ln n + k), n = len(intervals in struct) and k = returned intervals
func (t *IntervalTree) Containing(x int) []*Interval {
	return intersecting(t.tree.Root(), x)
}

// Intersecting returns all intervals intersecting the Interval given in parameter.
// Output sensitive: Complexity of O(ln n + k), n = len(intervals in struct) and k = returned intervals
func (t *IntervalTree) Intersecting(interval *Interval) []*Interval {
	// First search in the BST for all intersecting intervals
	intervalSearchResult := t.bst.IntervalSearch(&Point{x: interval.Start}, &Point{x: interval.End})
	// remove the duplicates, time depending on searchResult size as bst.IntervalSearch is output sensitive
	set := make(map[*Interval]bool) // uses of map prevent duplicates
	for _, p := range intervalSearchResult {
		p1 := p.(*Point)
		for _, i := range p1.ptrs {
			set[i] = true
		}
	}
	// query the IntervalTree to get all interval that intersect the query interval
	intersectSearchResult := t.Containing(interval.Start)
	for _, in := range intersectSearchResult {
		set[in] = true
	}
	result := make([]*Interval, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	return result
}

// -----------------------------------------------------
// 				INTERVAL TREE NODE
// -----------------------------------------------------

// elt structure representing an IntervalTree element containing a left sorted (start increment) list of intervals
// and right sorted (end decrement) list of intervals with the median point of all its intervals
type elt struct {
	leftSorted  []*Interval
	rightSorted []*Interval
	xMid        int
}

// newElt creates a new element with
// intervals must be a Slice
// This method is in O(n log n) as it uses sort.SliceStable to sort left and right lists
func newElt(intervals []*Interval, xMid int) *elt {
	length := len(intervals)
	intervalTreeElt := &elt{
		leftSorted:  make([]*Interval, length),
		rightSorted: make([]*Interval, length),
		xMid:        xMid,
	}
	// copy the array of intervals
	copy(intervalTreeElt.leftSorted[:], intervals)
	copy(intervalTreeElt.rightSorted[:], intervals)
	// sort start
	sort.SliceStable(
		intervalTreeElt.leftSorted[:], func(i, j int) bool {
			return intervalTreeElt.leftSorted[i].lessStart(intervalTreeElt.leftSorted[j])
		},
	)
	// sort end
	sort.SliceStable(
		intervalTreeElt.rightSorted[:], func(i, j int) bool {
			return intervalTreeElt.rightSorted[i].lessEnd(intervalTreeElt.rightSorted[j])
		},
	)
	return intervalTreeElt
}

// intersecting returns all the intervals that intersect the value "x".
// This method creates a new array of intervals
// Method in O(k) where k is the number of returned intervals
func (e *elt) intersecting(x int) []*Interval {
	if len(e.rightSorted) != len(e.leftSorted) {
		log.Fatalln("MUST HAVE SAME LENGTH")
	}
	var res []*Interval
	if x > e.xMid {
		// begin to check from the end
		for _, in := range e.rightSorted {
			if in.End < x {
				break
			}
			res = append(res, in)
		}
	} else if x < e.xMid {
		// begin to check from the start
		for _, in := range e.leftSorted {
			if in.Start > x {
				break
			}
			res = append(res, in)
		}
	} else {
		// return all the intervals
		res = append(res, e.leftSorted...)
	}
	return res
}

// -----------------------------------------------------
// 				INTERVAL
// -----------------------------------------------------

// Interval structure used to store an interval
type Interval struct {
	Start   int // Start <= End
	End     int
	Payload interface{}
}

// lessStart method used to sort Interval in ascending order of the Interval.Start value, if equals,
// use ascending comparison on Interval.End
// This is used to build the elt.leftSorted array
func (interval *Interval) lessStart(than *Interval) bool {
	if interval.Start == than.Start {
		return interval.End < than.End
	}
	return interval.Start < than.Start
}

// lessEnd method used to sort Interval in descending order of the Interval.End value, if equals,
// use descending comparison on Interval.Start
// This is used to build the elt.rightSorted array
func (interval *Interval) lessEnd(than *Interval) bool {
	if interval.End == than.End {
		return interval.Start > than.Start
	}
	return interval.End > than.End
}

// String prints an interval
func (interval *Interval) String() string {
	return fmt.Sprintf("[ %d - %d ]", interval.Start, interval.End)
}

// -----------------------------------------------------
// 				BST FOR INTERVAL SEARCH
// -----------------------------------------------------

// Point struct representing a point and linked to one or more interval
type Point struct {
	x    int
	ptrs []*Interval
}

// CompareTo implementation of the compare to method from bst.Comparable
func (p *Point) CompareTo(other bst.Comparable) int {
	switch v := other.(type) {
	case *Point:
		if p.x < v.x {
			return -1
		} else if p.x > v.x {
			return 1
		}
		return 0
	default:
		log.Fatalf("POINT COMPARETO: Trying to compare to %s", v)
	}
	return -1 // never reached
}

func buildBST(intervals []*Interval) *bst.BST {
	length := len(intervals)
	if length == 0 {
		return bst.NewBST()
	}
	// Create the array for the BST
	allPoints := make([]bst.Comparable, length*2)
	for i, in := range intervals {
		allPoints[i] = &Point{in.Start, []*Interval{in}}
		allPoints[length+i] = &Point{in.End, []*Interval{in}}
	}

	sort.Slice(
		allPoints, func(i, j int) bool {
			return allPoints[i].CompareTo(allPoints[j]) < 0
		},
	)

	allPoints = removeDuplicateByFusion(allPoints)
	return bst.NewBSTReady(allPoints)
}

// fusion payload of a point with another
// POST: p.ptrs = [p.ptrs +  p2.ptrs]
func (p *Point) fusion(p2 bst.Comparable) {
	p.ptrs = append(p.ptrs, p2.(*Point).ptrs...) // must be *Point, else panic
}

// removeDuplicateByFusion returns a new list of points without duplicates. It merge duplicates points to save all
// the Interval pointer to a same point if multiple intervals shared the same point
func removeDuplicateByFusion(points []bst.Comparable) []bst.Comparable {
	var res []bst.Comparable
	for _, p := range points {
		if res != nil && res[len(res)-1].CompareTo(p) == 0 {
			res[len(res)-1].(*Point).fusion(p) // must be *Point, else panic
		} else {
			res = append(res, p)
		}
	}
	return res
}
