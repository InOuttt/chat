package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// For stripping comments from JSON config
	jcr "github.com/DisposaBoy/JsonConfigReader"
	af "github.com/abaron/chat/server/adiraFinance"
	"github.com/abaron/chat/server/auth"
	"github.com/hako/branca"
)

type credValidator struct {
	// AuthLevel(s) which require this validator.
	requiredAuthLvl []auth.Level
	addToTags       bool
}

var globals struct {
	// Indicator that shutdown is in progress
	shuttingDown bool

	// Credential validators.
	validators map[string]credValidator
	// Validators required for each auth level.
	authValidators map[auth.Level][]string

	// Salt used for signing API key.
	apiKeySalt []byte
	// Tag namespaces (prefixes) which are immutable to the client.
	immutableTagNS map[string]bool
	// Tag namespaces which are immutable on User and partially mutable on Topic:
	// user can only mutate tags he owns.
	maskedTagNS map[string]bool

	// Add Strict-Transport-Security to headers, the value signifies age.
	// Empty string "" turns it off
	tlsStrictMaxAge string
	// Listen for connections on this address:port and redirect them to HTTPS port.
	tlsRedirectHTTP string
	// Maximum message size allowed from peer.
	maxMessageSize int64
	// Maximum number of group topic subscribers.
	maxSubscriberCount int
	// Maximum number of indexable tags.
	maxTagCount int

	// Maximum allowed upload size.
	maxFileUploadSize int64

	// Ldap server base_url
	ldapServer map[string]string
}

type validatorConfig struct {
	// TRUE or FALSE to set
	AddToTags bool `json:"add_to_tags"`
	//  Authentication level which triggers this validator: "auth", "anon"... or ""
	Required []string `json:"required"`
	// Validator params passed to validator unchanged.
	Config json.RawMessage `json:"config"`
}

type mediaConfig struct {
	// The name of the handler to use for file uploads.
	UseHandler string `json:"use_handler"`
	// Maximum allowed size of an uploaded file
	MaxFileUploadSize int64 `json:"max_size"`
	// Garbage collection timeout
	GcPeriod int `json:"gc_period"`
	// Number of entries to delete in one pass
	GcBlockSize int `json:"gc_block_size"`
	// Individual handler config params to pass to handlers unchanged.
	Handlers map[string]json.RawMessage `json:"handlers"`
}

// Contentx of the configuration file
type configType struct {
	// HTTP(S) address:port to listen on for websocket and long polling clients. Either a
	// numeric or a canonical name, e.g. ":80" or ":https". Could include a host name, e.g.
	// "localhost:80".
	// Could be blank: if TLS is not configured, will use ":80", otherwise ":443".
	// Can be overridden from the command line, see option --listen.
	Listen string `json:"listen"`
	// Base URL path where the streaming and large file API calls are served, default is '/'.
	// Can be overriden from the command line, see option --api_path.
	ApiPath string `json:"api_path"`
	// Cache-Control value for static content.
	CacheControl int `json:"cache_control"`
	// Address:port to listen for gRPC clients. If blank gRPC support will not be initialized.
	// Could be overridden from the command line with --grpc_listen.
	GrpcListen string `json:"grpc_listen"`
	// Enable handling of gRPC keepalives https://github.com/grpc/grpc/blob/master/doc/keepalive.md
	// This sets server's GRPC_ARG_KEEPALIVE_TIME_MS to 60 seconds instead of the default 2 hours.
	GrpcKeepalive bool `json:"grpc_keepalive_enabled"`
	// URL path for mounting the directory with static files (usually TinodeWeb).
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
	// Ldap exchange token server
	LdapServer map[string]string `json:"ldap_server"`

	// Configs for subsystems
	Cluster json.RawMessage `json:"cluster_config"`
	Plugin  json.RawMessage `json:"plugins"`
	// Store       json.RawMessage             `json:"store_config"`
	Push        json.RawMessage             `json:"push"`
	TLS         json.RawMessage             `json:"tls"`
	Auth        map[string]json.RawMessage  `json:"auth_config"`
	Validator   map[string]*validatorConfig `json:"acc_validation"`
	Media       *mediaConfig                `json:"media"`
	StoreConfig struct {
		Adapters struct {
			MSSQL struct {
				Salt string `json:"salt"`
			} `json:"mssql"`
		} `json:"adapters"`
	} `json:"store_config"`
}

var (
	config configType
	b      *branca.Branca
)

func main() {
	fmt.Println("##### ENCRYPTION TOOLS #####")

	separator := string(os.PathSeparator)
	rootpath := os.Getenv("GOPATH") + separator + "src" + separator + "github.com" + separator + "abaron" + separator + "chat" + separator + "server"

	var configfile = flag.String("config", "tinode.conf", "Path to config file.")
	flag.Parse()

	*configfile = toAbsolutePath(rootpath, *configfile)

	if file, err := os.Open(*configfile); err != nil {
		log.Fatal("Failed to read config file: ", err)
	} else if err = json.NewDecoder(jcr.New(file)).Decode(&config); err != nil {
		log.Fatal("Failed to parse config file: ", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	if len(config.StoreConfig.Adapters.MSSQL.Salt) >= 32 {
		b = branca.NewBranca(config.StoreConfig.Adapters.MSSQL.Salt[:32]) // This key must be exactly 32 bytes long.
	} else {
		b = branca.NewBranca("@ry0b4RoNpRoJ3ct4d!r@SmG4BaRoKAH") // This key must be exactly 32 bytes long.
	}

	fmt.Println("Choose option 1, 2, or 3")
	fmt.Println("1. Free Hash")
	fmt.Println("2. Formatted")
	fmt.Println("3. Exit")
	fmt.Print("Your choice: ")

	scanner.Scan()
	if scanner.Text() == "1" {
		fmt.Println("\nEnter your string what you wan't to hash like hostname, username, or password")
		fmt.Println("Or type \"q\" for quit")
		var text string
		for {
			fmt.Print("Enter string : ")
			scanner.Scan()
			text = scanner.Text()

			if text == "q" || text == "Q" || text == "quit" { // break the loop
				break
			}

			// Encode String to Branca Token.
			token, err := b.EncodeToString(text)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Your hash was:", token)
			}
		}
	} else if scanner.Text() == "2" {
		var (
			err       error
			hostname  string
			username  string
			password  string
			port      string
			database  string
			dsnFormat = "\"dsn\": \"server=%s;user id=%s;password=%s;port=%s;database=%s\","
			dbFormat  = "\"database\": \"%s\","
		)

		fmt.Print("\nEnter hostname: ")
		scanner.Scan()
		hostname, err = b.EncodeToString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("Enter username: ")
		scanner.Scan()
		username, err = b.EncodeToString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("Enter password: ")
		scanner.Scan()
		password, err = b.EncodeToString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("Enter port: ")
		scanner.Scan()
		port, err = b.EncodeToString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Print("Enter database: ")
		scanner.Scan()
		database, err = b.EncodeToString(scanner.Text())
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("\nYour DSN:")
		fmt.Printf(dsnFormat, hostname, username, password, port, database)
		fmt.Println("\n\nYour Database:")
		fmt.Printf(dbFormat, database)
		fmt.Println("\n")
		log.Println(af.GetStringInBetween(dsnFormat, "server=", ";user id"))
	} else {
		fmt.Print("Exit. Thank you")
	}

}

func toAbsolutePath(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Clean(filepath.Join(base, path))
}
