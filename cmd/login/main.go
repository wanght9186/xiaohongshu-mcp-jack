package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/sirupsen/logrus"
	"github.com/xpzouying/xiaohongshu-mcp/browser"
	"github.com/xpzouying/xiaohongshu-mcp/cookies"
	"github.com/xpzouying/xiaohongshu-mcp/xiaohongshu"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("登录工具发生异常: %v", r)
			fmt.Println("登录工具发生异常，请重新运行。")
			os.Exit(1)
		}
	}()

	var (
		binPath string // 浏览器二进制文件路径
	)
	flag.StringVar(&binPath, "bin", "", "浏览器二进制文件路径")
	flag.Parse()

	// 登录的时候，需要界面，所以不能无头模式
	b := browser.NewBrowser(false, browser.WithBinPath(binPath))
	// 不自动关闭浏览器，让用户可以在登录后继续操作

	page := b.NewPage()
	// 不自动关闭页面，保持浏览器打开

	action := xiaohongshu.NewLogin(page)

	checkCtx, cancelCheck := context.WithTimeout(context.Background(), 45*time.Second)
	status, err := action.CheckLoginStatus(checkCtx)
	cancelCheck()
	if err != nil {
		logrus.Warnf("首次检查登录状态失败，将继续登录流程: %v", err)
		status = false
	}

	logrus.Infof("当前登录状态: %v", status)

	if status {
		logrus.Info("您已经登录，浏览器将保持打开以便您继续操作。")
		if err := saveCookies(page); err != nil {
			logrus.Fatalf("failed to save cookies: %v", err)
		}
		// 保持浏览器打开，等待用户按回车键后退出
		fmt.Println("\n========================================")
		fmt.Println("浏览器窗口将保持打开，您可以继续操作。")
		fmt.Println("操作完成后，按回车键退出程序并关闭浏览器。")
		fmt.Println("========================================")

		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')

		// 用户按回车后，清理资源
		page.Close()
		b.Close()
		return
	}

	// 开始登录流程
	logrus.Info("开始登录流程...")
	loginCtx, cancelLogin := context.WithTimeout(context.Background(), 8*time.Minute)
	if err = action.Login(loginCtx); err != nil {
		cancelLogin()
		logrus.Fatalf("登录失败: %v", err)
	} else {
		cancelLogin()
		if err := saveCookies(page); err != nil {
			logrus.Fatalf("failed to save cookies: %v", err)
		}
	}

	// 再次检查登录状态确认成功
	checkCtx2, cancelCheck2 := context.WithTimeout(context.Background(), 45*time.Second)
	status, err = action.CheckLoginStatus(checkCtx2)
	cancelCheck2()
	verifiedAfterLogin := true
	if err != nil {
		logrus.Warnf("登录后检查状态失败，请手动确认是否已登录: %v", err)
		verifiedAfterLogin = false
	}

	if verifiedAfterLogin && status {
		logrus.Info("登录成功！")
	} else if verifiedAfterLogin {
		logrus.Error("登录流程完成但仍未登录")
	} else {
		logrus.Warn("登录流程已完成，但无法自动验证登录状态，请在 MCP 中执行 check_login_status 再确认。")
	}

	// 保持浏览器打开，等待用户按回车键后退出
	fmt.Println("\n========================================")
	fmt.Println("浏览器窗口将保持打开，您可以继续操作。")
	fmt.Println("操作完成后，按回车键退出程序并关闭浏览器。")
	fmt.Println("========================================")

	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

	// 用户按回车后，清理资源
	page.Close()
	b.Close()
}

func saveCookies(page *rod.Page) error {
	cks, err := page.Browser().GetCookies()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cks)
	if err != nil {
		return err
	}

	cookieLoader := cookies.NewLoadCookie(cookies.GetCookiesFilePath())
	return cookieLoader.SaveCookies(data)
}
