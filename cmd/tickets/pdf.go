package main

import (
	"bytes"
	"image/png"
	"log"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/skip2/go-qrcode"
)

func pdfGenerator(inputPDF string, outputPDF string, name, url string) bool {

	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Println("Error generating QR code:", err)
		return false
	}

	// 2. Convert QR code to PNG bytes in memory
	var qrBuf bytes.Buffer
	if err := png.Encode(&qrBuf, qr.Image(500)); err != nil {
		log.Println("Error encoding QR code:", err)
		return false
	}

	conf := model.NewDefaultConfiguration()
	onTop, update := true, false

	qrDesc := "rot:0, scale:.25 abs, pos:c, offset:155 170"
	qrWm, err := api.ImageWatermarkForReader(&qrBuf, qrDesc, onTop, update, types.POINTS)
	if err != nil {
		log.Println("Error creating image watermark:", err)
		return false
	}
	textDesc := "font:Helvetica, scalefactor:.2, rot:0, pos:c, offset:150 250, color:#000000"
	textWM, err := api.TextWatermark(name, textDesc, true, false, types.POINTS)
	if err != nil {
		log.Println("Error creating text watermark:", err)
		return false
	}

	if err := api.AddWatermarksFile(inputPDF, outputPDF, nil, qrWm, conf); err != nil {
		log.Println("Error applying watermark:", err)
		return false
	}

	if err := api.AddWatermarksFile(outputPDF, outputPDF, nil, textWM, conf); err != nil {
		log.Println("Error applying watermark:", err)
		return false
	}

	log.Printf("✅ Ticket %s for %s generated successfully \n", outputPDF, name)
	return true
}
