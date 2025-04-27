.PHONY: run-all run-main run-validator run-verifier run-producer clean

run-all: run-main run-validator run-verifier run-producer

run-main:
	@echo "Starting main node..."
	@go run mainNode/main.go & echo $$! > .main.pid
	@sleep 2

run-validator:
	@echo "Starting validator..."
	@go run validator/validator.go & echo $$! > .validator.pid
	@sleep 1

run-verifier:
	@echo "Starting verifier..."
	@go run mainNode/tcp/main_tcp.go & echo $$! > .verifier.pid
	@sleep 1

run-producer:
	@echo "Starting producer..."
	@go run producer/main.go & echo $$! > .producer.pid
	@echo "All components running. Press Ctrl+C to stop."
	@trap 'make clean' INT TERM EXIT; \
	while true; do sleep 1; done

clean:
	@echo "Shutting down all components..."
	@-kill $$(cat .main.pid) 2>/dev/null || true
	@-kill $$(cat .validator.pid) 2>/dev/null || true
	@-kill $$(cat .verifier.pid) 2>/dev/null || true
	@-kill $$(cat .producer.pid) 2>/dev/null || true
	@-rm .*.pid 2>/dev/null || true
	@echo "Cleanup complete."