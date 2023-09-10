objects = src/main.go src/notes.go

run: $(objects)
	@go run $(objects)
