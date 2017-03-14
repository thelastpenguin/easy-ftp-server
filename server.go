package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"time"

	"github.com/koofr/graval"
)

type User struct {
	Username string
	Password string
	FsRoot   string
}

type FileSystemDriver struct {
	Users      []User
	AuthedUser *User
}

func (d *FileSystemDriver) resolvePath(path string) string {
	return filepath.Join(d.AuthedUser.FsRoot, filepath.Clean(path))
}

func (d *FileSystemDriver) Authenticate(username, password string) bool {
	for _, user := range d.Users {
		if user.Username == username && user.Password == password {
			d.AuthedUser = &user
			return true
		}
	}
	return false
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

type FSDriverFactory struct {
	Users []User
}

func (d *FSDriverFactory) NewDriver() (graval.FTPDriver, error) {
	driver := FileSystemDriver{
		Users:      d.Users,
		AuthedUser: nil,
	}
	return graval.FTPDriver(&driver), nil
}

type FTPServerJsonConfig struct {
	Host  string
	Port  int
	Users []User
}

func main() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	configFile, err := os.Open(path.Join(usr.HomeDir, ".easyftp"))

	if err != nil {
		fmt.Printf("Error! No config file at %s\n", path.Join(usr.HomeDir, ".easyftp"))
		return
	}

	var config FTPServerJsonConfig
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		fmt.Printf("Error! Error parsing JSON in config at ~/.easyftp. Error: %v\n", err)
		return
	}

	factory := &FSDriverFactory{
		Users: config.Users,
	}

	server := graval.NewFTPServer(&graval.FTPServerOpts{
		ServerName: "Example FTP server",
		Factory:    graval.FTPDriverFactory(factory),
		Hostname:   config.Host,
		Port:       config.Port,
		PassiveOpts: &graval.PassiveOpts{
			ListenAddress: config.Host,
			NatAddress:    config.Host,
			PassivePorts: &graval.PassivePorts{
				Low:  42000,
				High: 45000,
			},
		},
	})

	log.Printf("Example FTP server listening on %s:%d\n", config.Host, config.Port)
	log.Printf("Loaded %d users.\n", len(config.Users))

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
