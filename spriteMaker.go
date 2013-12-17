package main

import (
	"code.google.com/p/gcfg"
	"container/list"
	"fmt"
	"image"
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

func stringInSlice(a string, list []string) (bool, int) {
	for i, b := range list {
		if b == a {
			return true, i
		}
	}
	return false, 0
}

func my_misc() {
	for i := 37; i < 46; i++ {
		/*[image "zoover4"]
		name = zoover-008
		name = zoover-009
		width = 51
		height = 9
		startX = 0
		startY =  9*/
		s := fmt.Sprintf("[image \"zoover%d\"]\nname = zoover-0%d\nname = zoover-0%d\nwidth = 51\nheight = 9\nstartX = 0\nstartY =  %d\n\n", i, i*2+7, i*2+8, i*9)
		fmt.Println(s)
	}

}

func main() {
	cfg := struct {
		Files struct {
			Name  []string
			Count []int
			Size  []int
			MaxH  []int
		}
		Image map[string]*struct {
			Name   []string
			Width  int
			Height int
			StartX int
			StartY int
			Hover  bool
		}
	}{}
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")

	dirSrc := "buttons"

	fmt.Println("start...")
	fSrc, err := os.Open(dirSrc)
	if err != nil {
		log.Fatal(err)
	}
	files, err := fSrc.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dirSrc, " contains", len(files), " files")

	images := list.New()
	elem := Element{}
	maxH, maxW := 0, 0
	for i := range files {
		imgSrc, err := os.Open(fmt.Sprintf("%s/%s", dirSrc, files[i].Name()))
		defer imgSrc.Close()
		img, _, err := image.Decode(imgSrc)
		if err != nil {
			continue
		}
		elem.name = strings.Split(files[i].Name(), ".")[0]
		elem.img = img
		images.PushFront(elem)
		if ok, c := stringInSlice(elem.name, cfg.Files.Name); ok {
			maxW += cfg.Files.Size[c]
			if cfg.Files.MaxH[c] > maxH {
				maxH = cfg.Files.MaxH[c]
			}
		} else {
			maxW += img.Bounds().Size().X
			if img.Bounds().Size().Y > maxH {
				maxH = img.Bounds().Size().Y
			}
		}

	}
	fmt.Printf("%d x %d\n", maxW, maxH)
	if err != nil {
		log.Fatal(err)
	}
	cssFile, err := os.Create("goStyle.css")
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
		if ok, c := stringInSlice(elem.name, cfg.Files.Name); ok {
			for i := 0; i < cfg.Files.Count[c]; i++ {
				spritePart := cfg.Image[fmt.Sprintf("%s%d", elem.name, i)]
				for id, n := range spritePart.Name {
					if spritePart.Hover {
						spritePart.Name[id] = fmt.Sprintf("a.%s:hover", n)
					} else {
						spritePart.Name[id] = fmt.Sprintf(".%s", n)
					}
				}
				fmt.Fprintf(cssFile, "%s\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(goSprite.png) -%dpx 0px;\n }\n \n", strings.Join(spritePart.Name, ", "), spritePart.Width, spritePart.Height, curX)
				curX += spritePart.Width
				r := image.Rectangle{p, p.Add(image.Point{spritePart.Width, spritePart.Height})}
				draw.Draw(m, r, img, image.Point{spritePart.StartX, spritePart.StartY}, draw.Src)
				p = image.Point{curX, 0}
			}
		} else {

			if strings.HasSuffix(elem.name, "_hover") {
				fmt.Fprintf(cssFile, "a.%s:hover\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(goSprite.png) -%dpx 0px;\n }\n \n", strings.Replace(elem.name, "_hover", "", 1), sr.Size().X, sr.Size().Y, curX)
			} else {
				fmt.Fprintf(cssFile, ".%s\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(goSprite.png) -%dpx 0px;\n }\n \n", elem.name, sr.Size().X, sr.Size().Y, curX)
			}
			curX += sr.Size().X
			r := image.Rectangle{p, p.Add(sr.Size())}
			draw.Draw(m, r, img, sr.Min, draw.Src)
		}
	}

	w, _ := os.Create("goSprite.png")
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.
}
