package utils

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func GenerateAuthCode() (int, error) {
	min := 100000
	max := 999999

	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(max-min+1) + min

	fmt.Println("=========")
	fmt.Println(" code: ", num)
	fmt.Println("=========")
	log.Println("=========")
	log.Println("code: ", num)
	log.Println("=========")

	return num, nil
}
