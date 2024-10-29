var template = `
<div id="page-setting">
    <h1 class="page-header">设置</h1>
    <div class="setting-container">
        <details open class="setting-group" id="setting-display">
            <summary>显示</summary>
            <label>
                主题 &nbsp;
                <select v-model="appOptions.Theme" @change="saveSetting">
                <option value="follow">跟随系统</option>
                <option value="light">浅色主题</option>
                <option value="dark">深色主题</option>
                </select>
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.ShowId" @change="saveSetting">
                显示书签ID
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.ListMode" @change="saveSetting">
                将书签显示为列表
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.HideThumbnail" @change="saveSetting">
                隐藏缩略图
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.HideExcerpt" @change="saveSetting">
                隐藏书签的节选
            </label>
        </details>
        <details v-if="activeAccount.owner" open class="setting-group" id="setting-bookmarks">
            <summary>书签</summary>
            <label>
                <input type="checkbox" v-model="appOptions.KeepMetadata" @change="saveSetting">
                更新时保留书签的元数据
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.UseArchive" @change="saveSetting">
                默认创建存档
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.CreateEbook" @change="saveSetting">
                默认创建图书
            </label>
            <label>
                <input type="checkbox" v-model="appOptions.MakePublic" @change="saveSetting">
                默认情况下公开书签
            </label>
        </details>
        <details v-if="activeAccount.owner" open class="setting-group" id="setting-accounts">
            <summary>账户</summary>
            <ul>
                <li v-if="accounts.length === 0">没有注册账户</li>
                <li v-for="(account, idx) in accounts">
                    <p>{{account.username}}
                        <span v-if="account.owner" class="account-level">(管理员)</span>
                    </p>
                    <a title="更改密码" @click="showDialogChangePassword(account)">
                        <i class="fa fas fa-fw fa-key"></i>
                    </a>
                    <a title="删除帐户" @click="showDialogDeleteAccount(account, idx)">
                        <i class="fa fas fa-fw fa-trash-alt"></i>
                    </a>
                </li>
            </ul>
            <div class="setting-group-footer">
                <a @click="loadAccounts">刷新帐户列表</a>
                <a v-if="activeAccount.owner" @click="showDialogNewAccount">添加新帐户</a>
            </div>
        </details>
		<details v-if="activeAccount.owner" class="setting-group" id="setting-system-info">
			<summary>系统信息</summary>
			<ul>
				<li><b>Shiori 版本:</b> <span>{{system.version?.tag}}<span></li>
				<li><b>数据库引擎:</b> <span>{{system.database}}</span></li>
				<li><b>操作系统:</b> <span>{{system.os}}</span></li>
			</ul>
	</details>
        <details v-if="activeAccount.owner" open class="setting-group">
            <summary>关于</summary>
            <label>
                &nbsp;&nbsp;翻译：昭君
            </label>
            <label>
                博客：http://yuanfangblog.xyz
            </label>
        </details>
    </div>
    <div class="loading-overlay" v-if="loading"><i class="fas fa-fw fa-spin fa-spinner"></i></div>
    <custom-dialog v-bind="dialog"/>
</div>`;

import customDialog from "../component/dialog.js";
import basePage from "./base.js";

