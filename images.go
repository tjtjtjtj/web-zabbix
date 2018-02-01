package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
)

func concatenateSideways(r1, r2 io.Reader) (*bytes.Buffer, error) {

	img1, fm1, err := image.Decode(r1)
	if err != nil {
		log.Fatalf("image decode err dayo %v", err)
	}
	img2, fm2, err := image.Decode(r2)
	if err != nil {
		log.Fatalf("image decode err dayo %v", err)
	}

	fmt.Printf("format:%s\n", fm1)
	fmt.Printf("max x:%d,max y:%d\n", img1.Bounds().Max.X, img1.Bounds().Max.Y)
	fmt.Printf("min x:%d,min y:%d\n", img1.Bounds().Min.X, img1.Bounds().Min.Y)
	fmt.Printf("format:%s\n", fm2)
	fmt.Printf("max x:%d,max y:%d\n", img2.Bounds().Max.X, img2.Bounds().Max.Y)
	fmt.Printf("min x:%d,min y:%d\n", img2.Bounds().Min.X, img2.Bounds().Min.Y)

	//starting position of the second image
	sp2 := image.Point{img1.Bounds().Dx(), 0}

	//new rectangle for the second image
	rt2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}

	//rectangle for the big image
	var rt image.Rectangle
	if img1.Bounds().Dy() < rt2.Dy() {
		rt = image.Rectangle{image.Point{0, 0}, rt2.Max}
	} else {
		rt = image.Rectangle{image.Point{0, 0}, image.Point{rt2.Max.X, img1.Bounds().Dy()}}

	}

	fmt.Printf("big rectangle:%v\n", rt)

	rgba := image.NewRGBA(rt)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rt2, img2, image.Point{0, 0}, draw.Src)

	b := new(bytes.Buffer)

	png.Encode(b, rgba)

	return b, nil
}

// bytesBbuffer writeoutputとか
func concatenateVertically(r1, r2 io.Reader) (*bytes.Buffer, error) {

	img1, fm1, err := image.Decode(r1)
	if err != nil {
		log.Fatalf("image decode err dayo %v", err)
	}
	img2, fm2, err := image.Decode(r2)
	if err != nil {
		log.Fatalf("image decode err dayo %v", err)
	}

	fmt.Printf("format:%s\n", fm1)
	fmt.Printf("max x:%d,max y:%d\n", img1.Bounds().Max.X, img1.Bounds().Max.Y)
	fmt.Printf("min x:%d,min y:%d\n", img1.Bounds().Min.X, img1.Bounds().Min.Y)
	fmt.Printf("format:%s\n", fm2)
	fmt.Printf("max x:%d,max y:%d\n", img2.Bounds().Max.X, img2.Bounds().Max.Y)
	fmt.Printf("min x:%d,min y:%d\n", img2.Bounds().Min.X, img2.Bounds().Min.Y)

	//starting position of the second image
	sp2 := image.Point{0, img1.Bounds().Dy()}

	//new rectangle for the second image
	rt2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}

	//rectangle for the big image
	var rt image.Rectangle
	if img1.Bounds().Dy() < rt2.Dx() {
		rt = image.Rectangle{image.Point{0, 0}, rt2.Max}
	} else {
		rt = image.Rectangle{image.Point{0, 0}, image.Point{img1.Bounds().Dx(), rt2.Max.Y}}

	}

	fmt.Printf("big rectangle:%v\n", rt)

	rgba := image.NewRGBA(rt)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, rt2, img2, image.Point{0, 0}, draw.Src)

	b := new(bytes.Buffer)

	png.Encode(b, rgba)

	return b, nil
}

/*
func main() {
	var files []*os.File
	r1, err := os.Open("./slime.png")
	if err != nil {
		log.Fatalln("open err dayo")
	}
	files = append(files, r1)

	r2, err := os.Open("./doraemon.png")
	if err != nil {
		log.Fatalln("open err dayo")
	}
	files = append(files, r2)
	r3, err := os.Open("./kitsune.png")
	if err != nil {
		log.Fatalln("open err dayo")
	}
	files = append(files, r3)

	fmt.Printf("%v", files)

	b := new(bytes.Buffer)
	io.Copy(b, files[0])

	for _, f := range files[1:] {
		b, _ = concatenateVertically(b, f)
	}

	file, _ := os.Create("output4.png")
	io.Copy(file, b)

*/
