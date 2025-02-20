build:
	-rm onepage
	go build -o onepage cmd/main/*.go
	./onepage

cleanjson:
	cp ./storage/default_jsons/* ./storage/