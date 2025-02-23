.PHONY: build clean docker

IMAGE?=ronenlib/kube-event-alert:1.0.1

exec:=kube-event-alert

build: clean $(exec)

clean:
	rm -f $(exec)

docker:
	docker build -t $(IMAGE) .

e2e: docker
	docker run --rm $(IMAGE)

test:
	go test -count=1 --race -v ./...

$(exec):
	go build -o $(exec) .
