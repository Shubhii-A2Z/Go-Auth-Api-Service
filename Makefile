include .env
export

# Create a new migration:- make migrate-create name="entity_table_name"
migrate-create:
	goose -dir $(MIGRATIONS_FOLDER) create $(name) sql 

# make migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_FOLDER) mysql "$(DB_URL)" up
