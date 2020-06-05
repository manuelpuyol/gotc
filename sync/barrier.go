package sync

import "sync"

type Barrier struct {
	mutex   *sync.Mutex
	cond    *sync.Cond
	group   *sync.WaitGroup
	counter int
	Threads int
}

func NewBarrier(threads int) *Barrier {
	var mutex sync.Mutex
	var group sync.WaitGroup
	cond := sync.NewCond(&mutex)

	return &Barrier{
		mutex:   &mutex,
		group:   &group,
		cond:    cond,
		counter: 0,
		Threads: threads,
	}
}

func (b *Barrier) Start() {
	b.group.Add(b.Threads)
}

func (b *Barrier) Reset() {
	b.counter = 0
}

func (b *Barrier) Wait() {
	b.group.Wait()
}

func (b *Barrier) Done() {
	b.mutex.Lock()
	b.counter++
	if b.counter == b.Threads {
		b.cond.Broadcast()
	} else {
		for b.counter != b.Threads {
			b.cond.Wait()
		}
	}
	b.mutex.Unlock()

	b.group.Done()
}
