TARGETS = build/monotonic

all: $(TARGETS)

build/monotonic: cmd/bot/main.go
	go build -o build/monotonic cmd/bot/main.go

run: build/monotonic
	CONFIG_PATH="./config/dev.yaml" ./build/monotonic

clean:
	rm -f $(TARGETS)

.PHONY: all clean build/monotonic run
