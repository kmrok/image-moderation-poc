package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"google.golang.org/api/option"

	pb "google.golang.org/genproto/googleapis/cloud/vision/v1"
)

var rc *rekognition.Rekognition
var vc *vision.ImageAnnotatorClient

func main() {
	var err error
	ctx := context.Background()
	session := session.Must(session.NewSession())
	rc = rekognition.New(session, aws.NewConfig().WithRegion("ap-northeast-1"))
	vc, err = vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GCLOUD_CREDENTIALS"))))
	if err != nil {
		fmt.Printf("failed to create client: %s", err.Error())
		return
	}

	url1 := "https://dhei5unw3vrsx.cloudfront.net/images/family_picnic_resized.jpg"
	url2 := "https://dhei5unw3vrsx.cloudfront.net/images/yoga_swimwear_resized.jpg"

	printMeasureTime(func() { _, _ = detectModerationLabels(ctx, url1) }, "AWS / Picnic")
	printMeasureTime(func() { _, _ = detectModerationLabels(ctx, url2) }, "AWS / Swimwear")
	printMeasureTime(func() { _, _ = detectSafeSearchWithoutDownload(ctx, url1) }, "GCP(url) / Picnic")
	printMeasureTime(func() { _, _ = detectSafeSearchWithoutDownload(ctx, url2) }, "GCP(url) / Swimwear")
	printMeasureTime(func() { _, _ = detectSafeSearch(ctx, url1) }, "GCP(download) / Picnic")
	printMeasureTime(func() { _, _ = detectSafeSearch(ctx, url2) }, "GCP(download) / Swimwear")
}

func detectModerationLabels(ctx context.Context, url string) (*rekognition.DetectModerationLabelsOutput, error) {
	image, err := downloadImage(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	d, err := rc.DetectModerationLabels(&rekognition.DetectModerationLabelsInput{
		Image: &rekognition.Image{
			Bytes: image,
		},
		MinConfidence: aws.Float64(0),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to detect: %w", err)
	}

	// Uncomment the following if you want to see the results.
	// fmt.Printf("ModerationLabels:\n")
	// for _, v := range d.ModerationLabels {
	// 	fmt.Printf("%+v \n", v)
	// }
	// fmt.Printf("ModerationModelVersion: %s\n", *d.ModerationModelVersion)

	return d, nil
}

func detectSafeSearchWithoutDownload(ctx context.Context, url string) (*pb.SafeSearchAnnotation, error) {
	props, err := vc.DetectSafeSearch(ctx, vision.NewImageFromURI(url), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detect: %w", err)
	}

	// Uncomment the following if you want to see the results.
	// fmt.Printf("%+v\n", props)

	return props, nil
}

func detectSafeSearch(ctx context.Context, url string) (*pb.SafeSearchAnnotation, error) {
	image, err := downloadImage(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download: %w", err)
	}

	img, err := vision.NewImageFromReader(bytes.NewReader(image))
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	props, err := vc.DetectSafeSearch(ctx, img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to detect: %w", err)
	}

	// Uncomment the following if you want to see the results.
	// fmt.Printf("%+v\n", props)

	return props, nil
}

func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func printMeasureTime(fn func(), description string) {
	start := time.Now()
	fn()
	end := time.Now()
	fmt.Printf("%s -> %f seconds\n", description, (end.Sub(start)).Seconds())
}

// Batch of Cloud Vision SafeSearch.
func batchDetectSafeSearch(ctx context.Context, urlList []string) (*pb.BatchAnnotateImagesResponse, error) {
	var req []*pb.AnnotateImageRequest
	for _, u := range urlList {
		image, err := downloadImage(u)
		if err != nil {
			return nil, fmt.Errorf("failed to download: %w", err)
		}

		img, err := vision.NewImageFromReader(bytes.NewReader(image))
		if err != nil {
			return nil, fmt.Errorf("failed to read: %w", err)
		}

		req = append(req, &pb.AnnotateImageRequest{
			Image:    img,
			Features: []*pb.Feature{{Type: pb.Feature_SAFE_SEARCH_DETECTION, MaxResults: int32(0)}},
		})
	}

	resp, err := vc.BatchAnnotateImages(ctx, &pb.BatchAnnotateImagesRequest{Requests: req})
	if err != nil {
		return nil, fmt.Errorf("failed to detect: %w", err)
	}

	return resp, nil
}
