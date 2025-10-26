package main

import "fmt"

type Detection struct {
	Ra  float64
	Dec float64
}

func setRa(ra float64, detection Detection) Detection {
	detection.Ra = ra
	return detection
}

func main() {
	detection := Detection{Ra: 1.0, Dec: 2.0}
	newDetection := setRa(99.0, detection)

	fmt.Println(detection.Ra)
	fmt.Println(newDetection.Ra)
	fmt.Println(newDetection.Dec)
}
