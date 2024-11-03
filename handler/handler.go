package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"arif/service"
)

const uploadDir = "./uploads/"

// UploadPDFHandler handles PDF file uploads
func UploadPDFHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form, with a max upload size of 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the "pdf" field
	file, handler, err := r.FormFile("pdf")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Ensure the uploaded file is a PDF
	if filepath.Ext(handler.Filename) != ".pdf" {
		http.Error(w, "Only PDF files are allowed", http.StatusUnsupportedMediaType)
		return
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	res, err := service.ProcessPDF(r.Context(), fileBytes)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error processing file", http.StatusInternalServerError)
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
	}

	_, err = w.Write(resBytes)
	if err != nil {
		http.Error(w, "Error writing response", http.StatusInternalServerError)
	}

	// Save the file to the uploads directory
	//filePath := filepath.Join(uploadDir, handler.Filename)
	//dst, err := os.Create(filePath)
	//if err != nil {
	//	http.Error(w, "Unable to save the file", http.StatusInternalServerError)
	//	return
	//}
	//defer dst.Close()
	//
	//if _, err := io.Copy(dst, file); err != nil {
	//	http.Error(w, "Error saving the file", http.StatusInternalServerError)
	//	return
	//}

	//fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}
