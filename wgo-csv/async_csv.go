package wgo_csv

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Provides an asynchronous csv file generation to avoid interface timeouts.
// 1. The backend uses AsyncCSV or AsyncWriteCSV in the interface function and implements the corresponding input function.
// 2. The front end can then poll the interface with the same request parameters and finally get the corresponding file address.
// If an error occurs, the error during sync is returned immediately, and the error after async is written to the csv file.
var (
	download = make(map[string]string)
	mu       sync.Mutex
)

// AsyncCSV generates a CSV file using a key derived from the MD5 hash of the provided request.
// It calls AsyncWriteCSV with the hashed key to manage file creation and writing.
func AsyncCSV(req interface{}, filePrefix string, writerFunc func(*csv.Writer) (err error)) (filename string, err error) {
	return AsyncWriteCSV(Md5Hash(req), filePrefix, writerFunc)
}

// AsyncWriteCSV generates a CSV file and associates it with the given key.
// It checks if the file is already being created or has been created.
// If not, it starts the CSV file creation process and calls writerFunc to populate the file.
// The function returns the file path once the file is ready or an error if the process fails.
// If the file is being created or was already created, it returns the existing file path.
func AsyncWriteCSV(key string, filePrefix string, writerFunc func(*csv.Writer) (err error)) (filename string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("AsyncCSV Panic: %v", r)
		}
	}()
	if p := GetFilepath(key); p != "" {
		return p, nil
	}
	//检查CSV文件是否正在组装或已经组装完毕
	finish := make(chan struct{})
	errChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("writerfunc Panic: %v", r)
			}
		}()
		defer close(errChan)
		defer close(finish)
		fileName := fmt.Sprintf("%s_%s.csv", filePrefix, time.Now().Format("20060102_150405"))
		filePath := filepath.Join("/app/cp/", fileName)
		file, err := os.Create(filePath)
		if err != nil {
			errChan <- err
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()

		err = writerFunc(writer)
		if err != nil {
			writer.Write([]string{"err:", err.Error()})
			errChan <- err
			return
		}
		mu.Lock()
		download[key] = fileName //确认文件已经完成
		mu.Unlock()
		finish <- struct{}{}
	}()
	select {
	//阻塞两秒
	case <-time.After(time.Second * 2):
		return "", nil
	case err = <-errChan:
		delete(download, key)
		return "", err
	case <-finish:
		return GetFilepath(key), nil
	}
}

// GetFilepath retrieves the file path associated with the given key.
// If the file is complete, it returns the path and schedules its deletion after 5 seconds.
// If the file is still being assembled, it returns an empty string or "wait".
func GetFilepath(key string) string {
	mu.Lock()
	defer mu.Unlock()
	if download == nil {
		download = make(map[string]string)
	}
	if val, ok := download[key]; ok && val != "" { //组装完毕
		go func() { //todo:cancel async delete
			time.Sleep(time.Second * 5)
			delete(download, key)
		}()
		return val
	} else if !ok {
		download[key] = ""
	}
	return download[key]
}

func Md5Hash(obj interface{}) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%+v", obj)))
	return hex.EncodeToString(hash.Sum(nil))
}
