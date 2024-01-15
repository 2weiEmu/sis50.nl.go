OBJECTS = src/*.go 
CMD = go
TARGET = sis50.nl

default: $(OBJECTS)
	$(CMD) run $(OBJECTS)

build: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS) 
