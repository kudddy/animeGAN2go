PROJECT_NAME ?= animegan2go
VERSION = $(shell python3 setup.py --version | tr '+' '-')
PROJECT_NAMESPACE ?= kudddy
REGISTRY_IMAGE ?= docker.io/$(PROJECT_NAMESPACE)/$(PROJECT_NAME)
DEPLOYMENTS_NAME ?= simplebackend-deployment


postgres:
	docker run -d --rm --name some-postgres -p 5434:5432 \
		-e POSTGRES_PASSWORD=pass -e POSTGRES_USER=user -e \
		POSTGRES_DB=db postgres:9.6

build:
	docker build -t $(REGISTRY_IMAGE) .

run:
	docker run \
		-p 9000:9000 \
		-e db_name="db" \
		-e db_pass="pass" \
		-e db_user="user" \
		-e db_type="postgres" \
		-e db_host="localhost" \
		-e db_port=5434 \
		-e bot_token="YOURE_TOKEN" \
		-d $(REGISTRY_IMAGE):latest