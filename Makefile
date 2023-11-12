CLIENT_BINARY=frontend
BROKER_BINARY=brokerService
AUTH_BINARY=authService
LOGGER_BINARY=loggerService
MAILER_BINARY=mailerService

up: down build_broker build_auth build_logger build_mailer
	@echo "Starting containers..."
	sudo docker compose up --build -d
	@echo "Containers started!"

down:
	@echo "Shutting down containers..."
	sudo docker compose down
	@echo "Done!"

restart: down up
	@echo "Restarting continers..."

build_broker:
	@echo "Building broker service binary..."
	cd ./broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${BROKER_BINARY} ./cmd/api
	@echo "Done!"

build_auth:
	@echo "Building auth service binary..."
	cd ./authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${AUTH_BINARY} ./cmd/api
	@echo "Done!"

build_logger:
	@echo "Building logger service binary..."
	cd ./logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${LOGGER_BINARY} ./cmd/api
	@echo "Done!"

build_mailer:
	@echo "Building mail service binary..."
	cd ./mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${MAILER_BINARY} ./cmd/api
	@echo "Done!"

build_client:
	@echo "Building client binary..."
	cd ./client && env GOOS=linux CGO_ENABLED=0 go build -o ./bin/${CLIENT_BINARY} ./cmd/web
	@echo "Done!"
