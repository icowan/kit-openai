/**
 * @Time : 2023/3/17 12:04 PM
 * @Author : solacowa@gmail.com
 * @File : service_test
 * @Software: GoLand
 */
package kitopenai

import (
	"context"
	"fmt"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
	"net/http/httputil"
	"testing"
)

func initSvc() Service {
	return New("", []kithttp.ClientOption{
		kithttp.ClientBefore(func(ctx context.Context, request *http.Request) context.Context {
			dump, err := httputil.DumpRequest(request, true)
			if err != nil {
				return ctx
			}
			fmt.Println(string(dump))
			return ctx
		}),
		kithttp.ClientAfter(func(ctx context.Context, response *http.Response) context.Context {
			dump, err := httputil.DumpResponse(response, true)
			if err != nil {
				return ctx
			}
			fmt.Println(string(dump))

			return ctx
		}),
	})
}

func TestService_Models(t *testing.T) {
	svc := initSvc()

	models, err := svc.Models(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range models.Data {
		t.Log(v.ID)
	}
}

func TestService_Model(t *testing.T) {
	svc := initSvc()

	model, err := svc.Model(context.Background(), "text-davinci-003")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(model.ID)
	for _, v := range model.Permissions {
		t.Log(v.ID)
	}
}

func TestService_ChatCompletions(t *testing.T) {
	svc := initSvc()
	res, err := svc.ChatCompletions(context.Background(), ChatCompletionsRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello!",
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res.ID)
	for _, v := range res.Choices {
		t.Log(v.Message.Content)
	}
}
