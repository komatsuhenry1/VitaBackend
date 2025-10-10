package utils

import (
	"math/rand"
	"time"
	"fmt"
	"log"
)

func GenerateAuthCode()(int, error){
	min := 100000
	max := 999999

	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(max-min+1) + min

	fmt.Println("=========")
	fmt.Println("2 factor code: ", num)
	fmt.Println("=========")
	log.Println("=========")
	log.Println("2 factor code: ", num)
	log.Println("=========")

	return num, nil
}