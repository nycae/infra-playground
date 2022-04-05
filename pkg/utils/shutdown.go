package utils

import (
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	switch sig := <-sigChan; sig {
	case syscall.SIGINT, syscall.SIGKILL:
		return
	case syscall.SIGTERM:
		os.Exit(1)
	}
}

func GetEnvWithDefault(key, def string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return def
	}

	return envVar
}

func MustGetDial(addr string, opts ...grpc.DialOption) *grpc.ClientConn {
	dial, err := grpc.Dial(addr, opts...)
	if err != nil {
		panic(err)
	}

	return dial
}
