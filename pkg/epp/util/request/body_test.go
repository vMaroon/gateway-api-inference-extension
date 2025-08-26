/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package request

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"sigs.k8s.io/gateway-api-inference-extension/pkg/epp/scheduling/types"
)

func TestExtractRequestData(t *testing.T) {
	tests := []struct {
		name    string
		body    map[string]any
		want    *types.LLMRequestData
		wantErr bool
	}{
		{
			name: "completions request body",
			body: map[string]any{
				"model":  "test",
				"prompt": "test prompt",
			},
			want: &types.LLMRequestData{
				Completions: &types.CompletionsRequest{
					Prompt: "test prompt",
				},
			},
		},
		{
			name: "chat completions request body",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{
						"role": "system", "content": "this is a system message",
					},
					map[string]any{
						"role": "user", "content": "hello",
					},
				},
			},
			want: &types.LLMRequestData{
				ChatCompletions: &types.ChatCompletionsRequest{
					Messages: []types.Message{
						{Role: "system", Content: "this is a system message"},
						{Role: "user", Content: "hello"},
					},
				},
			},
		},
		{
			name: "chat completions with all optional fields",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"tools":                        []any{map[string]any{"type": "function"}},
				"documents":                    []any{map[string]any{"content": "doc"}},
				"chat_template":                "custom template",
				"return_assistant_tokens_mask": true,
				"continue_final_message":       true,
				"add_generation_prompt":        true,
				"chat_template_kwargs":         map[string]any{"key": "value"},
			},
			want: &types.LLMRequestData{
				ChatCompletions: &types.ChatCompletionsRequest{
					Messages:                  []types.Message{{Role: "user", Content: "hello"}},
					Tools:                     []any{map[string]any{"type": "function"}},
					Documents:                 []any{map[string]any{"content": "doc"}},
					ChatTemplate:              "custom template",
					ReturnAssistantTokensMask: true,
					ContinueFinalMessage:      true,
					AddGenerationPrompt:       true,
					ChatTemplateKWArgs:        map[string]any{"key": "value"},
				},
			},
		},
		{
			name:    "nil body",
			body:    nil,
			wantErr: true,
		},
		{
			name: "invalid prompt format",
			body: map[string]any{
				"model":  "test",
				"prompt": 123,
			},
			wantErr: true,
		},
		{
			name: "invalid messages format",
			body: map[string]any{
				"model":    "test",
				"messages": "invalid",
			},
			wantErr: true,
		},
		{
			name: "neither prompt nor messages",
			body: map[string]any{
				"model": "test",
			},
			wantErr: true,
		},
		{
			name: "empty messages array",
			body: map[string]any{
				"model":    "test",
				"messages": []any{},
			},
			wantErr: true,
		},
		{
			name: "message with non-string role",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": 123, "content": "hello"},
				},
			},
			wantErr: true,
		},
		{
			name: "message with non-string content",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": 123},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid tools format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"tools": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid documents format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"documents": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid chat_template format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"chat_template": 123,
			},
			wantErr: true,
		},
		{
			name: "invalid return_assistant_tokens_mask format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"return_assistant_tokens_mask": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid continue_final_message format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"continue_final_message": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid add_generation_prompt format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"add_generation_prompt": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid chat_template_kwargs format",
			body: map[string]any{
				"model": "test",
				"messages": []any{
					map[string]any{"role": "user", "content": "hello"},
				},
				"chat_template_kwargs": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractRequestData(tt.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractRequestData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ExtractRequestData() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Benchmark tests for performance comparison
func BenchmarkExtractRequestData_Completions(b *testing.B) {
	body := map[string]any{
		"model":  "test",
		"prompt": "test prompt",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractRequestData(body)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExtractRequestData_ChatCompletions(b *testing.B) {
	body := map[string]any{
		"model": "test",
		"messages": []any{
			map[string]any{"role": "user", "content": "hello"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractRequestData(body)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkExtractRequestData_ChatCompletionsWithOptionals(b *testing.B) {
	body := map[string]any{
		"model": "test",
		"messages": []any{
			map[string]any{"role": "user", "content": "hello"},
		},
		"tools":                        []any{map[string]any{"type": "function"}},
		"documents":                    []any{map[string]any{"content": "doc"}},
		"chat_template":                "custom template",
		"return_assistant_tokens_mask": true,
		"continue_final_message":       true,
		"add_generation_prompt":        true,
		"chat_template_kwargs":         map[string]any{"key": "value"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractRequestData(body)
		if err != nil {
			b.Fatal(err)
		}
	}
}
