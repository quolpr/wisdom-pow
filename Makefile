lint:
	golangci-lint run ./...

test:
	go test -v ./...

vet:
	go vet ./...

run-server:
	@echo "Building Docker Image..."
	@docker_image_id=$$(docker build -q -f ./build/server/Dockerfile .) ; \
	if [ -z "$$docker_image_id" ]; then \
		echo "Docker build failed, image ID not found." ; \
		exit 1 ; \
	else \
		echo "Running Docker Container using Image ID $$docker_image_id" ; \
		docker run -p 8080:8080 --rm -it $$docker_image_id ; \
	fi

run-client:
	@echo "Building Docker Image..."
	@docker_image_id=$$(docker build -q -f ./build/client/Dockerfile .) ; \
	if [ -z "$$docker_image_id" ]; then \
		echo "Docker build failed, image ID not found." ; \
		exit 1 ; \
	else \
		echo "Running Docker Container using Image ID $$docker_image_id" ; \
		docker run --network host --rm -it $$docker_image_id ; \
	fi
