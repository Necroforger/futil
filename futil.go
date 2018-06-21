package futil

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"
)

var (
	// ErrSkipDir can be returned by a walk function to skip walking a directory
	ErrSkipDir = errors.New("Skip directory")
)

// Ls lists the contents of a directory
// And sorts them with directories coming first
//    dir : directory to list the contents of
func Ls(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	sort.Sort(FileInfoByType(files))
	return files, err
}

// Walk recursively walks through a directory
//    dir   : directory to walk through
//    fn    : function called for every file
//            in the directory tree
func Walk(dir string, fn func(string, os.FileInfo) error) error {
	info, err := Ls(dir)
	if err != nil {
		return err
	}

	for _, v := range info {
		if v.IsDir() {
			err = fn(dir, v)
			if err != nil {
				if err == ErrSkipDir {
					continue
				}
				return err
			}
			err = Walk(path.Join(dir, v.Name()), fn)
			if err != nil {
				return err
			}
		} else {
			err := fn(dir, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WalkFromTo compares two directory trees
//    from  :  from directory
//    to    :  to directory
//    fn    :  walk function
func WalkFromTo(from string, to string, fn func(from string, to string, info os.FileInfo) error) error {
	return Walk(from, func(source string, info os.FileInfo) error {
		return fn(source, path.Join(to, strings.TrimPrefix(source, from)), info)
	})
}

// Cp copies a file
//    from  : location to copy from
//    to    : destination path for the new copy
func Cp(from, to string) error {
	fa, err := os.OpenFile(from, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer fa.Close()
	stat, err := fa.Stat()
	if err != nil {
		return err
	}
	fb, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE, stat.Mode())
	if err != nil {
		return err
	}
	defer fb.Close()
	_, err = io.Copy(fb, fa)
	return err
}

// CpDir recursively copies a directory
//    from  : directory to copy from
//    to    : location to copy to
func CpDir(from, to string) error {
	return WalkFromTo(from, to, func(f, t string, info os.FileInfo) error {
		if info.IsDir() {
			return nil
		}
		return Cp(path.Join(from, info.Name()), path.Join(to, info.Name()))
	})
}

// Mv moves a file from one location to another
// if it fails to move the file, it will attempt to copy
//    from  : location to move from
//    to    : location to move to
func Mv(from, to string) error {
	err := os.Rename(from, to)
	if err != nil {
		return Cp(from, to)
	}
	return nil
}

// MvDir moves a directory from one location to another
//     from : location to move from
//     to   : location to move to
func MvDir(from, to string) error {
	// Attempt to rename the directory
	err := os.Rename(from, to)
	if err != nil {
		// If renaming the directory fails, fall back to
		// copying or moving the files individually
		return WalkFromTo(from, to, func(f, t string, info os.FileInfo) error {
			if info.IsDir() {
				return nil
			}
			return Mv(path.Join(from, info.Name()), path.Join(to, info.Name()))
		})
	}
	return nil
}

// ZipDir zips the directory dir to dest
//     path : path to zip. if a directory, recursively zip it.
//     dest : ouput zip file
func ZipDir(source, dest string) error {
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	// Add a slash to the end of the path
	// So the prefix is trimmed properly later on
	source = path.Clean(source) + "/"

	zwr := zip.NewWriter(f)
	err = Walk(source, func(p string, info os.FileInfo) error {

		// remove the root folder name from the archive
		npath := strings.TrimPrefix(p, source)

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// Remove the root directory name from the archive
		header.Name = path.Join(npath, info.Name())

		// List the file as a directory in the archive
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		hdr, err := zwr.CreateHeader(header)
		if err != nil {
			return err
		}

		// The file is a directory, we do not need to copy anything into it
		if info.IsDir() {
			return nil
		}

		f, err := os.OpenFile(path.Join(p, info.Name()), os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(hdr, f)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if err = zwr.Close(); err != nil {
		return err
	}

	return nil
}

// Unzip unzips a directory
//     from : location of zip file
//     to   : destination folder
func Unzip(from, to string) error {
	f, err := os.OpenFile(from, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	rd, err := zip.NewReader(f, stat.Size())
	if err != nil {
		return err
	}

	os.MkdirAll(to, 0666)

	for _, v := range rd.File {
		// Do not unzip directories
		if v.FileInfo().IsDir() {
			os.MkdirAll(path.Join(to, v.Name), 0666)
			continue
		}

		zf, err := v.Open()
		if err != nil {
			return err
		}
		defer zf.Close()
		fpath := path.Join(to, v.Name)

		df, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer df.Close()

		_, err = io.Copy(df, zf)
		if err != nil {
			return err
		}
	}

	return nil
}
