# go-concurrency
Sample code to demonstrate concurrent programming implementing Dijkstra's solution

# Description

Go, often referred to as Golang, is well-known for making it remarkably easy to work with concurrency. In order to make a particular function run concurrently, all we have to do is prepend the word "go" to the function call, and it cheerfully runs in the background, as a GoRoutine. Go's built in scheduler takes are of making sure that a given GoRoutine runs when it should, and as efficiently as it can.

However, this does not mean that working with concurrency is simple in Go—thread safe programming takes careful planning, and most importantly it requires that developers have an absolutely solid understanding of how Go deals with concurrency.

In the standard library, Go offers us several ways of dealing with concurrently running parts of our program, right in the standard library: sync.WaitGroup, which lets us wait for tasks to finish; sync.Mutex, which allows us to lock and unlock resources, so that no two GoRoutines can access the same memory location at the same time; and finally, Channels, which allow GoRoutines to send and receive data to and from each other.

Go's approach to concurrency is fairly straightforward, and is more or less summed up this mantra: Don't communicate by sharing memory; instead, share memory by communicating. Channels are the means by which we usually share memory by communicating.

This repo try to cover concurrency by solving some of the classic problem including the Dining Philosophers, the Producer/Consumer problem, and the Sleeping Barber.
