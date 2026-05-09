package model

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookmarkToDTO(t *testing.T) {
	// Create a test bookmark
	bookmark := Bookmark{
		ID:         123,
		URL:        "https://example.com",
		Title:      "Example Title",
		Excerpt:    "This is an excerpt",
		Author:     "John Doe",
		Public:     1,
		CreatedAt:  "2023-01-01 12:00:00",
		ModifiedAt: "2023-01-02 12:00:00",
		HasContent: true,
	}

	// Convert to DTO
	dto := bookmark.ToDTO()

	// Verify all fields are correctly transferred
	assert.Equal(t, bookmark.ID, dto.ID, "ID应该匹配")
	assert.Equal(t, bookmark.URL, dto.URL, "URL应该匹配")
	assert.Equal(t, bookmark.Title, dto.Title, "标题应该匹配")
	assert.Equal(t, bookmark.Excerpt, dto.Excerpt, "摘录内容应该匹配")
	assert.Equal(t, bookmark.Author, dto.Author, "作者应该匹配")
	assert.Equal(t, bookmark.Public, dto.Public, "公共内容应该匹配")
	assert.Equal(t, bookmark.CreatedAt, dto.CreatedAt, "创建时间应该匹配")
	assert.Equal(t, bookmark.ModifiedAt, dto.ModifiedAt, "修改时间应该匹配")
	assert.Equal(t, bookmark.HasContent, dto.HasContent, "包含目录应该匹配")

	// Verify default values for fields not in Bookmark
	assert.Empty(t, dto.Content, "Content应该为空")
	assert.Empty(t, dto.HTML, "HTML应该为空")
	assert.Empty(t, dto.ImageURL, "图片网址应该为空")
	assert.Empty(t, dto.Tags, "标签应该为空")
	assert.False(t, dto.HasArchive, "已存档应该为否")
	assert.False(t, dto.HasEbook, "电子书应该为否")
	assert.False(t, dto.CreateArchive, "创建存档应该为否")
	assert.False(t, dto.CreateEbook, "创建电子书应该为否")
}

func TestBookmarkDTOToBookmark(t *testing.T) {
	// Create a test BookmarkDTO with all fields populated
	dto := BookmarkDTO{
		ID:            123,
		URL:           "https://example.com",
		Title:         "Example Title",
		Excerpt:       "This is an excerpt",
		Author:        "John Doe",
		Public:        1,
		CreatedAt:     "2023-01-01 12:00:00",
		ModifiedAt:    "2023-01-02 12:00:00",
		Content:       "This is the content",
		HTML:          "<p>This is HTML</p>",
		ImageURL:      "https://example.com/image.jpg",
		HasContent:    true,
		Tags:          []TagDTO{{Tag: Tag{ID: 1, Name: "tag1"}}, {Tag: Tag{ID: 2, Name: "tag2"}}},
		HasArchive:    true,
		HasEbook:      true,
		CreateArchive: true,
		CreateEbook:   true,
	}

	// Convert to Bookmark
	bookmark := dto.ToBookmark()

	// Verify all fields are correctly transferred
	assert.Equal(t, dto.ID, bookmark.ID, "ID应该匹配")
	assert.Equal(t, dto.URL, bookmark.URL, "URL应该匹配")
	assert.Equal(t, dto.Title, bookmark.Title, "标题应该匹配")
	assert.Equal(t, dto.Excerpt, bookmark.Excerpt, "摘录内容应该匹配")
	assert.Equal(t, dto.Author, bookmark.Author, "作者应该匹配")
	assert.Equal(t, dto.Public, bookmark.Public, "公共内容应该匹配")
	assert.Equal(t, dto.CreatedAt, bookmark.CreatedAt, "创建时间应该匹配")
	assert.Equal(t, dto.ModifiedAt, bookmark.ModifiedAt, "修改时间应该匹配")
	assert.Equal(t, dto.HasContent, bookmark.HasContent, "包含目录应该匹配")

	// Fields that should not be transferred
	// These fields are only in BookmarkDTO and not in Bookmark
}

func TestGetThumbnailPath(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		bookmark BookmarkDTO
		expected string
	}{
		{
			name: "With ID",
			bookmark: BookmarkDTO{
				ID: 123,
			},
			expected: filepath.Join("thumb", "123"),
		},
		{
			name: "With zero ID",
			bookmark: BookmarkDTO{
				ID: 0,
			},
			expected: filepath.Join("thumb", "0"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := GetThumbnailPath(&tc.bookmark)
			assert.Equal(t, tc.expected, path, "缩略图路径应与预期值匹配")
		})
	}
}

func TestGetEbookPath(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		bookmark BookmarkDTO
		expected string
	}{
		{
			name: "With ID",
			bookmark: BookmarkDTO{
				ID: 123,
			},
			expected: filepath.Join("ebook", "123.epub"),
		},
		{
			name: "With zero ID",
			bookmark: BookmarkDTO{
				ID: 0,
			},
			expected: filepath.Join("ebook", "0.epub"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := GetEbookPath(&tc.bookmark)
			assert.Equal(t, tc.expected, path, "电子书路径应与预期值匹配")
		})
	}
}

func TestGetArchivePath(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		bookmark BookmarkDTO
		expected string
	}{
		{
			name: "With ID",
			bookmark: BookmarkDTO{
				ID: 123,
			},
			expected: filepath.Join("archive", "123"),
		},
		{
			name: "With zero ID",
			bookmark: BookmarkDTO{
				ID: 0,
			},
			expected: filepath.Join("archive", "0"),
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			path := GetArchivePath(&tc.bookmark)
			assert.Equal(t, tc.expected, path, "归档路径应与预期值匹配")
		})
	}
}

func TestBookmarkRoundTrip(t *testing.T) {
	// Test that converting from Bookmark to DTO and back preserves data
	original := Bookmark{
		ID:         123,
		URL:        "https://example.com",
		Title:      "Example Title",
		Excerpt:    "This is an excerpt",
		Author:     "John Doe",
		Public:     1,
		CreatedAt:  "2023-01-01 12:00:00",
		ModifiedAt: "2023-01-02 12:00:00",
		HasContent: true,
	}

	// Convert to DTO and back
	dto := original.ToDTO()
	roundTrip := dto.ToBookmark()

	// Verify all fields are preserved
	assert.Equal(t, original.ID, roundTrip.ID, "ID应该保留")
	assert.Equal(t, original.URL, roundTrip.URL, "URL应该保留")
	assert.Equal(t, original.Title, roundTrip.Title, "标题应该保留")
	assert.Equal(t, original.Excerpt, roundTrip.Excerpt, "摘录内容应该保留")
	assert.Equal(t, original.Author, roundTrip.Author, "作者应该保留")
	assert.Equal(t, original.Public, roundTrip.Public, "公共内容应该保留")
	assert.Equal(t, original.CreatedAt, roundTrip.CreatedAt, "创建时间信息应该保留")
	assert.Equal(t, original.ModifiedAt, roundTrip.ModifiedAt, "修改时间应该保留")
	assert.Equal(t, original.HasContent, roundTrip.HasContent, "包含目录应该保留")
}
