OBJECTS = src/*.go 
CMD = go
TARGET = build/sis50.nl

DEPLOY_PORT = 80
DEPLOY_BASE = "sis50.nl"

default: $(OBJECTS)
	$(CMD) run $(OBJECTS)

build: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS) 

test: $(OBJECTS)
	$(CMD) test -coverprofile -v ./...

deploy: $(OBJECTS)
	$(CMD) build -o $(TARGET) $(OBJECTS)
	killall sis50.nl
	./$(TARGET) -p $(DEPLOY_PORT) -base=$(DEPLOY_BASE) & disown
