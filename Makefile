all : documentation lint

documentation : doc/index.html doc.go README.md 

lint : example/ps.js.ok example/app.js.ok example/util.js.ok example/index.html.ok

%.js.ok : %.js
	jshint $<
	touch $@

build :
	cd ../caddy-custom
	go build -v
	sudo setcap cap_net_bind_service=+ep ./caddy
	./caddy -plugins | grep pubsub
	./caddy -version

%.html.ok : %.html
	tidy -quiet -output /dev/null $<
	touch $@

cov : all
	go test -v -coverprofile=coverage && go tool cover -html=coverage -o=coverage.html

check :
	golint .
	go vet -all .
	gofmt -s -l .

README.md : doc/document.md
	pandoc --read=markdown --write=gfm < $< > $@

doc/index.html : doc/document.md doc/html.txt doc/caddy.xml
	pandoc --read=markdown --write=html --template=doc/html.txt \
		--metadata pagetitle="Pubsub for Caddy" --syntax-definition=doc/caddy.xml < $< > $@

doc.go : doc/document.md doc/cleandoc.awk
	pandoc --read=markdown --write=plain $< | awk -f doc/cleandoc.awk > $@
	gofmt -s -w $@

perm :
	chgrp -R caddy example
	chmod -R u=rwX,g=rX,o= example

clean :
	rm -f coverage.html coverage doc/*.ok example/*.ok doc/index.html doc.go README.md
