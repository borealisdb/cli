CLUSTER_NAME ?= borealis-example
HOSTNAME = borealisdb.io
ENVIRONMENT = development
CHART_URL ?= ../charts/borealis

bind:
	go-bindata -o pkg/templates/templates_bind.go -pkg templates pkg/templates/

clean:
	rm -rf ${HOME}/.borealis/ || true
	rm -r ./borealis* || true

run.help:
	go run -modfile=localgo.mod main.go --help

run.init: clean bind
	go run -modfile=localgo.mod main.go init --log debug --host "https://$(HOSTNAME):8443" --chart $(CHART_URL)

run.create: bind
	go run -modfile=localgo.mod main.go create --cluster-name $(CLUSTER_NAME) --environment $(ENVIRONMENT) --host "https://$(HOSTNAME):8443" --log debug

run.deploy: bind
	go run -modfile=localgo.mod main.go deploy --environment $(ENVIRONMENT) --chart $(CHART_URL)

run.login:
	go run -modfile=localgo.mod main.go login --log debug

run.cluster.deploy: bind
	go run -modfile=localgo.mod main.go cluster deploy --cluster-name $(CLUSTER_NAME) --environment $(ENVIRONMENT) --log debug

run.cluster.token:
	go run -modfile=localgo.mod main.go cluster token --cluster-name $(CLUSTER_NAME) --log debug

run.cluster.connect:
	go run -modfile=localgo.mod main.go cluster connect --cluster-name $(CLUSTER_NAME) --username admin

precommit: bind

release:
	git tag $(shell svu next)
	git push --tags
	goreleaser release --clean