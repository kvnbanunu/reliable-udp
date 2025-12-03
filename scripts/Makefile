SERVER_IP = 192.168.0.53
SERVER_PORT = 9091
PROXY_IP = 192.168.0.73
PROXY_PORT = 9090
TIMEOUT = 3
MAX_RETRIES = 4
CLIENT_DROP = 0
SERVER_DROP = 0
CLIENT_DELAY = 0
SERVER_DELAY = 0
CLIENT_DELAY_MIN = 0
CLIENT_DELAY_MAX = 0
SERVER_DELAY_MIN = 0
SERVER_DELAY_MAX = 0

rserver:
	./server \
		--listen-ip $(SERVER_IP) \
		--listen-port $(SERVER_PORT)

rclient:
	./client \
		--target-ip $(SERVER_IP) \
		--target-port $(SERVER_PORT) \
		--timeout $(TIMEOUT) \
		--max-retries $(MAX_RETRIES)

rclient2:
	./client \
		--target-ip $(PROXY_IP) \
		--target-port $(PROXY_PORT) \
		--timeout $(TIMEOUT) \
		--max-retries $(MAX_RETRIES)

rproxy:
	./proxy \
		--listen-ip $(PROXY_IP) \
		--listen-port $(PROXY_PORT) \
		--target-ip $(SERVER_IP) \
		--target-port $(SERVER_PORT) \
		--client-drop $(CLIENT_DROP) \
		--server-drop $(SERVER_DROP) \
		--client-delay $(CLIENT_DELAY) \
		--server-delay $(SERVER_DELAY) \
		--client-delay-time-min $(CLIENT_DELAY_MIN) \
		--client-delay-time-max $(CLIENT_DELAY_MAX) \
		--server-delay-time-min $(SERVER_DELAY_MIN) \
		--server-delay-time-max $(SERVER_DELAY_MAX)
