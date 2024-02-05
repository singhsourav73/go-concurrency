package main

import (
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Philosopher is a struct which stores information about a philosopher
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// list out all 5 philosopher
var philosophers = []Philosopher{
	{name: "Debojit", leftFork: 4, rightFork: 0},
	{name: "Ankit", leftFork: 0, rightFork: 1},
	{name: "Soumit", leftFork: 1, rightFork: 2},
	{name: "Sanket", leftFork: 2, rightFork: 3},
	{name: "Pritam", leftFork: 3, rightFork: 4},
}

// Variable definition used during dining
var (
	hunger    = 3
	eatTime   = 3 * time.Second
	thinkTime = 1 * time.Second
	sleepTime = 1 * time.Second
)

// variable to display the order in which philosopher finish dining
var (
	orderMutex    sync.Mutex
	orderFinished []string
)

func main() {
	// print out welcome message
	color.Cyan("Dining Philosophers Problem!!")
	color.Cyan("-----------------------------")
	color.Cyan("Dining started!! The table is empty...")

	// Initial Sleep time
	time.Sleep(sleepTime)

	// Start the meal
	dine()

	// print out finish dining
	color.Cyan("Dining ended!! The table is empty!!")

	// Print the order of meal completion
	time.Sleep(sleepTime)
	color.Cyan("Order finished: %s.\n", strings.Join(orderFinished, ","))
}

func dine() {
	// eatTime = 0 * time.Second
	// sleepTime = 0 * time.Second
	// thinkTime = 0 * time.Second

	// wg is the WaitGroup that keeps track of how many philosophers are still at the table. When
	// it reaches zero, everyone is finished eating and has left. We add 5 (the number of philosophers) to this
	// wait group.
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// We want everyone to be seated before they start eating, so create a WaitGroup for that.
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks. Forks are assigned using the fields leftFork and rightFork in the Philosopher
	// type. Each fork, then, can be found using the index (an integer), and each fork has a unique mutex.
	forks := make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	for i := 0; i < len(philosophers); i++ {
		go dinningProblem(philosophers[i], wg, forks, seated)
	}

	// Wait for the philosophers to finish. This blocks until the wait group is 0.
	wg.Wait()
}

// diningProblem is the function fired off as a goroutine for each of our philosophers. It takes one
// philosopher, our WaitGroup to determine when everyone is done, a map containing the mutexes for every
// fork on the table, and a WaitGroup used to pause execution of every instance of this goroutine
// until everyone is seated at the table.
func dinningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	color.Yellow("%s is seated at the table.\n", philosopher.name)
	seated.Done() // Marks the philosopher as seated
	seated.Wait() // Wait till everyone is seated

	for i := hunger; i > 0; i-- {

		// Get a lock on the left and right forks. We have to choose the lower numbered fork first in order
		// to avoid a logical race condition, which is not detected by the -race flag in tests; if we don't do this,
		// we have the potential for a deadlock, since two philosophers will wait endlessly for the same fork.
		// Note that the goroutine will block (pause) until it gets a lock on both the right and left forks.
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			color.Yellow("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			color.Yellow("\t%s takes the right fork.\n", philosopher.name)
		}

		color.White("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)
		color.White("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)

		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()
		color.White("\t%s put down the fork.\n", philosopher.name)
	}

	color.Green("%s is done eating", philosopher.name)
	color.Green("%s left the table", philosopher.name)

	orderMutex.Lock()
	orderFinished = append(orderFinished, philosopher.name)
	orderMutex.Unlock()
}
