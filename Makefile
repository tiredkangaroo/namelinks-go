
build: #build go code
	@ cp .git/refs/heads/main version.txt
	@ go build -o urlshortner.out *.go
