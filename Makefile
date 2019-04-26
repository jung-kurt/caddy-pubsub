all : doc example/ps.js.ok example/app.js.ok example/util.js.ok example/index.html.ok

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

%.js.ok : %.js
	jshint $<
	touch $@

%.html.ok : %.html
	tidy -quiet -output /dev/null $<
	touch $@

doc : README.md doc.go doc/index.html.ok

README.md doc.go doc/index.html : doc/doc.txt

doc/pubsub.md README.md doc.go : doc/doc.txt
	cd doc; lua doc.lua
	gofmt -s -w doc.go

doc/index.html : doc/hdr.html doc/body.html doc/ftr.html
	cat doc/hdr.html doc/body.html doc/ftr.html > $@

doc/body.html : doc/pubsub.md
	markdown -f +links,+image,+smarty,+ext,+divquote -o $@ $<

perm :
	chgrp -R caddy example
	chmod -R u=rwX,g=rX,o= example

clean :
	rm -f coverage.html coverage doc/*.ok example/*.ok doc/pubsub.md README.md doc.go doc/index.html doc/body.html
