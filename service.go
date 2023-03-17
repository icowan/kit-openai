/**
 * @Time : 2023/3/17 11:11 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package kitopenai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"
)

type ContextKey string

type Middleware func(Service) Service

const (
	ContextKeyOpenApiModel   ContextKey = "ctx-open-api-model"
	ContextKeyOpenApiToken   ContextKey = "ctx-open-api-token"
	ContextKeyOpenApiBaseUrl ContextKey = "ctx-open-api-base-url"
)

type Service interface {
	// Models 获取模型列表
	Models(ctx context.Context) (res ResponseModel, err error)
	// Model 获取单个模型
	Model(ctx context.Context, modelName string) (res Model, err error)

	Completions(ctx context.Context, req CompletionsRequest) (res ResponseCompletions, err error)
	// ChatCompletions 聊天
	ChatCompletions(ctx context.Context, req ChatCompletionsRequest) (res ResponseChatCompletions, err error)
}

type service struct {
	openApiToken, openApiModel, openApiBaseUrl string
	kitOpts                                    []kithttp.ClientOption
}

func (s *service) ChatCompletions(ctx context.Context, req ChatCompletionsRequest) (res ResponseChatCompletions, err error) {
	token, _, baseUrl := s.getConfig(ctx)
	u, _ := url.Parse(fmt.Sprintf("%s/v1/chat/completions", baseUrl))
	_, err = s.endpoint(ctx, http.MethodPost, u, func(ctx context.Context, request *http.Request, i interface{}) error {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		b, _ := json.Marshal(req)
		request.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}, &res)
	return res, err
}

func (s *service) Completions(ctx context.Context, req CompletionsRequest) (res ResponseCompletions, err error) {
	token, _, baseUrl := s.getConfig(ctx)
	u, _ := url.Parse(fmt.Sprintf("%s/v1/completions", baseUrl))
	_, err = s.endpoint(ctx, http.MethodPost, u, func(ctx context.Context, request *http.Request, i interface{}) error {
		request.Header.Add("Content-Type", "application/json")
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		b, _ := json.Marshal(req)
		request.Body = io.NopCloser(bytes.NewReader(b))
		return nil
	}, &res)
	return res, err
}

func (s *service) endpoint(ctx context.Context, method string, tgt *url.URL, enc kithttp.EncodeRequestFunc, data interface{}) (res []byte, err error) {
	ep := kithttp.NewClient(method, tgt, enc, decodeResponse, s.kitOpts...).Endpoint()
	resp, err := ep(ctx, nil)
	if err != nil {
		err = errors.Wrap(err, "kithttp.NewClient.Endpoint")
		return
	}
	if resp != nil {
		if err = json.Unmarshal(resp.([]byte), &data); err != nil {
			return resp.([]byte), errors.Wrap(err, "json.NewDecoder(resp.(io.ReadCloser)).Decode(&res)")
		}
	}
	return nil, nil
}

func (s *service) Model(ctx context.Context, modelName string) (res Model, err error) {
	token, _, baseUrl := s.getConfig(ctx)
	u, _ := url.Parse(fmt.Sprintf("%s/v1/models/%s", baseUrl, modelName))
	_, err = s.endpoint(ctx, http.MethodGet, u, func(ctx context.Context, request *http.Request, i interface{}) error {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}, &res)
	return res, err
}

func (s *service) getConfig(ctx context.Context) (token, model, baseUrl string) {
	model = s.openApiModel
	token = s.openApiToken
	baseUrl = s.openApiBaseUrl
	if val, ok := ctx.Value(ContextKeyOpenApiModel).(string); ok && !strings.EqualFold(val, "") {
		model = val
	}
	if val, ok := ctx.Value(ContextKeyOpenApiToken).(string); ok && !strings.EqualFold(val, "") {
		token = val
	}
	if val, ok := ctx.Value(ContextKeyOpenApiBaseUrl).(string); ok && !strings.EqualFold(val, "") {
		baseUrl = val
	}

	if strings.EqualFold(baseUrl, "") {
		baseUrl = "https://api.openai.com"
	}

	return
}

func (s *service) Models(ctx context.Context) (res ResponseModel, err error) {
	token, _, baseUrl := s.getConfig(ctx)
	u, _ := url.Parse(fmt.Sprintf("%s/v1/models", baseUrl))
	_, err = s.endpoint(ctx, http.MethodGet, u, func(ctx context.Context, request *http.Request, i interface{}) error {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
		return nil
	}, &res)
	return res, err
}

func New(defaultOpenApiToken string, kitOpts []kithttp.ClientOption) Service {
	return &service{openApiToken: defaultOpenApiToken, kitOpts: kitOpts}
}

type HttpResponse struct {
	HttpHeaders http.Header `json:"-"`
	StatusCode  int         `json:"-"`
	Body        []byte      `json:"-"`
}

func decodeResponse(ctx context.Context, resp *http.Response) (res interface{}, err error) {
	if resp.StatusCode != http.StatusOK {
		var re Response
		if resp.Body != nil {
			if err = json.NewDecoder(resp.Body).Decode(&re); err != nil {
				return nil, errors.Wrap(err, "json.NewDecoder.Decode")
			}
			return re, errors.New(re.Error.Message)
		}
		return re, errors.New(resp.Status)
	}

	return io.ReadAll(resp.Body)
}
