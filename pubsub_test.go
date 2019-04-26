package pubsub

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func configure(expectErr bool, str string) (rules []ruleType, err error) {
	var c *caddy.Controller
	var mids []httpserver.Middleware
	var cfg *httpserver.SiteConfig
	var srvHnd httpserver.Handler
	var hnd handlerType
	var ok bool

	// printf("--- Directive begin --\n%s\n--- Directive end ---\n", str)
	c = caddy.NewTestController("http", str)
	// Caddy controller does not export its instance field, so it appears to be
	// impossible to test OnStartup and OnShutdown functions
	err = setup(c)
	if err == nil {
		if !expectErr {
			cfg = httpserver.GetConfig(c)
			mids = cfg.Middleware()
			midLen := len(mids)
			if midLen > 0 {
				for j := 0; j < midLen && err == nil; j++ {
					srvHnd = mids[j](httpserver.EmptyNext)
					hnd, ok = srvHnd.(handlerType)
					if ok {
						rules = append(rules, hnd.rules...)
					} else {
						err = fmt.Errorf("expected middleware handler to be pubsub handler")
					}
				}
			} else {
				err = fmt.Errorf("no middlewares present")
			}
		} else {
			err = fmt.Errorf("expected error but got succcess")
		}
	} else if expectErr {
		err = nil
	}
	return
}

func TestSetup(t *testing.T) {
	var err error
	// Each of the following directives is submitted for parsing. If the string is
	// prefixed with "0:", it is expected to parse successfully. If it is prefixed
	// with "1:", an error is expected. The prefix is removed before parsing.
	var directiveList = []string{
		`0:pubsub /publish /subscribe`,
		`1:pubsubx /publish /subscribe`,
		`1:pubsub /same /same`,
		`1:pubsub`,
		`1:pubsub /just_one`,
		`1:pubsub /publish /subscribe too_many`,
		`0:pubsub /publish /subscribe {
	MaxLongpollTimeoutSeconds 60
	MaxEventBufferSize 100
	EventTimeToLiveSeconds 300
	DeleteEventAfterFirstRetrieval
}`,
		`1:pubsub /publish /subscribe {
	MaxLongpollTimeoutSeconds 12 foo
}`,
		`1:pubsub /publish /subscribe {
	MaxEventBufferSize
}`,
		`1:pubsub /publish /subscribe {
	EventTimeToLiveSeconds 300 100
}`,
		`1:pubsub /publish /subscribe {
	EventTimeToLiveSeconds july
}`,
		`1:pubsub /publish /subscribe {
	DeleteEventAfterFirstRetrieval foo
}`,
		`1:pubsub /publish /subscribe {
	foo bar
}`,
	}

	for j := 0; j < len(directiveList) && err == nil; j++ {
		str := directiveList[j]
		_, err = configure(str[0:1] != "0", str[2:])
		if err != nil {
			fmt.Printf("error [%s], str [%s]\n", err.Error(), str)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

// handleGet returns a cgi handler (which implements ServeHTTP) based on the
// specified directive string. The directive can contain more than one block,
// but only the first is associated with the returned handler.
func handlerGet(directiveStr, rootStr string) (hnd handlerType, err error) {
	var c *caddy.Controller
	var mids []httpserver.Middleware
	var cfg *httpserver.SiteConfig
	var ok bool
	var midLen int

	c = caddy.NewTestController("http", directiveStr)
	cfg = httpserver.GetConfig(c)
	cfg.Root = rootStr
	err = configureServer(c, cfg)
	if err == nil {
		mids = cfg.Middleware()
		midLen = len(mids)
		if midLen > 0 {
			hnd, ok = mids[0](httpserver.EmptyNext).(handlerType)
			if ok {
				// .. success
			} else {
				err = fmt.Errorf("expected middleware handler to be CGI handler")
			}
		} else {
			err = fmt.Errorf("no middlewares present")
		}
	}
	return
}

func TestServe(t *testing.T) {
	var err error
	var hnd handlerType
	var buf strings.Builder
	var srv *httptest.Server
	directiveList := []string{
		`pubsub /publish /subscribe`,
	}

	requestList := []string{
		"/publish?category=demo&body=foo",
		"/publish?body=foo",
		"/publish?category=demo",
		"/not_a_pubsub",
		"/subscribe?timeout=1&category=test",
	}

	expectStr := `++-+-+++++`

	setErrorFlag := func(err error) {
		var c byte
		if err == nil {
			c = '+'
		} else {
			c = '-'
		}
		buf.WriteByte(c)
	}

	for dirJ := 0; dirJ < len(directiveList) && err == nil; dirJ++ {
		hnd, err = handlerGet(directiveList[dirJ], "./test")
		if err == nil {
			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, err = hnd.ServeHTTP(w, r)
				setErrorFlag(err)
			}))
			// In production, Caddy calls a middleware's startup and shutdown functions
			// appropriately, but not in a test. Call them here manually to remedy that.
			for reqJ := 0; reqJ < len(requestList) && err == nil; reqJ++ {
				var res *http.Response
				res, err = http.Get(srv.URL + requestList[reqJ])
				setErrorFlag(err)
				if err == nil {
					res.Body.Close()
				}
			}
			hnd.shutdown()
			srv.Close()
		}
	}

	if err == nil {
		gotStr := buf.String()
		if expectStr != gotStr {
			err = fmt.Errorf("expected %s, got %s", expectStr, gotStr)
		}
	}

	if err != nil {
		t.Fatalf("%s", err)
	}

}
