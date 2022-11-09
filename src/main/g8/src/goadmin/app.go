package goadmin

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	hocon "github.com/go-akka/configuration"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"main/src/cocostore"
	"main/src/utils"
)

const (
	defaultConfigFile = "./config/application.conf"
)

var (
	AppConfig *hocon.Config

	EchoServer           *echo.Echo
	echoServerListenAddr string
	echoServerListenPort int32
)

// Start bootstraps the application.
func Start(bootstrappers ...IBootstrapper) {
	var err error

	AppConfig = initAppConfig()
	utils.DevMode = AppConfig.GetBoolean("dev_mode", utils.DevMode)
	if utils.Location, err = time.LoadLocation(AppConfig.GetString("timezone")); err != nil {
		panic(err)
	}

	EchoServer, echoServerListenAddr, echoServerListenPort = initEchoServer()

	// map static resources
	if confV := AppConfig.GetValue("static_resources"); confV != nil && confV.IsObject() {
		for uri, dirO := range confV.GetObject().Items() {
			if dirO.IsString() {
				dir := dirO.GetString()
				log.Printf("Mapping static resources: %s -> %s", uri, dir)
				if !strings.HasPrefix(uri, "/") {
					uri = "/" + uri
				}
				EchoServer.Static(uri, dir)
			}
		}
	}

	// bootstrapping
	if bootstrappers != nil {
		for _, b := range bootstrappers {
			log.Println("Bootstrapping", b)
			if err := b.Bootstrap(AppConfig, EchoServer); err != nil {
				log.Println(err)
			}
		}
	}

	startEchoServer(EchoServer, echoServerListenAddr, echoServerListenPort)
}

func initAppConfig() *hocon.Config {
	configFile := os.Getenv("APP_CONFIG")
	if configFile == "" {
		log.Printf("No environment APP_CONFIG found, fallback to [%s]", defaultConfigFile)
		configFile = defaultConfigFile
	}
	return loadAppConfig(configFile)
}

func initEchoServer() (*echo.Echo, string, int32) {
	listenPort := AppConfig.GetInt32("http.listen_port", 0)
	if listenPort <= 0 {
		panic("No valid [http.listen_port] configured")
	}
	listenAddr := AppConfig.GetString("http.listen_addr", "127.0.0.1")

	e := echo.New()

	// register session middleware
	sessionKey := AppConfig.GetString("goadmin.session_key", "s3cr3t_s3ssion_2uth3ntic2tion_k3y")
	// e.Use(session.Middleware(sessions.NewCookieStore([]byte(sessionKey))))
	e.Use(session.Middleware(cocostore.NewCompressedCookieStore(cocostore.CompressionLevelBestCompression, []byte(sessionKey))))

	requestTimeout := AppConfig.GetTimeDuration("http.request_timeout", time.Duration(0))
	if requestTimeout > 0 {
		e.Server.ReadTimeout = requestTimeout
	}
	bodyLimit := AppConfig.GetByteSize("http.max_request_size")
	if bodyLimit != nil && bodyLimit.Int64() > 0 {
		e.Use(middleware.BodyLimit(bodyLimit.String()))
	}

	TemplateRenderer = newGoadminRenderer()
	e.Renderer = TemplateRenderer

	return e, listenAddr, listenPort
}

func startEchoServer(echoServer *echo.Echo, listenAddr string, listenPort int32) {
	log.Printf("Starting [%s] on [%s:%d]...\n", AppConfig.GetString("app.name")+" v"+AppConfig.GetString("app.version"), listenAddr, listenPort)
	go echoServer.Logger.Fatal(echoServer.Start(fmt.Sprintf("%s:%d", listenAddr, listenPort)))
}
