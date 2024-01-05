OBJECTS = src/*.go 
CMD = go

default: $(OBJECTS)
	$(CMD) run $(OBJECTS)
