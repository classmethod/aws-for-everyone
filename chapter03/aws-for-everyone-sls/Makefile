start-localstack:
	@TMPDIR=/private${TMPDIR} docker-compose up -d

test:
	@cd src/backend/persons && \
	go test -v ./...
