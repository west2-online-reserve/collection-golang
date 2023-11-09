package main

import (
	"fmt"
	"os"

	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func main() {
	// 启动 ChromeDriver /path/to/chromedriver
	wd, err := selenium.NewChromeDriverService("C:/Users/86181/AppData/Local/Google/Chrome/Application/chromedriver.exe", 5869)
	if err != nil {
		fmt.Println("无法启动 ChromeDriverfirst:", err)
		return
	}
	defer wd.Stop()

	// 设置 Chrome 浏览器选项
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 启动 Chrome 浏览器
	chromeCaps := chrome.Capabilities{
		Args: []string{
			// 这里可以添加 Chrome 浏览器的启动参数
		},
	}
	caps.AddChrome(chromeCaps)

	// 启动 Chrome 浏览器实例
	driver, err := selenium.NewRemote(caps, "http://localhost:5869/wd/hub")
	if err != nil {
		fmt.Println("无法启动 Chrome 浏览器second:", err)
		return
	}
	defer driver.Quit()

	// 打开网页
	err = driver.Get("https://www.bilibili.com/video/BV12341117rG/?p=2&spm_id_from=pageDriver&vd_source=6862116cc28d19480bad52bd3ab2308d")
	if err != nil {
		fmt.Println("无法打开网页：", err)
		return
	}

	// 等待一段时间,来进行手动登录
	time.Sleep(20 * time.Second)

	// 刷新页面
	err = driver.Refresh()
	if err != nil {
		fmt.Println("刷新页面失败：", err)
		return
	}

	fmt.Println("页面已加载")

	// 在这里可以继续编写其他操作，如查找元素、填写表单等
	// 获取页面初始高度
	js := "return action=document.body.scrollHeight"
	height, err := driver.ExecuteScript(js, nil)
	if err != nil {
		fmt.Println("获取页面高度失败：", err)
		return
	}
	initialHeight := int(height.(float64))

	// 初始化一个空的切片用于存储已经打印的文本
	printedText := make([]string, 0)
	for i := 0; i <= 30; i++ {
		// 将滚动条滚动到页面底部
		_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
		if err != nil {
			fmt.Println("滚动到页面底部失败：", err)
			return
		}

		// 等待加载更多评论
		time.Sleep(3 * time.Second)

		elements, err := driver.FindElements(selenium.ByXPATH, "//span[@class='reply-content-container root-reply']")
		if err != nil {
			fmt.Println("查找元素失败：", err)
			return
		}

		// 遍历获取到的元素
		for _, element := range elements {
			// 对每个元素进行操作
			// 获取元素的文本内容

			text, err := element.Text()
			text += "\n"
			if err != nil {
				fmt.Println("获取文本内容失败：", err)
				continue
			}

			// 检查文本是否已经打印过
			if contains(printedText, text) {
				continue
			}

			// 打开文件，以追加的方式写入数据
			file, err := os.OpenFile("comment.txt", os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("打开文件失败：", err)
				return
			}
			defer file.Close()

			// 将评论写入文件
			_, err = fmt.Fprintln(file, text)
			if err != nil {
				fmt.Println("写入文件失败：", err)
				return
			}

			// 将已打印的文本添加到列表中
			printedText = append(printedText, text)

		}

		// 判断滚动后页面高度是否发生变化
		js = "return action=document.body.scrollHeight"
		height, err = driver.ExecuteScript(js, nil)
		if err != nil {
			fmt.Println("获取页面高度失败：", err)
			return
		}
		newHeight := int(height.(float64))

		if newHeight == initialHeight {
			time.Sleep(3 * time.Second)
			height, _ = driver.ExecuteScript(js, nil)
			newHeight := int(height.(float64))
			if newHeight == initialHeight {
				continue
			} else {
				break
			}

		} else {
			fmt.Println("滚动条滚动后页面高度发生变化")
			initialHeight = newHeight
		}
	}

	// 获取网页响应代码
	//pageSource, err := driver.FindElements("By.XPath","span[@class="reply-content-container root-reply"]")

}
