package meta

// Filemeta ：文件元信息
type Filemeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]Filemeta

func init() {
	fileMetas = make(map[string]Filemeta)
}

// UpdateFileMeta : 新增/更新文件元信息
func UpdateFileMeta(fmeta Filemeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func GetFileMeta(filesha1 string) Filemeta {
	return fileMetas[filesha1]
}


func GetAllFile(count int)[]Filemeta  {
	fMetaArray := make([]Filemeta, len(fileMetas))
	for _, v := range fileMetas{
		fMetaArray = append(fMetaArray, v)
	}
	return fMetaArray[0:count]
}

// 删除文件，非物理删除，将唯一索引hash值从map中删除
func RemoveFileMeta(filehash string)  {
	delete(fileMetas, filehash)
}