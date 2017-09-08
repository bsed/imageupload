package main

import (
	"github.com/bsed/imageupload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var currentImage *imageupload.Image

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.GET("/image", func(c echo.Context) error {
		if currentImage == nil {
			return c.String(http.StatusNotFound, "404")
		}
		currentImage.Write(c.Response().Writer)

		return nil
	})

	e.GET("/thumbnail", func(c echo.Context) error {
		if currentImage == nil {
			return c.String(http.StatusNotFound, "404")
		}
		t, err := imageupload.ThumbnailJPEG(currentImage, 300, 300, 80)
		if err != nil {
			panic(err)
		}
		t.Write(c.Response().Writer)

		return nil
	})

	e.POST("/upload", func(c echo.Context) error {
		img, err := imageupload.Process(c.Request(), "file")
		if err != nil {
			panic(err)
		}

		currentImage = img
		c.Redirect(http.StatusMovedPermanently, "/")
		return nil
	})
	e.Logger.Fatal(e.Start(":8443"))
}
