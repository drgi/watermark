package watermark

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	fontPath = "wqy-zenhei.ttf"

	opt = &Options{
		WatermarkHPercent:             0.1,
		WatermarkWPercent:             0.9,
		WatermarkHorizontalOffset:     20,
		WatermarkVerticalOffset:       20,
		WatermarkTextHorizontalOffset: 20,
		WatermarkTextVerticalOffset:   20,
		WatermarkOpacity:              128,
		// WatermarkTextColor            *image.Uniform
		// WatermarkBackgroundColor      *image.Uniform

		// JpegOptions     *jpeg.Options
		AutoOrientation: true,

		RowSpacing: 1.5,
		DPI:        100,
	}

	texts = [][]string{{
		"Копия документа изготовлена",
		"в электронном виде и заверена 12:44 12.12.2021",
		"Иванов Иван Иваныч - 789056634345",
	},
		{
			"Копия документа изготовлена",
		},
		{
			"Копия документа изготовлена",
			"в электронном виде и заверена 12:44 12.12.2021",
			"Иванов Иван Иваныч - 789056634345",
			"Копия документа изготовлена",
			"в электронном виде и заверена 12:44 12.12.2021",
			"Иванов Иван Иваныч - 789056634345",
		},
		{
			"",
		},
		{},
	}
)

const (
	TestFileDir = "testfiles/dif-size"
	ResultDir   = "testfiles/results"
)

func TestWatermark(t *testing.T) {
	t.Log("start testing")

	wm, err := New(fontPath, opt)
	if err != nil {
		t.Error(err)
	}

	fileList, _ := ioutil.ReadDir(TestFileDir)
	for _, file := range fileList {
		path := filepath.Join(TestFileDir, file.Name())
		fmt.Println("path: ", path)
		targetImage, _ := ioutil.ReadFile(path)
		for i, text := range texts {
			resultImage, err := wm.AddFromText(targetImage, text)
			if err != nil {
				fmt.Printf("Watermark error: %v\n", err)
			}
			name := fmt.Sprintf("%s_%v.jpg", file.Name(), i)
			resultPath := filepath.Join(ResultDir, name)
			ioutil.WriteFile(resultPath, resultImage, 0644)
		}

	}

}
