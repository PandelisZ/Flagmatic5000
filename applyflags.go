package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gocarina/gocsv"

	"gopkg.in/yaml.v2"
)

type Expert struct {
	City string
}

type City struct {
	City string `csv:"city"`
	Iso2 string `csv:"iso2"`
}

type Flags []string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadCities() []City {
	citiesCsv, err := os.OpenFile("simplemaps-worldcities-basic.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer citiesCsv.Close()

	cities := []City{}

	if err := gocsv.UnmarshalFile(citiesCsv, &cities); err != nil { // Load clients from file
		panic(err)
	}

	for idx, c := range cities {
		cities[idx].City = strings.ToLower(c.City)
	}

	return cities
}

func removeDuplicates(elements []Expert) []Expert {
	encountered := map[Expert]bool{}
	result := []Expert{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func removeFlagDupes(elements Flags) Flags {
	encountered := map[string]bool{}
	result := Flags{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func removeCommas(e []Expert) []Expert {
	for idx, ex := range e {
		e[idx].City = strings.ToLower(strings.SplitN(ex.City, ",", 2)[0])
	}
	return removeDuplicates(e)
}

func getCountry(e []Expert, c []City) Flags {
	var result Flags

	cities := make(map[string]string)
	for _, cit := range c {
		cities[cit.City] = cit.Iso2
	}

	for _, ex := range e {
		isoCode, ok := cities[ex.City]
		if ok {
			result = append(result, strings.ToLower(isoCode))
		}
	}

	return removeFlagDupes(result)

}

func getImageDimension(img image.Image) (int, int) {
	return img.Bounds().Dx(), img.Bounds().Dy()
}

func loadFlagImg(country string) image.Image {
	infile, err := os.Open("./flags/" + country + ".png")
	if err != nil {
		check(err)
	}
	defer infile.Close()

	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(infile)
	if err != nil {
		check(err)
	}
	return src
}

func loadBanner() image.Image {
	infile, err := os.Open("experts-empty.png")
	if err != nil {
		check(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		check(err)
	}
	return src
}

func main() {
	args := os.Args[1:]
	e := []Expert{}

	cities := loadCities()

	yamlIn, err := ioutil.ReadFile(args[0])
	check(err)

	err2 := yaml.Unmarshal(yamlIn, &e)
	check(err2)

	e = removeCommas(e)
	isoFlags := getCountry(e, cities)
	fmt.Print(isoFlags)

	var flagImages []image.Image
	for _, iso := range isoFlags {
		flagImages = append(flagImages, loadFlagImg(iso))
	}

	fx, fy := getImageDimension(flagImages[0])
	fmt.Print(fx)
	cursor := 0
	yPos := 0
	maxWidth := 0
	padding := 10
	for _, img := range flagImages {
		maxWidth += img.Bounds().Dx() + padding
	}

	expertFlags := image.NewRGBA(image.Rect(0, 0, maxWidth, fy*2+padding))

	for idx, fImg := range flagImages {
		draw.Draw(expertFlags, expertFlags.Bounds(), fImg, image.Point{-cursor, yPos}, draw.Src)
		if idx == 10 {
			cursor = 0
			yPos = yPos - fImg.Bounds().Dy() - padding
		} else {
			cursor = cursor + fImg.Bounds().Dx() + padding
		}

	}

	toimg, _ := os.Create("expert-flags.png")
	defer toimg.Close()
	png.Encode(toimg, expertFlags)

	expertsBanner, _ := os.Create("expert-banner-flags.png")
	defer expertsBanner.Close()

	banner := loadBanner()
	bX, bY := getImageDimension(banner)
	outBanner := image.NewRGBA(image.Rect(0, 0, bX, bY))

	draw.Draw(outBanner, banner.Bounds(), banner, image.Point{0, 0}, draw.Src)
	draw.Draw(outBanner, banner.Bounds(), expertFlags, image.Point{-550, -300}, draw.Over)
	png.Encode(expertsBanner, outBanner)

}
