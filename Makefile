OBJECTS = src/*.go 
CMD = go
TARGET = build/sis50.nl

default: $(OBJECTS)
	$(CMD) run $(OBJECTS)

build: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS) 

deploy: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS)
	killall sis50.nl
	./build/sis50.nl -p 80 -base="sis50.nl"
