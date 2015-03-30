// targiz project targiz.go
package targz

import (
	//"fmt"
	"io"
	"os"
	// "time"
	"archive/tar"
	//"bytes"
	"compress/gzip"
	"path"
	"path/filepath"
)

func Uncompress(fr io.Reader, basepath string) error {
	// file read
	//fr, err := os.Open(srcpath)
	//if err != nil {
	//	return err
	//}
	//defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}
	defer gr.Close()

	// tar read
	tr := tar.NewReader(gr)

	// 读取文件
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// 显示文件
		//fmt.Println(h.Name, h.FileInfo().IsDir())
		if h.FileInfo().IsDir() {
			//fmt.Println(h.Name, ">>>>", h.Size)
			if _, err := os.Stat(path.Join(basepath, h.Name)); os.IsNotExist(err) {
				errs := os.MkdirAll(path.Join(basepath, h.Name), 0766)
				if errs != nil {
					return errs
				}
			}
			//err := os.MkdirAll(path.Join(basepath, h.Name), 0766)
			//if err != nil {
			//	fmt.Println(err)
			//	return err
			//}

			continue
		}

		// 打开文件
		if _, err := os.Stat(path.Dir(path.Join(basepath, h.Name))); os.IsNotExist(err) {
			errs := os.MkdirAll(path.Dir(path.Join(basepath, h.Name)), 0766)
			if errs != nil {
				return errs
			}
		}
		fw, err := os.OpenFile(path.Join(basepath, h.Name), os.O_CREATE|os.O_WRONLY, 0644 /*os.FileMode(h.Mode)*/)
		if err != nil {
			return err
		}

		// 写文件
		_, err = io.Copy(fw, tr)
		fw.Close()
		if err != nil {
			return err
		}

	}

	return nil
}

func Compress(distWr io.Writer, srcpath ...string) error {
	// Create a buffer to write our archive to.
	//buf := new(bytes.Buffer)
	gw := gzip.NewWriter(distWr)
	defer gw.Close()
	// Create a new tar archive.
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, v := range srcpath {
		//fmt.Println(v)
		fw, err := os.OpenFile(v, os.O_RDONLY, 0766 /*os.FileMode(h.Mode)*/)
		if err != nil {
			return err
		}

		flInfo, err := fw.Stat()
		if os.IsNotExist(err) {
			return err
		}

		if flInfo.IsDir() {

			filepath.Walk(v, func(filename string, fi os.FileInfo, err error) error {
				//fmt.Println(filename)
				target, _ := os.Readlink(filename)

				hdr, err := tar.FileInfoHeader(fi, target)
				if err != nil {
					return err
				}
				hdr.Name = filename
				//if fi.IsDir() {
				//	hdr.Name += "/"
				//}
				if err := tw.WriteHeader(hdr); err != nil {
					return err
				}

				if fi.IsDir() {
					//fmt.Println(hdr.Name)
					return nil
				}

				tempfw, err := os.OpenFile(filename, os.O_RDONLY, 0644 /*os.FileMode(h.Mode)*/)
				if err != nil {
					return err
				}

				io.Copy(tw, tempfw)
				tempfw.Close()

				return nil
			})
		} else {
			target, _ := os.Readlink(v)

			hdr, err := tar.FileInfoHeader(flInfo, target)
			if err != nil {
				return err
			}

			//hdr := &tar.Header{
			//	Name: v,
			//	Size: flInfo.Size(),
			//}
			//fmt.Println(hdr.Size, hdr.Name)
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			io.Copy(tw, fw)

		}
		fw.Close()

	}
	return nil
}
