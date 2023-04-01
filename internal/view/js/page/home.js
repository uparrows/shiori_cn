var template = `
<div id="page-home">
    <div class="page-header">
        <input type="text" placeholder="搜索网址、关键字或标签" v-model.trim="search" @focus="$event.target.select()" @keyup.enter="searchBookmarks"/>
        <a title="刷新内容" @click="reloadData">
            <i class="fas fa-fw fa-sync-alt" :class="loading && 'fa-spin'"></i>
        </a>
        <a v-if="activeAccount.owner" title="添加新书签" @click="showDialogAdd">
            <i class="fas fa-fw fa-plus-circle"></i>
        </a>
        <a v-if="tags.length > 0" title="显示标签" @click="showDialogTags">
            <i class="fas fa-fw fa-tags"></i>
        </a>
        <a v-if="activeAccount.owner" title="批量编辑" @click="toggleEditMode">
            <i class="fas fa-fw fa-pencil-alt"></i>
        </a>
    </div>
    <div class="page-header" id="edit-box" v-if="editMode">
        <p>{{selection.length}} 个项目已选择</p>
        <a title="删除书签" @click="showDialogDelete(selection)">
            <i class="fas fa-fw fa-trash-alt"></i>
        </a>
        <a title="添加标签" @click="showDialogAddTags(selection)">
            <i class="fas fa-fw fa-tags"></i>
        </a>
        <a title="更新档案" @click="showDialogUpdateCache(selection)">
            <i class="fas fa-fw fa-cloud-download-alt"></i>
        </a>
        <a title="取消" @click="toggleEditMode">
            <i class="fas fa-fw fa-times"></i>
        </a>
    </div>
    <p class="empty-message" v-if="!loading && listIsEmpty">尚未保存书签 :(</p>
    <div id="bookmarks-grid" ref="bookmarksGrid" :class="{list: appOptions.listMode}">
        <pagination-box v-if="maxPage > 1" 
            :page="page" 
            :maxPage="maxPage" 
            :editMode="editMode"
            @change="changePage">
        </pagination-box>
        <bookmark-item v-for="(book, index) in bookmarks" 
            :id="book.id"
            :url="book.url"
            :title="book.title"
            :excerpt="book.excerpt"
            :public="book.public"
            :imageURL="book.imageURL"
            :hasContent="book.hasContent"
            :hasArchive="book.hasArchive"
            :tags="book.tags"
            :index="index"
            :key="book.id" 
            :editMode="editMode"
            :showId="appOptions.showId"
            :listMode="appOptions.listMode"
            :hideThumbnail="appOptions.hideThumbnail"
            :hideExcerpt="appOptions.hideExcerpt"
            :selected="isSelected(book.id)"
            :menuVisible="activeAccount.owner"
            @select="toggleSelection"
            @tag-clicked="bookmarkTagClicked"
            @edit="showDialogEdit"
            @delete="showDialogDelete"
            @update="showDialogUpdateCache">
        </bookmark-item>
        <pagination-box v-if="maxPage > 1" 
            :page="page" 
            :maxPage="maxPage" 
            :editMode="editMode"
            @change="changePage">
        </pagination-box>
    </div>
    <div class="loading-overlay" v-if="loading"><i class="fas fa-fw fa-spin fa-spinner"></i></div>
    <custom-dialog id="dialog-tags" v-bind="dialogTags">
        <a @click="filterTag('*')">(全部已标记)</a>
        <a @click="filterTag('*', true)">(全部未标记)</a>
        <a v-for="tag in tags" @click="dialogTagClicked($event, tag)">
            #{{tag.name}}<span>{{tag.nBookmarks}}</span>
        </a>
    </custom-dialog>
    <custom-dialog v-bind="dialog"/>
</div>`

import paginationBox from "../component/pagination.js";
import bookmarkItem from "../component/bookmark.js";
import customDialog from "../component/dialog.js";
import basePage from "./base.js";

