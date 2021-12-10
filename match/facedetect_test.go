package match

import (
	_ "embed"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	_ "golang.org/x/image/webp"
)

func TestFaceDetector_DetectFaces(t *testing.T) {

	const testImageDir = "face_testdata"

	if _, err := os.Stat(testImageDir); os.IsNotExist(err) {
		t.Skip("Skipping test because the test image directory does not exist")
		return
	} else if err != nil {
		panic(err)
	}

	ts := httptest.NewServer(http.FileServer(http.FS(os.DirFS(testImageDir))))

	var detector = NewFaceDetector()

	var (
		fails int
		total int
	)

	err := filepath.Walk(testImageDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		total++

		t.Run(path, func(t *testing.T) {
			var condChar = filepath.Base(path)[0]

			gotCount, err := detector.DetectFaces(ts.URL + "/" + info.Name())
			if err != nil && condChar != 'e' {
				t.Errorf("FaceDetector.DetectFaces(%s) error = %v, but wanted no error", path, err)
				return
			}

			switch condChar {
			case 'e':
				if err == nil {
					// files starting with "e" should return an error, so this didn't work
					t.Errorf("FaceDetector.DetectFaces(%s) error = %v, but wanted error", path, err)
				}
			case 'y':
				if gotCount == 0 {
					t.Logf("FaceDetector.DetectFaces(%s) = %v, but want faces", path, gotCount)
					fails++
				}
			case 'n':
				if gotCount > 0 {
					t.Logf("FaceDetector.DetectFaces(%s) = %v, but want none", path, gotCount)
					fails++
				}
			default:
				t.Fatalf("unexpected character %c as first char in file name", condChar)
			}
		})

		return nil
	})
	if err != nil {
		panic("unexpected walk failure: " + err.Error())
	}

	var failRate = float32(fails) / float32(total)
	if failRate > .30 {
		t.Errorf("More than 30%% of face checks failed: %d of %d (%.2f)", fails, total, failRate*100)
	} else {
		t.Logf("Less than 30%% of face checks failed")
	}
}
