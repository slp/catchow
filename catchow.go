package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var hashes map[string]int
const letters = "TRWAGMYFPDXBNJZSQVHLCKE"
const ncpus = 4

func calcHash(i int, letter rune, date string, zip string) string {
	key := fmt.Sprintf("%05d%c%s%s", i, letter, date, zip)
	sum := []byte(key)

	for i := 0; i <= 1715; i++ {
		tmp := sha256.Sum256(sum)
		sum = []byte(fmt.Sprintf("%x", tmp[:]))
	}
	return fmt.Sprintf("%s", sum)
}

func genHashes(wg *sync.WaitGroup, start int, end int, year int, zip string) {
	defer wg.Done()

	for month := time.January; month <= 12; month++ {
		t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)

		for day := 1; day <= t.Day(); day++ {
			date := fmt.Sprintf("%d%02d%02d", year, month, day)

			for _, l := range letters {
				for i := start; i <= end; i++ {
					hash := calcHash(i, l, date, zip)
					if _, ok := hashes[hash]; ok {
						fmt.Printf("%05d-%c %s %s\n", i, l, date, zip)
						fmt.Println(hash)
					}
				}
			}
		}
	}
}

func main() {
	file, err := os.Open("all.db")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hashes = make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hashes[scanner.Text()[:64]] = 1
	}

	start := 0
	end := 99999
	step := 100000 / ncpus

	var wg sync.WaitGroup
	wg.Add(4)

	for i := 1; i <= 4; i++ {
		if i == ncpus {
			go genHashes(&wg, start, end, 1969, "08008")
		} else {
			next := start + step
			go genHashes(&wg, start, next, 1969, "08008")
			start = next
		}
	}

	wg.Wait()
}
