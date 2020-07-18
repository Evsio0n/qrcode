package qrcode

import (
	"../log"
	"bytes"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
)

func ToPng(contentString string, width int, height int) *bytes.Buffer{
	qrCode, err := qr.Encode(contentString, qr.M, qr.Auto)
	qrCode, err = barcode.Scale(qrCode, width, height)
	log.Info(contentString)
	// create the output file
	bytePng  := new(bytes.Buffer)
	err=png.Encode(bytePng,qrCode)
	if err != nil {
		log.Error(err)
	}

	if err != nil {
		log.Error(err)
	}
	return bytePng
}
func NoToPng(contentString string) *bytes.Buffer{
	qrCode, err := qr.Encode(contentString, qr.M, qr.Auto)
	// create the output file
	bytePng  := new(bytes.Buffer)
	err=png.Encode(bytePng,qrCode)
	if err != nil {
		log.Error(err)
	}

	if err != nil {
		log.Error(err)
	}
	return bytePng
}