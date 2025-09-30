package tickets

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/types"
	"github.com/skip2/go-qrcode"
)

// PdfGenerator generates a ticket PDF with QR code and name/company/ticketType rendered as PNG.
func PdfGenerator(inputPDF, outputPDF, name, url, company, ticketType string) bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic recovered in PdfGenerator: %v", r)
			log.Printf("Parameters: name=%s, company=%s, ticketType=%s", name, company, ticketType)
		}
	}()

	conf := model.NewDefaultConfiguration()
	onTop, update := true, false

	// 1️⃣ QR code watermark
	qr, err := qrcode.New(url, qrcode.Medium)
	if err != nil {
		log.Println("Error generating QR code:", err)
		return false
	}
	var qrBuf bytes.Buffer
	if err := png.Encode(&qrBuf, qr.Image(500)); err != nil {
		log.Println("Error encoding QR code:", err)
		return false
	}
	qrWM, err := api.ImageWatermarkForReader(&qrBuf, "rot:0, scale:.25 abs, pos:c, offset:155 170", onTop, update, types.POINTS)
	if err != nil {
		log.Println("Error creating QR watermark:", err)
		return false
	}

	tempPDF := outputPDF + ".tmp"
	if err := api.AddWatermarksFile(inputPDF, tempPDF, nil, qrWM, conf); err != nil {
		log.Println("Error applying QR watermark:", err)
		return false
	}

	// 2️⃣ Render name/company/ticketType as PNG (regular CJK font)
	const (
		imgWidth  = 600
		imgHeight = 220
	)

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.SetColor(color.Transparent)
	dc.Clear()
	dc.SetColor(color.Black)

	// Paths to installed CJK font (regular)
	fontCJK := "NotoSerifCJK-VF.ttf.ttc"

	y := 50.0
	lineSpacing := 50.0

	// TicketType → largest
	if ticketType != "" {
		if err := dc.LoadFontFace(fontCJK, 40); err != nil {
			log.Println("Error loading CJK font:", err)
			return false
		}
		offsets := []struct{ x, y float64 }{
			{0, 0}, {1, 0}, {0, 1}, {1, 1},
		}

		for _, off := range offsets {
			dc.DrawStringAnchored(ticketType, float64(imgWidth)/2+off.x, y+off.y, 0.5, 0.5)
		}
		dc.DrawStringAnchored(ticketType, float64(imgWidth)/2, y, 0.5, 0.5)
		y += lineSpacing
	}

	// Name → medium
	if name != "" {
		if err := dc.LoadFontFace(fontCJK, 30); err != nil {
			log.Println("Error loading CJK font:", err)
			return false
		}
		offsets := []struct{ x, y float64 }{
			{0, 0}, {1, 0}, {0, 1}, {1, 1},
		}

		for _, off := range offsets {
			dc.DrawStringAnchored(name, float64(imgWidth)/2+off.x, y+off.y, 0.5, 0.5)
		}
		dc.DrawStringAnchored(name, float64(imgWidth)/2, y, 0.5, 0.5)
		y += lineSpacing + 15.0
	}

	// Company → smaller
	if company != "" {
		if err := dc.LoadFontFace(fontCJK, 23); err != nil {
			log.Println("Error loading CJK font:", err)
			return false
		}
		offsets := []struct{ x, y float64 }{
			{0, 0}, {1, 0}, {0, 1}, {1, 1},
		}

		for _, off := range offsets {
			dc.DrawStringWrapped(company, float64(imgWidth)/2+off.x, y+off.y, 0.5, 0.5, 400, 1.5, gg.AlignCenter)
		}

		y += 18*1.5 + lineSpacing
	}

	var textBuf bytes.Buffer
	img := dc.Image()

	// Ensure alpha channel
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	if err := png.Encode(&textBuf, rgba); err != nil {
		log.Println("Error encoding text image:", err)
		return false
	}

	textWM, err := api.ImageWatermarkForReader(&textBuf, "rot:0, scale:.5 abs, pos:c, offset:150 285", onTop, update, types.POINTS)
	if err != nil {
		log.Println("Error creating text PNG watermark:", err)
		return false
	}

	// 3️⃣ Apply text PNG watermark
	if err := api.AddWatermarksFile(tempPDF, outputPDF, nil, textWM, conf); err != nil {
		log.Println("Error applying text PNG watermark:", err)
		return false
	}

	_ = os.Remove(tempPDF)

	log.Printf("✅ Ticket %s for %s generated successfully\n", outputPDF, name)
	return true
}
