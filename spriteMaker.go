package main

import (
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type Element struct {
	name string
	img  image.Image
}

var (
	white color.Color = color.RGBA{255, 255, 255, 255}
	black color.Color = color.RGBA{0, 0, 0, 255}
	blue  color.Color = color.RGBA{0, 0, 255, 255}
)

// ref) http://golang.org/doc/articles/image_draw.html
func main() {

	fmt.Println("start...")
	fSrc, err := os.Open("buttons")
	if err != nil {
		log.Fatal(err)
	}
	files, err := fSrc.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("stars contains", len(files), " files")

	images := list.New()
	//images := make([]image.Image, len(files))
	elem := Element{}
	maxH, maxW := 0, 0
	for i := range files {
		//fmt.Println(files[i].Name())
		imgSrc, err := os.Open(fmt.Sprintf("%s/%s", "buttons", files[i].Name()))
		defer imgSrc.Close()
		img, _, err := image.Decode(imgSrc)
		if err != nil {
			continue
			//fmt.Println(imgForm)
			//log.Fatal(err)
		}
		elem.name = strings.Split(files[i].Name(), ".")[0]
		elem.img = img
		images.PushFront(elem)
		//images[i] = img
		maxW += img.Bounds().Size().X
		if img.Bounds().Size().Y > maxH {
			maxH = img.Bounds().Size().Y
		}
	}
	fmt.Printf("%d x %d", maxW, maxH)
	fmt.Printf("opening %s\n", "stars/2_stern1.gif")
	fSrc, err = os.Open("stars/2_stern1.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer fSrc.Close()

	if err != nil {
		log.Fatal(err)
	}

	//buf := make([]byte, 1024)
	cssFile, err := os.Create("style.css")
	if err != nil {
		log.Fatal(err)
	}

	defer cssFile.Close()
	m := image.NewRGBA(image.Rect(0, 0, maxW, maxH)) //*NRGBA (image.Image interface)
	curX := 0
	for e := images.Front(); e != nil; e = e.Next() {
		p := image.Point{curX, 0}
		elem = e.Value.(Element)
		img := elem.img
		sr := img.Bounds()
		fmt.Fprintf(cssFile, ".%s\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(buttons.png) -%dpx 0px;\n }\n \n", elem.name, sr.Size().X, sr.Size().Y, curX)
		curX += sr.Size().X
		r := image.Rectangle{p, p.Add(sr.Size())}
		draw.Draw(m, r, img, sr.Min, draw.Src)
	}

	w, _ := os.Create("buttons.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.
}
