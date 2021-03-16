package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

func NewLogger(writer io.Writer, dir string) (*Logger, error) {

	// init dir
	if !fileExists(dir, true) {
		err := os.Mkdir(dir, 0777)
		if err != nil {
			return nil, err
		}
	}

	// init file
	filename := time.Now().Format(`2006-01-02`) + `.log`
	if !fileExists(filename, false) {
		file, err := os.OpenFile(path.Join(dir, filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			return nil, err
		}
		writer = io.MultiWriter(writer, file)
	}

	return &Logger{
		l:   log.New(writer, ``, log.LstdFlags),
		dir: dir,
	}, nil
}

type Logger struct {
	l   *log.Logger
	dir string
}

func (l *Logger) GetLogger() *log.Logger {
	return l.l
}

// Удаляет старые записи, можно было бы чекать по названию ( парсить дату из названия)
func (l *Logger) RemoveOldFile(days int) error {
	files, err := ioutil.ReadDir(path.Base(l.dir))
	if err != nil {
		return err
	}

	for _, file := range files {
		hours := int(time.Since(file.ModTime()).Hours())
		if hours > days*24 {
			if err = os.Remove("./" + l.dir + `/` + file.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

func fileExists(filename string, isDir bool) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir() == isDir
}
