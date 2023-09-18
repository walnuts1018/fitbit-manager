.PHONY: run
run:
	$(BINARIES)
	docker compose up -d
	sleep 5
	go run main.go

.PHONY: clean
clean:
	docker compose down

