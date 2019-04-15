CADDYDIR=/home/kurt/go/src/github.com/mholt/caddy/caddy
CADDY=${CADDYDIR}/caddy

${CADDY} : pubsub.go example/ps.js.ok example/app.js.ok example/util.js.ok example/index.html.ok
	./build

cov : ${CADDY}
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

run : ${CADDY}
	${CADDY} -conf example/Caddyfile

%.js.ok : %.js
	jshint $<
	touch $@

%.html.ok : %.html
	tidy -quiet -output /dev/null $<
	touch $@

publish :
	wget --output-document=- --quiet "http://192.168.1.20/demo/publish?category=demo&body=from%20makefile"

doc : README.md doc.go doc/index.html

README.md doc.go doc/index.html : doc/doc.txt

doc/pubsub.md README.md doc.go : doc/doc.txt
	cd doc; lua doc.lua
	gofmt -s -w doc.go

doc/index.html : doc/hdr.html doc/body.html doc/ftr.html
	cat doc/hdr.html doc/body.html doc/ftr.html > $@
	tidy -q -o /dev/null $@

doc/body.html : doc/pubsub.md
	markdown -f +links,+image,+smarty,+ext,+divquote -o $@ $<

clean :
	rm -f coverage.html coverage example/*.ok doc/pubsub.md README.md doc.go doc/index.html doc/body.html
