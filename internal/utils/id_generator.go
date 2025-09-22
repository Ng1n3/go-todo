package utils

import (
	"math/rand"
	"sync"
	"time"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var (
  rng *rand.Rand
  once sync.Once
)

// generate the random number once
func initRNG() {
  rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}


func GenerateID(length int) string {
  once.Do(initRNG)

  b := make([]byte, length)
  for i := range b {
    b[i] = letters[rng.Intn(len(letters))]
  }
  return string(b)
}


