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

env-logs:
	@sudo docker compose logs -f pg

env-shell:
	@sudo docker exec -it todoapp-pg-1 psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)