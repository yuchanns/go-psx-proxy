package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func download(w http.ResponseWriter, req *http.Request) {
	pathSplit := strings.Split(req.URL.Path, "/")
	fileName := pathSplit[len(pathSplit) - 1]
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Cannot get running file directory:", err)
		return
	}
	filePath := strings.Join([]string{dir, "downloads", fileName}, "/")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File %s is not exists, need to be cached.\nPlease download file by link:\n%s\nAnd move the file to path %s\n", fileName, req.URL, filePath)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Cannot open file: %s\n", err)
		return
	}
	fileHeader := make([]byte, 512)
	if _, err := file.Read(fileHeader); err != nil {
		fmt.Printf("Cannot read fileHeader: %s\n", err)
		return
	}

	fileStat, _ := file.Stat()
	w.Header().Set("Content-Disposition", "attachment; filename=" + fileName)
	w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))

	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("Cannot seek file offset 0: %s\n", err)
		return
	}
	if _, err := io.Copy(w, file); err != nil {
		fmt.Printf("Cannot copy file: %s\n", err)
		return
	}
}

func main() {
	http.HandleFunc("/", download)
	http.ListenAndServe("0.0.0.0:1088", nil)
}
