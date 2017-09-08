package imageupload

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Image struct {
	Filename    string
	ContentType string // file contenet type
	Data        []byte
	Size        int // file size
}

// Save save image to file.
func (i *Image) Save(filename string) error {
	return ioutil.WriteFile(filename, i.Data, 0600)
}

// DataURI convert image to base64 data uri.
func (i *Image) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))
}

// Write Write image to HTTP response.
func (i *Image) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", i.ContentType)
	w.Header().Set("Content-Length", strconv.Itoa(i.Size))
	w.Write(i.Data)
}

// ThumbnaimJPEG create JPEG thumbnail from image.
func (i *Image) ThumbnailJPEG(width int, height int, quality int) (*Image, error) {
	return ThumbnailJPEG(i, width, height, quality)
}

// ThumbnailPNG create PNG thumbnail from image.
func (i *Image) ThumbnailPNG(width int, height int) {
	return ThumbnailPNG(i, width, height)
}

// Limit lime the size of uploaded files.
func LimitFileSize(maxsize int64, w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxsize)
}

func okContentType(contentType string) bool {
	return contentType == "image/png" || contentType == "image/jpeg" || contentType == "image/gif"
}

// Process process image upload info.
func Process(r *http.Request, field string) (*Image, error) {
	file, fh, err := r.FormFile(field)
	if err != nil {
		return nil, err
	}

	contentType := fh.Header.Get("Content-Type")
	if !okContentType(contentType) {
		return nil, errors.New(fmt.Sprintf("Wrong content type: %s", contentType))
	}

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	_, _, err = image.Decode(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}

	i := &Image{
		Filename:    fh.Filename,
		ContentType: contentType,
		Data:        bs,
		Size:        len(bs),
	}

	return i, nil
}

// ThumbnailJPEG create JPEG thumbnail
func ThumbnailJPEG(i *Image, width int, height int, quality int) (*Image, error) {
	img, _, err := image.Decode(bytes.NewBuffer(i.Data))
	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = jpeg.Encode(data, thumbnail, &jpeg.Options{
		Quality: quality,
	})

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()

	t := &Image{
		Filename:    "thumbnail.jpg",
		ContentType: "image/jpeg",
		Data:        bs,
		Size:        len(bs),
	}

	return t, nil
}

// ThumbnailPNG create PNG thumnail
func ThumbnailPNG(i *Image, width int, height int) (*Image, error) {
	img, _, err := image.Decode(bytes.NewReader(i.Data))

	thumbnail := resize.Thumbnail(uint(width), uint(height), img, resize.Lanczos3)

	data := new(bytes.Buffer)
	err = png.Encode(data, thumbnail)

	if err != nil {
		return nil, err
	}

	bs := data.Bytes()
	t := &Image{
		Filename:    "thumbnail.png",
		ContentType: "image/png",
		Data:        bs,
		Size:        len(bs),
	}
	return t, nil
}
