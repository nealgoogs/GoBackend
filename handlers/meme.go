package handlers

import (
	"image"
	"image/png"
	"net/http"

	"github.com/fogleman/gg"
)

// GenerateMeme handles the meme generation request
func GenerateMeme(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the uploaded file
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get the meme text
	memeText := r.FormValue("memeText")

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		http.Error(w, "Unable to decode image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a new context for drawing
	const W = 800 // width of the meme image
	const H = 600 // height of the meme image
	dc := gg.NewContext(W, H)

	// Draw the original image
	dc.DrawImage(img, 0, 0)

	// Load a common system font (make sure the path is correct)
	if err := dc.LoadFontFace("C:\\Windows\\Fonts\\Arial.ttf", 36); err != nil {
		http.Error(w, "Unable to load font: "+err.Error(), http.StatusInternalServerError)
		return
	}
	dc.SetRGB(1, 1, 1) // White color

	// Draw the text on the image
	dc.DrawStringAnchored(memeText, W/2, H-50, 0.5, 0.5)

	// Set the response header to indicate an image is being returned
	w.Header().Set("Content-Type", "image/png")

	// Write the image directly to the response
	if err := png.Encode(w, dc.Image()); err != nil {
		http.Error(w, "Unable to encode meme image: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
