package main

import (
	"fmt"
	"sync"
)

type Task struct {
	ID     int
	Result string
	f      func() string
}

func NewTask(id int, f func() string) *Task {
	return &Task{ID: id, f: f}
}

func (t *Task) Run(wg *sync.WaitGroup) {
	t.Result = t.f()
	wg.Done()
}

type Pool struct {
	Tasks []*Task

	concurrency int
	taskChan    chan *Task
	wg          sync.WaitGroup
}

func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:       tasks,
		concurrency: concurrency,
		taskChan:    make(chan *Task),
	}
}

func (p *Pool) Run() {
	for workerID := 0; workerID < p.concurrency; workerID++ {
		go p.work(workerID)
	}

	p.wg.Add(len(p.Tasks))
	for _, task := range p.Tasks {
		p.taskChan <- task
	}

	close(p.taskChan)

	p.wg.Wait()
}

func (p *Pool) work(workerID int) {
	for task := range p.taskChan {
		task.Run(&p.wg)
		fmt.Println("Worker #", workerID, " processed task #", task.ID, " with result ", task.Result)
	}
}
