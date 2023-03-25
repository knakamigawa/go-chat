MAKEFILE_DIR:=$(dir $(abspath $(lastword $(MAKEFILE_LIST))))

wire:
	cd $(MAKEFILE_DIR)app/di; wire; cd $(MAKEFILE_DIR)

sqlc:
	cd $(MAKEFILE_DIR)infra/database; sqlc generate; cd $(MAKEFILE_DIR)

migrate-create:
	migrate create -ext sql -dir infra/database/migrations -seq $(NAME)