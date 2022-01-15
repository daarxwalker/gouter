package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Gouter interface {
	run(config)
	readConfig() config
}

type gouter struct {
	router *mux.Router
}

func newGouter() Gouter {
	return &gouter{}
}

func (g gouter) handler(p *httputil.ReverseProxy, remote *url.URL) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Header.Set("X-Origin-Host", remote.Host)
		r.URL.Path = mux.Vars(r)["rest"]
		r.Header.Del("Server")
		p.ServeHTTP(w, r)
		r.Header.Del("Server")
	}
}

func (g gouter) handleCancelError(w http.ResponseWriter, r *http.Request, err error) {
	if err != context.Canceled {
		log.Printf("http: proxy error: %v", err)
	}
	w.WriteHeader(http.StatusBadGateway)
}

func (g gouter) readConfig() config {
	fmt.Println("-- GOUTER SIMPLE REVERSE PROXY --")

	bytes, err := ioutil.ReadFile("gouter.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	var c config
	if err := yaml.Unmarshal(bytes, &c); err != nil {
		log.Fatalln(err)
	}

	return c
}

func (g *gouter) run(c config) {
	g.router = mux.NewRouter()

	if len(c.Hosts) == 0 {
		log.Fatalln("ERROR: Config hosts does not exist.")
	}

	fmt.Println("Config loaded:")

	for _, item := range c.Hosts {
		parts := strings.Split(item, "::")
		fmt.Println(fmt.Sprintf("- [%s] FORWARDED TO ORIGIN [%s]", parts[1], parts[0]))
		remote, _ := url.Parse(parts[0])
		proxy := httputil.NewSingleHostReverseProxy(remote)

		proxy.ErrorHandler = g.handleCancelError

		g.router.HandleFunc(parts[1]+"{rest:.*}", g.handler(proxy, remote))
		g.router.HandleFunc(strings.TrimSuffix(parts[1], "/")+"{rest:.*}", g.handler(proxy, remote))
	}

	http.Handle("/", g.router)
	log.Println("Running with no errors...")
	log.Fatalln(http.ListenAndServe(":80", g.router))
}
