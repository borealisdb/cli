CLUSTER_NAME ?= borealis-example
HOSTNAME = borealisdb.io
ENVIRONMENT = development

bind:
	go-bindata -o pkg/templates/templates_bind.go -pkg templates pkg/templates/

clean:
	rm -rf ${HOME}/.borealis/ || true
	rm -r ./borealis* || true

run.help:
	go run main.go --help

run.init: clean bind
	go run main.go init --log debug --host "https://$(HOSTNAME):8443"

run.create: bind
	go run main.go create --cluster-name $(CLUSTER_NAME) --environment $(ENVIRONMENT) --host "https://$(HOSTNAME):8443" --log debug

run.deploy: bind
	go run main.go deploy --environment $(ENVIRONMENT)

run.login:
	go run main.go login --log debug

run.cluster.deploy: bind
	go run main.go cluster deploy --cluster-name $(CLUSTER_NAME) --environment $(ENVIRONMENT) --log debug

run.cluster.token:
	go run main.go cluster token --cluster-name $(CLUSTER_NAME) --log debug

run.cluster.connect:
	go run main.go cluster connect --cluster-name $(CLUSTER_NAME) --username admin

precommit: bind

install-dep:
	go install github.com/goreleaser/goreleaser@latest
	go install github.com/caarlos0/svu@latest