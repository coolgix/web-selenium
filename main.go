package main

const (
	chromeDriver = "/Users/yanghaipeng/go/src/github.com/web-selenium/ChromeDriver/chromedriver"
	port         = 10001
)

func main() {
	/* 开启WebDriver服务 */
	// 设置Selenium的服务配置
	opts := []selenium.ServiceOption{
		// 开启Selenium的执行记录
		//selenium.Output(os.Stderr),
	}

}
