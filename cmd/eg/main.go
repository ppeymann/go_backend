package main

import (
	"fmt"
	"time"
)

func main() {
	now := time.Now().UTC()

	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	start := time.Date(now.Year(), now.Month(), now.Day(), 7, 35, 0, 0, time.UTC).Unix()
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 30, 0, 0, time.UTC).Unix()

	fmt.Println("date:", base, "start:", start, "end:", end)
}
