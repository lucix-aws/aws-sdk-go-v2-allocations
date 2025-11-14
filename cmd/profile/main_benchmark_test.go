package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	besting "testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

func BenchmarkListQueues(b *besting.B) {
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

	for i := 0; i < b.N; i++ {
		_, err = svc.ListQueues(context.Background(), nil)
		if err != nil {
			panic(err)
		}
	}
}
