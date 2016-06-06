package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func readJsonFromFile(v interface{}, name string) (err error) {
	bs, err := ioutil.ReadFile(name)
	if err != nil {
		return
	}
	err = json.Unmarshal(bs, v)
	return
}
func writeJsonToFile(v interface{}, name string) (err error) {
	bs, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return
	}

	err = ioutil.WriteFile(name, bs, 0x666)
	return
}

const fname = "./conf.json"

type Conf struct {
	CertFilePath string
	KeyFilePath  string
	Router       []ProxyRouter
}
type ProxyRouter struct {
	Src string
	Dst string
}

var conf Conf

func loadData() {
	err := readJsonFromFile(&conf, fname)
	if err != nil {
		writeJsonToFile(conf, fname)
		panic("not proxy router")
	}
	fmt.Printf("proxyRouter:%#v", conf)
}

func main() {
	loadData()

	var hostMap = map[string]http.Handler{}
	for _, router := range conf.Router {
		u, err := url.Parse(router.Dst)
		if err != nil {
			panic(err)
		}
		hostMap[router.Src] = httputil.NewSingleHostReverseProxy(u)
	}
	fh := func(rw http.ResponseWriter, req *http.Request) {

		if host, ok := hostMap[req.Host]; ok {
			host.ServeHTTP(rw, req)
		} else {
			fmt.Printf("url:%#v\n", req.Host)
			http.Error(rw, "no proxy", 500)
		}
	}
	if len(conf.CertFilePath) > 0 && len(conf.KeyFilePath) > 0 {
		go func() {
			err := http.ListenAndServeTLS(":443", conf.CertFilePath, conf.KeyFilePath, http.HandlerFunc(fh))
			panic(err)
		}()
	}
	http.ListenAndServe(":80", http.HandlerFunc(fh))
}
