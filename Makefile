.NOTPARALLEL:

workTenantId := 72f988bf-86f1-41af-91ab-2d7cd011db47
workSubscriptionId := 27b750cd-ed43-42fd-9044-8d75e124ae55
workClientId := undefined
workClientSecret := undefined

# personal client id is 'azkube-ci'
# this application runs as 'azkube'
personalTenantId := 13de0a15-b5db-44b9-b682-b4ba82afbd29
personalSubscriptionId := aff271ee-e9be-4441-b9bb-42f5af4cbaeb
personalClientId := 20f97fda-60b5-4557-9100-947b9db06ec0
personalClientSecret := $(shell cat clientsecret)

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

deploy-cs-personal:
	go run main.go deploy \
		--tenant-id="$(personalTenantId)" \
		--subscription-id="$(personalSubscriptionId)" \
		--auth-method="clientsecret" \
		--client-id="$(personalClientId)" \
		--client-secret="$(personalClientSecret)"

deploy-device-personal:
	go run main.go deploy \
		--tenant-id="$(personalTenantId)" \
		--subscription-id="$(personalSubscriptionId)" \
		--auth-method="device"

deploy-cs-work:
	go run main.go deploy \
		--tenant-id="$(workTenantId)" \
		--subscription-id="$(workSubscriptionId)" \
		--auth-method="clientsecret" \
		--client-id="$(workClientId)" \
		--client-secret="$(workClientSecret)"

deploy-device-work:
	go run main.go deploy \
		--tenant-id="$(workTenantId)" \
		--subscription-id="$(workSubscriptionId)" \
		--auth-method="device"

clean:
	rm -f azkube
	rm -f state.json

clean-run: clean run
