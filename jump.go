package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var (
	errorScan = errors.New("找不到对象")
)

func screenShot() (string, error) {
	log.Print("adb shell screencap -p /sdcard/jump.png")
	output, err := runCMD("adb", "shell", "screencap", "-p", "/sdcard/jump.png")
	if err != nil {
		return "", err
	}
	log.Print(output)
	png := fmt.Sprintf("jump_%d.png", time.Now().Unix())
	png = filepath.Join(cfg.BasePATH, "img", png)
	log.Printf("adb pull /sdcard/jump.png %s", png)
	output, err = runCMD("adb", "pull", "/sdcard/jump.png", png)
	if err != nil {
		return "", err
	}
	return png, nil
}
func calDistance(f, t *image.Point) float64 {
	fs := f.X - t.X
	ts := f.Y - t.Y
	return math.Sqrt(float64(ts*ts + fs*fs))
}
func jump(from, to *image.Point) error {
	distance := calDistance(from, to)
	log.Printf("jump %s -> %s (%f)", from, to, distance)
	//log.Printf("jump : %d", distance)
	pt := distance * cfg.PressCoefficient
	if pt < 200 {
		pt = 200
	}
	t := int(pt)
	x1, y1 := strconv.Itoa(from.X), strconv.Itoa(from.Y)
	x2, y2 := strconv.Itoa(to.X), strconv.Itoa(to.Y)
	ts := strconv.Itoa(t)
	log.Printf("adb shell input swipe %s %s %s %s %s", x1, y1, x2, y2, ts)
	r, err := runCMD("adb", "shell", "input", "swipe", x1, y1, x2, y2, ts)
	if err != nil {
		log.Print("jump error : ", err)
		return err
	}
	log.Print(r)
	return err
}

func findPoints(source string) (*image.Point, *image.Point, error) {
	input, err := os.Open(source)
	if err != nil {
		return nil, nil, err
	}
	defer input.Close()
	imgCfg, err := png.DecodeConfig(input)
	if err != nil {
		return nil, nil, err
	}
	h, w := imgCfg.Height, imgCfg.Width
	input.Seek(0, 0)
	img, _ := png.Decode(input)
	start, end, step := h/3, h*2/3, 50
	var startY int
	for i := start; i < end; i += step {
		lastPx := toRGBA(img.At(0, i))
		for j := 1; j < w; j++ {
			pix := toRGBA(img.At(j, i))
			if !sameColor(lastPx, pix) {
				startY = i - step
				break
			}
		}
		if startY != 0 {
			break
		}
	}
	sxb := w / 8
	log.Printf("从%d开始扫描", startY)
	var pxs, pxc, pym int
	for i := startY; i < end; i++ {
		for j := sxb; j < w-sxb; j++ {
			pix := toRGB(img.At(j, i))
			if 50 < pix.R && pix.R < 60 && 53 < pix.G && pix.G < 60 && 95 < pix.B && pix.B < 110 {
				pxs += j
				pxc += 1
				pym = maxInt(i, pym)
			}
		}
	}
	if pxs == 0 || pxc == 0 {
		return nil, nil, errorScan
	}
	px := pxs / pxc
	py := pym - cfg.PieceBaseHeightHalf
	var bxs, bxe, bxc, bx, by int
	if px < w/2 {
		bxs = px
		bxe = w
	} else {
		bxs = 0
		bxe = px
	}
	start, end = h/3, h*2/3
	i := start
	for ; i < end; i++ {
		lp := toRGB(img.At(0, i))
		if bx != 0 || by != 0 {
			break
		}
		bxs, bxc = 0, 0
		for j := bxs; j < bxe; j++ {
			pix := toRGB(img.At(j, i))
			if absInt(j-px) < cfg.PieceBodyWidth {
				continue
			}
			if absInt(pix.R-lp.R)+absInt(pix.G-lp.G)+absInt(pix.B-lp.B) > 10 {
				bxs += j
				//log.Printf("bxs(%d)+=j(%d)", bxs, j)
				bxc += 1
			}
		}
		if bxs != 0 {
			bx = bxs / bxc
			//log.Printf("bx(%d)=bxs(%d)/bxc(%d)", bx, bxs, bxc)
		}
	}
	lp := toRGB(img.At(bx, i))
	var k int
	for k = i + 274; k > i; k -= 1 {
		pix := toRGB(img.At(bx, k))
		if absInt(pix.R-lp.R)+absInt(pix.G-lp.G)+absInt(pix.B-lp.B) < 10 {
			break
		}
	}
	by = (i + k) / 2
	for j := i; j < i+200; j++ {
		pix := toRGB(img.At(bx, j))
		if absInt(pix.R-245)+absInt(pix.G-245)+absInt(pix.B-245) == 0 {
			by = j + 10
			break
		}
	}
	if bx == 0 || by == 0 {
		return nil, nil, errorScan
	}
	return &image.Point{px, py}, &image.Point{bx, by}, nil
}

func btnPosition() (*image.Point, *image.Point, error) {
	h, w, err := deviceSize()
	if err != nil {
		return nil, nil, err
	}
	l := w / 2
	t := int(1584 * (h / 1920.0))
	return &image.Point{l, t}, &image.Point{l, t}, err
}
