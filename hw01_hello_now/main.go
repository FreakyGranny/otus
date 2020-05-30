package main

import (
	"fmt"
	"log"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	fmt.Printf("current time: %s\n", time.Now().Round(time.Second))

	ntpTime, err := ntp.Time("0.ru.pool.ntp.org")
	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Printf("exact time: %s\n", ntpTime.Round(time.Second))
}
