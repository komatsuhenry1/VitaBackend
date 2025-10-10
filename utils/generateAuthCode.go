package utils

import (
	"math/rand"
	"time"
	"fmt"
)

func GenerateAuthCode()(int, error){
	min := 100000
	max := 999999

	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(max-min+1) + min

	fmt.Println("=========")
	fmt.Println("2 factor code: ", num)
	fmt.Println("=========")

	return num, nil
}