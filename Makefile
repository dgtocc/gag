build:
	go build -o gag ./src

clean_release:
	- rm -rf ./dist/*

release: clean_release
	GOOS=windows go build -ldflags  "-s -w" -o ./dist/$(VER)/windows/gag.exe ./cmd
	GOOS=linux go build -ldflags  "-s -w" -o ./dist/$(VER)/linux/gag ./cmd
	GOOS=darwin go build -ldflags  "-s -w" -o ./dist/$(VER)/darwin/gag ./cmd
	zip ./dist/gag-windows-$(VER).zip ./dist/$(VER)/windows/gag.exe
	zip ./dist/gag-linux-$(VER).zip ./dist/$(VER)/linux/gag
	zip ./dist/gag-darwin-$(VER).zip ./dist/$(VER)/darwin/gag

