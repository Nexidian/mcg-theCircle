GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=main
NAME=listenQuestionTrigger
ZIP_NAME=listenQuestionTrigger.zip

all: test build
build:
	@rm -rf dist
	@mkdir dist
	@echo "===== Building $(BINARY_NAME) ($(NAME)) ======"
	env GOOS=linux $(GOCMD) build -o $(BINARY_NAME)
	zip $(ZIP_NAME) main
	@echo "\n"
	@mv $(ZIP_NAME) ./dist/
	@rm -f main
	@rm -f $(ZIP_NAME)
test: 
	@echo "===== Running Tests ($(NAME)) ======"
	$(GOTEST)
	@echo "\n"
clean:
	rm -rf dist
	rm -f main
	rm -f $(ZIP_NAME)
