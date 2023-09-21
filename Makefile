objects = src/main.go src/notes.go src/grid.go

run: $(objects)
	@go run $(objects)
