package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	coordTransform "github.com/qichengzx/coordtransform"
)

var traccarserveraddr string

func main() {
	flag.StringVar(&traccarserveraddr, "u", "", "traccar server addr")
	flag.Parse()
	if traccarserveraddr == "" {
		log.Fatalln("must have traccar server addr!")
	}
	router := httprouter.New()
	router.HandlerFunc("POST", "/", Index)
	transserveraddr := "0.0.0.0:8006"

	log.Printf("start server at %v ...\n", transserveraddr)

	log.Println(http.ListenAndServe(transserveraddr, router))

}

func Index(w http.ResponseWriter, r *http.Request) {

	v := r.URL.Query()
	slon := v.Get("lon")
	slat := v.Get("lat")

	lon, _ := strconv.ParseFloat(slon, 64)
	lat, _ := strconv.ParseFloat(slat, 64)

	clon, clat := coordTransform.WGS84toGCJ02(lon, lat)

	// /?id=295182&timestamp=1650858135&lat=31.096157&lon=121.232793&speed=0&bearing=0&altitude=6.753916263580322&accuracy=35&batt=99
	data := make(url.Values)
	data["id"] = []string{v.Get("id")}
	data["timestamp"] = []string{v.Get("timestamp")}
	data["speed"] = []string{v.Get("speed")}
	data["bearing"] = []string{v.Get("bearing")}
	data["altitude"] = []string{v.Get("altitude")}
	data["accuracy"] = []string{v.Get("accuracy")}
	data["batt"] = []string{v.Get("batt")}
	data["lat"] = []string{fmt.Sprintf("%v", clat)}
	data["lon"] = []string{fmt.Sprintf("%v", clon)}

	proxyReq, _ := http.NewRequest("POST", traccarserveraddr, strings.NewReader(data.Encode())) // URL-encoded payload

	proxyReq.Header.Set("Host", r.Host)
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	for header, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(header, value)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

}
