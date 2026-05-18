# Загрузка переменных из .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

env-up:
	@sudo docker compose up -d pg

env-down:
	@sudo docker compose down -v

env-cleanup:
	@read -p "Do you really wanna delete all files? [y/N]: " ans; \
	case "$$ans" in \
		y|Y|yes|YES) \
			sudo docker compose down -v && \
			sudo rm -rv out/pgdata && \
			echo "All done!" ;; \
		*) \
			echo "Cleanup cancelled" ;; \
	esac

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Missing parameter seq, usage: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	sudo docker compose run --rm pg-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)" 

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Missing action parameter, usage: make migrate-action action=up"; \
		exit 1; \
	fi; \
	sudo docker compose run --rm pg-migrate \
		-path /migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@pg:5432/$(POSTGRES_DB)?sslmode=disable" \
		$(action)

env-port-forward:
	@sudo docker compose up -d port-forwarder

env-port-close:
	@sudo docker compose down port-forwarder

logs-cleanup:
	@read -p "Do you really wanna delete all log files? [y/N]: " ans; \
	case "$$ans" in \
		y|Y|yes|YES) \
			sudo rm -rv out/logs && \
			echo "All done!" ;; \
		*) \
			echo "Cleanup cancelled" ;; \
	esac

env-logs:
	@sudo docker compose logs -f pg

env-shell:
	@sudo docker exec -it todoapp-pg-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

todoapp-run:
	@ export LOGGER_FOLDER=./out/logs &&  export POSTGRES_HOST=localhost && go mod tidy && 	go run cmd/todoapp/main.go

swagger-gen:
	@docker compose run --rm --remove-orphans swagger \
		init \
		-g cmd/todoapp/main.go \
		-o docs \
		--parseInternal \
		--parseDependency