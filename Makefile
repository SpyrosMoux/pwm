build:
	go mod tidy
	go build -o ${PWD}/bin/pwm main.go

install: build
	sudo cp ${PWD}/bin/pwm /usr/local/bin/
	chmod +x /usr/local/bin/pwm
