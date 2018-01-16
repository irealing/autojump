package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const (
	adbPATH    = "lib"
	screenSize = `(?P<width>\d+)x(?P<height>\d+)`
)

var (
	regexpScreen = regexp.MustCompile(screenSize)
	cfg          = &Config{}
)

func init() {
	cwd, err := initPATH()
	if err != nil {
		log.Fatal(err)
	}
	cfg.BasePATH = cwd
	h, w, err := deviceSize()
	if err != nil {
		log.Fatal(err)
	}
	jf := fmt.Sprintf("config/%dx%d/config.json", h, w)
	jf = filepath.Join(cwd, jf)
	if _, err = os.Stat(jf); err != nil {
		jf = filepath.Join(cwd, "config", "default.json")
	}
	cfg.Height, cfg.Width = h, w
	cf, err := os.Open(jf)
	if err != nil {
		log.Fatal(err)
	}
	defer cf.Close()
	var buf []byte
	if buf, err = ioutil.ReadAll(cf); err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(buf, cfg)
}

//initPATH 初始化PATH
//添加adb到PATH
func initPATH() (string, error) {
	log.Println("初始化PATH")
	filename, e := os.Executable()
	if e != nil {
		return "", e
	}
	cwd := filepath.Dir(filename)
	path := os.Getenv("PATH")
	path = fmt.Sprintf("%s;%s", path, filepath.Join(cwd, adbPATH))
	os.Setenv("PATH", path)
	return cwd, nil
}

func deviceSize() (int, int, error) {
	r, err := runCMD("adb", "shell", "wm", "size")
	log.Println(r)
	if err != nil {
		return -1, -1, err
	}
	si := regexpScreen.FindStringSubmatch(r)
	if si != nil && len(si) != 3 {
		log.Println(si)
		return -1, -1, errors.New("获取屏幕尺寸异常")
	}
	w, _ := strconv.Atoi(si[1])
	h, _ := strconv.Atoi(si[2])
	return h, w, nil
}

type Config struct {
	BasePATH            string  `json:"-"`
	Height              int     `json:"-"`
	Width               int     `json:"-"`
	UnderGameScoreY     int     `json:"under_game_score_y"`
	PressCoefficient    float64 `json:"press_coefficient"`
	PieceBaseHeightHalf int     `json:"piece_base_height_1_2"`
	PieceBodyWidth      int     `json:"piece_body_width"`
}
