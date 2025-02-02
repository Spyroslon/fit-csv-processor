package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const uploadDir = "uploads"
const processedDir = "processed"

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	os.MkdirAll(uploadDir, os.ModePerm)
	os.MkdirAll(processedDir, os.ModePerm)

	// Register specific handlers first
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	// Use a custom handler for the homepage
	http.HandleFunc("/", serveIndex)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe("localhost:8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Ensure only POST requests are handled
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received file upload request")

	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		log.Println("Failed to parse form:", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println("Failed to read file:", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		log.Println("Failed to save file:", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	// Process the file with Python
	outputFile := filepath.Join(processedDir, handler.Filename[:len(handler.Filename)-4]+"_summary.csv")
	cmd := exec.Command("python", "fit_processor.py", "--i", savePath, "--o", processedDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Python processing failed:", string(output))
		http.Error(w, fmt.Sprintf("Failed to process file: %s\n%s", err, output), http.StatusInternalServerError)
		return
	}

	log.Println("Python script output:", string(output))

	// Check if the output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		log.Println("Processed file not found")
		http.Error(w, "Processed file not found", http.StatusInternalServerError)
		return
	}

	// ✅ Make sure we respond with JSON
	response := map[string]string{
		"message":       "File processed successfully",
		"download_link": fmt.Sprintf("/download?file=%s", filepath.Base(outputFile)),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // ✅ Ensures correct response
	json.NewEncoder(w).Encode(response)

	log.Println("JSON response sent successfully")
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "File not specified", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(processedDir, file)
	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, filePath)
}
