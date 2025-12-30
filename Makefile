
# Variables
ID ?= 1
TITLE ?= Task
DESCRIPTION ?= make via makefile
COMPLETED ?= false

# Docker
build:
	docker build -t skooby-todo .

run: build
	docker run -p 8080:8080 skooby-todo


# API COMMANDS -----------

# POST /todos (make test-create TITLE="New Task")
test-create:
	@echo "--- Creating Task ---"
	@curl -i -X POST localhost:8080/todos \
		-d '{"title": "$(TITLE)", "description": "$(DESCRIPTION)"}'

# GET /todos
test-getall:
	@echo "--- Getting All Tasks ---"
	@curl -i localhost:8080/todos

# GET /todos/{id} (make test-getbyid ID=2)
test-getbyid:
	@echo "--- Getting Task ID=$(ID) ---"
	@curl -i localhost:8080/todos/$(ID)

# PUT /todos/{id} (make test-update ID=2)
test-update:
	@echo "--- Updating Task ID=$(ID) ---"
	@curl -i -X PUT localhost:8080/todos/$(ID) \
		-d '{"title": "$(TITLE)", "description": "$(DESCRIPTION)", "completed": $(COMPLETED)}'

# DELETE /todos/{id} (make test-delete ID=2)
test-delete:
	@echo "--- Deleting Task ID=$(ID) ---"
	@curl -i -X DELETE localhost:8080/todos/$(ID)