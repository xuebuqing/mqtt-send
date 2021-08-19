.PHONY: mqtt-send

demo:
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	go build mqtt-send.go
	@chmod 777 mqtt-send


.PHONY: clean
clean:
	@rm -rf mqtt-send
	@echo "[clean Done]"
