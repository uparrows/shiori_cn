package cmd

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-shiori/shiori/internal/database"
)

func Test_parseCsvExport_old_format(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
	}{
		{
			name:     "测试旧文件格式",
			fileName: "pocket-old.csv",
		},
		{
			name:     "测试新文件格式",
			fileName: "pocket-new.csv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.Open("../../testdata/" + tt.fileName)
			if err != nil {
				t.Error(err.Error())
			}
			defer file.Close()
			ctx := context.TODO()

			tmpDir, err := os.MkdirTemp("", "shiori-test-*")
			if err != nil {
				t.Fatalf("创建临时目录失败: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			dbPath := filepath.Join(tmpDir, "shiori.db")
			db, err := database.OpenSQLiteDatabase(ctx, dbPath)
			if err != nil {
				t.Fatalf("无法打开 SQLite 数据库: %v", err)
			}

			if err := db.Migrate(ctx); err != nil {
				t.Fatalf("迁移 SQLite 数据库失败: %v", err)
			}

			bookmarks := parseCsvExport(ctx, db, file)
			if len(bookmarks) != 1 {
				t.Errorf("预期有 1 个书签, 实际 %d", len(bookmarks))
			}
			bm := bookmarks[0]
			if bm.Title != "Shiori" {
				t.Errorf("预期标题 Shiori 得到 %s", bm.URL)
			}
			if bm.URL != "https://github.com/go-shiori/shiori" {
				t.Errorf("预期网址 https://github.com/go-shiori/shiori, 得到 %s", bm.URL)
			}
			if len(bm.Tags) != 1 {
				t.Errorf("预期1个标签, 得到 %d", len(bm.Tags))
			}
			if bm.Tags[0].Name != "shiori" {
				t.Errorf("预期标签 shiori, 得到 %s", bm.Tags[0].Name)
			}
			if bm.CreatedAt == "" {
				t.Error("预期创建时间不为空")
			}
			if bm.ModifiedAt == "" {
				t.Error("预期创建时间不为空")
			}
		})
	}
}
