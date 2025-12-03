SERVER = cmd/server/main.go
CLIENT = cmd/client/main.go
PROXY = cmd/proxy/main.go
SERVER_TARGET = bin/server
CLIENT_TARGET = bin/client
PROXY_TARGET = bin/proxy
BUILD = go build -o
RUN = go run
HELP = -h
SERVER_IP = 127.0.0.1
PROXY_IP = 127.0.0.1
SERVER_PORT = 8080
PROXY_PORT = 8081
TIMEOUT = 2
MAX_RETRIES = 5
SERVER_ARGS = --listen-ip $(SERVER_IP) --listen-port $(SERVER_PORT)
CLIENT_ARGS = --target-ip $(SERVER_IP) --target-port $(SERVER_PORT) --timeout $(TIMEOUT) --max-retries $(MAX_RETRIES)
CLIENT_ARGS_PROXY = --target-ip $(PROXY_IP) --target-port $(PROXY_PORT) --timeout $(TIMEOUT) --max-retries $(MAX_RETRIES)
PROXY_PARAMS = --client-drop 10 --server-drop 5 --client-delay 20 --server-delay 15 --client-delay-time-min 100 --client-delay-time-max 200 --server-delay-time-min 150 --server-delay-time-max 300
PROXY_ARGS = --listen-ip $(PROXY_IP) --listen-port $(PROXY_PORT) --target-ip $(SERVER_IP) --target-port $(SERVER_PORT) $(PROXY_PARAMS)
COPY_CONFIG = cp config.json bin/

all: clean buildserver buildclient buildproxy copymakefile

server: clear
	@$(RUN) $(SERVER) $(SERVER_ARGS)

client: clear
	@$(RUN) $(CLIENT) $(CLIENT_ARGS)

clientp:
	$(RUN) $(CLIENT) $(CLIENT_ARGS_PROXY)

proxy:
	$(RUN) $(PROXY) $(PROXY_ARGS)

buildserver:
	$(BUILD) $(SERVER_TARGET) $(SERVER)
	@$(COPY_CONFIG)

buildclient:
	$(BUILD) $(CLIENT_TARGET) $(CLIENT)
	@$(COPY_CONFIG)

buildproxy:
	$(BUILD) $(PROXY_TARGET) $(PROXY)
	@$(COPY_CONFIG)

clear:
	@clear

clean:
	rm -rf bin

copymakefile:
	@cp scripts/Makefile bin/
