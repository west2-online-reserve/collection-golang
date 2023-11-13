package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	//创建文件
	file, err := os.Create("comment.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
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

	// 初始化一个空的map用于存储已经打印的文本
	printedText := make(map[string]bool)

	for {
		// 将滚动条滚动到页面底部
		_, err = driver.ExecuteScript("window.scrollTo(0, document.body.scrollHeight)", nil)
		if err != nil {
			fmt.Println("滚动到页面底部失败：", err)
			return
		}

		// 等待加载更多评论
		time.Sleep(3 * time.Second)

		elements, err := driver.FindElements(selenium.ByXPATH, "//div[@class='reply-item']")
		if err != nil {
			fmt.Println("查找元素失败：", err)
			return
		}

		// 遍历获取到的元素
		for _, element := range elements {
			// 对每个元素进行操作
			// 首先获取到主评论elementfirst,并添加到txt中								".//span[@class='reply-content-container root-reply']"
			firstUsername, _ := element.FindElement(selenium.ByXPATH, ".//div[@class='user-name']")
			Pritext, err := firstUsername.Text()
			elementfirst, _ := element.FindElement(selenium.ByXPATH, ".//span[@class='reply-content-container root-reply']")
			s, _ := elementfirst.Text()
			Pritext = Pritext + ": " + s + "\n"
			Sumtext := Pritext
			if err != nil {
				fmt.Println("获取文本内容失败：", err)
				continue
			}
			// 检查文本是否已经打印过
			if printedText[Pritext] {
				continue
			}

			//获取当前节点的子评论
			/*
				1. 首先尝试找到（点击查看）按钮，如果没找到就继续
				2.如果找到，点击按钮，第一次添加当前文本到Sumtext
			*/
			button, err := element.FindElement(selenium.ByXPATH, ".//span[@class='view-more-btn']")
			if err == nil && button != nil {
				fmt.Println("ooooooooooooooooooooooooooooooooo")
				driver.ExecuteScript("arguments[0].scrollIntoView(true);", []interface{}{button})
				button.Click()
				time.Sleep(1 * time.Second)
				//点击之后会更新子评论

				//当页子评论的集合
				subreplyitems, _ := element.FindElements(selenium.ByXPATH, ".//div[@class='sub-reply-item']")

				for _, subreplyitem := range subreplyitems {

					replyuser, err := subreplyitem.FindElement(selenium.ByXPATH, ".//div[@class='sub-user-name']")
					if err != nil {
						continue
					}
					replycontent, err := subreplyitem.FindElement(selenium.ByXPATH, ".//span[@class='reply-content']")
					if err != nil {
						continue
					}
					username, _ := replyuser.Text()
					srr, _ := replycontent.Text()
					Sumtext = Sumtext + "\n" + username + ": " + srr
				}
				Sumtext += "\n"

			}

			// 打开文件，以追加的方式写入数据
			file, err := os.OpenFile("comment.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Println("打开文件失败：", err)
				return
			}
			defer file.Close()
			// 将评论写入文件
			_, err = fmt.Fprintln(file, Sumtext)
			if err != nil {
				fmt.Println("写入文件失败：", err)
				return
			}

			// 将已打印的文本添加到map中
			printedText[Pritext] = true
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

			} else {
				continue
			}

		} else {
			fmt.Println("滚动条滚动后页面高度发生变化")
			initialHeight = newHeight
		}
	}

	// 获取网页响应代码
	//pageSource, err := driver.FindElements("By.XPath","span[@class="reply-content-container root-reply"]")

}
