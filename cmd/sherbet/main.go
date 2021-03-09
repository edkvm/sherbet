package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	"net/url"

	"flag"

	"github.com/edkvm/sherbet"

	"github.com/edkvm/sherbet/entryparser"
)

type logWriter struct{}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s [DEBUG] %s", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"), string(bytes))
}

type Env struct {
	enbaleSecurity bool
	key            string
	salt           string
}

const Name = "sherbet"
const Vesrion = "0.1"

func main() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))

	env := initEnv()

	var securityMiddleware func(next http.Handler) http.Handler

	if env.enbaleSecurity {
		securityMiddleware = securityMiddlewareBuilder(env.key, env.salt)
	}

	middlewareChain := []Middleware{serverDetailsMiddleware, securityMiddleware, cacheHeadersMiddleware}

	srvs := []*http.Server{
		sherbet.StartServer(6060, applyMiddleware(middlewareChain, entryparser.NewOpEngineMiddelware())),
		sherbet.StartServer(6061, applyMiddleware(middlewareChain, entryparser.NewHTTPRubyParserMiddelware())),
	}

	log.Println("server started")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for sig := range c {

		if sig == os.Interrupt {

		}
		log.Printf("%s %s server is shutting down\n", Name, Vesrion)
		ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		for _, srv := range srvs {
			if srv != nil {
				srv.Shutdown(ctx)
			}
		}

		return
	}
}

func initEnv() *Env {
	enableSecurityPtr := flag.Bool("security", false, "enable security")
	keyPtr := flag.String("key", "", "This key will be used to verify the origin of the request")
	saltPtr := flag.String("salt", "", "Salt will be added to the request for security reasons")

	flag.Parse()

	// TODO(ekiselman): Verify key and salt are not nil if security enabled
	return &Env{
		enbaleSecurity: *enableSecurityPtr,
		key:            *keyPtr,
		salt:           *saltPtr,
	}
}

func applyMiddleware(chain Chain, h http.Handler) http.Handler {
	for i := range chain {
		idx := (len(chain) - 1) - i
		m := chain[idx]
		if m != nil {
			h = m(h)
		}
	}

	return h
}

type Middleware func(http.Handler) http.Handler

type Chain []Middleware

func securityMiddlewareBuilder(key, salt string) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		rgxp := regexp.MustCompile(`^/(?P<secret>[\d\w]{10,})/`)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			validSign := false
			match := rgxp.FindAllStringSubmatch(r.URL.Path, -1)
			if match != nil && len(match) == 1 {
				mod := strings.TrimPrefix(r.URL.Path, match[0][0])
				hash := match[0][1]
				if sherbet.VerfiySignature(mod, hash, key, salt) {
					validSign = true
					modUrl, _ := url.Parse(mod)
					r.URL.Path = modUrl.Path
					r.URL.RawPath = modUrl.RawPath
				}
			}

			//TODO(ekiselman): Switch the ifs
			if !validSign {
				http.Error(w, "", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: Move to conf file
func cacheHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentTime := time.Now().UTC()
		w.Header().Set("Expires", currentTime.Add(24*30*time.Hour).Format(time.RFC1123))
		w.Header().Set("Date", currentTime.Format(time.RFC1123))
		maxAge := time.Hour * 24 * 365
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		next.ServeHTTP(w, r)
	})
}

func serverDetailsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", fmt.Sprintf("%s/%s", Name, Vesrion))
		next.ServeHTTP(w, r)
	})
}
