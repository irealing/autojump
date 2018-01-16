package main

import "image/color"

type RGB struct {
	R, G, B int
}

func toRGB(c color.Color) *RGB {
	rgba := toRGBA(c)
	return &RGB{int(rgba.R), int(rgba.G), int(rgba.B)}
}
