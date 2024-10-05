package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"

	hp "github.com/dirodriguezm/healpix"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go [ang2pix|conesearch] Nside Iterations")
		os.Exit(1)
	}
	_nside, _ := strconv.ParseInt(os.Args[2], 10, 64)
	nside := int(_nside)
	_iterations, _ := strconv.ParseInt(os.Args[3], 10, 64)
	iterations := int(_iterations)
	if os.Args[1] == "ang2pix" {
		testAng2pix(nside, iterations)
	}
	if os.Args[1] == "conesearch" {
		testConesearch(nside, iterations)
	}
	if os.Args[1] == "conesearch2" {
		testConesearch2(nside, iterations)
	}
	if os.Args[1] == "conesearch2-goroutine" {
		testConesearch2Goroutine(nside, iterations)
	}
	if os.Args[1] == "conesearch-append" {
		testConesearchWithAppend(nside, iterations)
	}
	if os.Args[1] == "conesearch-goroutine" {
		testConesearchGoroutine(nside, iterations)
	}
	if os.Args[1] == "conesearch-single" {
		testConesearchSingleOp(nside, iterations)
	}
	if os.Args[1] == "conesearch-single-goroutine" {
		testConesearchSingleOpGoroutine(nside, iterations)
	}
}

func testAng2pix(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	results := make([]int64, 360*180*iterations, 360*180*iterations)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				point := hp.RADec(float64(ra), float64(dec))
				results[i*(360*180)+(90+dec)] = mapper.PixelAt(point)
			}
		}
	}
	fmt.Printf("Results: %d\n", len(results))
	fmt.Println(results[0:10])
}

func testConesearch(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	results := make([][]hp.PixelRange, 360*180*iterations, 360*180*iterations)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				point := hp.RADec(float64(ra), float64(dec))
				results[i*(360*180)+(90+dec)] = mapper.QueryDiscInclusive(point, radius, 4)
			}
		}
	}
	fmt.Printf("Results: %d\n", len(results))
	fmt.Println(results[0])
}

func testConesearch2(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				point := hp.RADec(float64(ra), float64(dec))
				mapper.QueryDiscInclusive(point, radius, 4)
			}
		}
	}
}

func testConesearch2Goroutine(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	nGoroutines := 2
	ch := make(chan int8, nGoroutines)
	var wg sync.WaitGroup
	radius := arcsecToRadians(10)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				wg.Add(1)
				ch <- 1
				go func() {
					defer func() {
						wg.Done()
						<-ch
					}()
					point := hp.RADec(float64(ra), float64(dec))
					mapper.QueryDiscInclusive(point, radius, 4)
				}()
			}
		}
	}
}

func testConesearchGoroutine(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	results := make([][]hp.PixelRange, 360*180*iterations, 360*180*iterations)
	var wg sync.WaitGroup
	maxGoroutines := 4
	limiter := make(chan int8, maxGoroutines)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				wg.Add(1)
				limiter <- 1
				go func() {
					defer func() {
						wg.Done()
						<-limiter
					}()
					point := hp.RADec(float64(ra), float64(dec))
					results[i*(360*180)+(90+dec)] = mapper.QueryDiscInclusive(point, radius, 4)
				}()
			}
		}
	}
	wg.Wait()
	fmt.Printf("Results: %d\n", len(results))
	fmt.Println(results[0])
}

func testConesearchWithAppend(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	results := make([][]hp.PixelRange, 0)
	for i := 0; i < iterations; i++ {
		for ra := 0; ra < 360; ra++ {
			for dec := -90; dec < 90; dec++ {
				point := hp.RADec(float64(ra), float64(dec))
				results = append(results, mapper.QueryDiscInclusive(point, radius, 4))
			}
		}
	}
	fmt.Printf("Results: %d\n", len(results))
	fmt.Println(results[0:10])
}

func testConesearchSingleOp(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	for i := 0; i < iterations; i++ {
		ra := rand.Intn(361)
		dec := rand.Intn(91)
		point := hp.RADec(float64(ra), float64(dec))
		mapper.QueryDiscInclusive(point, radius, 4)
	}
}

func testConesearchSingleOpGoroutine(nside, iterations int) {
	mapper, err := hp.NewHEALPixMapper(nside, hp.Nest)
	if err != nil {
		panic(err)
	}
	radius := arcsecToRadians(10)
	ch := make(chan int8, 10)
	var wg sync.WaitGroup
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		ch <- 1
		go func() {
			defer func() {
				wg.Done()
				<-ch
			}()
			ra := rand.Intn(361)
			dec := rand.Intn(91)
			point := hp.RADec(float64(ra), float64(dec))
			mapper.QueryDiscInclusive(point, radius, 4)
		}()
	}
}

func arcsecToRadians(arcsec float64) float64 {
	return (arcsec / 3600) * (math.Pi / 180)
}
