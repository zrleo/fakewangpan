package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"lzr.com/wangpan/meta"
	"lzr.com/wangpan/util"
)

// UploadHandler： 上传文件处理
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回前端页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal error %s")
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 上传文件，存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("internal error %s", err.Error())
			return
		}
		defer file.Close()
		fileMeta := meta.Filemeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf(" create file failed err %s", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Printf("======================filehash is %s====================\n", fileMeta.FileSha1)
		meta.UpdateFileMeta(fileMeta)
		if err != nil {
			fmt.Printf("failed to save data to new file, error %s", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSuccessHandler: 上传成功提示
func UploadSuccessHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload finished")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// 获取参数
	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("ContentType", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// FileDelHandler: 删除文件
func FileDelHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	fm := meta.GetFileMeta(filehash)
	err := os.Remove(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	meta.RemoveFileMeta(filehash)
	fmt.Printf("删除成功， %s", filehash)
}

// FilenameUpdateHandler修改文件名字接口
func FilenameUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	filehash := r.Form.Get("filehash")
	fileNewName := r.Form.Get("newName")
	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fm := meta.GetFileMeta(filehash)
	fm.FileName = fileNewName
	meta.UpdateFileMeta(fm)
	data, err := json.Marshal(fm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
