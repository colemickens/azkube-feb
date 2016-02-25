.NOTPARALLEL:

subscriptionId := "aff271ee-e9be-4441-b9bb-42f5af4cbaeb"
tenantId := "13de0a15-b5db-44b9-b682-b4ba82afbd29"
clientId := "20f97fda-60b5-4557-9100-947b9db06ec0"
clientSecret := $(shell cat clientsecret)

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

deploy-cs:
	go run main.go deploy \
		--tenant-id="$(tenantId)" \
		--subscription-id="$(subscriptionId)" \
		--auth-method="clientsecret" \
		--client-id="$(clientId)" \
		--client-secret="$(clientSecret)"

deploy-device:
	go run main.go deploy \
		--tenant-id="$(tenantId)" \
		--subscription-id="$(subscriptionId)" \
		--auth-method="device"

clean:
	rm -f azkube
	rm -f state.json

clean-run: clean run
