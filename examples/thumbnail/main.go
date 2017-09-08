package main

import (
	"fmt"
	"github.com/bsed/imageupload"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"time"
)

var currentImage *imageupload.Image

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.POST("/upload", func(c echo.Context) error {
		img, err := imageupload.Process(c.Request(), "file")
		if err != nil {
			panic(err)
		}

		thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
		if err != nil {
			panic(err)
		}

		thumb.Save(fmt.Sprintf("%d.png", time.Now().Unix()))
		thumb.Write(c.Response().Writer)
		c.Redirect(http.StatusMovedPermanently, "/")
		return nil
	})
	e.Logger.Fatal(e.Start(":8443"))
}
