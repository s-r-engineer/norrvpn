package main

import "testing"

func TestRandPrime(t *testing.T) {
	prime, _ := generateLargePrime()
	encodedPrime := encodeLargePrime(prime)
	newPrime := parseLargePrime(encodedPrime)
	if prime.Cmp(newPrime) != 0 {
		t.Fail()
	}
}
