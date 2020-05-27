package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/browser"
)

type fileSystem struct {
	RootDir string
}

func (fs *fileSystem) Open(name string) (http.File, error) {
	log.Printf("debug: serving name=%q", name)

	name = filepath.Join(fs.RootDir, name)
	return os.Open(name)
}

func doServe(cfg Config, dir string, openBrowser bool) error {
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
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Serve(l)
	}()

	if openBrowser {
		var u url.URL
		u.Scheme = "http"
		u.Host = l.Addr().String()
		if err := browser.OpenURL(u.String()); err != nil {
			log.Printf("warn: failed to open url %q: %v", u, err)
		}
	}

	return <-errCh
}
