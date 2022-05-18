This repo is a fork of go-shiori/shiori ,what i've done was to translate this repo into chinese .
DO NOT UPLOADED TO CSDN

本仓库是go-shiori/shiori的分支，我所做的只是将其汉化为中文版。
请勿上传到csdn。

教程详见：https://yuanfangblog.xyz/technology/601.html




# Shiori

[![IC](https://github.com/go-shiori/shiori/actions/workflows/push.yml/badge.svg?branch=master)](https://github.com/go-shiori/shiori/actions/workflows/push.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-shiori/shiori)](https://goreportcard.com/report/github.com/go-shiori/shiori)
[![#shiori@libera.chat](https://img.shields.io/badge/irc-%23shiori-orange)](https://web.libera.chat/#shiori)
[<img src="https://img.shields.io/docker/pulls/dezhao/shiori_cn.svg">](https://hub.docker.com/r/dezhao/shiori_cn)

**查看我们最新的[公告](https://github.com/go-shiori/shiori/discussions/categories/announcements)**

Shiori 是一个用 Go 语言编写的简单书签管理器。 旨在作为 [Pocket][pocket] 的简单克隆。 您可以将其用作命令行应用程序或 Web 应用程序。 此应用程序作为单个二进制文件分发，这意味着它可以轻松安装和使用。

![1](https://user-images.githubusercontent.com/38988286/169091978-da3b42d4-283d-4ae0-87ca-025910e90144.jpg)


![2](https://user-images.githubusercontent.com/38988286/169091998-1bbc7c19-7795-42e3-b5bd-90813b782dd8.jpg)

![3](https://user-images.githubusercontent.com/38988286/169092022-62ad0e3a-4c6d-46ad-92ee-799fc6e5b5d7.jpg)



## 功能

- 基本书签管理，即添加、编辑、删除和搜索。
- 从 Netscape 书签文件导入和导出书签。
- 从 Pocket 导入书签。
- 简单干净的命令行界面。
- 对于那些不想使用命令行应用程序的人来说，简单而漂亮的网络界面。
- 便携，得益于其单一的二进制格式。
- 支持 sqlite3、PostgreSQL 和 MySQL 作为其数据库。
- 在可能的情况下，默认情况下 `shiori` 将解析可读内容并创建网页的离线存档。
- 浏览器插件支持 Firefox 和 Chrome。

![Comparison of reader mode and archive mode][mode-comparison]

## 文档

所有文档都可以在 [wiki][wiki] 中找到。 如果您认为信息不完整或不正确，请随时编辑。

## 许可

Shiori 是根据 [MIT license][mit] 的条款分发的，这意味着您可以使用它并根据需要进行修改。 但是，如果您对其进行了增强，如果可能，请发送拉取请求。

[wiki]: https://github.com/go-shiori/shiori/wiki
[mit]: https://choosealicense.com/licenses/mit/
[web-extension]: https://github.com/go-shiori/shiori-web-ext
[screenshot]: https://raw.githubusercontent.com/go-shiori/shiori/master/docs/readme/cover.png
[mode-comparison]: https://raw.githubusercontent.com/go-shiori/shiori/master/docs/readme/comparison.png
[pocket]: https://getpocket.com/
[256]: https://github.com/go-shiori/shiori/issues/256