export default {
	template: template,
	mixins: [basePage],
	components: {
		customDialog,
	},
	data() {
		return {
			loading: false,
			accounts: [],
			system: {},
		};
	},
	methods: {
		saveSetting() {
			let options = {
				ShowId: this.appOptions.ShowId,
				ListMode: this.appOptions.ListMode,
				HideThumbnail: this.appOptions.HideThumbnail,
				HideExcerpt: this.appOptions.HideExcerpt,
				Theme: this.appOptions.Theme,
			};

			if (this.activeAccount.owner) {
				options = {
					...options,
					KeepMetadata: this.appOptions.KeepMetadata,
					UseArchive: this.appOptions.UseArchive,
					CreateEbook: this.appOptions.CreateEbook,
					MakePublic: this.appOptions.MakePublic,
				};
			}

			this.$emit("setting-changed", options);
			//request
			fetch(new URL("api/v1/auth/account", document.baseURI), {
				method: "PATCH",
				body: JSON.stringify({
					config: this.appOptions,
				}),
				headers: {
					"Content-Type": "application/json",
					Authorization: "Bearer " + localStorage.getItem("shiori-token"),
				},
			})
				.then((response) => {
					if (!response.ok) throw response;
					return response.json();
				})
				.then((responseData) => {
					const responseString = JSON.stringify(responseData.message);
					localStorage.setItem("shiori-account", responseString);
				})
				.catch((err) => {
					this.getErrorMessage(err).then((msg) => {
						this.showErrorDialog(msg);
					});
				});
		},
		loadAccounts() {
			if (this.loading) return;

			this.loading = true;
			fetch(new URL("api/accounts", document.baseURI), {
				headers: {
					"Content-Type": "application/json",
					Authorization: "Bearer " + localStorage.getItem("shiori-token"),
				},
			})
				.then((response) => {
					if (!response.ok) throw response;
					return response.json();
				})
				.then((json) => {
					this.loading = false;
					this.accounts = json;
				})
				.catch((err) => {
					this.loading = false;
					this.getErrorMessage(err).then((msg) => {
						this.showErrorDialog(msg);
					});
				});
		},
		loadSystemInfo() {
			if (this.system.version !== undefined) return;

			fetch(new URL("api/v1/system/info", document.baseURI), {
				headers: {
					"Content-Type": "application/json",
					Authorization: "Bearer " + localStorage.getItem("shiori-token"),
				},
			})
				.then((response) => {
					if (!response.ok) throw response;
					return response.json();
				})
				.then((json) => {
					this.system = json.message;
				})
				.catch((err) => {
					this.getErrorMessage(err).then((msg) => {
						this.showErrorDialog(msg);
					});
				});
		},
		showDialogNewAccount() {
			this.showDialog({
				title: "新账户",
				content: "输入新账户的数据 :",
				fields: [
					{
						name: "username",
						label: "用户名",
						value: "",
					},
					{
						name: "password",
						label: "密码",
						type: "password",
						value: "",
					},
					{
						name: "repeat",
						label: "重复密码",
						type: "password",
						value: "",
					},
					{
						name: "visitor",
						label: "此帐户供访客使用",
						type: "check",
						value: false,
					},
				],
				mainText: "确认",
				secondText: "取消",
				mainClick: (data) => {
					if (data.username === "") {
						this.showErrorDialog("用户名不能为空");
						return;
					}

					if (data.password === "") {
						this.showErrorDialog("密码不能为空");
						return;
					}

					if (data.password !== data.repeat) {
						this.showErrorDialog("密码不匹配");
						return;
					}

					var request = {
						username: data.username,
						password: data.password,
						owner: !data.visitor,
					};

					this.dialog.loading = true;
					fetch(new URL("api/accounts", document.baseURI), {
						method: "post",
						body: JSON.stringify(request),
						headers: {
							"Content-Type": "application/json",
							Authorization: "Bearer " + localStorage.getItem("shiori-token"),
						},
					})
						.then((response) => {
							if (!response.ok) throw response;
							return response;
						})
						.then(() => {
							this.dialog.loading = false;
							this.dialog.visible = false;

							this.accounts.push({
								username: data.username,
								owner: !data.visitor,
							});
							this.accounts.sort((a, b) => {
								var nameA = a.username.toLowerCase(),
									nameB = b.username.toLowerCase();

								if (nameA < nameB) {
									return -1;
								}

								if (nameA > nameB) {
									return 1;
								}

								return 0;
							});
						})
						.catch((err) => {
							this.dialog.loading = false;
							this.getErrorMessage(err).then((msg) => {
								this.showErrorDialog(msg);
							});
						});
				},
			});
		},
		showDialogChangePassword(account) {
			this.showDialog({
				title: "更改密码",
				content: "输入新密码 :",
				fields: [
					{
						name: "oldPassword",
						label: "旧密码",
						type: "password",
						value: "",
					},
					{
						name: "password",
						label: "新密码",
						type: "password",
						value: "",
					},
					{
						name: "repeat",
						label: "重复密码",
						type: "password",
						value: "",
					},
				],
				mainText: "确认",
				secondText: "取消",
				mainClick: (data) => {
					if (data.oldPassword === "") {
						this.showErrorDialog("旧密码不能为空");
						return;
					}

					if (data.password === "") {
						this.showErrorDialog("新密码不能为空");
						return;
					}

					if (data.password !== data.repeat) {
						this.showErrorDialog("密码不匹配");
						return;
					}

					var request = {
						username: account.username,
						oldPassword: data.oldPassword,
						newPassword: data.password,
						owner: account.owner,
					};

					this.dialog.loading = true;
					fetch(new URL("api/accounts", document.baseURI), {
						method: "put",
						body: JSON.stringify(request),
						headers: {
							"Content-Type": "application/json",
							Authorization: "Bearer " + localStorage.getItem("shiori-token"),
						},
					})
						.then((response) => {
							if (!response.ok) throw response;
							return response;
						})
						.then(() => {
							this.dialog.loading = false;
							this.dialog.visible = false;
						})
						.catch((err) => {
							this.dialog.loading = false;
							this.getErrorMessage(err).then((msg) => {
								this.showErrorDialog(msg);
							});
						});
				},
			});
		},
		showDialogDeleteAccount(account, idx) {
			this.showDialog({
				title: "删除帐户",
				content: `删除帐户 "${account.username}" ?`,
				mainText: "是",
				secondText: "否",
				mainClick: () => {
					this.dialog.loading = true;
					fetch(`api/accounts`, {
						method: "delete",
						body: JSON.stringify([account.username]),
						headers: {
							"Content-Type": "application/json",
							Authorization: "Bearer " + localStorage.getItem("shiori-token"),
						},
					})
						.then((response) => {
							if (!response.ok) throw response;
							return response;
						})
						.then(() => {
							this.dialog.loading = false;
							this.dialog.visible = false;
							this.accounts.splice(idx, 1);
						})
						.catch((err) => {
							this.dialog.loading = false;
							this.getErrorMessage(err).then((msg) => {
								this.showErrorDialog(msg);
							});
						});
				},
			});
		},
	},
	mounted() {
		if (this.activeAccount.owner) {
			this.loadAccounts();
			this.loadSystemInfo();
		}
	},
};
