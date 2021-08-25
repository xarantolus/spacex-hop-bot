package match

import (
	_ "embed"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "golang.org/x/image/webp"
)

func TestFaceDetector_DetectFaces(t *testing.T) {

	const testImageDir = "face_testdata"

	if _, err := os.Stat(testImageDir); os.IsNotExist(err) {
		t.Skip()
		return
	} else if err != nil {
		panic(err)
	}

	ts := httptest.NewServer(http.FileServer(http.FS(os.DirFS(testImageDir))))

	var tests = []struct {
		filename      string
		containsFaces bool
		wantErr       bool
	}{
		{"a.txt", false, true},

		{"1.jpg", true, false},
		{"2.jpg", true, false},

		{"3.jpg", false, false},
		{"4.jpg", false, false},
		{"5.jpg", false, false},
		{"6.jpg", false, false},
	}

	var detector = NewFaceDetector()

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			gotCount, err := detector.DetectFaces(ts.URL + "/" + tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("FaceDetector.DetectFaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (gotCount == 0) == tt.containsFaces {
				if tt.containsFaces {
					t.Errorf("FaceDetector.DetectFaces() = %v, but want faces", gotCount)
				} else {
					t.Errorf("FaceDetector.DetectFaces() = %v, but want none", gotCount)
				}
			}
		})
	}
}
