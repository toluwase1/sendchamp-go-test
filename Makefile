up:
	docker compose up --build

generate-mock:
	 mockgen -destination=mocks/user_mock.go -package=mocks sendchamp-go-test/services UserService
	 mockgen -destination=mocks/user_repo_mock.go -package=mocks sendchamp-go-test/db UserRepository

test: generate-mock
	 MEDDLE_ENV=test go test ./...
