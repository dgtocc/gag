build:
	go build -o gag ./src

clean_release:
	- rm -rf ./dist/*

release: clean_release
	GOOS=windows go build -o ./dist/$(VER)/windows/gag.exe ./cmd
	GOOS=linux go build -o ./dist/$(VER)/linux/gag ./cmd
	GOOS=darwin go build -o ./dist/$(VER)/darwin/gag ./cmd
