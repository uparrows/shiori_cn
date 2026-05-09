package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListTagsOptions_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		options ListTagsOptions
		wantErr bool
	}{
		{
			name: "valid options with search",
			options: ListTagsOptions{
				Search:            "test",
				WithBookmarkCount: true,
				OrderBy:           DBTagOrderByTagName,
			},
			wantErr: false,
		},
		{
			name: "valid options with bookmark ID",
			options: ListTagsOptions{
				BookmarkID:        123,
				WithBookmarkCount: true,
				OrderBy:           DBTagOrderByTagName,
			},
			wantErr: false,
		},
		{
			name: "invalid options with both search and bookmark ID",
			options: ListTagsOptions{
				Search:            "test",
				BookmarkID:        123,
				WithBookmarkCount: true,
				OrderBy:           DBTagOrderByTagName,
			},
			wantErr: true,
		},
		{
			name: "valid options with neither search nor bookmark ID",
			options: ListTagsOptions{
				WithBookmarkCount: true,
				OrderBy:           DBTagOrderByTagName,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.options.IsValid()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "搜索和书签 ID 过滤不能同时使用")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
