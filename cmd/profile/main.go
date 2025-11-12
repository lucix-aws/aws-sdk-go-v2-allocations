package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"net/http"
	_ "net/http/pprof"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

const runs = 1024 * 1024

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		panic(err)
	}

	svc := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.Region = "us-east-1"
		o.Credentials = credentials.NewStaticCredentialsProvider("akid", "secret", "session")
		o.HTTPClient = smithyhttp.ClientDoFunc(func(r *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader("{}")),
			}, nil
		})
	})

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	for i := 0; i < runs; i++ {
		fmt.Printf("\033[2K\r%d / %d (%.2f%%)", i+1, runs, float64(i+1)/runs*100)
		_, err = svc.ListQueues(context.Background(), nil)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println()

	cmd := exec.Command("pprof", "-http=localhost:9090", "-sample_index=alloc_objects", "localhost:6060/debug/pprof/allocs")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
