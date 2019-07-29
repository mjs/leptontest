// Programme de test divers en golang

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/TheCacophonyProject/lepton3"
	"github.com/TheCacophonyProject/periph/host"
)

func dumpToPNG(path string, frame *lepton3.Frame) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	defer func() {
		w.Flush()
		f.Close()
	}()
	return png.Encode(w, reduce(frame))
}

var dst = image.NewGray16(image.Rect(0, 0, lepton3.FrameCols, lepton3.FrameRows))

func reduce(src *lepton3.Frame) *image.Gray16 {
	minVal := uint16(math.MaxUint16)
	maxVal := uint16(0)
	for y := 0; y < lepton3.FrameRows; y++ {
		for x := 0; x < lepton3.FrameCols; x++ {
			i := src.Pix[y][x]
			if i > maxVal {
				maxVal = i
			}
			if i < minVal {
				minVal = i
			}
		}
	}

	var norm = math.MaxUint16 / (maxVal - minVal)
	for y, row := range src.Pix {
		for x, val := range row {
			dst.SetGray16(x, y, color.Gray16{Y: (val - minVal) * norm})
		}
	}
	return dst
}

func main() {

	// Init host
	_, err := host.Init()
	if err != nil {
		fmt.Println("Error in init host")
	}

	speed := int64(20000000)
	camera, err := lepton3.New(speed)
	if err != nil {
		fmt.Println("Error in lepton3.New(speed)")
	}

	// if err := camera.SetRadiometry(true); err != nil {
	// 	fmt.Println("Error in camera.SetRadiometry(true)")
	// }

	err = camera.Open()
	if err != nil {
		fmt.Println("Error in camera.Open()")
	}
	defer camera.Close()

	rawFrame := new(lepton3.RawFrame)
	frame := new(lepton3.Frame)
	i := 0
	for {
		err := camera.NextFrame(rawFrame)
		if err != nil {
			fmt.Println("Error in camera.NextFrame")
		}

		rawFrame.ToFrame(frame)

		filename := filepath.Join("/home/pi", fmt.Sprintf("%05d.png", i))
		err = dumpToPNG(filename, frame)
		if err != nil {
			fmt.Println("Error in dumpToPNG")
		}
		i++
	}
}
