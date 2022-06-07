package image_process

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

type ImageProcess struct {
	imgRGBA     *image.RGBA64
	BinaryImage image.Image
	OriginImage image.Image
	workDir     string
}

func NewImageProcess(originImage, binaryImage image.Image, workDir string) *ImageProcess {
	img := &ImageProcess{
		BinaryImage: binaryImage,
		OriginImage: originImage,
		workDir:     workDir,
	}
	return img
}

// 扫描图片 percentage为透明度，取值0-1
func (i *ImageProcess) GradientImage(percentage float32) (string, error) {
	//copy原图生成rgb图
	i.imgRGBA = image.NewRGBA64(i.OriginImage.Bounds())
	draw.Draw(i.imgRGBA, i.OriginImage.Bounds(), i.OriginImage, image.Point{}, draw.Over)

	rect := i.BinaryImage.Bounds()
	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X-1; x++ {
			r, g, b, _ := i.BinaryImage.At(x, y).RGBA()
			nr, ng, nb, _ := i.BinaryImage.At(x+1, y).RGBA()
			if r*g*b != nr*ng*nb {
				or, og, ob, oa := i.OriginImage.At(x, y).RGBA()
				opacity := uint16(float32(oa) * percentage)
				v := i.OriginImage.ColorModel().Convert(color.NRGBA64{R: uint16(or), G: uint16(og), B: uint16(ob), A: opacity})
				rr, gg, bb, aa := v.RGBA()
				i.imgRGBA.SetRGBA64(x, y, color.RGBA64{uint16(rr), uint16(gg), uint16(bb), uint16(aa)})
			}
		}
	}
	path := i.workDir + "/output.png"
	err := i.SaveImage(path)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (i *ImageProcess) SaveImage(filePath string) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()
	return png.Encode(out, i.imgRGBA)
}
