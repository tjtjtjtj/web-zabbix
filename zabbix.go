package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/sclevine/agouti"
)

// 読み込み対象スクリーンはyamlファイルで定義する
// それを構造体読み込みとして処理

type Zabbix struct {
	// UserID/Password/Proxyはjsonファイルに記載しない
	UserID   string
	Password string
	Proxy    string
	URL      string `json:"url"`
	Page     *agouti.Page
}

type screenimage struct {
	x, y   int
	w, h   int
	src    string
	reader io.Reader
}

func (z *Zabbix) SetupEnv() {
	env := os.Getenv("ZABBIX_ENV")
	file, err := os.Open(fmt.Sprintf("conf/%s/zabbix_env.json", env))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(z)
	if err != nil {
		log.Fatal(err)
	}

	z.UserID = os.Getenv("ZABBIX_USER")
	z.Password = os.Getenv("ZABBIX_PASSWORD")
	z.Proxy = os.Getenv("ZABBIX_PROXY")

	fmt.Println(*z)
}

func (z *Zabbix) Login() {
	fmt.Println("login start")
	fmt.Printf("page:%v", z.Page)

	err := z.Page.Navigate(z.URL)
	if err != nil {
		log.Fatalf("Failed to navigate:%v", err)
	}

	fmt.Println("login start2")
	html, _ := z.Page.HTML()
	fmt.Printf("html:%s", html)
	z.Page.Screenshot("/tmp/outputs/zabbix1.png")
	fmt.Println("login start3")

	userid := z.Page.FindByID("name")
	password := z.Page.FindByID("password")
	userid.Fill(z.UserID)
	password.Fill(z.Password)
	z.Page.Screenshot("/tmp/outputs/zabbix2.png")
	if err := z.Page.FindByID("enter").Click(); err != nil {
		log.Fatal("Failed to set password", err)
	}

}

// 対象スクリーンを指定できるようにする
func (z *Zabbix) ScreenshotALL() {

	if err := z.Page.FindByLink("スクリーン").Click(); err != nil {
		log.Fatal("Failed to click", err)
	}

	if err := z.Page.FindByID("elementid").Select("testscreen"); err != nil {
		log.Fatal("Failed to click", err)
	}

	//e, err := z.Page.FindByClass("screen_view").AllByClass("flickerfreescreen").Elements()
	trcount, err := z.Page.FindByClass("screen_view").All("tr").Count()
	if err != nil {
		log.Fatal("Failed to click", err)
	}
	fmt.Printf("trcount:%d", trcount)

	//tr tdで位置を特定

	screenimages := make([][]*screenimage, 0)
	for i := 0; i < trcount; i++ {
		ff := z.Page.FindByClass("screen_view").All("tr").At(i).AllByClass("flickerfreescreen")
		ffcount, err := ff.Count()
		if err != nil {
			log.Fatal("Failed to count", err)
		}
		if ffcount == 0 {
			continue
		}
		fmt.Println("kkkkkkkkkkkkkkkokoni nnn")
		screenimages = append(screenimages, make([]*screenimage, 0))
		for j := 0; j < ffcount; j++ {
			s := new(screenimage)
			s.y, s.x = i, j
			s.src, err = ff.At(j).Find("img").Attribute("src")
			if err != nil {
				log.Fatal("Failed to src\n", err)
			}
			fmt.Println("kokoni nnn")
			screenimages[i] = append(screenimages[i], s)
		}
	}

	allb := new(bytes.Buffer)
	for i, ss := range screenimages {
		fmt.Printf("aaaaaaaallb %d\n", i)
		sb := new(bytes.Buffer)
		for k, v := range ss {
			fmt.Printf("screenimages:%v\n", v)

			err := z.Page.Navigate(v.src)
			if err != nil {
				log.Fatalf("Failed to navigate:%v", err)
			}
			time.Sleep(5 * time.Second)
			i := z.Page.Find("img")
			body_ele, _ := i.Elements()
			w, h, err := body_ele[0].GetSize()
			if err != nil {
				log.Fatalf("Failed to element:%v", err)
			}
			fmt.Printf("size:width %d,heigh %d¥n", w, h)
			z.Page.Size(w, h)
			b, err := z.Page.Session().GetScreenshot()
			if err != nil {
				log.Fatalf("Failed to getscreen:%v", err)
			}
			v.reader = bytes.NewReader(b)
			if k == 0 {
				io.Copy(sb, v.reader)
				continue
			}
			sb, _ = concatenateSideways(sb, v.reader)

			z.Page.Screenshot(fmt.Sprintf("/tmp/outputs/zabbix00%d.png", k))
			//	z.Screenshot(fmt.Sprintf("/tmp/outputs/zabbix000%d.png", k))
		}
		if i == 0 {
			io.Copy(allb, sb)
		} else {
			allb, _ = concatenateVertically(allb, sb)

		}
	}

	file, _ := os.Create("/tmp/outputs/output4.png")
	io.Copy(file, allb)

}

func (z *Zabbix) Screenshot(filepath string) error {

	time.Sleep(5 * time.Second)
	i := z.Page.Find("img")
	body_ele, _ := i.Elements()
	w, h, err := body_ele[0].GetSize()
	if err != nil {
		return err
	}
	fmt.Printf("size:width %d,heigh %d¥n", w, h)
	z.Page.Size(w, h)
	z.Page.Screenshot(filepath)
	return nil
}

func concatinate(i1, i2 io.Reader) {

	img1, type1, err := image.Decode(i1)
	img2, type2, err := image.Decode(i2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("type1:%s,type2:%s", type1, type2)

	sp2 := image.Point{img1.Bounds().Dx(), 0}
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
	r := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)

	out, err := os.Create("/tmp/outputs/zabbixcon.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)
}

func concati() {
	imgFile1, err := os.Open("/tmp/outputs/zabbix0000.png")
	imgFile2, err := os.Open("/tmp/outputs/zabbix0001.png")
	if err != nil {
		fmt.Println(err)
	}
	img1, type1, err := image.Decode(imgFile1)
	img2, type2, err := image.Decode(imgFile2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("type1:%s,type2:%s", type1, type2)

	sp2 := image.Point{img1.Bounds().Dx(), 0}
	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
	r := image.Rectangle{image.Point{0, 0}, r2.Max}
	rgba := image.NewRGBA(r)

	draw.Draw(rgba, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
	draw.Draw(rgba, r2, img2, image.Point{0, 0}, draw.Src)

	out, err := os.Create("/tmp/outputs/zabbixcon.png")
	if err != nil {
		fmt.Println(err)
	}

	png.Encode(out, rgba)
}
