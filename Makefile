objects = src/*.go 

run: $(objects)
	@go run $(objects)

deploy:
	@go build -o server $(objects) 
	@./server -deploy