export default {
	template: template,
	mixins: [basePage],
	components: {
		bookmarkItem,
		paginationBox,
		customDialog
	},
	data() {
		return {
			loading: false,
			editMode: false,
			selection: [],

			search: "",
			page: 0,
			maxPage: 0,
			bookmarks: [],
			tags: [],

			dialogTags: {
				visible: false,
				editMode: false,
				title: '现有标签',
				mainText: '确认',
				secondText: '重命名标签',
				mainClick: () => {
					if (this.dialogTags.editMode) {
						this.dialogTags.editMode = false;
					} else {
						this.dialogTags.visible = false;
					}
				},
				secondClick: () => {
					this.dialogTags.editMode = true;
				},
				escPressed: () => {
					this.dialogTags.visible = false;
					this.dialogTags.editMode = false;
				}
			},
		}
	},
	computed: {
		listIsEmpty() {
			return this.bookmarks.length <= 0;
		}
	},
	watch: {
		"dialogTags.editMode"(editMode) {
			if (editMode) {
				this.dialogTags.title = "重命名标签";
				this.dialogTags.mainText = "取消";
				this.dialogTags.secondText = "";
			} else {
				this.dialogTags.title = "现有标签";
				this.dialogTags.mainText = "确认";
				this.dialogTags.secondText = "重命名标签";
			}
		}
	},
	methods: {
		reloadData() {
			if (this.loading) return;
			this.page = 1;
			this.search = "";
			this.loadData(true, true);
		},
		loadData(saveState, fetchTags) {
			if (this.loading) return;

			// Set default args
			saveState = (typeof saveState === "boolean") ? saveState : true;
			fetchTags = (typeof fetchTags === "boolean") ? fetchTags : false;

			// Parse search query
			var keyword = this.search,
				rxExcludeTagA = /(^|\s)-tag:["']([^"']+)["']/i, // -tag:"with space"
				rxExcludeTagB = /(^|\s)-tag:(\S+)/i, // -tag:without-space
				rxIncludeTagA = /(^|\s)tag:["']([^"']+)["']/i, // tag:"with space"
				rxIncludeTagB = /(^|\s)tag:(\S+)/i, // tag:without-space
				tags = [],
				excludedTags = [],
				rxResult;

			// Get excluded tag first, while also removing it from keyword
			while (rxResult = rxExcludeTagA.exec(keyword)) {
				keyword = keyword.replace(rxResult[0], "");
				excludedTags.push(rxResult[2]);
			}

			while (rxResult = rxExcludeTagB.exec(keyword)) {
				keyword = keyword.replace(rxResult[0], "");
				excludedTags.push(rxResult[2]);
			}

			// Get included tags
			while (rxResult = rxIncludeTagA.exec(keyword)) {
				keyword = keyword.replace(rxResult[0], "");
				tags.push(rxResult[2]);
			}

			while (rxResult = rxIncludeTagB.exec(keyword)) {
				keyword = keyword.replace(rxResult[0], "");
				tags.push(rxResult[2]);
			}

			// Trim keyword
			keyword = keyword.trim().replace(/\s+/g, " ");

			// Prepare URL for API
			var url = new URL("api/bookmarks", document.baseURI);
			url.search = new URLSearchParams({
				keyword: keyword,
				tags: tags.join(","),
				exclude: excludedTags.join(","),
				page: this.page
			});

			// Fetch data from API
			var skipFetchTags = Error("跳过获取标签");

			this.loading = true;
			fetch(url)
				.then(response => {
					if (!response.ok) throw response;
					return response.json();
				})
				.then(json => {
					// Set data
					this.page = json.page;
					this.maxPage = json.maxPage;
					this.bookmarks = json.bookmarks;

					// Save state and change URL if needed
					if (saveState) {
						var history = {
							activePage: "page-home",
							search: this.search,
							page: this.page
						};

						var url = new Url(document.baseURI);
						url.hash = "home";
						url.clearQuery();
						if (this.page > 1) url.query.page = this.page;
						if (this.search !== "") url.query.search = this.search;

						window.history.pushState(history, "page-home", url);
					}

					// Fetch tags if requested
					if (fetchTags) {
						return fetch(new URL("api/tags", document.baseURI));
					} else {
						this.loading = false;
						throw skipFetchTags;
					}
				})
				.then(response => {
					if (!response.ok) throw response;
					return response.json();
				})
				.then(json => {
					this.tags = json;
					this.loading = false;
				})
				.catch(err => {
					this.loading = false;

					if (err !== skipFetchTags) {
						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					}
				});
		},
		searchBookmarks() {
			this.page = 1;
			this.loadData();
		},
		changePage(page) {
			this.page = page;
			this.$refs.bookmarksGrid.scrollTop = 0;
			this.loadData();
		},
		toggleEditMode() {
			this.selection = [];
			this.editMode = !this.editMode;
		},
		toggleSelection(item) {
			var idx = this.selection.findIndex(el => el.id === item.id);
			if (idx === -1) this.selection.push(item);
			else this.selection.splice(idx, 1);
		},
		isSelected(bookId) {
			return this.selection.findIndex(el => el.id === bookId) > -1;
		},
		dialogTagClicked(event, tag) {
			if (!this.dialogTags.editMode) {
				this.filterTag(tag.name, event.altKey);
			} else {
				this.dialogTags.visible = false;
				this.showDialogRenameTag(tag);
			}
		},
		bookmarkTagClicked(event, tagName) {
			this.filterTag(tagName, event.altKey);
		},
		filterTag(tagName, excludeMode) {
			// Set default parameter
			excludeMode = (typeof excludeMode === "boolean") ? excludeMode : false;

			if (this.dialogTags.editMode) {
				return;
			}

			if (tagName === "*") {
				this.search = excludeMode ? "-tag:*" : "tag:*";
				this.page = 1;
				this.loadData();
				return;
			}

			var rxSpace = /\s+/g,
				includeTag = rxSpace.test(tagName) ? `tag:"${tagName}"` : `tag:${tagName}`,
				excludeTag = "-" + includeTag,
				rxIncludeTag = new RegExp(`(^|\\s)${includeTag}`, "ig"),
				rxExcludeTag = new RegExp(`(^|\\s)${excludeTag}`, "ig"),
				search = this.search;

			search = search.replace("-tag:*", "");
			search = search.replace("tag:*", "");
			search = search.trim();

			if (excludeMode) {
				if (rxExcludeTag.test(search)) {
					return;
				}

				if (rxIncludeTag.test(search)) {
					this.search = search.replace(rxIncludeTag, "$1" + excludeTag);
				} else {
					search += ` ${excludeTag}`;
					this.search = search.trim();
				}
			} else {
				if (rxIncludeTag.test(search)) {
					return;
				}

				if (rxExcludeTag.test(search)) {
					this.search = search.replace(rxExcludeTag, "$1" + includeTag);
				} else {
					search += ` ${includeTag}`;
					this.search = search.trim();
				}
			}

			this.page = 1;
			this.loadData();
		},
		showDialogAdd() {
			this.showDialog({
				title: "新书签",
				content: "创建新书签",
				fields: [{
					name: "url",
					label: "输入链接,以 http://开头",
				}, {
					name: "title",
					label: "自定义标题（可选）"
				}, {
					name: "excerpt",
					label: "自定义摘录（可选）",
					type: "area"
				}, {
					name: "tags",
					label: "标签以逗号分隔（可选）",
					separator: ",",
					dictionary: this.tags.map(tag => tag.name)
				}, {
					name: "createArchive",
					label: "创建存档",
					type: "check",
					value: this.appOptions.useArchive,
				}, {
					name: "makePublic",
					label: "公开存档",
					type: "check",
					value: this.appOptions.makePublic,
				}],
				mainText: "确认",
				secondText: "取消",
				mainClick: (data) => {
					// Make sure URL is not empty
					if (data.url.trim() === "") {
						this.showErrorDialog("网址不能为空");
						return;
					}

					// Prepare tags
					var tags = data.tags
						.toLowerCase()
						.replace(/\s+/g, " ")
						.split(/\s*,\s*/g)
						.filter(tag => tag.trim() !== "")
						.map(tag => {
							return {
								name: tag.trim()
							};
						});

					// Send data
					var data = {
						url: data.url.trim(),
						title: data.title.trim(),
						excerpt: data.excerpt.trim(),
						public: data.makePublic ? 1 : 0,
						tags: tags,
						createArchive: data.createArchive,
					};

					this.dialog.loading = true;
					fetch(new URL("api/bookmarks", document.baseURI), {
						method: "post",
						body: JSON.stringify(data),
						headers: { "Content-Type": "application/json" }
					}).then(response => {
						if (!response.ok) throw response;
						return response.json();
					}).then(json => {
						this.dialog.loading = false;
						this.dialog.visible = false;
						this.bookmarks.splice(0, 0, json);
					}).catch(err => {
						this.dialog.loading = false;
						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				}
			});
		},
		showDialogEdit(item) {
			// Check the item
			if (typeof item !== "object") return;

			var id = (typeof item.id === "number") ? item.id : 0,
				index = (typeof item.index === "number") ? item.index : -1;

			if (id < 1 || index < 0) return;

			// Get the existing bookmark value
			var book = JSON.parse(JSON.stringify(this.bookmarks[index])),
				strTags = book.tags.map(tag => tag.name).join(", ");

			this.showDialog({
				title: "编辑书签",
				content: "编辑书签的数据",
				showLabel: true,
				fields: [{
					name: "url",
					label: "链接",
					value: book.url,
				}, {
					name: "title",
					label: "标题",
					value: book.title,
				}, {
					name: "excerpt",
					label: "摘录",
					type: "area",
					value: book.excerpt,
				}, {
					name: "tags",
					label: "标签",
					value: strTags,
					separator: ",",
					dictionary: this.tags.map(tag => tag.name)
				}, {
					name: "makePublic",
					label: "使存案公开可用",
					type: "check",
					value: book.public >= 1,
				}],
				mainText: "确认",
				secondText: "取消",
				mainClick: (data) => {
					// Validate input
					if (data.title.trim() === "") return;

					// Prepare tags
					var tags = data.tags
						.toLowerCase()
						.replace(/\s+/g, " ")
						.split(/\s*,\s*/g)
						.filter(tag => tag.trim() !== "")
						.map(tag => {
							return {
								name: tag.trim()
							};
						});

					// Set new data
					book.url = data.url.trim();
					book.title = data.title.trim();
					book.excerpt = data.excerpt.trim();
					book.public = data.makePublic ? 1 : 0;
					book.tags = tags;

					// Send data
					this.dialog.loading = true;
					fetch(new URL("api/bookmarks", document.baseURI), {
						method: "put",
						body: JSON.stringify(book),
						headers: { "Content-Type": "application/json" }
					}).then(response => {
						if (!response.ok) throw response;
						return response.json();
					}).then(json => {
						this.dialog.loading = false;
						this.dialog.visible = false;
						this.bookmarks.splice(index, 1, json);
					}).catch(err => {
						this.dialog.loading = false;
						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				}
			});
		},
		showDialogDelete(items) {
			// Check and filter items
			if (typeof items !== "object") return;
			if (!Array.isArray(items)) items = [items];

			items = items.filter(item => {
				var id = (typeof item.id === "number") ? item.id : 0,
					index = (typeof item.index === "number") ? item.index : -1;

				return id > 0 && index > -1;
			});

			if (items.length === 0) return;

			// Split ids and indices
			var ids = items.map(item => item.id),
				indices = items.map(item => item.index).sort((a, b) => b - a);

			// Create title and content
			var title = "删除书签",
				content = "删除选定的书签？ 这个动作是不可逆的。";

			if (items.length === 1) {
				title = "删除书签";
				content = "你确定吗 ？ 这个动作是不可逆的。";
			}

			// Show dialog
			this.showDialog({
				title: title,
				content: content,
				mainText: "是",
				secondText: "否",
				mainClick: () => {
					this.dialog.loading = true;
					fetch(new URL("api/bookmarks", document.baseURI), {
						method: "delete",
						body: JSON.stringify(ids),
						headers: { "Content-Type": "application/json" },
					}).then(response => {
						if (!response.ok) throw response;
						return response;
					}).then(() => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;
						this.dialog.visible = false;
						indices.forEach(index => this.bookmarks.splice(index, 1))

						if (this.bookmarks.length < 20) {
							this.loadData(false);
						}
					}).catch(err => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;

						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				}
			});
		},
		showDialogUpdateCache(items) {
			// Check and filter items
			if (typeof items !== "object") return;
			if (!Array.isArray(items)) items = [items];

			items = items.filter(item => {
				var id = (typeof item.id === "number") ? item.id : 0,
					index = (typeof item.index === "number") ? item.index : -1;

				return id > 0 && index > -1;
			});

			if (items.length === 0) return;

			// Show dialog
			var ids = items.map(item => item.id);

			this.showDialog({
				title: "更新缓存",
				content: "更新所选书签的缓存？ 这个动作是不可逆的。",
				fields: [{
					name: "keepMetadata",
					label: "保留旧标题和摘录",
					type: "check",
					value: this.appOptions.keepMetadata,
				}, {
					name: "createArchive",
					label: "同样更新存档",
					type: "check",
					value: this.appOptions.useArchive,
				}],
				mainText: "是",
				secondText: "否",
				mainClick: (data) => {
					var data = {
						ids: ids,
						createArchive: data.createArchive,
						keepMetadata: data.keepMetadata,
					};

					this.dialog.loading = true;
					fetch(new URL("api/cache", document.baseURI), {
						method: "put",
						body: JSON.stringify(data),
						headers: { "Content-Type": "application/json" },
					}).then(response => {
						if (!response.ok) throw response;
						return response.json();
					}).then(json => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;
						this.dialog.visible = false;

						json.forEach(book => {
							var item = items.find(el => el.id === book.id);
							this.bookmarks.splice(item.index, 1, book);
						});
					}).catch(err => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;

						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				}
			});
		},
		showDialogAddTags(items) {
			// Check and filter items
			if (typeof items !== "object") return;
			if (!Array.isArray(items)) items = [items];

			items = items.filter(item => {
				var id = (typeof item.id === "number") ? item.id : 0,
					index = (typeof item.index === "number") ? item.index : -1;

				return id > 0 && index > -1;
			});

			if (items.length === 0) return;

			// Show dialog
			this.showDialog({
				title: "添加新标签",
				content: "为选定的书签添加新标签",
				fields: [{
					name: "tags",
					label: "以逗号分隔标签",
					value: "",
					separator: ",",
					dictionary: this.tags.map(tag => tag.name)
				}],
				mainText: '确定',
				secondText: '取消',
				mainClick: (data) => {
					// Validate input
					var tags = data.tags
						.toLowerCase()
						.replace(/\s+/g, ' ')
						.split(/\s*,\s*/g)
						.filter(tag => tag.trim() !== '')
						.map(tag => {
							return {
								name: tag.trim()
							};
						});

					if (tags.length === 0) return;

					// Send data
					var request = {
						ids: items.map(item => item.id),
						tags: tags
					}

					this.dialog.loading = true;
					fetch(new URL("api/bookmarks/tags", document.baseURI), {
						method: "put",
						body: JSON.stringify(request),
						headers: { "Content-Type": "application/json" },
					}).then(response => {
						if (!response.ok) throw response;
						return response.json();
					}).then(json => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;
						this.dialog.visible = false;

						json.forEach(book => {
							var item = items.find(el => el.id === book.id);
							this.bookmarks.splice(item.index, 1, book);
						});
					}).catch(err => {
						this.selection = [];
						this.editMode = false;
						this.dialog.loading = false;

						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				}
			});
		},
		showDialogTags() {
			this.dialogTags.visible = true;
			this.dialogTags.editMode = false;
			this.dialogTags.secondText = this.activeAccount.owner ? "重命名标签" : "";
		},
		showDialogRenameTag(tag) {
			this.showDialog({
				title: "重命名标签",
				content: `更改标签的名称 "#${tag.name}"`,
				fields: [{
					name: "newName",
					label: "新标签名",
					value: tag.name,
				}],
				mainText: "确定",
				secondText: "取消",
				secondClick: () => {
					this.dialog.visible = false;
					this.dialogTags.visible = true;
				},
				escPressed: () => {
					this.dialog.visible = false;
					this.dialogTags.visible = true;
				},
				mainClick: (data) => {
					// Save the old query
					var rxSpace = /\s+/g,
						oldTagQuery = rxSpace.test(tag.name) ? `"#${tag.name}"` : `#${tag.name}`,
						newTagQuery = rxSpace.test(data.newName) ? `"#${data.newName}"` : `#${data.newName}`;

					// Send data
					var newData = {
						id: tag.id,
						name: data.newName,
					};

					this.dialog.loading = true;
					fetch(new URL("api/tag", document.baseURI), {
						method: "PUT",
						body: JSON.stringify(newData),
						headers: { "Content-Type": "application/json" },
					}).then(response => {
						if (!response.ok) throw response;
						return response.json();
					}).then(() => {
						tag.name = data.newName;

						this.dialog.loading = false;
						this.dialog.visible = false;
						this.dialogTags.visible = true;
						this.dialogTags.editMode = false;
						this.tags.sort((a, b) => {
							var aName = a.name.toLowerCase(),
								bName = b.name.toLowerCase();

							if (aName < bName) return -1;
							else if (aName > bName) return 1;
							else return 0;
						});

						if (this.search.includes(oldTagQuery)) {
							this.search = this.search.replace(oldTagQuery, newTagQuery);
							this.loadData();
						}
					}).catch(err => {
						this.dialog.loading = false;
						this.dialogTags.visible = false;
						this.dialogTags.editMode = false;
						this.getErrorMessage(err).then(msg => {
							this.showErrorDialog(msg);
						})
					});
				},
			});
		},
	},
	mounted() {
		// Prepare history state watcher
		var stateWatcher = (e) => {
			var state = e.state || {},
				activePage = state.activePage || "page-home",
				search = state.search || "",
				page = state.page || 1;

			if (activePage !== "page-home") return;

			this.page = page;
			this.search = search;
			this.loadData(false);
		}

		window.addEventListener('popstate', stateWatcher);
		this.$once('hook:beforeDestroy', () => {
			window.removeEventListener('popstate', stateWatcher);
		})

		// Set initial parameter
		var url = new Url;
		this.search = url.query.search || "";
		this.page = url.query.page || 1;

		this.loadData(false, true);
	}
}