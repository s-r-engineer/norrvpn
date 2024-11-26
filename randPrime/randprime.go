package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const defaultBits = 4096
const defaultBase = 62

func generateLargePrime() (*big.Int, error) {
	return rand.Prime(rand.Reader, defaultBits)
}

func encodeLargePrime(b *big.Int) string {
	return b.Text(defaultBase)
}

func parseLargePrime(s string) *big.Int {
	b := new(big.Int)
	b.SetString(s, defaultBase)
	return b
}

func main() {
	s, _ := generateLargePrime()
	fmt.Print(s)
}
