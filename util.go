package futil

import "os"

// FileInfoByType allows the sorting of a slice of file infos by type
type FileInfoByType []os.FileInfo

func (f FileInfoByType) Less(a, b int) bool {
	var aa, bb = 1, 1
	if f[a].IsDir() {
		aa = 0
	}
	if f[b].IsDir() {
		bb = 0
	}
	return aa < bb
}

func (f FileInfoByType) Len() int {
	return len(f)
}

func (f FileInfoByType) Swap(a, b int) {
	f[a], f[b] = f[b], f[a]
}

// SplitFileInfo splits a slice of file infos into separate slices of files and
// directories. Should pass the results of Ls or ioutil.ReadDir to this function
// To sort the output
//    info : slice of os.FileInfo
func SplitFileInfo(infos []os.FileInfo) (dirs []os.FileInfo, files []os.FileInfo) {
	dirs = []os.FileInfo{}
	files = []os.FileInfo{}

	for _, v := range infos {
		if v.IsDir() {
			dirs = append(dirs, v)
		} else {
			files = append(files, v)
		}
	}

	return
}
