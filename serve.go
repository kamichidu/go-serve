package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type fileSystem struct {
	RootDir string
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	log.Printf("debug: serving name=%q", name)

	name = filepath.Join(fs.RootDir, name)
	return os.Open(name)
}

func doServe(cfg Config, dir string) error {
	func(wd string, err error) {
		log.Printf("debug: wd=%q err=%q", wd, err)
	}(os.Getwd())
	log.Printf("info: serving directory %s ...", dir)
	l, err := net.Listen("tcp", cfg.ListenAddress)
	if err != nil {
		return err
	}
	defer l.Close()
	log.Printf("info: listening on %s ...", l.Addr())

	var hdlr http.Handler = http.FileServer(&fileSystem{
		RootDir: dir,
	})

	srv := &http.Server{
		Handler:      hdlr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return srv.Serve(l)
}
