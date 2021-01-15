package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"strings"
)

// DirUpLevel returns dir name of upper `level`
func DirUpLevel(p string, level int) string {
	if level > 0 {
		panic("level should <= 0")
	}
	dirs := strings.Split(p, "/")
	return strings.Join(dirs[:len(dirs)+level], "/")
}

// GetCronFilePath returns the cron file path
func GetCronFilePath(user string) string {
	f := os.Getenv("CRON_FILE_PATH")
	if f == "" {
		return fmt.Sprintf("/var/spool/cron/%s", user)
	}
	return f
}

// GetWhoAmI returns current user name
func GetWhoAmI() string {
	u := os.Getenv("USER")
	if u == "" {
		return "automation"
	}
	return u
}

// IsFileExist returns true if a file `filename` exist, returns false otherwise
func IsFileExist(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirExist returns true if a dir `pathname` exist, returns false otherwise
func IsDirExist(pathname string) bool {
	info, err := os.Stat(pathname)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// LsDir list files under pathname and return files that is a dir. return error if pathname is not a dir
func LsDir(pathname string) ([]string, error) {
	if !IsDirExist(pathname) {
		return nil, fmt.Errorf("pathname %s is not a directory", pathname)
	}
	dirs := []string{}
	f, err := os.OpenFile(pathname, os.O_RDONLY, 0755)
	defer f.Close()

	if err != nil {
		log.Println(err)
		return nil, err
	}
	subfiles, err := f.Readdir(0)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, fileinfo := range subfiles {
		if fileinfo.IsDir() {
			dirs = append(dirs, fileinfo.Name())
		}
	}
	return dirs, nil
}

// GetCacheDir returns goat cache dir. If env GOAT_CACHE_DIR exist, use its value. Otherwise return /tmp/goat
func GetCacheDir() string {
	d := os.Getenv("GOAT_CACHE_DIR")
	if d == "" {
		return "/tmp/goat"
	}
	return d
}

// Sha256Sum use SHA256 hash function to digest the input string and return the hash string
func Sha256Sum(input string) string {
	shaSum := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", shaSum)
}
