SERVER_IP = 192.168.0.53
SERVER_PORT = 9091
PROXY_IP = 192.168.0.73
PROXY_PORT = 9090
TIMEOUT = 3
MAX_RETRIES = 4

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
		--client-drop 0 \
		--server-drop 0 \
		--client-delay 0 \
		--server-delay 0 \
		--client-delay-time-min 0 \
		--client-delay-time-max 0 \
		--server-delay-time-min 0 \
		--server-delay-time-max 0
