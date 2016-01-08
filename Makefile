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

docker: build
	docker build -t azkube .

docker-push: docker
	docker tag -f azkube "colemickens/azkube:latest"
	docker push "colemickens/azkube"



# temporary for dev because I'm lazy

create-common:
	./azkube create-common \
		--location "westus" \
		--subscription-id "aff271ee-e9be-4441-b9bb-42f5af4cbaeb" \
		--tenant-id "13de0a15-b5db-44b9-b682-b4ba82afbd29"

create-app:
	./azkube create-app

create-pki:
	./azkube create-pki

create-ssh:
	./azkube create-ssh

deploy-vault:
	./azkube deploy-vault

upload-secrets:
	./azkube upload-secrets

deploy-myriad:
	./azkube deploy-myriad \
		--option \
		--option \
		--option

clean:
	rm -f azkube
	rm -f state.json

run: quick create-common create-app create-pki create-ssh deploy-vault upload-secrets deploy-myriad

clean-run: clean run
