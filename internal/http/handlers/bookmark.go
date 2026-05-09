package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-shiori/shiori/internal/http/response"
	"github.com/go-shiori/shiori/internal/model"
	"github.com/gofrs/uuid/v5"
)

// getBookmark retrieves and validates a bookmark by ID from the request
func getBookmark(deps model.Dependencies, c model.WebContext) (*model.BookmarkDTO, error) {
	bookmarkID, err := strconv.Atoi(c.Request().PathValue("id"))
	if err != nil {
		return nil, response.SendError(c, http.StatusNotFound, "无效的书签 ID")
	}

	if bookmarkID == 0 {
		return nil, response.SendError(c, http.StatusNotFound, "未找到书签")
	}

	// Get bookmark from database
	bookmark, err := deps.Domains().Bookmarks().GetBookmark(c.Request().Context(), model.DBID(bookmarkID))
	if err != nil {
		return nil, response.SendError(c, http.StatusNotFound, "未找到书签")
	}

	// Check access permissions
	if bookmark.Public != 1 && !c.UserIsLogged() {
		response.RedirectToLogin(c, deps.Config().Http.RootPath, c.Request().URL.String())
		return nil, nil
	}

	return bookmark, nil
}

// HandleBookmarkContent serves the bookmark content page
func HandleBookmarkContent(deps model.Dependencies, c model.WebContext) {
	bookmark, err := getBookmark(deps, c)
	if err != nil || bookmark == nil {
		return
	}

	data := map[string]any{
		"RootPath": deps.Config().Http.RootPath,
		"Version":  model.BuildVersion,
		"Book":     bookmark,
		"HTML":     template.HTML(bookmark.HTML),
	}

	if err := response.SendTemplate(c, "content.html", data); err != nil {
		deps.Logger().WithError(err).Error("内容模板渲染失败")
	}
}

// HandleBookmarkArchive serves the bookmark archive page
func HandleBookmarkArchive(deps model.Dependencies, c model.WebContext) {
	bookmark, err := getBookmark(deps, c)
	if err != nil || bookmark == nil {
		return
	}

	if !deps.Domains().Bookmarks().HasArchive(bookmark) {
		response.NotFound(c)
		return
	}

	data := map[string]any{
		"RootPath": deps.Config().Http.RootPath,
		"Version":  model.BuildVersion,
		"Book":     bookmark,
	}

	if err := response.SendTemplate(c, "archive.html", data); err != nil {
		deps.Logger().WithError(err).Error("渲染存档模板失败")
	}
}

// HandleBookmarkArchiveFile serves files from the bookmark archive
func HandleBookmarkArchiveFile(deps model.Dependencies, c model.WebContext) {
	bookmark, err := getBookmark(deps, c)
	if err != nil || bookmark == nil {
		return
	}

	if !deps.Domains().Bookmarks().HasArchive(bookmark) {
		response.NotFound(c)
		return
	}

	resourcePath := c.Request().PathValue("path")

	archive, err := deps.Domains().Archiver().GetBookmarkArchive(bookmark)
	if err != nil {
		deps.Logger().WithError(err).Error("打开存档时出错")
		response.SendInternalServerError(c)
		return
	}
	defer archive.Close()

	if !archive.HasResource(resourcePath) {
		response.NotFound(c)
		return
	}

	content, resourceContentType, err := archive.Read(resourcePath)
	if err != nil {
		deps.Logger().WithError(err).Error("读取存档文件时出错")
		response.SendInternalServerError(c)
		return
	}

	// Generate weak ETAG
	shioriUUID := uuid.NewV5(uuid.NamespaceURL, model.ShioriURLNamespace)
	etag := fmt.Sprintf("W/%s", uuid.NewV5(shioriUUID, fmt.Sprintf("%x-%x-%x", bookmark.ID, resourcePath, len(content))))

	c.ResponseWriter().Header().Set("Etag", etag)
	c.ResponseWriter().Header().Set("Cache-Control", "max-age=31536000")
	c.ResponseWriter().Header().Set("Content-Encoding", "gzip")
	c.ResponseWriter().Header().Set("Content-Type", resourceContentType)
	c.ResponseWriter().WriteHeader(http.StatusOK)
	c.ResponseWriter().Write(content)
}

// HandleBookmarkThumbnail serves the bookmark thumbnail
func HandleBookmarkThumbnail(deps model.Dependencies, c model.WebContext) {
	bookmark, err := getBookmark(deps, c)
	if err != nil || bookmark == nil {
		return
	}

	if !deps.Domains().Bookmarks().HasThumbnail(bookmark) {
		response.NotFound(c)
		return
	}

	etag := "w/" + model.GetThumbnailPath(bookmark) + "-" + bookmark.ModifiedAt

	// Check if the client's ETag matches
	if c.Request().Header.Get("If-None-Match") == etag {
		c.ResponseWriter().WriteHeader(http.StatusNotModified)
		return
	}

	options := &response.SendFileOptions{
		Headers: []http.Header{
			{"Cache-Control": {"no-cache, must-revalidate"}},
			{"Last-Modified": {bookmark.ModifiedAt}},
			{"ETag": {etag}},
		},
	}

	response.SendFile(c, deps.Domains().Storage(), model.GetThumbnailPath(bookmark), options)
}

// HandleBookmarkEbook serves the bookmark's ebook file
func HandleBookmarkEbook(deps model.Dependencies, c model.WebContext) {
	bookmark, err := getBookmark(deps, c)
	if err != nil || bookmark == nil {
		return
	}

	ebookPath := model.GetEbookPath(bookmark)
	if !deps.Domains().Storage().FileExists(ebookPath) {
		response.SendError(c, http.StatusNotFound, "未找到电子书")
		return
	}

	c.ResponseWriter().Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.epub"`, bookmark.Title))
	response.SendFile(c, deps.Domains().Storage(), ebookPath, nil)
}
