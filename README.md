# Pubsub for Caddy

[![MIT licensed][badge-mit]][license]
[![Report][badge-report]][report]

Package pubsub implements a longpoll-based publish and subscribe middleware for 
[Caddy][caddy], a modern, full-featured, easy-to-use web server.

This plugin lets you easily push event notifications to any practical number of
web clients. To publish an event, content that includes a category and body is
posted to the "publish" URL configured in the Caddyfile. To subscribe to
published events, a web client simply connects to the "subscribe" URL
configured in the Caddyfile. This connection can be managed by including the
small, dependency-free file ps.js in your web application and using the
non-blocking methods of a Subscriber instance.

This plugin uses longpolling (specifically, [golongpoll][longpoll]) to connect
clients to the server. The advantages of this are significant. Longpoll
connections

* are straightforward HTTP/HTTPS
* are not thwarted by firewalls and proxies
* are supported by virtually all browsers

Additionally, this plugin provides a simple web-based interface to publish
events, so any software capable of posting content to the Caddy server, such as
wget and web browsers, can dispatch information to listening clients. This
flexibility allows short-lived applications, suchs as crontab scripts and CGI
scripts, to publish events that are of interest to subscribing web clients.

On the downside, longpoll connections are one-direction only. Published events
flow only from the server to clients. However, because of this plugin's simple
publishing interface, a web client that receives an event can immediately
publish its own response.

## Security Considerations

As with websockets, longpolling requires special care to protect both the
server and all connected clients.

**Longpolling consumes resources on the server.** Too many connections to
clients can impact server operations. It is important to protect the
configured "subscribe" path with some form of authentication such as [basic
authentication][auth] (be sure to use HTTPS!) or [JWT][jwt] in order to
manage the number of connections that your system will maintain.

**Published events can instantly reach a large number of clients.** Be sure
to require authorization in order to access the configured "publish" path to
prevent rogue publishers from dispatching unexpected content to clients or
flooding the subscription channels.

## Basic Syntax

The basic pubsub directive lets you specify a "publish" path and a
corresponding "subscribe" path. This directive can be repeated. Each pubsub
block is managed by its own longpoll instance so categories are effectively
scoped by directive.

	pubsub publish_path subscribe_path

For example:

	cgi /chat/publish /chat/subscribe

The specified paths are virtual; they do not refer to any filesystem resources.

### Publishing

When the Caddy server receives a call that matches the publish_path URL, the
pubsub plugin responds by checking the request for the url-encoded form fields
"category" and "body". If these form values are sent to the server in a POST
request rather than included in the tail of the URL in a GET request, the
Content-Type must be "application/x-www-form-urlencoded". The body value is
then dispatched verbatim to all clients that are currently subscribed to the
specified category. Structured data is easily dispatched by sending a
JSON-encoded value. The included JavaScript file ps.js has functions that
handle this encoding and decoding automatically.

At its simplest, a publish call might look like 

	https://example.com/chat/publish?category=team&body=Hello%20world

In this example, the body "Hello world" is dispatched to all subscribers of the
"team" category. 

### Subscribing

When the Caddy server receives a call that matches the subscribe_path URL, the
pubsub plugin keeps the connection alive until a publication event of the
correct category is returned or the configured time limit is reached. In either
case, the client then makes another similar request of the server. This cycle
continues until the client page is dismissed. When the longpoll instance
detects that the client is no longer responsive it gracefully drops the client
from its subscription list.

## Advanced Syntax

The basic syntax shown above is likely all you will need to configure the
pubsub plugin. If some control over the underlying golongpoll package is
needed, you can use all or part of the advanced syntax shown here.

	pubsub publish_path subscribe_path {
		MaxLongpollTimeoutSeconds timeout
		MaxEventBufferSize count
		EventTimeToLiveSeconds timeout
		DeleteEventAfterFirstRetrieval
	}

Any missing fields are replaced with their default values; see the
[golongpoll documentation][golongpoll-doc] for more details.

The `MaxLongpollTimeoutSeconds` subdirective specifies the maximum number of
seconds that the longpoll server will keep a client connection alive.

The `MaxEventBufferSize` subdirective specifies the maximum number of events of
a particular category that will be kept by the longpoll server. Beyond this
limit, events will be dropped even if they have not expired.

The `EventTimeToLiveSeconds` subdirective specifies how long events will be
retained by the longpoll server.

If the `DeleteEventAfterFirstRetrieval` subdirective is present then events
will be deleted right after they have been dispatched to current subscribers.

[auth]: https://caddyserver.com/docs/basicauth
[badge-author]: https://img.shields.io/badge/author-Kurt_Jung-blue.svg
[badge-github]: https://img.shields.io/badge/project-Git_Hub-blue.svg
[badge-mit]: https://img.shields.io/badge/license-MIT-blue.svg
[badge-report]: https://goreportcard.com/badge/github.com/jung-kurt/caddy-pubsub
[caddy]: https://caddyserver.com/
[github]: https://github.com/jung-kurt/caddy-pubsub
[golongpoll-doc]: https://godoc.org/github.com/jcuga/golongpoll
[jung]: https://github.com/jung-kurt/
[jwt]: https://github.com/BTBurke/caddy-jwt
[key]: class:key
[license]: https://raw.githubusercontent.com/jung-kurt/caddy-pubsub/master/LICENSE
[longpoll]: https://github.com/jcuga/golongpoll
[report]: https://goreportcard.com/report/github.com/jung-kurt/caddy-pubsub
[subkey]: class:subkey