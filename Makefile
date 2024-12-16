.PHONY: run
run:
	$(BINARIES)
	docker compose up -d
	sleep 5
	go run main.go

.PHONY: clean
clean:
	docker compose down

certs/fitbit-manager.local.walnuts.dev.pem:
	mkcert -cert-file ./certs/fitbit-manager.local.walnuts.dev.pem -key-file ./certs/fitbit-manager.local.walnuts.dev-key.pem fitbit-manager.local.walnuts.dev
	mkcert -install

.PHONY: up
up: certs/fitbit-manager.local.walnuts.dev.pem
	docker-compose up -d
