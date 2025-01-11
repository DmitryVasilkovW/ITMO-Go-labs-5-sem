//go:build !solution

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	height     = 12
	widthColon = 4
	width      = 8
)

func main() {
	portPtr := flag.Int("port", 8080, "port")

	flag.Parse()
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *portPtr), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	urlPath, err := buildURLPath(r)
	if err != nil {
		http.Error(w, "failed to parse URL", http.StatusBadRequest)
		return
	}

	timeRequest, err := getTimeRequest(urlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	kRequest, err := getKRequest(urlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	img := drawPicture(timeRequest, kRequest)

	if err := sendImageResponse(w, img); err != nil {
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
	}
}

func buildURLPath(r *http.Request) (string, error) {
	urlPath := "http://" + r.Host + r.URL.String()
	u, err := url.Parse(urlPath)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func getTimeRequest(urlPath string) (string, error) {
	u, _ := url.Parse(urlPath)
	q := u.Query()

	timeRequest := q.Get("time")
	if timeRequest == "" {
		timeRequest = time.Now().Format("15:04:05")
	}

	if len(timeRequest) != 8 || !isValidTimeFormat(timeRequest) {
		return "", fmt.Errorf("incorrect time format")
	}

	return timeRequest, nil
}

func isValidTimeFormat(timeStr string) bool {
	_, err := time.Parse("15:04:05", timeStr)
	return err == nil
}

func getKRequest(urlPath string) (int, error) {
	u, _ := url.Parse(urlPath)
	q := u.Query()

	kStr := q.Get("k")
	if kStr == "" {
		return 1, nil
	}

	k, err := strconv.Atoi(kStr)
	if err != nil || k < 1 || k > 30 {
		return 0, fmt.Errorf("invalid k")
	}

	return k, nil
}

func sendImageResponse(w http.ResponseWriter, img image.Image) error {
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	return png.Encode(w, img)
}

func getDigitPixels(digit int) string {
	switch digit {
	case 0:
		return Zero
	case 1:
		return One
	case 2:
		return Two
	case 3:
		return Three
	case 4:
		return Four
	case 5:
		return Five
	case 6:
		return Six
	case 7:
		return Seven
	case 8:
		return Eight
	case 9:
		return Nine
	default:
		return ""
	}
}

func drawDigit(img *image.RGBA, digit int32, xStart int, yStart int, k int) (*image.RGBA, int, int) {
	drawPixels(img, digit, xStart, yStart, k)
	x := xStart + width
	return img, x, 0
}

func drawPixels(img *image.RGBA, digit int32, xStart, yStart int, k int) {
	x := xStart
	y := yStart
	digitPixels := getDigitPixels(int(digit - '0'))
	for i, sampleItem := range digitPixels {
		if digitPixels[i] == 10 {
			continue
		}
		img = drawPixel(img, string(sampleItem), x*k, y*k, k)
		if x-xStart == width-1 {
			x = xStart
			y++
		} else {
			x++
		}
	}

}

func drawColon(img *image.RGBA, xStart int, yStart int, k int) (*image.RGBA, int, int) {
	drawImg(img, xStart, yStart, k)
	x := xStart + widthColon
	return img, x, 0
}

func drawImg(img *image.RGBA, xStart, yStart int, k int) {
	x := xStart
	y := yStart
	for i, sampleItem := range Colon {
		if Colon[i] == 10 {
			continue
		}
		img = drawPixel(img, string(sampleItem), x*k, y*k, k)
		if x-xStart == widthColon-1 {
			x = xStart
			y++
		} else {
			x++
		}
	}

}

func drawPixel(img *image.RGBA, sign string, xStart int, yStart int, k int) *image.RGBA {
	for y := yStart; y < yStart+k; y++ {
		for x := xStart; x < xStart+k; x++ {
			if sign == "1" {
				img.Set(x, y, Cyan)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}
	return img
}

func drawPicture(time string, k int) *image.RGBA {
	resultWidth := (2*widthColon + 6*width) * k
	resultHeight := height * k
	img := image.NewRGBA(image.Rect(0, 0, resultWidth, resultHeight))

	return getImg(time, img, k)
}

func getImg(time string, img *image.RGBA, k int) *image.RGBA {
	x := 0
	y := 0
	for _, value := range time {
		if value == ':' {
			img, x, y = drawColon(img, x, y, k)
		} else {
			img, x, y = drawDigit(img, value, x, y, k)
		}
	}

	return img
}
