package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, customerNumber, total int

type Producer struct {
	curRequest chan int
	data       chan PizzaOrder
	quit       chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		customerNumber++
		fmt.Printf("Receive order #%d!\n", pizzaNumber)
		delay := rand.Intn(5) + 1
		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d. It will take %d seconds...\n", pizzaNumber, delay)
		// delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 1 {
			msg = fmt.Sprintf("*** Failure reason #A for pizza #%d!", pizzaNumber)
		} else if rnd <= 2 {
			msg = fmt.Sprintf("*** Failure reason #B for pizza #%d!", pizzaNumber)
		} else if rnd <= 3 {
			msg = fmt.Sprintf("*** Failure reason #C for pizza #%d!", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** Failure reason #D for pizza #%d!", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza Order #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
		message:     fmt.Sprintf("*** All Pizza resources are busy to service #%d!", pizzaNumber),
	}
}

func pizzaShop(pizzaMaker *Producer) {
	i := 0
	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}

	}
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out message of pizza business opening
	color.Cyan("The Pizza Shop is open for business!")
	color.Cyan("------------------------------------")

	// create a producer
	pizzaJob := &Producer{
		curRequest: make(chan int),
		data:       make(chan PizzaOrder),
		quit:       make(chan chan error),
	}

	// run the producer
	go pizzaShop(pizzaJob)

	// create and run the consumer
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green(" Customer is served with Order #%d!", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer with order #%d is really mad and leaving!", i.pizzaNumber)
			}
		} else {
			color.Yellow("The Pizza shop is done making pizza!!")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error closing channel!", err)
			}
		}
	}

	// Print the Ending message
	color.Cyan("-----------------")
	color.Cyan("Done for the day.")

	color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts in total.", pizzasMade, pizzasFailed, total)
}
