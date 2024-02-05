package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

// Declare some variable to define time spend on activities
// like cutting, waiting, etc
var (
	seatingCapacity = 10
	arrivalRate     = 100
	cutDuration     = 1000 * time.Millisecond
	timeOpen        = 10 * time.Second
)

func main() {
	// seed a random number generator
	rand.Seed(time.Now().UnixNano())

	// print welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels for client and completion flag
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create a barber shop
	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarberDoneChan:  doneChan,
		Open:            true,
	}
	color.Green("The Shop is open for the day!")

	// Add barber/barbers
	shop.addBarber("Saha")
	shop.addBarber("Aditya")
	shop.addBarber("Khetan")

	// add client
	shopClosing := make(chan bool)
	closed := make(chan bool)

	// block until the barbershop is closed
	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <- true
	}()

	i := 1

	go func() {
		for {
			randomMilliSecond := rand.Int() % (2 * arrivalRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliSecond)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()

	<-closed
}
