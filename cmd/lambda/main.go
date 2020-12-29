package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
)

func OnSchedule(ctx context.Context, event events.CloudWatchEvent) (string, error) {
	sess := session.Must(session.NewSession())
	svc := cloudwatchevents.New(sess)

	result, err := svc.ListRules(&cloudwatchevents.ListRulesInput{})
	if err != nil {
		return "", err
	}

	var requestData request

	for _, er := range event.Resources {
		for _, r := range result.Rules {
			if *r.Arn != er {
				continue
			}

			err = json.Unmarshal([]byte(*r.Description), &requestData)
			if err != nil {
				return "", err
			}
		}
	}

	req, err := http.NewRequest(requestData.Method, requestData.URL, bytes.NewBuffer([]byte(requestData.Body)))
	for k, v := range requestData.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	defer resp.Body.Close()

	return "", nil
}

type request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Body    string            `json:"body"`
	Headers map[string]string `json:"headers`
}

func main() {
	lambda.Start(OnSchedule)
}
