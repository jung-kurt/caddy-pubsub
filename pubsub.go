/*
 * Copyright (c) 2019 Kurt Jung (Gmail: kurt.w.jung)
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package pubsub

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/jcuga/golongpoll"
)

type handlerType struct {
	next              httpserver.Handler
	rules             []ruleType
	startup, shutdown func() error
}

var (
	errNoBody     = errors.New("publication body missing")
	errNoCategory = errors.New("publication category missing")
)

// ruleType represents a pubsub handling rule; it is parsed from the pubsub directive
// in the Caddyfile
type ruleType struct {
	// Publication path
	publishPath string
	// Subscription path
	subscribePath string
	// golongpoll options
	opt golongpoll.Options
	// longpoll instance for this block
	manager *golongpoll.LongpollManager
}

func init() {
	caddy.RegisterPlugin("pubsub", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

// configureServer processes the tokens collected from the Caddy configuration
// file for the "pubsub" directives and, if successful, instantiates a longpoll
// manager and inserts the pubsub handler into the middleware chain.
func configureServer(ctrl *caddy.Controller, cfg *httpserver.SiteConfig) (err error) {
	var hnd handlerType

	hnd.shutdown = func() (err error) {
		for j := 0; j < len(hnd.rules) && err == nil; j++ {
			rule := &hnd.rules[j]
			if rule.manager != nil {
				rule.manager.Shutdown()
				rule.manager = nil
			}
		}
		return
	}

	hnd.rules, err = pubsubParse(ctrl)
	if err == nil {
		for j := 0; j < len(hnd.rules) && err == nil; j++ {
			rule := &hnd.rules[j]
			rule.manager, err = golongpoll.StartLongpoll(rule.opt)
		}
		if err == nil {
			ctrl.OnShutdown(hnd.shutdown)
			cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
				hnd.next = next
				return hnd
			})
		}
	}
	return
}

// setup configures a new pubsub middleware instance with the specified filesystem
// root
func setup(c *caddy.Controller) (err error) {
	return configureServer(c, httpserver.GetConfig(c))
}

// pubsubParse parses one or more "pubsub" configuration directives
func pubsubParseAdvanced(c *caddy.Controller, opt *golongpoll.Options) (err error) {
	for err == nil && c.NextBlock() {
		val := c.Val()
		args := c.RemainingArgs()
		argCount := len(args)
		switch val {
		// If the "LoggingEnabled" flag is true, golongpoll logs messages to the
		// global log instance which interferes with caddy's logging mechanism.
		case "MaxLongpollTimeoutSeconds":
			if argCount == 1 {
				opt.MaxLongpollTimeoutSeconds, err = strconv.Atoi(args[0])
			} else {
				err = fmt.Errorf("expecting 1 argument after \"MaxLongpollTimeoutSeconds\", got %d", argCount)
			}
		case "MaxEventBufferSize":
			if argCount == 1 {
				opt.MaxEventBufferSize, err = strconv.Atoi(args[0])
			} else {
				err = fmt.Errorf("expecting 1 argument after \"MaxEventBufferSize\", got %d", argCount)
			}
		case "EventTimeToLiveSeconds":
			if argCount == 1 {
				opt.EventTimeToLiveSeconds, err = strconv.Atoi(args[0])
			} else {
				err = fmt.Errorf("expecting 1 argument after \"EventTimeToLiveSeconds\", got %d", argCount)
			}
		case "DeleteEventAfterFirstRetrieval":
			if argCount == 0 {
				opt.DeleteEventAfterFirstRetrieval = true
			} else {
				err = fmt.Errorf("unexpected arguments after \"DeleteEventAfterFirstRetrieval\"")
			}
		default:
			err = fmt.Errorf("unexpected subdirective \"%s\"", val)
		}
	}
	return
}

// pubsubParse parses one or more "pubsub" configuration directives
func pubsubParse(c *caddy.Controller) (rules []ruleType, err error) {
	for err == nil && c.Next() {
		var rule ruleType
		val := c.Val()
		args := c.RemainingArgs()
		if val == "pubsub" {
			switch {
			case len(args) == 2:
				if args[0] != args[1] {
					rule.publishPath = args[0]
					rule.subscribePath = args[1]
					err = pubsubParseAdvanced(c, &rule.opt)
				} else {
					err = fmt.Errorf("publish path and subscribe path must be different")
				}
			default:
				err = fmt.Errorf("expecting 2 arguments to \"pubsub\", got %d", len(args))
			}
		} else {
			err = fmt.Errorf("expecting \"pubsub\", got \"%s\"", val)
		}
		if err == nil {
			rules = append(rules, rule)
		}
	}
	return
}

// ServeHTTP satisfies the httpserver.Handler interface.
func (h handlerType) ServeHTTP(w http.ResponseWriter, r *http.Request) (code int, err error) {
	for _, rule := range h.rules {
		if httpserver.Path(r.URL.Path).Matches(rule.subscribePath) {
			// The following call blocks until an event is published or the call times out
			rule.manager.SubscriptionHandler(w, r)
			return
		} else if httpserver.Path(r.URL.Path).Matches(rule.publishPath) {
			err = r.ParseForm()
			if err == nil {
				category := r.Form.Get("category")
				if category != "" {
					body := r.Form.Get("body")
					if body != "" {
						err = rule.manager.Publish(category, body)
					} else {
						err = errNoBody
					}
				} else {
					err = errNoCategory
				}
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			if err == nil {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "OK")
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Not OK")
			}
			return
		}
	}
	return h.next.ServeHTTP(w, r)
}
