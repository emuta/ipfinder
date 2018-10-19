package main

import (
    "flag"
    "log"
    "net"
    "net/http"
    "encoding/json"

    "github.com/ipipdotnet/datx-go"
)

var addr   = flag.String("addr", ":80", "service listen address")
var dbfile = flag.String("dbfile", "./17monipdb.datx", "db file of 17monipdb.datx path")

var city *datx.City

func init() {
    flag.Parse()

    var err error
    if city, err = datx.NewCity(*dbfile); err != nil {
        log.Fatal(err)
    }
}

type Location struct {
    Country  string `json:"country"`
    Province string `json:"province"`
    City     string `json:"city"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query()
    ip := query.Get("ip")

    if ip == "" {
        ip, _, _ = net.SplitHostPort(r.RemoteAddr)
    }

    // loc, err := city.Find(ip)
    loc, err := city.FindLocation(ip)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte(err.Error()))
        return
    }
    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    l := &Location{loc.Country, loc.Province, loc.City}
    json.NewEncoder(w).Encode(l)

    go func() {
        log.Printf("[%s] %s, %#v", r.RemoteAddr, r.URL, l)
        }()
}

func main() {
    http.HandleFunc("/", Handler)
    log.Fatal(http.ListenAndServe(*addr, nil))
}