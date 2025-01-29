package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const uploadDir = "uploads"
const processedDir = "processed"

func main() {
	// Ensure directories exist
	os.MkdirAll(uploadDir, os.ModePerm)
	os.MkdirAll(processedDir, os.ModePerm)

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	savePath := filepath.Join(uploadDir, handler.Filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	// Process the file with Python
	outputFile := filepath.Join(processedDir, handler.Filename[:len(handler.Filename)-4]+"_summary.csv")
	cmd := exec.Command("python3", "fit_processor.py", savePath)
	if err := cmd.Run(); err != nil {
		http.Error(w, "Failed to process file", http.StatusInternalServerError)
		return
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to process file: %s\n%s", err, output), http.StatusInternalServerError)
		return
	}

	// Return processed file link
	fmt.Fprintf(w, "File processed. <a href='/download?file=%s'>Download</a>", filepath.Base(outputFile))
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
