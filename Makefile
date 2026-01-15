include .env
export 

run-worker:
	@go run backend/cmd/robot-workflow/worker/main.go