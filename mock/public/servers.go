package public

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

type Server interface {
	SetID(id string)
	GetRequestRate() int
	GetLogMessage(ts int64) map[string]interface{}
	GetSystemMessage(ts int64) map[string]interface{}
}

type BaseServer struct {
	ID string
	SpeedFactor int
	RequestRate int // Request rate in thousands per minute. 2 for non-web servers, 5 for web
}

func (b BaseServer) SetID(id string) {
	b.ID = id
}

func (b BaseServer) GetRequestRate() int {
	return b.RequestRate * 1000
}

func (b BaseServer) GetSystemMessage(ts int64) map[string]interface{} {
	return map[string]interface{}{
		"timestamp":       ts,
		"type":            "system",
		"server_id":       b.ID,
		"cpu_usage":       math.Round(((rand.Float64() * 30) + 40) * 10) / 10, // CPU usages are between 40 and 70 with one decimal place
		"memory_usage":    math.Round((rand.Float64() * 20) + 10) * 100, // Memory usage in MB, between 1000 and 3000
		"memory_capacity": math.Round((rand.Float64() * 5) + 3) * 1000,    // Memory capacity in MB, between 3000 and 8000
		"disk_space":      math.Round(rand.Float64() * 100) * 100000, // Available disk space in bytes
	}
}

var urlPaths = [7]string{"/", "/contact", "/users", "/bi-dashboard", "/login", "/upload", "/search",}

func GetServerCluster() []Server {
	rand.Seed(time.Now().Unix())

	var servers []Server

	var webNames = [15]string{
		"Natty Narwhal",
		"Oneiric Ocelot",
		"Precise Pangolin",
		"Quantal Quetzal",
		"Raring Ringtail",
		"Saucy Salamander",
		"Trusty Tahr",
		"Utopic Unicorn",
		"Vivid Vervet",
		"Wily Werewolf",
		"Xenial Xerus",
		"Yakkety Yak",
		"Zesty Zapus",
		"Artful Aardvark",
		"Bionic Beaver",
	}

	// Add 15 web servers
	for i := 32; i < 47; i++ {
		servers = append(servers, WebServer{BaseServer{"app-0" + strconv.Itoa(i) + ": " + webNames[i - 32], (rand.Intn(40) + 1) * 10, 4}})
	}

	// 3 DB servers
	for i := 0; i < 3; i++ {
		servers = append(servers, DBServer{BaseServer{"db-00" + strconv.Itoa(i), (rand.Intn(40) + 1) * 10, 2}})
	}

	// 2 search servers
	for i := 0; i < 2; i++ {
		servers = append(servers, SearchServer{BaseServer{"search-00" + strconv.Itoa(i), (rand.Intn(40) + 1) * 10, 2}})
	}

	// 2 reverse proxies
	for i := 0; i < 2; i++ {
		servers = append(servers, ReverseProxy{BaseServer{"proxy-00" + strconv.Itoa(i), (rand.Intn(40) + 1) * 10, 5}})
	}

	return servers
}