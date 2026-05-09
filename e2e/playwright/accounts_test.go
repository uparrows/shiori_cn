package playwright

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-shiori/shiori/e2e/e2eutil"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestE2EAccounts(t *testing.T) {
	// Start a new Shiori container
	container := e2eutil.NewShioriContainer(t, "")
	baseURL := fmt.Sprintf("http://localhost:%s", container.GetPort())

	mainTestHelper, err := NewTestHelper(t, "main")
	require.NoError(t, err)
	defer mainTestHelper.Close()

	t.Run("001 login as admin", func(t *testing.T) {
		// Navigate to the login page
		_, err = mainTestHelper.page.Goto(baseURL)
		mainTestHelper.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := mainTestHelper.page.Locator("#username")
		passwordLocator := mainTestHelper.page.Locator("#password")
		buttonLocator := mainTestHelper.page.Locator(".button")

		// Wait for and fill the login form
		mainTestHelper.Require().NoError(t, usernameLocator.WaitFor(), "Wait for username field")
		mainTestHelper.Require().NoError(t, usernameLocator.Fill("shiori"), "Fill username field")
		mainTestHelper.Require().NoError(t, passwordLocator.Fill("gopher"), "Fill password field")

		// Click login and wait for success
		mainTestHelper.Require().NoError(t, buttonLocator.Click(), "Click login button")
		mainTestHelper.Require().NoError(t, mainTestHelper.page.Locator("#bookmarks-grid").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待书签栏显示")
	})

	t.Run("002 create new admin account", func(t *testing.T) {
		// Navigate to settings page
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`[title="Settings"]`).Click(),
			"点击设置按钮")
		mainTestHelper.Require().NoError(t, mainTestHelper.page.Locator(".setting-container").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待设置页面显示")

		// Click on "Add new account" <a> element
		mainTestHelper.page.Locator(`[title="Add new account"]`).Click()
		mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
			State: playwright.WaitForSelectorStateVisible,
		})

		// Fill modal
		mainTestHelper.page.Locator(`[name="username"]`).Fill("admin2")
		mainTestHelper.page.Locator(`[name="password"]`).Fill("admin2")
		mainTestHelper.page.Locator(`[name="repeat_password"]`).Fill("admin2")
		mainTestHelper.page.Locator(`[name="admin"]`).Check()

		// Click on "Ok" button
		mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click()

		// Wait for modal to disappear
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(1000),
			}),
			"等待模态框消失")

		// Refresh account list
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`a[title="Refresh accounts"]`).Click(),
			"点击账户刷新按钮")
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(".loading-overlay").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(1000),
			}),
			"等待加载提示消失")

		// Check if new account is created
		accountsCount, err := mainTestHelper.page.Locator(".accounts-list li").Count()
		mainTestHelper.Require().NoError(t, err, "统计列表中的帐户数量")
		mainTestHelper.Require().Equal(t, 2, accountsCount, "创建新管理员帐户后，请验证是否存在 2 个管理员帐户")
	})

	t.Run("003 create new user account", func(t *testing.T) {
		// Click on "Add new account" <a> element
		mainTestHelper.page.Locator(`[title="添加新帐户"]`).Click()
		mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		})

		// Fill modal
		mainTestHelper.page.Locator(`[name="username"]`).Fill("user1")
		mainTestHelper.page.Locator(`[name="password"]`).Fill("user1")
		mainTestHelper.page.Locator(`[name="repeat_password"]`).Fill("user1")

		// Click on "Ok" button
		mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click()

		// Wait for modal to disappear
		mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateHidden,
			Timeout: playwright.Float(1000),
		})

		// Refresh account list
		mainTestHelper.page.Locator(`a[title="刷新账户"]`).Click()
		mainTestHelper.page.Locator(".loading-overlay").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateHidden,
			Timeout: playwright.Float(1000),
		})

		// Check if new account is created
		accountsCount, err := mainTestHelper.page.Locator(".accounts-list li").Count()
		mainTestHelper.Require().NoError(t, err, "统计列表中的帐户数量失败")
		mainTestHelper.Require().Equal(t, 3, accountsCount, "创建用户帐户后，预计会生成 3 个帐户。")
	})

	t.Run("004 check admin account created successfully", func(t *testing.T) {
		th, err := NewTestHelper(t, t.Name())
		require.NoError(t, err, "创建测试助手")
		defer th.Close()

		// Navigate to the login page
		_, err = th.page.Goto(baseURL)
		th.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := th.page.Locator("#username")
		passwordLocator := th.page.Locator("#password")
		buttonLocator := th.page.Locator(".button")

		// Wait for and fill the login form
		th.Require().NoError(t, usernameLocator.WaitFor(), "等待用户名输入")
		th.Require().NoError(t, usernameLocator.Fill("admin2"), "请填写用户名")
		th.Require().NoError(t, passwordLocator.Fill("admin2"), "请填写密码")

		// Click login and wait for success
		th.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		th.Require().NoError(t, th.page.Locator("#bookmarks-grid").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待书签栏显示")

		// Navigate to settings
		th.Require().NoError(t,
			th.page.Locator(`[title="Settings"]`).Click(),
			"点击设置按钮")
		th.Require().NoError(t,
			th.page.Locator(".setting-container").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待设置页面显示")

		// Check if can see system info (admin only)
		visible, err := th.page.Locator(`#setting-system-info`).IsVisible()
		th.Require().NoError(t, err, "检查系统信息的可见性")
		th.Require().True(t, visible, "验证管理员用户对系统信息的可见性")
	})

	t.Run("005 check user account created successfully", func(t *testing.T) {
		th, err := NewTestHelper(t, t.Name())
		require.NoError(t, err, "创建测试助手")

		defer th.Close()

		// Navigate to the login page
		_, err = th.page.Goto(baseURL)
		th.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := th.page.Locator("#username")
		passwordLocator := th.page.Locator("#password")
		buttonLocator := th.page.Locator(".button")

		// Wait for and fill the login form
		th.Require().NoError(t, usernameLocator.WaitFor(), "等待用户名输入")
		th.Require().NoError(t, usernameLocator.Fill("user1"), "请填写用户名")
		th.Require().NoError(t, passwordLocator.Fill("user1"), "请填写密码")

		// Click login and wait for success
		th.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		th.Require().NoError(t, th.page.Locator("#bookmarks-grid").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待书签栏显示")

		// Navigate to settings
		th.Require().NoError(t,
			th.page.Locator(`[title="设置"]`).Click(),
			"点击设置按钮")
		th.Require().NoError(t, th.page.Locator(".setting-container").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待设置页面显示")

		// Check if can see system info (admin only)
		visible, err := th.page.Locator(`#setting-system-info`).IsVisible()
		th.Require().NoError(t, err, "检查系统信息的可见性")
		th.Require().False(t, visible, "验证系统信息部分对普通用户不可见")

		// My account settings is visible
		visible, err = th.page.Locator(`#setting-my-account`).IsVisible()
		th.Require().NoError(t, err, "检查帐户设置的可见性")
		th.Require().True(t, visible, "验证用户帐户设置的可见性")

		// Check change password requires current password
		th.Require().NoError(t,
			th.page.Locator(`li[shiori-username="user1"] a[title="更改密码"]`).Click(),
			"点击“更改密码”按钮")
		th.Require().NoError(t,
			th.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待修改密码弹窗出现")
		visible, err = th.page.Locator(`[name="old_password"]`).IsVisible()
		th.Require().NoError(t, err, "检查旧密码字段的可见性")
		th.Require().True(t, visible, "更改密码时，请检查旧密码字段是否可见")

		// Fill modal
		th.Require().NoError(t,
			th.page.Locator(`[name="old_password"]`).Fill("user1"),
			"填写旧密码字段")
		th.Require().NoError(t,
			th.page.Locator(`[name="new_password"]`).Fill("new_user1"),
			"填写新密码字段")
		th.Require().NoError(t,
			th.page.Locator(`[name="repeat_password"]`).Fill("new_user1"),
			"重复新密码字段")

		// Click on "Ok" button
		th.Require().NoError(t,
			th.page.Locator(`.custom-dialog-button.main`).Click(),
			"点击确定按钮")

		// Wait for modal to display text: "密码已更改."
		dialogContent := th.page.Locator(".custom-dialog-content")
		th.Require().NoError(t,
			dialogContent.WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待对话框内容显示")

		contentText, err := dialogContent.TextContent()
		th.Require().NoError(t, err, "获取对话框内容文本")
		th.Require().Equal(t, "密码已更改.", contentText, "验证密码更改确认消息")
	})

	t.Run("006 delete user account", func(t *testing.T) {
		// Click on "Delete" button
		mainTestHelper.page.Locator(`li[shiori-username="user1"] a[title="删除帐户"]`).Click()
		mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		})

		// Click on "Ok" button
		mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click()

		// Wait for modal to disappear
		mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateHidden,
			Timeout: playwright.Float(1000),
		})

		// Refresh account list
		mainTestHelper.page.Locator(`a[title="刷新账户"]`).Click()
		mainTestHelper.page.Locator(".loading-overlay").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateHidden,
			Timeout: playwright.Float(1000),
		})

		// Check if account is deleted
		accountsCount, err := mainTestHelper.page.Locator(".accounts-list li").Count()
		mainTestHelper.Require().NoError(t, err, "统计列表中的帐户数量")
		mainTestHelper.Require().Equal(t, 2, accountsCount, "创建管理员帐户后，请验证是否存在 2 个管理员帐户")

		time.Sleep(5 * time.Second)
	})

	t.Run("007 change password for admin account", func(t *testing.T) {
		// Click on "Change password" button
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`li[shiori-username="admin2"] a[title="更改密码"]`).Click(),
			"点击更改密码按钮")
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待密码对话框出现")

		// Fill modal
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`[name="new_password"]`).Fill("admin3"),
			"填写新密码")
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`[name="repeat_password"]`).Fill("admin3"),
			"重复新密码")

		// Click on "Ok" button
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click(),
			"点击确定按钮")

		// Wait for modal to display text: "密码已更改."
		dialogContent := mainTestHelper.page.Locator(".custom-dialog-content")
		mainTestHelper.Require().NoError(t,
			dialogContent.WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待对话框内容显示")

		contentText, err := dialogContent.TextContent()
		mainTestHelper.Require().NoError(t, err, "获取对话框内容文本")
		mainTestHelper.Require().Equal(t, "密码已更改.", contentText, "验证密码更改确认消息")

		// Click on "Ok" button
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click(),
			"点击确定按钮")

		// Wait for modal to disappear
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(".custom-dialog").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(2000),
			}),
			"等待对话框关闭")

		// Refresh account list
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`a[title="刷新账户"]`).Click(),
			"点击刷新帐户")
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(".loading-overlay").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateHidden,
				Timeout: playwright.Float(1000),
			}),
			"等待刷新完成")

		t.Run("0071 login with new password", func(t *testing.T) {
			th, err := NewTestHelper(t, t.Name())
			require.NoError(t, err, "创建测试助手失败")
			defer th.Close()

			// Navigate to the login page
			_, err = th.page.Goto(baseURL)
			th.Require().NoError(t, err, "导航至基础 URL")

			// Wait for login page
			th.Require().NoError(t,
				th.page.Locator("#username").WaitFor(playwright.LocatorWaitForOptions{
					State:   playwright.WaitForSelectorStateVisible,
					Timeout: playwright.Float(1000),
				}),
				"等待登录页面")
			th.Require().NoError(t, th.page.Locator("#username").Fill("admin2"), "请填写用户名")
			th.Require().NoError(t, th.page.Locator("#password").Fill("admin3"), "请填写密码字段")
			th.Require().NoError(t, th.page.Locator(".button").Click(), "点击登录按钮")
			th.Require().NoError(t, th.page.Locator("#bookmarks-grid").WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}), "等待书签栏显示")
		})
	})

	t.Run("008 logout", func(t *testing.T) {
		// Click on "Logout" button
		mainTestHelper.Require().NoError(t, mainTestHelper.page.Locator(`a[title="注销"]`).Click(), "点击注销按钮")

		// Wait for modal to display text
		dialogContent := mainTestHelper.page.Locator(".custom-dialog-content")
		mainTestHelper.Require().NoError(t,
			dialogContent.WaitFor(playwright.LocatorWaitForOptions{
				State:   playwright.WaitForSelectorStateVisible,
				Timeout: playwright.Float(1000),
			}),
			"等待对话框内容显示")

		contentText, err := dialogContent.TextContent()
		mainTestHelper.Require().NoError(t, err, "获取对话框内容文本")
		mainTestHelper.Require().Equal(t, "您确定要退出登录吗 ?", contentText, "验证注销确认消息")

		// Click on "Yes" button
		mainTestHelper.Require().NoError(t,
			mainTestHelper.page.Locator(`.custom-dialog-button.main`).Click(),
			"点击“是”按钮")

		// Wait for login page
		err = mainTestHelper.page.Locator("#login-scene").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		})
		mainTestHelper.Require().NoError(t, err, "等待登录页面")
	})
}
