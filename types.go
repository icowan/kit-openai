/**
 * @Time : 2023/3/17 11:40 AM
 * @Author : solacowa@gmail.com
 * @File : types
 * @Software: GoLand
 */

package kitopenai

type (
	Response struct {
		Error *Error `json:"error"`
	}
	ResponseModel struct {
		Response
		Data   []Model `json:"data"`
		Object string  `json:"object"`
	}
	ResponseCompletions struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Model   string `json:"model"`
		Choices []struct {
			Text         string      `json:"text"`
			Index        int         `json:"index"`
			Logprobs     interface{} `json:"logprobs"`
			FinishReason string      `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}
	ResponseChatCompletions struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		Choices []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
	}

	Error struct {
		Message string      `json:"message"`
		Type    string      `json:"type"`
		Param   interface{} `json:"param"`
		Code    interface{} `json:"code"`
	}
	Model struct {
		ID          string       `json:"id"`
		Object      string       `json:"object"`
		OwnedBy     string       `json:"owned_by"`
		Permissions []Permission `json:"permission"`
		Root        string       `json:"root"`
	}
	Permission struct {
		ID                 string      `json:"id"`
		Object             string      `json:"object"`
		Created            int         `json:"created"`
		AllowCreateEngine  bool        `json:"allow_create_engine"`
		AllowSampling      bool        `json:"allow_sampling"`
		AllowLogprobs      bool        `json:"allow_logprobs"`
		AllowSearchIndices bool        `json:"allow_search_indices"`
		AllowView          bool        `json:"allow_view"`
		AllowFineTuning    bool        `json:"allow_fine_tuning"`
		Organization       string      `json:"organization"`
		Group              interface{} `json:"group"`
		IsBlocking         bool        `json:"is_blocking"`
	}
)

type (
	CompletionsRequest struct {
		Model       string      `json:"model,omitempty"`
		Prompt      string      `json:"prompt,omitempty"`
		MaxTokens   int         `json:"max_tokens,omitempty"`
		Temperature int         `json:"temperature,omitempty"`
		TopP        int         `json:"top_p,omitempty"`
		N           int         `json:"n,omitempty"`
		Stream      bool        `json:"stream,omitempty"`
		Logprobs    interface{} `json:"logprobs,omitempty"`
		Stop        string      `json:"stop,omitempty"`
	}
	ChatCompletionsRequest struct {
		Model    string    `json:"model"`
		Messages []Message `json:"messages"`
	}
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
)
