package domains_test

import (
	"context"
	"testing"

	"github.com/go-shiori/shiori/internal/model"
	"github.com/go-shiori/shiori/internal/testutil"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBookmarkTagOperations(t *testing.T) {
	ctx := context.Background()
	logger := logrus.New()

	// Setup using the test configuration and dependencies
	_, deps := testutil.GetTestConfigurationAndDependencies(t, ctx, logger)
	bookmarksDomain := deps.Domains().Bookmarks()
	tagsDomain := deps.Domains().Tags()
	db := deps.Database()

	// Create a test bookmark
	bookmark := model.BookmarkDTO{
		URL:   "https://example.com/bookmark-tags-test",
		Title: "Bookmark Tags Test",
	}
	savedBookmarks, err := db.SaveBookmarks(ctx, true, bookmark)
	require.NoError(t, err)
	require.Len(t, savedBookmarks, 1)
	bookmarkID := savedBookmarks[0].ID

	// Create a test tag
	tagDTO := model.TagDTO{
		Tag: model.Tag{
			Name: "test-tag",
		},
	}
	createdTag, err := tagsDomain.CreateTag(ctx, tagDTO)
	require.NoError(t, err)
	tagID := createdTag.ID

	// Test BookmarkExists
	t.Run("BookmarkExists", func(t *testing.T) {
		// Test with existing bookmark
		exists, err := bookmarksDomain.BookmarkExists(ctx, bookmarkID)
		require.NoError(t, err)
		assert.True(t, exists, "书签应该存在")

		// Test with non-existent bookmark
		exists, err = bookmarksDomain.BookmarkExists(ctx, 9999)
		require.NoError(t, err)
		assert.False(t, exists, "不存在的书签不应该存在")
	})

	// Test TagExists
	t.Run("TagExists", func(t *testing.T) {
		// Test with existing tag
		exists, err := tagsDomain.TagExists(ctx, tagID)
		require.NoError(t, err)
		assert.True(t, exists, "标签应该存在")

		// Test with non-existent tag
		exists, err = tagsDomain.TagExists(ctx, 9999)
		require.NoError(t, err)
		assert.False(t, exists, "不存在的标签不应该存在")
	})

	// Test AddTagToBookmark
	t.Run("AddTagToBookmark", func(t *testing.T) {
		// Add tag to bookmark
		err := bookmarksDomain.AddTagToBookmark(ctx, bookmarkID, tagID)
		require.NoError(t, err)

		// Verify tag was added by listing tags for the bookmark
		tags, err := tagsDomain.ListTags(ctx, model.ListTagsOptions{
			BookmarkID: bookmarkID,
		})
		require.NoError(t, err)
		require.Len(t, tags, 1, "应该只有一个标签")
		assert.Equal(t, tagID, tags[0].ID, "标签 ID 应匹配")
		assert.Equal(t, "test-tag", tags[0].Name, "标签名称应匹配")

		// Test adding the same tag again (should not error)
		err = bookmarksDomain.AddTagToBookmark(ctx, bookmarkID, tagID)
		require.NoError(t, err, "再次添加相同的标签不应该出错。")

		// Test adding tag to non-existent bookmark
		err = bookmarksDomain.AddTagToBookmark(ctx, 9999, tagID)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrBookmarkNotFound, "应返回“找不到书签”错误")

		// Test adding non-existent tag to bookmark
		err = bookmarksDomain.AddTagToBookmark(ctx, bookmarkID, 9999)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrTagNotFound, "应返回“未找到标签”错误")
	})

	// Test RemoveTagFromBookmark
	t.Run("RemoveTagFromBookmark", func(t *testing.T) {
		// Remove tag from bookmark
		err := bookmarksDomain.RemoveTagFromBookmark(ctx, bookmarkID, tagID)
		require.NoError(t, err)

		// Verify tag was removed by listing tags for the bookmark
		tags, err := tagsDomain.ListTags(ctx, model.ListTagsOptions{
			BookmarkID: bookmarkID,
		})
		require.NoError(t, err)
		require.Len(t, tags, 0, "移除后不应有任何标签")

		// Test removing a tag that's not associated with the bookmark (should not error)
		err = bookmarksDomain.RemoveTagFromBookmark(ctx, bookmarkID, tagID)
		require.NoError(t, err, "删除未关联的标签不应出错。")

		// Test removing tag from non-existent bookmark
		err = bookmarksDomain.RemoveTagFromBookmark(ctx, 9999, tagID)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrBookmarkNotFound, "应返回“找不到书签”错误")

		// Test removing non-existent tag from bookmark
		err = bookmarksDomain.RemoveTagFromBookmark(ctx, bookmarkID, 9999)
		require.Error(t, err)
		assert.ErrorIs(t, err, model.ErrTagNotFound, "应返回“未找到标签”错误")
	})
}
