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
SERVER_ARGS = -i $(SERVER_IP) -p $(SERVER_PORT)
CLIENT_ARGS = -i $(SERVER_IP) -p $(SERVER_PORT)
CLIENT_ARGS_PROXY = -i $(PROXY_IP) -p $(PROXY_PORT)
PROXY_ARGS = $(SERVER_ARGS) -I $(PROXY_IP) -P $(PROXY_PORT)
COPY_CONFIG = cp config.json bin/

all: clean buildserver buildclient buildproxy

server:
	@$(RUN) $(SERVER) $(SERVER_ARGS)

client:
	@$(RUN) $(CLIENT) $(CLIENT_ARGS)

clientp:
	@$(RUN) $(CLIENT) $(CLIENT_ARGS_PROXY)

proxy:
	@$(RUN) $(PROXY) $(PROXY_ARGS)

buildserver:
	@$(BUILD) $(SERVER_TARGET) $(SERVER)
	@$(COPY_CONFIG)

buildclient:
	@$(BUILD) $(CLIENT_TARGET) $(CLIENT)
	@$(COPY_CONFIG)

buildproxy:
	@$(BUILD) $(PROXY_TARGET) $(PROXY)
	@$(COPY_CONFIG)

clean:
	@rm -rf bin
