// movingminmax_test.go, jpad 2015

package movingminmax

import (
	//	"fmt"
	"math/rand"
	"testing"

	"github.com/notnot/container/deque_int"
)

const (
	N = 1000
	W = 10
)

var (
	values       []float64 // input sample values
	mmins, mmaxs []float64 // reference results
)

// init initializes the test values.
func init() {
	values = make([]float64, N)
	mmins = make([]float64, N)
	mmaxs = make([]float64, N)

	// generate random sample values in [0.0, 1.0)
	rand.Seed(123)
	for i := 0; i < N; i++ {
		values[i] = rand.float64()
	}
	/*
		// display reference results
		fmt.Printf("moving minmax, reference code (offline), W = %d:\n", W)
		MovingMinMax_offline_a(W)
		for i := 0; i < N; i++ {
			fmt.Printf("value[%02d] %.3f : min %.3f, max %.3f\n",
				i, values[i], mmins[i], mmaxs[i])
		}
	*/
}

//// tests /////////////////////////////////////////////////////////////////////

func TestMovingMinMax(t *testing.T) {
	// test with various window widths
	for w := uint(1); w < N; w++ {
		minmax := NewMovingMinMax(w)
		MovingMinMax_offline(int(w))

		for i := range values[:len(values)-1] {
			minmax.Update(values[i])

			min := minmax.Min()
			if min != mmins[i] {
				t.Errorf("W %d: values[%d] Min() got: %.3f, want: %.3f",
					w, i, min, mmins[i])
			}

			max := minmax.Max()
			if max != mmaxs[i] {
				t.Errorf("W %d: values[%d] Max() got: %.3f, want: %.3f",
					w, i, max, mmaxs[i])
			}
		}
	}
}

func TestMovingMinMax0(t *testing.T) {
	// test with various window widths
	for w := uint(1); w < N; w++ {
		minmax := NewMovingMinMax0(w)
		MovingMinMax_offline(int(w))

		for i := range values[:len(values)-1] {
			minmax.Update(values[i])

			min := minmax.Min()
			if min != mmins[i] {
				t.Errorf("W %d: values[%d] Min() got: %.3f, want: %.3f",
					w, i, min, mmins[i])
			}

			max := minmax.Max()
			if max != mmaxs[i] {
				t.Errorf("W %d: values[%d] Max() got: %.3f, want: %.3f",
					w, i, max, mmaxs[i])
			}
		}
	}
}

func TestMovingMin(t *testing.T) {
	// test with various window widths
	for w := uint(1); w < N; w++ {
		mmin := NewMovingMin(w)
		MovingMinMax_offline(int(w))

		for i := range values[:len(values)-1] {
			mmin.Update(values[i])

			min := mmin.Min()
			if min != mmins[i] {
				t.Errorf("values[%d] Min() got: %.3f, want: %.3f",
					i, min, mmins[i])
			}
		}
	}
}

func TestMovingMax(t *testing.T) {
	// test with various window widths
	for w := uint(1); w < N; w++ {
		mmax := NewMovingMax(w)
		MovingMinMax_offline(int(w))

		for i := range values[:len(values)-1] {
			mmax.Update(values[i])

			max := mmax.Max()
			if max != mmaxs[i] {
				t.Errorf("values[%d] Max() got: %.3f, want: %.3f",
					i, max, mmaxs[i])
			}
		}
	}
}

func TestMovingMean(t *testing.T) {
	mmean := NewMovingMean(W)
	for i := range values {
		mmean.Update(values[i])
		/*
			fmt.Printf("values[%2d] %.3f: mean %.3f\n",
				i, values[i], mmean.Mean())
		*/
	}
}

func TestDequePushPop(t *testing.T) {
	deque := newDeque_IV(10)
	iv1 := _IV{1, 1.0}
	iv2 := _IV{2, 2.0}

	deque.PushFront(iv1.i, iv1.v)
	deque.PushBack(iv2.i, iv2.v)

	item := deque.PopFront()
	if *item != iv1 {
		t.Errorf("got: %v, want: %v", item, iv1)
	}
	item = deque.PopBack()
	if *item != iv2 {
		t.Errorf("got: %v, want: %v", item, iv2)
	}
}

func TestDequeFrontBackItem(t *testing.T) {
	deque := newDeque_IV(10)
	iv1 := _IV{1, 1.0}
	iv2 := _IV{2, 2.0}

	deque.PushFront(iv1.i, iv1.v)
	deque.PushBack(iv2.i, iv2.v)

	front := deque.FrontItem()
	if *front != iv1 {
		t.Errorf("got: %v, want: %v", front, iv1)
	}
	back := deque.BackItem()
	if *back != iv2 {
		t.Errorf("got: %v, want: %v", back, iv2)
	}
}

//// benchmarks ////////////////////////////////////////////////////////////////

func BenchmarkReference(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MovingMinMax_offline(W)
	}
}

func BenchmarkMovingMinMax0(b *testing.B) {
	minmax := NewMovingMinMax0(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			minmax.Update(values[i])
		}
	}
}

func BenchmarkMovingMinMax(b *testing.B) {
	minmax := NewMovingMinMax(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			minmax.Update(values[i])
		}
	}
}

func BenchmarkMovingMin(b *testing.B) {
	mmin := NewMovingMin(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			mmin.Update(values[i])
		}
	}
}

func BenchmarkMovingMax(b *testing.B) {
	mmax := NewMovingMax(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			mmax.Update(values[i])
		}
	}
}

func BenchmarkMovingMean(b *testing.B) {
	mmean := NewMovingMean(W)
	for j := 0; j < b.N; j++ {
		for i := range values {
			mmean.Update(values[i])
		}
	}
}

//// reference code ////////////////////////////////////////////////////////////

func MovingMinMax_offline(w int) {
	U := deque_int.New() // upper -> indices of maxima
	L := deque_int.New() // lower -> indices of minima

	// initial minimum and maximum
	mmins[0] = values[0]
	mmaxs[0] = values[0]

	for i := 1; i < len(values); i++ {
		if i < w {
			// 'absolute' minimum and maximum
			if values[i] < mmins[i-1] {
				mmins[i] = values[i]
			} else {
				mmins[i] = mmins[i-1]
			}
			if values[i] > mmaxs[i-1] {
				mmaxs[i] = values[i]
			} else {
				mmaxs[i] = mmaxs[i-1]
			}
		} else {
			// 'moving' minimum and maximum
			mmins[i-1] = values[iExtreme(L, i)]
			mmaxs[i-1] = values[iExtreme(U, i)]
		}

		// update monotonic wedge
		if values[i] > values[i-1] {
			L.PushBack(i - 1)
			if i == w+L.FrontItem() {
				L.PopFront()
			}
			for U.Size() > 0 {
				if values[i] <= values[U.BackItem()] {
					if i == w+U.FrontItem() {
						U.PopFront()
					}
					break
				}
				U.PopBack()
			}
		} else {
			U.PushBack(i - 1)
			if i == w+U.FrontItem() {
				U.PopFront()
			}
			for L.Size() > 0 {
				if values[i] >= values[L.BackItem()] {
					if i == w+L.FrontItem() {
						L.PopFront()
					}
					break
				}
				L.PopBack()
			}
		}
	}
	// final minimum and maximum
	i := len(values) - 1
	mmins[i] = values[iExtreme(L, i)]
	mmaxs[i] = values[iExtreme(U, i)]
}

func iExtreme(d *deque_int.Deque, i int) int {
	if d.Size() > 0 {
		return d.FrontItem()
	} else {
		return i - 1
	}
}
