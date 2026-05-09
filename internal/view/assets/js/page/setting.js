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
        <details v-if="activeAccount.owner" open class="setting-group setting-accounts" id="setting-accounts">
            <summary>账户</summary>
            <ul class="accounts-list">
                <li v-if="accounts.length === 0">没有注册账户</li>
                <li v-for="(account, idx) in accounts" :shiori-username="account.username">
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
                <a @click="loadAccounts" title="Refresh accounts">刷新帐户列表</a>
                <a v-if="activeAccount.owner" @click="showDialogNewAccount" title="Add new account">添加新帐户</a>
            </div>
        </details>
        <details v-if="!activeAccount.owner" open class="setting-group setting-accounts" id="setting-my-account">
            <summary>我的账户</summary>
            <ul>
                <li v-for="(account, idx) in [this.activeAccount]" :shiori-username="account.username">
                    <p>{{account.username}}
                        <span v-if="account.owner" class="account-level">(owner)</span>
                    </p>
                    <a title="修改密码" @click="showDialogChangePassword(account)">
                        <i class="fa fas fa-fw fa-key"></i>
                    </a>
                </li>
            </ul>
            <div class="setting-group-footer">
                <a @click="showDialogChangePassword(this.activeAccount)" title="Change password">修改密码</a>
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
import { apiRequest } from "../utils/api.js";

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
					const responseString = JSON.stringify(responseData);
					localStorage.setItem("shiori-account", responseString);
				})
				.catch((err) => {
					this.getErrorMessage(err).then((msg) => {
						this.showErrorDialog(msg);
					});
				});
		},
		async loadAccounts() {
			if (this.loading) return;

			this.loading = true;
			try {
				const json = await apiRequest(
					new URL("api/v1/accounts", document.baseURI),
				);
				this.loading = false;
				this.accounts = json;
			} catch (err) {
				this.loading = false;
				this.showErrorDialog(err.message);
			}
		},
		async loadSystemInfo() {
			if (this.system.version !== undefined) return;

			try {
				const json = await apiRequest(
					new URL("api/v1/system/info", document.baseURI),
				);
				this.system = json;
			} catch (err) {
				this.showErrorDialog(err.message);
			}
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
						name: "repeat_password",
						label: "重复密码",
						type: "password",
						value: "",
					},
					{
						name: "admin",
						label: "这是一个管理员帐户.",
						type: "check",
						value: false,
					},
				],
				mainText: "确认",
				secondText: "取消",
				mainClick: async (data) => {
					if (data.username === "") {
						return;
					}

					var request = {
						username: data.username,
						password: data.password,
						owner: data.admin,
					};

					this.dialog.loading = true;
					try {
						const json = await apiRequest(
							new URL("api/v1/accounts", document.baseURI),
							{
								method: "post",
								body: JSON.stringify(request),
							},
						);

						this.dialog.loading = false;
						this.dialog.visible = false;

						this.accounts.push(json);
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
					} catch (err) {
						this.dialog.loading = false;
						this.showErrorDialog(err.message);
					}
				},
			});
		},
		showDialogChangePassword(account) {
			let fields = [
				{
					name: "new_password",
					label: "新密码",
					type: "password",
					value: "",
				},
				{
					name: "repeat_password",
					label: "重复密码",
					type: "password",
					value: "",
				},
			];

			const requiresOldPassword =
				!this.activeAccount.owner || this.activeAccount.id === account.id;

			// Only owners can update user passwords without
			// providing the old password

			if (requiresOldPassword) {
				fields.unshift({
					name: "old_password",
					label: "当前密码",
					type: "password",
					value: "",
				});
			}

			this.showDialog({
				title: "更改密码",
				content: "",
				fields: fields,
				mainText: "确认",
				secondText: "取消",
				mainClick: async (data) => {
					if (requiresOldPassword) {
						if (data.old_password === "") {
							this.showErrorDialog("您必须提供当前密码.");
							return;
						}
					}

					if (data.new_password === "") {
						this.showErrorDialog("新密码不能为空");
						return;
					}

					if (data.new_password !== data.repeat_password) {
						this.showErrorDialog("密码不匹配");
						return;
					}

					var request = {
						old_password: data.old_password,
						new_password: data.new_password,
					};

					// Determine which URL to use depending if the user is updating its own
					// account or another user's account.
					let url = `api/v1/accounts/${account.id}`;
					if (this.activeAccount.id === account.id) {
						url = "api/v1/auth/account";
					}

					this.dialog.loading = true;
					try {
						await apiRequest(new URL(url, document.baseURI), {
							method: "PATCH",
							body: JSON.stringify(request),
						});

						this.showDialog({
							title: "密码已更改",
							content: "密码已更改.",
							mainText: "确认",
							mainClick: () => {
								this.dialog.visible = false;
							},
						});
					} catch (err) {
						this.dialog.loading = false;
						this.showErrorDialog(err.message);
					}
				},
			});
		},
		showDialogDeleteAccount(account, idx) {
			this.showDialog({
				title: "删除帐户",
				content: `删除帐户 "${account.username}" ?`,
				mainText: "是",
				secondText: "否",
				mainClick: async () => {
					this.dialog.loading = true;
					try {
						await apiRequest(`api/v1/accounts/${account.id}`, {
							method: "DELETE",
						});

						this.dialog.loading = false;
						this.dialog.visible = false;
						this.accounts.splice(idx, 1);
					} catch (err) {
						this.dialog.loading = false;
						this.showErrorDialog(err.message);
					}
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
