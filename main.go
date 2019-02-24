package main

import (
    "flag"
    "net"
    "net/http"
    "encoding/json"
    
    "github.com/sirupsen/logrus"
    "github.com/ipipdotnet/datx-go"
)

var addr   = flag.String("addr", ":80", "service listen address")
var dbfile = flag.String("dbfile", "./17monipdb.datx", "db file of 17monipdb.datx path")

var city *datx.City
var log *logrus.Logger

func init() {
    flag.Parse()

    // initailize logger
    log = logrus.New()
    log.SetFormatter(&logrus.TextFormatter{
        FullTimestamp:   true,
        TimestampFormat: "2006-01-02 15:04:05",
    })

    // load dbfile
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

    loc, err := city.FindLocation(ip)
    if err != nil {
        w.WriteHeader(404)
        w.Write([]byte(err.Error()))
        log.WithError(err).WithField("ip", ip).Error(err)
        return
    }

    w.Header().Set("Content-Type", "application/json;charset=UTF-8")
    l := &Location{loc.Country, loc.Province, loc.City}
    json.NewEncoder(w).Encode(l)

    log.WithFields(logrus.Fields{
        "remote": r.RemoteAddr, 
        "country": l.Country,
        "province": l.Province,
        "city": l.City}).Info(r.URL)
}

func main() {
    http.HandleFunc("/", Handler)
    log.Fatal(http.ListenAndServe(*addr, nil))
}