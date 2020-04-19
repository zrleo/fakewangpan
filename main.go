package main

import (
	"fmt"
	"net/http"

	"lzr.com/wangpan/handler"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSuccessHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/delete", handler.FileDelHandler)
	http.HandleFunc("/file/update", handler.FilenameUpdateHandler)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Printf("failed to start server, err %s", err.Error())
	}
}
