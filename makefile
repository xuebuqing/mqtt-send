.PHONY: ctrlapp

demo:
	@export GO111MODULE=on && \
	export GOPROXY=https://goproxy.io && \
	go build ctrlapp.go
	@chmod 777 ctrlapp


.PHONY: clean
clean:
	@rm -rf ctrlapp
	@echo "[clean Done]"
