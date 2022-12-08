package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

// 设置常量
const (
	// ChromeDriver路径，存放在E:\mygo，因此使用相对路径表示
	chromeDriver = "./Chromedriver/chromedriver"
	// ChromeDriver运行端口
	port = 8080
)

// 定义结构体，用于存储数据
type Job struct {
	Name   string `json:"name"`
	Area   string `json:"area"`
	Pays   string `json:"pays"`
	Exp    string `json:"exp"`
	Tags   string `json:"tags"`
	Desc   string `json:"desc"`
	Publis string `json:"publis"`
	Cmp    string `json:"cmp"`
	Scale  string `json:"scale"`
}

// 创建浏览器对象
func get_wd() (selenium.WebDriver, *selenium.Service) {
	// 开启Selenium服务
	s, _ := selenium.NewChromeDriverService(chromeDriver, port)
	/* 连接WebDriver服务 */
	// 设置浏览器功能
	caps := selenium.Capabilities{}
	// 设置chrome特定功能
	chromeCaps := chrome.Capabilities{
		// 使用开发者调试模式
		ExcludeSwitches: []string{"enable-automation"},
		// 基本功能
		Args: []string{
			"--no-sandbox",
			// 设置请求头
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; " +
				"x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
				"Chrome/94.0.4606.61 Safari/537.36",
		},
	}
	// 将谷歌浏览器特定功能chromeCaps添加到caps
	caps.AddChrome(chromeCaps)
	// 根据浏览器功能连接
	urlPrefix := fmt.Sprintf("http://localhost:%d/wd/hub", port)
	wd, _ := selenium.NewRemote(caps, urlPrefix)
	return wd, s
}

// 获取当前页数的所有职位信息
func get_jobs(wd selenium.WebDriver) []Job {
	var jobs []Job
	jf, _ := wd.FindElements(selenium.ByClassName, "search-job-result")
	for _, v := range jf {
		j := Job{}
		// 获取职位名称
		name, _ := v.FindElement(selenium.ByClassName, "job-name")
		j.Name, _ = name.Text()
		// 获取工作地点
		area, _ := v.FindElement(selenium.ByClassName, "job-area")
		j.Area, _ = area.Text()
		// 获取薪资
		pays, _ := v.FindElement(selenium.ByClassName, "salary")
		j.Pays, _ = pays.Text()
		// 获取经验学历
		exp, _ := v.FindElement(selenium.ByClassName, "tag-list")
		j.Exp, _ = exp.Text()
		// // 获取职位标签
		tags, _ := v.FindElement(selenium.ByClassName, "job-card-footer")
		j.Tags, _ = tags.Text()
		// // 获取公司福利
		desc, _ := v.FindElement(selenium.ByClassName, "info-desc")
		j.Desc, _ = desc.Text()
		// 获取公司人事信息
		publis, _ := v.FindElement(selenium.ByClassName, "info-public")
		j.Publis, _ = publis.Text()
		// 获取公司名称
		cmp, _ := v.FindElement(selenium.ByClassName, "company-name")
		j.Cmp, _ = cmp.Text()
		// 获取公司行业和规模
		scale, _ := v.FindElement(selenium.ByClassName, "company-tag-list")
		j.Scale, _ = scale.Text()

		jobs = append(jobs, j)
	}
	return jobs

}

// 保存数据
func save_data(jobs []Job) {
	// 将变量jobs的数据写入JSON文件

	f2, _ := os.OpenFile("output.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	encoder := json.NewEncoder(f2)
	err := encoder.Encode(jobs)
	// 如果err不为空值nil，则说明写入错误
	if err != nil {
		fmt.Printf("JSON写入失败：%v\n", err.Error())
	} else {
		fmt.Printf("JSON写入成功\n")
	}
}

func main() {
	// 获取浏览器对象
	wd, s := get_wd()
	// 关闭服务
	defer s.Stop()
	// 关闭浏览器对象
	defer wd.Quit()
	// 访问网址
	wd.Get("https://www.zhipin.com/")
	// 最大化窗口
	wd.MaximizeWindow("")
	time.Sleep(5 * time.Second)
	// 输入查询职位
	query, _ := wd.FindElement(selenium.ByName, "query")
	query.SendKeys("go语言")
	time.Sleep(2 * time.Second)
	// 点击搜索按钮
	search, _ := wd.FindElement(selenium.ByCSSSelector, `[class="btn btn-search"]`)
	search.Click()
	time.Sleep(2 * time.Second)

	// 获取第一页的职位信息
	jobs := get_jobs(wd)
	// 使用死循环实现翻页
	for {
		np, err := wd.FindElement(selenium.ByCSSSelector, `[class="page"]>[class="next"]`)
		// err不等于nil说明无法点击下一页，终止死循环
		if err != nil {
			break
		} else {
			// 点击下一页按钮
			np.Click()
			time.Sleep(2 * time.Second)
			// 获取当前页的职位信息
			// 将当前页所有职业合并到切片jobs
			jobs = append(jobs, get_jobs(wd)...)
		}
	}
	// 保存数据
	save_data(jobs)
}
