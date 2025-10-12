package main

import "fmt"

//返回的是一个匿名函数，但是这个匿名函数引用到函数外的变量/参数，
//因此这个匿名函数就和这个变量构成一个整体，即闭包
//闭包特性1：记忆状态
func getSum() func(int) int {
	var sum int = 10
	return func(num int) int {
		sum = sum + num
		return sum
	}
}

//闭包特性2：数据封存
func bankAccount(initialBalance float64) (func(float64) bool, func() float64) {
	balance := initialBalance // 私有变量

	// 存款函数
	deposit := func(amount float64) bool {
		if amount > 0 {
			balance += amount
			return true
		}
		return false
	}

	// 查询余额函数
	getBalance := func() float64 {
		return balance
	}

	return deposit, getBalance
}

//应用：
type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewDatabaseConfig() func(...func(*DatabaseConfig)) *DatabaseConfig {
	config := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "password",
	}

	return func(opts ...func(*DatabaseConfig)) *DatabaseConfig {
		for _, opt := range opts {
			opt(config)
		}
		return config
	}
}
func (D *DatabaseConfig) String() string {
	str := fmt.Sprintf("%v,%v,%v,%v", D.Host, D.Port, D.Username, D.Password)
	return str
}

//闭包是指一个函数（通常是匿名函数）与其相关的引用环境组合而成的实体。简单来说，闭包是一个能够捕获并记住其创建时作用域中变量的函数。
//闭包本质依旧是一个匿名函数，只是这个函数引入外界的变量/参数
//闭包中使用的变量/参数会一直保存在内存中

func main() {
	// 外部变量
	message := "Hello"
	// 闭包：捕获了外部的 message 变量
	greeting := func() {
		fmt.Println(message) // 访问外部变量
	}
	greeting() // 输出: Hello
	// 修改外部变量，闭包能看到变化
	message = "World"
	greeting() // 输出: World

	//====================================================================

	f := getSum()
	fmt.Println(f(1))
	fmt.Println(f(2))
	fmt.Println(f(3))
	fmt.Println(getSum()(4)) //重新创建新的闭包

	fmt.Println("-------------------------")

	fmt.Println(getSum01(1))
	fmt.Println(getSum01(2))
	fmt.Println(getSum01(3))

	//====================================================================

	deposit, getBalance := bankAccount(1000)
	fmt.Println("初始余额:", getBalance()) // 输出: 初始余额: 1000
	deposit(500)
	fmt.Println("存款后余额:", getBalance()) // 输出: 存款后余额: 1500
	// 无法直接访问 balance 变量，实现了数据封装

	//===================================================================

	configurator := NewDatabaseConfig()
	config := configurator(
		func(c *DatabaseConfig) { c.Host = "192.168.1.100" },
		func(c *DatabaseConfig) { c.Port = 3306 },
	)

	fmt.Printf("数据库配置: %+v\n", *config)
	fmt.Printf("数据库配置: %+v\n", config) //直接调用String方法

}

//不使用闭包的时候，函数中的值不能反复使用

func getSum01(num int) int {
	var sum int = 0
	sum = sum + num
	return sum
}
