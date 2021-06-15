package IntervalTree

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestIntervalTree(t *testing.T) {
	iterations := 20_000
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := 500
	valueSearchedRange := 50000
	maxDistance := 20000
	for i := 0; i < iterations; i++ {
		rand.Seed(time.Now().UnixNano())
		valueSearched := rnd.Intn(valueSearchedRange)
		intsBigger := rand.Intn(n)
		intsSmaller := rand.Intn(n)
		intsInSearch := rand.Intn(n)
		totalInts := intsBigger + intsSmaller + intsInSearch
		name := fmt.Sprintf("Iteration %d", i)
		t.Run(name, func(t *testing.T) {
			// creates the intervals
			log.Printf("Generating %d intervals with %d smaller, %d in search,  %d bigger",
				totalInts, intsSmaller, intsInSearch, intsBigger)
			ints := make([]*Interval, totalInts)
			for i := 0; i < totalInts; i++ {
				var start int
				var end int
				if i < intsBigger {
					// interval bigger
					start = valueSearched + rand.Intn(maxDistance) + 1
					end = start + rand.Intn(maxDistance) + 1 // end > start
				} else if i < intsBigger+intsSmaller {
					// interval smaller
					end = valueSearched - rand.Intn(maxDistance) - 1 // end > start
					start = end - rand.Intn(maxDistance)
				} else {
					// interval intersecting
					start = valueSearched - rand.Intn(maxDistance)
					end = valueSearched + rand.Intn(maxDistance) + 1 // end > start
				}
				ints[i] = &Interval{
					Start: start,
					End:   end,
				}
			}
			log.Println("Generation of the the intervals done")
			log.Println("Creating the the interval tree...")
			now := time.Now()
			intTree := NewIntervalTree(ints[:])
			log.Printf("IntervalTree created in %d nanoseconds", time.Now().Sub(now).Nanoseconds())
			log.Printf("Query question with the value %d", valueSearched)
			now = time.Now()
			containing := intTree.Containing(valueSearched)
			log.Printf("Query in %d nanoseconds", time.Now().Sub(now).Nanoseconds())
			if len(containing) != intsInSearch {
				t.Fatalf("Number of returned interval is not the good one, wanted: %d / received: %d", intsInSearch, len(containing))
			}
		})
	}
}

func TestIntervalTree_Intersecting(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	// Bounds
	upper := rand.Intn(10_000) + 10_100 // from 10_100 to 20_099
	lower := rand.Intn(10_000)          // from 0 to 9_999
	maxDistance := rand.Intn(500) + 100 // from 100 to 599
	// max number of intervals by category
	maxIntervals := rand.Intn(1000) + 500
	// number of intervals intersecting lower bound
	inLow := rand.Intn(maxIntervals)
	// number of intervals intersecting upper bound
	inUpper := rand.Intn(maxIntervals)
	// number of intervals fully enclosed
	inEncl := rand.Intn(maxIntervals)
	// number of intervals fully out closed
	inOutcl := rand.Intn(maxIntervals)
	// number of intervals outside upper
	inOutUppper := rand.Intn(maxIntervals)
	// number of intervals outside lower
	inOutLower := rand.Intn(maxIntervals)
	// total good response
	totalIntersect := inLow + inUpper + inEncl + inOutcl
	log.Printf("Generating intervals with %d inLow, %d inUpper, %d inEncl, %d inOutcl, %d inOutUpper, %d inOutLower",
		inLow, inUpper, inEncl, inOutcl, inOutUppper, inOutLower)

	// create few intervals
	var intervals []*Interval
	//inLow intervals
	for i := 0; i < inLow; i++ {
		rand.Seed(time.Now().UnixNano())
		intervals = append(intervals, &Interval{
			Start:   lower - rand.Intn(maxDistance),
			End:     rand.Intn(upper-lower) + lower,
			Payload: nil,
		})
	}
	// inUpper intervals
	for i := 0; i < inUpper; i++ {
		rand.Seed(time.Now().UnixNano())
		intervals = append(intervals, &Interval{
			Start:   rand.Intn(upper-lower) + lower,
			End:     rand.Intn(maxDistance) + upper,
			Payload: nil,
		})
	}
	// inEncl intervals
	for i := 0; i < inEncl; i++ {
		rand.Seed(time.Now().UnixNano())
		start := rand.Intn((upper-lower)/2) + lower
		end := rand.Intn((upper-lower)/2-1) + start + 1
		intervals = append(intervals, &Interval{
			Start:   start,
			End:     end,
			Payload: nil,
		})
	}
	// inOutcl intervals
	for i := 0; i < inOutcl; i++ {
		rand.Seed(time.Now().UnixNano())
		intervals = append(intervals, &Interval{
			Start:   lower - rand.Intn(maxDistance),
			End:     rand.Intn(maxDistance) + upper,
			Payload: nil,
		})
	}
	// inOutLower intervals
	for i := 0; i < inOutLower; i++ {
		rand.Seed(time.Now().UnixNano())
		end := lower - rand.Intn(maxDistance) - 1
		start := end - rand.Intn(maxDistance) - 1
		intervals = append(intervals, &Interval{
			Start:   start,
			End:     end,
			Payload: nil,
		})
	}
	// inOutUpper intervals
	for i := 0; i < inOutUppper; i++ {
		rand.Seed(time.Now().UnixNano())
		start := upper + rand.Intn(maxDistance) + 1
		end := start + rand.Intn(maxDistance) + 1
		intervals = append(intervals, &Interval{
			Start:   start,
			End:     end,
			Payload: nil,
		})
	}

	tree := NewIntervalTree(intervals)
	result := tree.Intersecting(&Interval{Start: lower, End: upper})
	if len(result) != totalIntersect {
		t.Fatalf("EXPECTING %d VALUES, GOT %d", totalIntersect, len(result))
	}
}
