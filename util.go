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
