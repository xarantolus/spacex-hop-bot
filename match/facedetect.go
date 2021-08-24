package match

import (
	"bufio"
	_ "embed"
	"fmt"
	"image"
	"log"
	"math/rand"
	"net/http"
	"time"

	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/webp"

	"github.com/dghubble/go-twitter/twitter"
	pigo "github.com/esimov/pigo/core"
)

type FaceDetector struct {
	client             http.Client
	possibleUserAgents []string

	faceClassifier *pigo.Pigo
}

// pigo facefinder (https://github.com/esimov/pigo/blob/master/cascade/facefinder), MIT licensed
//go:embed facefinder
var faceFinderCascadeBytes []byte

func NewFaceDetector() *FaceDetector {
	pigo := pigo.NewPigo()

	classifier, err := pigo.Unpack(faceFinderCascadeBytes)
	if err != nil {
		log.Fatalf("Error unpacking cascade file: %s", err)
	}

	return &FaceDetector{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		possibleUserAgents: []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:86.0) Gecko/20100101 Firefox/86.0",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
			"Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
		},
		faceClassifier: classifier,
	}
}

// fetchImage downloads the image at the given URL for use with pigo
func (f *FaceDetector) fetchImage(url string) (image *image.NRGBA, err error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	// Set a few headers to look like a browser
	req.Header.Set("User-Agent", f.possibleUserAgents[rand.Intn(len(f.possibleUserAgents))])
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US;q=0.7,en;q=0.3")

	resp, err := f.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		err = fmt.Errorf("unexpected StatusCode %d", resp.StatusCode)
		return
	}

	return pigo.DecodeImage(bufio.NewReader(resp.Body))
}

// DetectFaces counts faces on an image at the given URL
func (f *FaceDetector) DetectFaces(url string) (count int, err error) {
	src, err := f.fetchImage(url)
	if err != nil {
		return
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     cols / 16,
		MaxSize:     1000,
		ShiftFactor: 0.25,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	dets := f.faceClassifier.RunCascade(cParams, 0.0)

	dets = f.faceClassifier.ClusterDetections(dets, 0.2)

	// Now count those with a certain score
	for _, d := range dets {
		if d.Q > 5.0 {
			count++
		}
	}

	return count, nil
}

// FaceRatio calculates the face ratio of a tweet. It is
// "Equation": (number of faces in all images) / (number of images)
func (f *FaceDetector) FaceRatio(tweet *twitter.Tweet) float32 {
	var imageUrls []string

	if tweet.ExtendedEntities != nil {
		for _, m := range tweet.ExtendedEntities.Media {
			imageUrls = append(imageUrls, m.MediaURLHttps)
		}
	} else if tweet.Entities != nil {
		for _, m := range tweet.Entities.Media {
			imageUrls = append(imageUrls, m.MediaURLHttps)
		}
	}

	if len(imageUrls) == 0 {
		return 0
	}

	var faceCount int

	for _, iurl := range imageUrls {
		count, err := f.DetectFaces(iurl)
		if err != nil {
			log.Printf("[FaceDetector] Error detecting face in %q: %s\n", iurl, err.Error())
			continue
		}

		faceCount += count
	}

	return float32(faceCount) / float32(len(imageUrls))
}
