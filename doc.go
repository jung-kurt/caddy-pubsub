/*
Package pubsub implements a longpoll-based publish and subscribe
middleware for Caddy, a modern, full-featured, easy-to-use web server.

This plugin lets you easily push event notifications to any practical
number of web clients. To publish an event, post content that includes a
category and body to the “publish” URL configured in the Caddyfile. To
subscribe to published events, simply connect your web client to the
“subscribe” URL configured in the Caddyfile. To manage the subscription
connection, include the small, dependency-free file ps.js in your web
application and use the non-blocking methods of a Subscriber instance.

This plugin uses longpolling (specifically, golongpoll) to connect
clients to the server. The advantages of this are significant. Longpoll
connections


-   are straightforward HTTP/HTTPS

-   are not thwarted by firewalls and proxies

-   are supported by virtually all browsers

Additionally, this plugin provides a simple web-based interface to
publish events, so any software capable of posting content to the Caddy
server, such as wget and web browsers, can dispatch information to
listening clients. This flexibility allows short-lived applications,
such as crontab scripts and CGI scripts, to publish events that are of
interest to subscribing web clients.

On the downside, longpoll connections are one-direction only. Published
events flow only from the server to clients. However, because of this
plugin’s simple publishing interface, a web client that receives an
event can immediately publish its own response.


Security Considerations

As with websockets, longpolling requires special care to protect both
the server and all connected clients.

LONGPOLLING CONSUMES RESOURCES ON THE SERVER. Too many connections to
clients can impact server operations. It is important to protect the
configured “subscribe” path with some form of authentication such as
basic authentication or JWT in order to manage the number of connections
that your system will maintain. Be sure to use HTTPS!

PUBLISHED EVENTS CAN INSTANTLY REACH A LARGE NUMBER OF CLIENTS. Be sure
to require authorization in order to access the configured “publish”
path to prevent rogue publishers from dispatching unexpected content to
clients or flooding the subscription channels.


Basic Syntax

The basic pubsub directive lets you specify a “publish” path and a
corresponding “subscribe” path. This directive can be repeated. Each
pubsub block is managed by its own longpoll instance so categories are
effectively scoped by directive.

    pubsub publish_path subscribe_path

For example:

    pubsub /chat/publish /chat/subscribe

The specified paths are virtual; they do not refer to any filesystem
resources.

Publishing

When the Caddy server receives a call that matches the publish_path URL,
the pubsub plugin responds by checking the request for the url-encoded
form fields “category” and “body”. If these form values are sent to the
server in a POST request rather than included in the tail of the URL in
a GET request, the Content-Type must be
“application/x-www-form-urlencoded”. The body value is then dispatched
verbatim to all clients that are currently subscribed to the specified
category. Structured data is easily dispatched by sending a JSON-encoded
value.

At its simplest, a publish call might look like

    https://example.com/chat/publish?category=team&body=Hello%20world

In this example, the body “Hello world” is dispatched to all subscribers
of the “team” category.

Subscribing

When the Caddy server receives a call that matches the subscribe_path
URL, the pubsub plugin keeps the connection alive until a publication
event of the correct category is returned or the configured time limit
is reached. In either case, the client then makes another similar
request of the server. This cycle continues until the client page is
dismissed. When the longpoll instance detects that the client is no
longer responsive it gracefully drops the client from its subscription
list.


Advanced Syntax

The basic syntax shown above is likely all you will need to configure
the pubsub plugin. If some control over the underlying golongpoll
package is needed, you can use all or part of the advanced syntax shown
here.

    pubsub publish_path subscribe_path {
        MaxLongpollTimeoutSeconds timeout
        MaxEventBufferSize count
        EventTimeToLiveSeconds timeout
        DeleteEventAfterFirstRetrieval
    }

Any missing fields are replaced with their default values; see the
golongpoll documentation for more details.

The MaxLongpollTimeoutSeconds subdirective specifies the maximum number
of seconds that the longpoll server will keep a client connection alive.

The MaxEventBufferSize subdirective specifies the maximum number of
events of a particular category that will be kept by the longpoll
server. Beyond this limit, events will be dropped even if they have not
expired.

The EventTimeToLiveSeconds subdirective specifies how long events will
be retained by the longpoll server.

If the DeleteEventAfterFirstRetrieval subdirective is present then
events will be deleted right after they have been dispatched to current
subscribers.


Running the example

Here is a sample Caddyfile that can be modified for use in the following
example:

    http://127.0.0.1 {
        bind 127.0.0.1
        basicauth /psdemo/publish publish 123
        basicauth /psdemo/subscribe subscribe 123
        pubsub /psdemo/publish /psdemo/subscribe
    }


-   Obtain or build a Caddy web server with the pubsub plugin.


-   Create an example directory and download the files in the repository
example directory to it. Make this your default working directory.


-   Edit the file named Caddyfile. Modify the site address to an
appropriate local network address. Ideally, you will use an address
that can be accessed from a number of devices on your network. (A
local interface, such as 127.0.0.5, will work, but then you will
have to simulate multiple devices by opening multiple tabs in your
browser.) You may wish to use a non-standard port if you need to
avoid interfering with another server. HTTP is used here for the
purposes of a local, insecure demonstration. USE HTTPS IN
PRODUCTION. The Caddyfile has values for the authorization fields
and paths that are used as defaults in the example script, so
leaving them as they are will simplify the demonstration.


-   Launch the Caddy web server from the directory that contains your
example Caddyfile. The command might be simply caddy or, if the
caddy executable you wish to run is not on the default search path,
you will need to qualify the name with its location.


-   Access the server with the web browsers of a number of devices.
Alternatively, open a number of tabs in a single browser and point
them to the server.


-   On each open page, click the “Configure” button and make appropriate
changes. Most fields, such as “Auth password” and “URL” are
pre-filled to match the values in the sample Caddyfile. You may
enter a name in the “Publisher name” field or leave it empty. A
blank value will be replaced with a random name like “user_42” the
first time you publish an event.


-   On each open page, click the “Run” button and then the “Start”
button. Simulate the publication of events by clicking the “A”, “B”,
and “C” buttons from various devices or tab pages. Additionally, you
may click the “Auto start” button. This will begin an automatic
sequence in which an event is published at random intervals.
Published events will be sent to the web server and dispatched to
all subscribing pages. These events will be displayed beneath the
buttons in a list. A page can publish events even if it does not
subscribe to them.


Using pubsub in your own web applications

The JavaScript file ps.js is included in the example shown above. This
script may be included in a web page with the following line:

    <script src="ps.js"></script>

The script is small and dependency-free and should be easy to modify if
needed. Within your application code, instantiate a Subscriber instance
as follows:

    subscribe = ps.Subscriber(category, url, callback, authorization, options);

The parameters are:


-   category: a short string that identifies the event category to which
to subscribe


-   url: the subscribe_path configured in the Caddyfile (in the example
above, this is “/psdemo/subscribe”)


-   callback: this is a function that is called (with the published body
and server timestamp) for each event of the specified category


-   authorization: a string like “Basic c3Vic2NyaWJlOjEyMw==” that will
be sent as an authorization header.


-   options: an object that contains the fields timeout (seconds,
default 45), successDelay (milliseconds, default 10), errorDelay
(milliseconds, default 3000), and json (boolean, true if event
bodies are JSON-encoded and should be automatically decoded, default
false).

More details can be found in the comments in the ps.js file.

To start the subscription, call

    subscribe.start();

To end the subscription, call

    subscribe.stop();

Publish an event as follows:

    ps.publish(category, url, body, authorization);

The parameters are:


-   category: a short string that identifies the event category of the
published event


-   url: the publish_path configured in the Caddyfile (in the example
above, this is “/psdemo/publish”)


-   body: this is the text that will be dispatched to all subscribers of
events with the specified category; this text is often a
JSON-encoded object


-   authorization: a string like “Basic cHVibGlzaDoxMjM=” that will be
sent as an authorization header.

A variation of the ps.publish() function is

    ps.publishObj(category, url, body, authorization);

In this case, body can be any JSON-encodable value such as an array or
an object. When subscribing to events that are published this way, it is
convenient to set the “json” field of the options argument of
ps.Subscriber() to true so that the event body is automatically decoded.
*/
package pubsub
