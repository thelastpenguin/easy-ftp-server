package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/koofr/graval"
)

type FileSystemDriver struct {

	// TODO: add multiple users
	Username string
	Password string
	FsRoot   string
}

func (d *FileSystemDriver) resolvePath(path string) string {
	return filepath.Join(d.FsRoot, filepath.Clean(path))
}

func (d *FileSystemDriver) Authenticate(username, password string) bool {
	return username == d.Username && password == d.Password
}

func (d *FileSystemDriver) Bytes(path string) int64 {
	realPath := d.resolvePath(path)

	if info, err := os.Stat(realPath); err == nil {
		return info.Size()
	}
	return -1
}

func (d *FileSystemDriver) ModifiedTime(path string) (time.Time, bool) {
	realPath := d.resolvePath(path)

	if info, err := os.Stat(realPath); err == nil {
		return info.ModTime(), true
	}
	return time.Now(), false
}

func (d *FileSystemDriver) ChangeDir(path string) bool {
	realPath := d.resolvePath(path)
	if info, err := os.Stat(realPath); err == nil {
		return info.IsDir()
	}
	return false
}

func (d *FileSystemDriver) DirContents(path string) ([]os.FileInfo, bool) {
	realPath := d.resolvePath(path)
	files, err := ioutil.ReadDir(realPath)
	fmt.Printf("listing directory: %s\n", realPath)
	if err == nil {
		return files, true
	}
	return nil, false
}

func (d *FileSystemDriver) DeleteDir(path string) bool {
	if err := os.RemoveAll(d.resolvePath(path)); err != nil {
		return true
	}
	return false
}

func (d *FileSystemDriver) DeleteFile(path string) bool {
	if err := os.Remove(d.resolvePath(path)); err == nil {
		return true
	}
	return false
}

func (d *FileSystemDriver) Rename(from_path, to_path string) bool {
	from := d.resolvePath(from_path)
	to := d.resolvePath(to_path)

	return os.Rename(from, to) == nil
}

func (d *FileSystemDriver) MakeDir(path string) bool {
	realPath := d.resolvePath(path)

	err := os.MkdirAll(realPath, 0700)
	fmt.Println(err)
	return err == nil
}

func (d *FileSystemDriver) GetFile(path string, position int64) (io.ReadCloser, bool) {
	realPath := d.resolvePath(path)

	f, err := os.OpenFile(realPath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, false
	}

	f.Seek(position, 0)

	return ioutil.NopCloser(f), true
}

func (d *FileSystemDriver) PutFile(path string, reader io.Reader) bool {
	realPath := d.resolvePath(path)

	f, err := os.Create(realPath)
	if err != nil {
		return false
	}
	defer f.Close()

	io.Copy(f, reader)

	return true
}

type DriverFactory struct{}

func (d *DriverFactory) NewDriver() (graval.FTPDriver, error) {
	driver := FileSystemDriver{
		Username: "username",
		Password: "password",
		FsRoot:   ".",
	}
	return graval.FTPDriver(&driver), nil
}

func main() {
	fmt.Println("HI!")
	host := "127.0.0.1"
	port := 8021
	username := "test"
	password := "test"

	factory := &DriverFactory{}

	server := graval.NewFTPServer(&graval.FTPServerOpts{
		ServerName: "Example FTP server",
		Factory:    graval.FTPDriverFactory(factory),
		Hostname:   host,
		Port:       port,
		PassiveOpts: &graval.PassiveOpts{
			ListenAddress: host,
			NatAddress:    host,
			PassivePorts: &graval.PassivePorts{
				Low:  42000,
				High: 45000,
			},
		},
	})

	log.Printf("Example FTP server listening on %s:%d", host, port)
	log.Printf("Access: ftp://%s:%s@%s:%d/", username, password, host, port)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
