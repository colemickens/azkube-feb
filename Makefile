.NOTPARALLEL:

all: build

glide:
	GO15VENDOREXPERIMENT=1 glide up
	GO15VENDOREXPERIMENT=1 glide rebuild

build:
	GO15VENDOREXPERIMENT=1 \
	CGO_ENABLED=0 \
	go build -a -tags netgo -installsuffix nocgo -ldflags '-w' .

docker: build
	docker build -t azkube .

docker-push: docker
	docker tag -f azkube "colemickens/azkube:latest"
	docker push "colemickens/azkube"

clean:
	rm azkube
