package playwright

import (
	"fmt"
	"testing"

	"github.com/go-shiori/shiori/e2e/e2eutil"
	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	// Start a new Shiori container
	container := e2eutil.NewShioriContainer(t, "")
	baseURL := fmt.Sprintf("http://localhost:%s", container.GetPort())

	mainTestHelper, err := NewTestHelper(t, "main")
	require.NoError(t, err)
	defer mainTestHelper.Close()

	t.Run("successful login with default credentials", func(t *testing.T) {
		// Navigate to the login page
		_, err = mainTestHelper.page.Goto(baseURL)
		mainTestHelper.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := mainTestHelper.page.Locator("#username")
		passwordLocator := mainTestHelper.page.Locator("#password")
		buttonLocator := mainTestHelper.page.Locator(".button")

		// Wait for and fill the login form
		mainTestHelper.Require().NoError(t, usernameLocator.WaitFor(), "等待用户名输入")
		mainTestHelper.Require().NoError(t, usernameLocator.Fill("shiori"), "请填写用户名")
		mainTestHelper.Require().NoError(t, passwordLocator.Fill("gopher"), "请填写密码字段")

		// Click login and wait for success
		mainTestHelper.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		mainTestHelper.Require().NoError(t, mainTestHelper.page.Locator("#bookmarks-grid").WaitFor(playwright.LocatorWaitForOptions{
			State:   playwright.WaitForSelectorStateVisible,
			Timeout: playwright.Float(1000),
		}), "等待书签栏显示")
	})

	t.Run("failed login with wrong username", func(t *testing.T) {
		th, err := NewTestHelper(t, t.Name())
		require.NoError(t, err)
		defer th.Close()

		// Navigate to the login page
		_, err = th.page.Goto(baseURL)
		th.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := th.page.Locator("#username")
		passwordLocator := th.page.Locator("#password")
		buttonLocator := th.page.Locator(".button")
		errorLocator := th.page.Locator(".error-message")

		// Wait for and fill the login form
		th.Require().NoError(t, usernameLocator.WaitFor(), "等待用户名输入")
		th.Require().NoError(t, usernameLocator.Fill("wrong_user"), "请填写用户名")
		th.Require().NoError(t, passwordLocator.Fill("gopher"), "请填写密码字段")

		// Click login and verify error
		th.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		errorText, err := errorLocator.TextContent()
		th.Require().NoError(t, err, "获取错误信息文本")
		th.Require().Contains(t, errorText, "用户名或密码不匹配")
	})

	t.Run("failed login with wrong password", func(t *testing.T) {
		th, err := NewTestHelper(t, t.Name())
		require.NoError(t, err)
		defer th.Close()

		// Navigate to the login page
		_, err = th.page.Goto(baseURL)
		th.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := th.page.Locator("#username")
		passwordLocator := th.page.Locator("#password")
		buttonLocator := th.page.Locator(".button")
		errorLocator := th.page.Locator(".error-message")

		// Wait for and fill the login form
		th.Require().NoError(t, usernameLocator.WaitFor(), "等待用户名输入")
		th.Require().NoError(t, usernameLocator.Fill("shiori"), "请填写用户名")
		th.Require().NoError(t, passwordLocator.Fill("wrong_password"), "请填写密码字段")

		// Click login and verify error
		th.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		errorText, err := errorLocator.TextContent()
		th.Require().NoError(t, err, "获取错误信息文本")
		th.Require().Contains(t, errorText, "用户名或密码不匹配")
	})

	t.Run("empty username validation", func(t *testing.T) {
		th, err := NewTestHelper(t, t.Name())
		require.NoError(t, err)
		defer th.Close()

		// Navigate to the login page
		_, err = th.page.Goto(baseURL)
		th.Require().NoError(t, err, "导航至基础 URL")

		// Get locators for form elements
		usernameLocator := th.page.Locator("#username")
		passwordLocator := th.page.Locator("#password")
		buttonLocator := th.page.Locator(".button")
		errorLocator := th.page.Locator(".error-message")

		// Wait for form and fill only password
		th.Require().NoError(t, usernameLocator.WaitFor(), "请填写用户名")
		th.Require().NoError(t, passwordLocator.Fill("gopher"), "请填写密码字段")

		// Click login and verify error
		th.Require().NoError(t, buttonLocator.Click(), "点击登录按钮")
		errorText, err := errorLocator.TextContent()
		th.Require().NoError(t, err, "获取错误信息文本")
		th.Require().Contains(t, errorText, "用户名不能为空")
	})
}
