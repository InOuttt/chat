package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// For stripping comments from JSON config
	jcr "github.com/DisposaBoy/JsonConfigReader"

	af "github.com/abaron/chat/server/adiraFinance"
	"github.com/gorilla/mux"
)

type configType struct {
	// Default HTTP(S) address:port to listen on for websocket and long polling clients. Either a
	// numeric or a canonical name, e.g. ":80" or ":https". Could include a host name, e.g.
	// "localhost:80".
	// Could be blank: if TLS is not configured, will use ":80", otherwise ":443".
	// Can be overridden from the command line, see option --listen.
	Listen string `json:"listen"`
	// Cache-Control value for static content.
	CacheControl int `json:"cache_control"`
	// Address:port to listen for gRPC clients. If blank gRPC support will not be initialized.
	// Could be overridden from the command line with --grpc_listen.
	GrpcListen string `json:"grpc_listen"`
	// Enable handling of gRPC keepalives https://github.com/grpc/grpc/blob/master/doc/keepalive.md
	// This sets server's GRPC_ARG_KEEPALIVE_TIME_MS to 60 seconds instead of the default 2 hours.
	GrpcKeepalive bool `json:"grpc_keepalive_enabled"`
	// URL path for mounting the directory with static files.
	StaticMount string `json:"static_mount"`
	// Local path to static files. All files in this path are made accessible by HTTP.
	StaticData string `json:"static_data"`
	// Salt used in signing API keys
	APIKeySalt []byte `json:"api_key_salt"`
	// Maximum message size allowed from client. Intended to prevent malicious client from sending
	// very large files inband (does not affect out of band uploads).
	MaxMessageSize int `json:"max_message_size"`
	// Maximum number of group topic subscribers.
	MaxSubscriberCount int `json:"max_subscriber_count"`
	// Masked tags: tags immutable on User (mask), mutable on Topic only within the mask.
	MaskedTagNamespaces []string `json:"masked_tags"`
	// Maximum number of indexable tags
	MaxTagCount int `json:"max_tag_count"`
	// URL path for exposing runtime stats. Disabled if the path is blank.
	ExpvarPath string `json:"expvar"`
}

var (
	config configType
)

func main() {
	af.LogInfo("##### ARYO BARON REST #####")

	executable, _ := os.Executable()

	rootpath, _ := filepath.Split(executable)

	var configfile = flag.String("config", "tinode.conf", "Path to config file.")
	flag.Parse()

	*configfile = toAbsolutePath(rootpath, *configfile)
	log.Printf("Using config from '%s'", *configfile)

	if file, err := os.Open(*configfile); err != nil {
		log.Fatal("Failed to read config file: ", err)
	} else if err = json.NewDecoder(jcr.New(file)).Decode(&config); err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	// Rest for SSO
	router := mux.NewRouter()
	router.HandleFunc("/test", createUser).Methods("GET")
	router.HandleFunc("/ping", ping).Methods("GET")
	router.HandleFunc("/v1/user", createUser).Methods("POST")
	router.HandleFunc("/v1/user/{id}", createUser).Methods("GET")
	router.HandleFunc("/v1/user/{id}", createUser).Methods("PUT")
	router.HandleFunc("/v1/user/{id}", createUser).Methods("DELETE")

	if err := http.ListenAndServe(":2097", router); err == nil {
		af.Log.Ln("Running Rest on :2097")
	} else {
		af.Log.Error(err)
	}
}

func toAbsolutePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(base, path))
}
