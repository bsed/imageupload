# imageupload 文件上传缩略图
image upload thumbnail for go

## Process 处理文件上传
```go
img, err := imageupload.Process(c.Request(), "file")
```

## Save 保存文件
```go
img, err := imageupload.Process(c.Request(), "file")
		if err != nil {
			panic(err)
		}

		thumb, err := imageupload.ThumbnailPNG(img, 300, 300)
		if err != nil {
			panic(err)
		}

		thumb.Save(fmt.Sprintf("%d.png", time.Now().Unix()))
```
## Write 写入到文件中
```go
var currentImage *imageupload.Image
currentImage.Write(c.Response().Writer)
```
## LimitFileSize 设置文件尺寸


## ThumbnailJPEG JPEG缩略图

```go
t, err := imageupload.ThumbnailJPEG(currentImage, 300, 300, 80)
		if err != nil {
			panic(err)
		}
t.Write(c.Response().Writer)
```

## ThumbnailPNG PNG缩略图

```go
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
```