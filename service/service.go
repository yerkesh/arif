package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"image/png"
	"strconv"
	"time"

	"github.com/unidoc/unipdf/v3/model"
	"github.com/unidoc/unipdf/v3/render"

	"arif/clients"
	"arif/clients/aws_s3"
	"arif/entity"
	"arif/repo"
)

const (
	bucketName = "arifs3"
)

// statuses
const (
	statusSplittingImages = "splitting_images"
	statusExtractingText  = "extracting_text"
	statusMakingTranslate = "making_translate"
	statusDone            = "done"
)

func ProcessPDF(ctx context.Context, fileBytes []byte) (entity.UploadResult, error) {
	md5Hash := GenerateMD5Hash(len(fileBytes))

	ctx = context.Background()
	go func() {
		locationURL, err := aws_s3.UploadImageToS3(ctx, fileBytes, bucketName, fmt.Sprintf("books/%s.pdf", md5Hash))
		if err != nil {
			fmt.Println(err)
		}

		err = repo.CreateRequest(ctx, md5Hash, locationURL)
		if err != nil {
			fmt.Println(err)
		}

		err = repo.UpdateRequestStatus(ctx, md5Hash, statusSplittingImages)
		if err != nil {
			fmt.Println(err)
		}

		urlsMap, err := PdfToImages(ctx, fileBytes, md5Hash)
		if err != nil {
			fmt.Println(err)
		}

		err = repo.CreateEntry(ctx, md5Hash, urlsMap)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(locationURL)
		fmt.Println(urlsMap)

		err = repo.UpdateRequestStatus(ctx, md5Hash, statusExtractingText)
		if err != nil {
			fmt.Println(err)
		}

		extractedTexts := make(map[int]string, len(urlsMap))
		for page, url := range urlsMap {
			text, err := clients.ExtractFromImage(ctx, url)
			if err != nil {
				fmt.Println(err)
			}
			extractedTexts[page] = text
			fmt.Println(text)
		}

		err = repo.InsertExtracted(ctx, md5Hash, extractedTexts)
		if err != nil {
			fmt.Println(err)
		}

		err = repo.UpdateRequestStatus(ctx, md5Hash, statusMakingTranslate)
		if err != nil {
			fmt.Println(err)
		}

		translated, err := clients.GenerateTranslateMessage(ctx, extractedTexts)
		if err != nil {
			fmt.Println(err)
		}

		err = repo.InsertTranslated(ctx, md5Hash, translated)
		if err != nil {
			fmt.Println(err)
		}

		err = repo.UpdateRequestStatus(ctx, md5Hash, statusDone)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(translated)
	}()

	return entity.UploadResult{Hash: md5Hash}, nil
}

// GenerateMD5Hash creates an MD5 hash from milliseconds and a string value.
func GenerateMD5Hash(numBytes int) string {
	// Convert milliseconds to string
	msString := strconv.FormatInt(time.Now().UnixMilli(), 10)
	strNum := strconv.FormatInt(int64(numBytes), 10)
	// Combine milliseconds and the input string
	combined := msString + strNum

	// Compute MD5 hash
	hash := md5.Sum([]byte(combined))

	// Convert hash bytes to hexadecimal string
	return hex.EncodeToString(hash[:])
}

// PdfToImages converts a PDF file into images, one image per page.
// pdfPath is the path to the input PDF file.
// outputDir is the directory where the images will be saved.
func PdfToImages(ctx context.Context, fileBytes []byte, hash string) (map[int]string, error) {
	result := make(map[int]string)

	// Create a bytes.Reader from the []byte
	reader := bytes.NewReader(fileBytes)

	// Create a new PDF reader from the bytes.Reader
	pdfReader, err := model.NewPdfReader(reader)
	if err != nil {
		return nil, fmt.Errorf("error creating PDF reader: %v", err)
	}

	// Get the total number of pages
	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return nil, fmt.Errorf("error getting page count: %v", err)
	}

	// Iterate through each page
	for n := 1; n <= numPages; n++ {
		// Get the page
		page, err := pdfReader.GetPage(n)
		if err != nil {
			return nil, fmt.Errorf("could not get page %d: %v", n, err)
		}

		// Create a new renderer
		renderer := render.NewImageDevice()

		// Render the page to an image
		img, err := renderer.Render(page)
		if err != nil {
			return nil, fmt.Errorf("could not render page %d: %v", n, err)
		}

		// Create a buffer to hold the encoded image
		var buf bytes.Buffer

		// Encode the image to PNG and write it to the buffer
		err = png.Encode(&buf, img)
		if err != nil {
			return nil, fmt.Errorf("could not encode image for page %d: %v", n, err)
		}

		// Get the byte slice from the buffer
		imageBytes := buf.Bytes()

		// Upload the image to S3
		locationURL, err := aws_s3.UploadImageToS3(ctx, imageBytes, bucketName, fmt.Sprintf("images/%s_page_%d.png", hash, n))
		if err != nil {
			return nil, fmt.Errorf("could not upload image for page %d: %v", n, err)
		}

		result[n] = locationURL

		fmt.Println("Uploading... page number:", n)
	}

	return result, nil
}
