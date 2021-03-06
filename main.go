package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type fallback struct {
	defaultPath string
	fs          http.FileSystem
}

func (fb fallback) Open(path string) (http.File, error) {
	f, err := fb.fs.Open(path)
	if os.IsNotExist(err) {
		log.Printf("Serving %s instead of %s: %v", fb.defaultPath, path, err)
		return fb.fs.Open(fb.defaultPath)
	}
	log.Printf("Serving %s", path)
	return f, err
}

func main() {
	AutoInjection(os.Args[1:]...)

	var fileSystem = fallback{
		defaultPath: "/index.html",
		fs:          http.Dir("/srv/http"),
	}

	handler := http.FileServer(fileSystem)

	http.Handle("/", handler)

	log.Printf("Listening at 0.0.0.0:80")
	log.Fatalln(http.ListenAndServe(":80", nil))
}

func readEnv() (env map[string]string) {
	env = make(map[string]string)
	for _, setting := range os.Environ() {
		pair := strings.SplitN(setting, "=", 2)
		env[pair[0]] = pair[1]
	}
	return
}

func AutoInjection(files ...string) {
	for _, patern := range files {
		fs, err := filepath.Glob(patern)
		if err != nil {
			log.Printf("[WARN] Injecting for the pattern %s failed : %v", patern, err)
		} else if len(fs) == 0 {
			log.Printf("[WARN] Injecting for the pattern %s returned no files", patern)
		}

		for _, f := range fs {
			err := InjectEnv(f)
			if err != nil {
				log.Printf("[FAIL] Injecting on file %s: error %v", f, err)
			} else {
				log.Printf("[DONE] Injecting on file %s", f)
			}
		}
	}
}

func InjectEnv(fileName string) error {
	env := readEnv()

	t := template.New(filepath.Base(fileName))
	if delims, ok := env["DELIM"]; ok {
		if delimEnd, ok := env["DELIM_END"]; ok {
			t = t.Delims(delims, delimEnd)
		} else {
			t = t.Delims(delims, delims)
		}
	}

	t, err := t.ParseFiles(fileName)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	file.Truncate(0)

	return t.Execute(file, env)
}
