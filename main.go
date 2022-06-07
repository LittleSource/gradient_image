package main

import (
	"fmt"
	"go1/clould"
	"go1/image_process"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
)

const dataDir = "testdata"

// testdata 目录下两个对应的文件夹目录
var (
	ImagesDir = filepath.Join(dataDir, "images")
)

func main() {
	workPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	//输入图片
	imagePath := filepath.Join(ImagesDir, "jimin4.jpg")
	orgFile, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	//生成二值图
	baiduCloud := clould.NewBaiduClould("", "")
	binPath := filepath.Join(ImagesDir, "bin.jpg")
	err = baiduCloud.PortraitSegmentation(orgFile, binPath)
	if err != nil {
		panic(err)
	}
	//原图和二值图片进行扫描
	orgFile, err = os.Open(imagePath)
	if err != nil {
		panic(err)
	}
	orgImage, _, err := image.Decode(orgFile)
	if err != nil {
		panic(err)
	}

	binFile, err := os.Open(binPath)
	if err != nil {
		panic(err)
	}
	binImage, _, err := image.Decode(binFile)
	if err != nil {
		panic(err)
	}
	//开始扫描并替换像素
	imageProcess := image_process.NewImageProcess(orgImage, binImage, workPath)
	path, err := imageProcess.GradientImage(0.9)
	if err != nil {
		panic(err)
	}
	orgFile.Close()
	binFile.Close()
	fmt.Println("save to:", path)
}
