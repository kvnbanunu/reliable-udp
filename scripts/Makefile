rserver:
	./server \
		--listen-ip 192.168.0.53 \
		--listen-port 9091

rclient:
	./client \
		--target-ip 192.168.0.53 \
		--target-port 9091 \
		--timeout 3 \
		--max-retries 4

rclient2:
	./client \
		--target-ip 192.168.0.73 \
		--target-port 9090 \
		--timeout 3 \
		--max-retries 4

rproxy:
	./proxy \
		--listen-ip 192.168.0.73 \
		--listen-port 9090 \
		--target-ip 192.168.0.53 \
		--target-port 9091 \
		--client-drop 0 \
		--server-drop 0 \
		--client-delay 0 \
		--server-delay 0 \
		--client-delay-time-min 0 \
		--client-delay-time-max 0 \
		--server-delay-time-min 0 \
		--server-delay-time-max 0
