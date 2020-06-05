package queue

import (
	"gotc/transaction"
	"sync"
)

type Queue struct {
	mutex        *sync.Mutex
	cond         *sync.Cond
	transactions []*transaction.Transaction
	Size         int
}

func NewQueue(cond *sync.Cond) *Queue {
	var mutex sync.Mutex
	queue := Queue{mutex: &mutex, cond: cond}

	return &queue
}

func (queue *Queue) Enqueue(t *transaction.Transaction) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.transactions = append(queue.transactions, t)
	queue.Size++
}

func (queue *Queue) Dequeue() *transaction.Transaction {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if queue.Size == 0 {
		return nil
	}

	t := queue.transactions[0]
	queue.transactions = queue.transactions[1:]
	queue.Size--

	return t
}
