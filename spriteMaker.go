package main

import (
	"code.google.com/p/gcfg"
	"container/list"
	"flag"
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

type Gcfg struct {
	Info struct {
		Sufix string
		Count int
		Width int
		MaxH  int
	}
	Image map[string]*struct {
		Name   []string
		Width  int
		Height int
		StartX int
		StartY int
		Hover  bool
	}
}

func (p *Gcfg) Clear() {
	p = &Gcfg{}
}

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

/*
(-d) - images directory
(-s) - sprite outfile
(-c) - css outfile
(-g) - config file (gcfg format)
*/

var dirSrc string
var spriteOutFile string
var cssOutFile string
var gcfgSrc string

func init() {
	flag.StringVar(&dirSrc, "d", "images", "Images directory")
	flag.StringVar(&spriteOutFile, "s", "goSprite.png", "Sprite out file")
	flag.StringVar(&cssOutFile, "c", "goStyle.css", "Css out file")
	flag.StringVar(&gcfgSrc, "g", "config", "Config files directory (gcfg format)")
	flag.Parse()
}

func main() {

	cfg := Gcfg{}
	configsId := make(map[string]int)
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
	configs := make([]Gcfg, len(files))
	images := list.New()
	elem := Element{}
	maxH, maxW := 0, 0
	for i := range files {
		imgSrc, err := os.Open(fmt.Sprintf("%s%c%s", dirSrc, os.PathSeparator, files[i].Name()))
		defer imgSrc.Close()
		img, _, err := image.Decode(imgSrc)
		if err != nil {
			continue
		}
		elem.name = strings.Split(files[i].Name(), ".")[0]
		elem.img = img
		images.PushFront(elem)
		if _, err := os.Stat(fmt.Sprintf("%s%c%s.gcfg", gcfgSrc, os.PathSeparator, elem.name)); os.IsNotExist(err) {
			maxW += img.Bounds().Size().X
			if img.Bounds().Size().Y > maxH {
				maxH = img.Bounds().Size().Y
			}
		} else {
			err = gcfg.ReadFileInto(&configs[i], fmt.Sprintf("%s%c%s.gcfg", gcfgSrc, os.PathSeparator, elem.name))
			if err != nil {
				fmt.Printf("Error while reading config(%s):%s.gcfg\n", elem.name, err)

			} else {
				fmt.Printf("config %s.gcfg loaded\n", elem.name)
				configsId[elem.name] = i
				for j := range configs[i].Image {
					for k := range configs[i].Image[j].Name {
						configs[i].Image[j].Name[k] = fmt.Sprintf("%s%s", configs[i].Info.Sufix, configs[i].Image[j].Name[k])
					}
				}
				maxW += configs[i].Info.Width
				if configs[i].Info.MaxH > maxH {
					maxH = configs[i].Info.MaxH
				}
			}
		}
	}
	fmt.Printf("%d x %d\n", maxW, maxH)
	cssFile, err := os.Create(cssOutFile)
	if err != nil {
		log.Fatal(err)
	}
	defer cssFile.Close()
	m := image.NewRGBA(image.Rect(0, 0, maxW, maxH))
	curX := 0
	for e := images.Front(); e != nil; e = e.Next() {
		p := image.Point{curX, 0}
		elem = e.Value.(Element)
		img := elem.img
		sr := img.Bounds()
		configId, ok := configsId[elem.name]
		if !ok {
			if strings.HasSuffix(elem.name, "_hover") {
				fmt.Fprintf(cssFile, "a.%s:hover\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(goSprite.png) -%dpx 0px;\n }\n \n", strings.Replace(elem.name, "_hover", "", 1), sr.Size().X, sr.Size().Y, curX)
			} else {
				fmt.Fprintf(cssFile, ".%s\n {\n  width:%dpx;\n  height:%dpx;\n  background:url(goSprite.png) -%dpx 0px;\n }\n \n", elem.name, sr.Size().X, sr.Size().Y, curX)
			}
			curX += sr.Size().X
			r := image.Rectangle{p, p.Add(sr.Size())}
			draw.Draw(m, r, img, sr.Min, draw.Src)
		} else {
			cfg = configs[configId]
			for i := 0; i < cfg.Info.Count; i++ {
				spritePart := cfg.Image[fmt.Sprintf("%d", i)]
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
			cfg.Clear()
		}
	}

	w, _ := os.Create(spriteOutFile)
	defer w.Close()
	png.Encode(w, m) //Encode writes the Image m to w in PNG format.
}
