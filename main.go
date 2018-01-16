package main

import (
	"log"
	"time"
)

func main() {
	for {
		img, err := screenShot()
		if err != nil {
			log.Fatal(err)
		}
		px, py, err := findPoints(img)
		if err != nil {
			log.Fatal(err)
		}
		err = jump(px, py)
		if err != nil {
			log.Fatal(err)
		}
		st := genRandI64(900, 1500)
		d := time.Duration(st) * time.Millisecond
		log.Print("sleep ", st, d)
		time.Sleep(d)
	}
}
