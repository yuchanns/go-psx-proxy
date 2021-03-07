package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func download(w http.ResponseWriter, req *http.Request) {
	pathSplit := strings.Split(req.URL.Path, "/")
	fileName := pathSplit[len(pathSplit)-1]
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Cannot get running file directory:", err)
		return
	}
	filePath := strings.Join([]string{dir, "downloads", fileName}, "/")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("File %s is not exists, need to be cached.\nPlease download file by link:\n%s\nAnd move the file to path %s\n", fileName, req.URL, filePath)
		return
	} else if err != nil {
		fmt.Printf("Cannot state file: %s\n", err)
		return
	}
	fmt.Printf("Cache hit: %s\n", fileName)
	http.ServeFile(w, req, filePath)
}

func main() {
	http.HandleFunc("/", download)
	http.ListenAndServe("0.0.0.0:1088", nil)
}
