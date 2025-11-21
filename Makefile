OBJECTS = .
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

required_files: $(OBJECTS)
	mkdir resources
	mkdir log
	mkdir build
	touch secret.conf
	touch log/sis50.log
	touch resources/calendar
	touch resources/centralDb
	touch resources/messages
	touch resources/shopping

	echo "0000\n0000\n0000\n0000\n0000\n0000\n0000\n" > resources/calendar
	
	sqlite3 -line ./resources/centralDb 'CREATE TABLE users(id integer primary key autoincrement, username not null unique, password_hash not null); CREATE TABLE sessions(user_id integer, session_token text not null);'

	# adding a random secret to a secret.conf file
	head -c 32 /dev/random > secret.conf
