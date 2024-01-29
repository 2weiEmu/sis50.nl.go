OBJECTS = src/*.go 
CMD = go
TARGET = build/sis50.nl

default: $(OBJECTS)
	$(CMD) run $(OBJECTS)

build: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS) 
