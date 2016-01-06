.NOTPARALLEL:

all: build

glide:
	GO15VENDOREXPERIMENT=1 glide up
	GO15VENDOREXPERIMENT=1 glide rebuild

build:
	GO15VENDOREXPERIMENT=1 \
	CGO_ENABLED=0 \
	go build -a -tags netgo -installsuffix nocgo -ldflags '-w' .

quick:
	GO15VENDOREXPERIMENT=1 \
	CGO_ENABLED=0 \
	go build .

run:
	./azkube create-common \
		--location "westus" \
		--subscription-id "aff271ee-e9be-4441-b9bb-42f5af4cbaeb" \
		--tenant-id "13de0a15-b5db-44b9-b682-b4ba82afbd29"

	#./azkube create-app \
	#	--this-doesnt-work

	./azkube create-pki

	./azkube create-ssh

	./azkube deploy-vault

	./azkube update-secrets

	./azkube deploy-myriad \
		--option \
		--option \
		--option

docker: build
	docker build -t azkube .

docker-push: docker
	docker tag -f azkube "colemickens/azkube:latest"
	docker push "colemickens/azkube"

