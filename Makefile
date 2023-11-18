objects = src/*.go 

run: $(objects)
	@go run $(objects)

