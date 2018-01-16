package main

import (
	"bytes"
	"image/color"
	"math/rand"
	"os/exec"
	"time"
)

//runCMD 执行命令
func runCMD(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	err := cmd.Run()
	return buf.String(), err
}
func sameColor(x, y *color.RGBA) bool {
	return x.R == y.R && x.G == y.G && x.B == y.B
}
func toRGBA(c color.Color) *color.RGBA {
	if rgba, ok := c.(color.RGBA); ok {
		return &rgba
	}
	tc := color.RGBAModel.Convert(c)
	rgba := tc.(color.RGBA)
	return &rgba
}
func absInt(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func maxInt(i, j int) int {
	t := i
	if j > t {
		t = j
	}
	return t
}
func genRandI64(min, max int64) int64 {
	seed := time.Now().Unix()
	rand.Seed(seed)
	return rand.Int63n(max-min) + min
}
