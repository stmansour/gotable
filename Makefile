SCSS_BIN := sass

gotable: css defaults *.go
	go clean
	go get -t -v ./...
	go vet
	golint
	go build
	go install

clean:
	go clean
	rm -rf *.out *.csv *.html *.txt *.pdf *.css .sass-cache
	rm -f defaults.go

defaults:
	./defaults.sh

css:
	${SCSS_BIN} ./scss/gotable.scss ./gotable.css --precision 5 --style=compressed --no-source-map
	@echo "Current working directory:"
	pwd
	@echo "scss completed.  ls -l ./gotable.css:"
	ls -l ./gotable.css

lint:
	golint

test:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

benchmark:
	go test -bench=. -benchmem

update:
	cp smoke_test.txt smoke_test.csv smoke_test.html smoke_test.pdf smoke_test_custom_template.html testdata/

all: clean gotable test

deps: wkhtmltopdf

wkhtmltopdf:
	./pdfinstall.sh
