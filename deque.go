// deque.go, jpad 2015

package movingminmax

//// deque_IV //////////////////////////////////////////////////////////////////

// deque_IV is a bounded deque with a fixed capacity to store index+value pairs.
// It is implemented as a power-of-two sized ring buffer for efficient wrapping
// around of the front and back indices.
type deque_IV struct {
	items []*_IV
	front uint // front item index
	back  uint // back item index
	size  uint // number of items
	mask  uint // wrap mask (capacity - 1, where capacity is a power of 2)
}

func newDeque_IV(capacity uint) *deque_IV {
	capacity = nextPowerOfTwo(capacity + 1)
	items := make([]*_IV, capacity)
	for i := range items {
		items[i] = &_IV{}
	}
	return &deque_IV{
		items: items,
		mask:  capacity - 1,
	}
}

func (d *deque_IV) PushFront(index uint, value float64) {
	d.front--
	d.front &= d.mask
	d.items[d.front].i = index
	d.items[d.front].v = value
	d.size++
}

func (d *deque_IV) PushBack(index uint, value float64) {
	d.items[d.back].i = index
	d.items[d.back].v = value
	d.back++
	d.back &= d.mask
	d.size++
}

func (d *deque_IV) PopFront() *_IV {
	item := d.items[d.front]
	d.size--
	d.front++
	d.front &= d.mask
	return item
}

func (d *deque_IV) PopBack() *_IV {
	d.size--
	d.back--
	d.back &= d.mask
	return d.items[d.back]
}

func (d *deque_IV) PruneFront() {
	d.size--
	d.front++
	d.front &= d.mask
}

func (d *deque_IV) PruneBack() {
	d.size--
	d.back--
	d.back &= d.mask
}

func (d *deque_IV) FrontItem() *_IV {
	return d.items[d.front]
}

func (d *deque_IV) BackItem() *_IV {
	return d.items[(d.back-1)&d.mask]
}

func (d *deque_IV) Size() uint {
	return d.size
}

//// deque_f32 /////////////////////////////////////////////////////////////////

// deque_f32 is a bounded deque with a fixed capacity to store float64 values.
// It is implemented as a power-of-two sized ring buffer for efficient wrapping
// around of the front and back indices.
type deque_f32 struct {
	items []float64
	front uint // front item index
	back  uint // back item index
	size  uint // number of items
	mask  uint // wrap mask (capacity - 1, where capacity is a power of 2)
}

func newDeque_f32(capacity uint) *deque_f32 {
	capacity = nextPowerOfTwo(capacity + 1)
	return &deque_f32{
		items: make([]float64, capacity),
		mask:  capacity - 1,
	}
}

func (d *deque_f32) PushFront(value float64) {
	d.front--
	d.front &= d.mask
	d.items[d.front] = value
	d.size++
}

func (d *deque_f32) PushBack(value float64) {
	d.items[d.back] = value
	d.back++
	d.back &= d.mask
	d.size++
}

func (d *deque_f32) PopFront() float64 {
	item := d.items[d.front]
	d.size--
	d.front++
	d.front &= d.mask
	return item
}

func (d *deque_f32) PopBack() float64 {
	d.size--
	d.back--
	d.back &= d.mask
	return d.items[d.back]
}

func (d *deque_f32) PruneFront() {
	d.size--
	d.front++
	d.front &= d.mask
}

func (d *deque_f32) PruneBack() {
	d.size--
	d.back--
	d.back &= d.mask
}

func (d *deque_f32) FrontItem() float64 {
	return d.items[d.front]
}

func (d *deque_f32) BackItem() float64 {
	return d.items[(d.back-1)&d.mask]
}

func (d *deque_f32) Size() uint {
	return d.size
}

//// utilities /////////////////////////////////////////////////////////////////

// nextPowerOfTwo returns a power of two that is larger or equal than x.
func nextPowerOfTwo(x uint) uint {
	result := uint(1)
	for result < x {
		result <<= 1
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////

// To implement the algorithm from Lemire 2006, we can use a fixed amount of
// memory and minimal memory allocation. This custom queue serves this purpose.
type intfloatqueue struct {
	nodes []*intfloatnode
	head  uint
	tail  uint
	mask  uint
}

func newintfloatqueue(size uint) *intfloatqueue {
	size = nextPowerOfTwo(size + 1)
	n := make([]*intfloatnode, size)
	for i := range n {
		n[i] = &intfloatnode{}
	}
	return &intfloatqueue{
		nodes: n,
		mask:  size - 1,
	}
}

func (q *intfloatqueue) empty() bool {
	return q.tail == q.head
}

func (q *intfloatqueue) nonempty() bool {
	return q.tail != q.head
}

func (q *intfloatqueue) count() uint {
	return (q.tail - q.head) & q.mask
}

func (q *intfloatqueue) push(index uint, value float64) {
	q.nodes[q.tail].index = index
	q.nodes[q.tail].value = value
	q.tail = (q.tail + 1) & q.mask
}

func (q *intfloatqueue) pop() *intfloatnode {
	node := q.nodes[q.head]
	q.head = (q.head + 1) & q.mask
	return node
}

func (q *intfloatqueue) tailnode() *intfloatnode {
	return q.nodes[(q.tail-1)&q.mask]
}

func (q *intfloatqueue) tailvalue() float64 {
	return q.nodes[(q.tail-1)&q.mask].value
}

func (q *intfloatqueue) poptail() *intfloatnode {
	q.tail = (q.tail - 1) & q.mask
	node := q.nodes[q.tail]
	return node
}

func (q *intfloatqueue) prunetail() {
	q.tail = (q.tail - 1) & q.mask
}

func (q *intfloatqueue) prunehead() {
	q.head = (q.head + 1) & q.mask
}

func (q *intfloatqueue) headnode() *intfloatnode {
	return q.nodes[q.head]
}

func (q *intfloatqueue) headindex() uint {
	return q.nodes[q.head].index
}

func (q *intfloatqueue) headvalue() float64 {
	return q.nodes[q.head].value
}

type intfloatnode struct {
	index uint
	value float64
}

func newintfloatnode(i uint, v float64) *intfloatnode {
	return &intfloatnode{
		index: i,
		value: v,
	}
}
