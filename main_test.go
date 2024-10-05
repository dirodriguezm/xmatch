package main

import "testing"

func TestConesearch(t *testing.T) {
	nside := 18
	iterations := 1
	testConesearch(nside, iterations)
}

func TestConesearchSingle(t *testing.T) {
	nside := 18
	iterations := 100000
	testConesearchSingleOp(nside, iterations)
}
