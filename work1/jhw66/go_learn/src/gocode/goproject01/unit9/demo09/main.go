// Go语言接口的主要作用
// 接口在Go语言中扮演着至关重要的角色，它们提供了抽象、多态和解耦的能力。

// // 1. 多态性（Polymorphism）
// // 作用：允许不同的类型以统一的方式被处理
package main

import "fmt"

// PaymentMethod 接口定义支付方式
type PaymentMethod interface {
	ProcessPayment(amount float64) bool
	GetName() string
}

// 信用卡支付实现
type CreditCard struct {
	CardNumber string
	HolderName string
}

func (c CreditCard) ProcessPayment(amount float64) bool {
	fmt.Printf("Processing credit card payment of $%.2f for %s\n", amount, c.HolderName)
	// 模拟支付处理逻辑
	return true
}

func (c CreditCard) GetName() string {
	return "Credit Card"
}

// PayPal支付实现
type PayPal struct {
	Email string
}

func (p PayPal) ProcessPayment(amount float64) bool {
	fmt.Printf("Processing PayPal payment of $%.2f for %s\n", amount, p.Email)
	// 模拟支付处理逻辑
	return true
}

func (p PayPal) GetName() string {
	return "PayPal"
}

// 加密货币支付实现
type CryptoCurrency struct {
	WalletAddress string
	CoinType      string
}

func (c CryptoCurrency) ProcessPayment(amount float64) bool {
	fmt.Printf("Processing %s payment of $%.2f to %s\n", c.CoinType, amount, c.WalletAddress)
	return true
}

func (c CryptoCurrency) GetName() string {
	return c.CoinType + " Crypto"
}

// 统一的支付处理函数
func ProcessOrder(paymentMethod PaymentMethod, amount float64) {
	fmt.Printf("Using %s for payment...\n", paymentMethod.GetName())
	success := paymentMethod.ProcessPayment(amount)
	if success {
		fmt.Println("Payment successful!")
	} else {
		fmt.Println("Payment failed!")
	}
	fmt.Println("---")
}

func main() {
	// 创建不同的支付方式
	creditCard := CreditCard{"1234-5678-9012-3456", "John Doe"}
	paypal := PayPal{"john@example.com"}
	bitcoin := CryptoCurrency{"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "Bitcoin"}

	// 统一处理不同的支付方式
	paymentMethods := []PaymentMethod{creditCard, paypal, bitcoin}
	amount := 99.99

	for _, method := range paymentMethods {
		ProcessOrder(method, amount)
	}
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
// // 2. 抽象与解耦（Abstraction & Decoupling）
// // 作用：隐藏具体实现细节，降低模块间的耦合度
// package main

// import "fmt"

// // Storage 接口定义数据存储操作
// type Storage interface {
// 	Save(data string) error
// 	Load(id string) (string, error)
// 	Delete(id string) error
// }

// // 文件存储实现
// type FileStorage struct {
// 	FilePath string
// }

// func (f FileStorage) Save(data string) error {
// 	fmt.Printf("Saving data to file: %s\n", f.FilePath)
// 	// 实际的文件保存逻辑
// 	return nil
// }

// func (f FileStorage) Load(id string) (string, error) {
// 	fmt.Printf("Loading data from file: %s, ID: %s\n", f.FilePath, id)
// 	return "data from file", nil
// }

// func (f FileStorage) Delete(id string) error {
// 	fmt.Printf("Deleting data from file: %s, ID: %s\n", f.FilePath, id)
// 	return nil
// }

// // 数据库存储实现
// type DatabaseStorage struct {
// 	ConnectionString string
// 	TableName        string
// }

// func (d DatabaseStorage) Save(data string) error {
// 	fmt.Printf("Saving data to database: %s, table: %s\n", d.ConnectionString, d.TableName)
// 	// 实际的数据库操作
// 	return nil
// }

// func (d DatabaseStorage) Load(id string) (string, error) {
// 	fmt.Printf("Loading data from database: %s, ID: %s\n", d.ConnectionString, id)
// 	return "data from database", nil
// }

// func (d DatabaseStorage) Delete(id string) error {
// 	fmt.Printf("Deleting data from database: %s, ID: %s\n", d.ConnectionString, id)
// 	return nil
// }

// // 业务逻辑层 - 不关心具体存储实现
// type UserService struct {
// 	storage Storage
// }

// func NewUserService(storage Storage) *UserService {
// 	return &UserService{storage: storage}
// }

// func (s *UserService) CreateUser(userData string) error {
// 	fmt.Println("Creating user...")
// 	return s.storage.Save(userData)
// }

// func (s *UserService) GetUser(userID string) (string, error) {
// 	fmt.Println("Getting user...")
// 	return s.storage.Load(userID)
// }

// func (s *UserService) DeleteUser(userID string) error {
// 	fmt.Println("Deleting user...")
// 	return s.storage.Delete(userID)
// }

// func main() {
// 	// 使用文件存储
// 	fileStorage := FileStorage{FilePath: "/data/users.txt"}
// 	userServiceWithFile := NewUserService(fileStorage)

// 	userServiceWithFile.CreateUser("John Doe, 30")
// 	userServiceWithFile.GetUser("user123")
//     userServiceWithFile.DeleteUser("user123")

// 	fmt.Println("=== Switching to database storage ===")

// 	// 切换到数据库存储 - 业务逻辑无需修改
// 	dbStorage := DatabaseStorage{
// 		ConnectionString: "localhost:5432/mydb",
// 		TableName:        "users",
// 	}
// 	userServiceWithDB := NewUserService(dbStorage)

// 	userServiceWithDB.CreateUser("Jane Smith, 25")
// 	userServiceWithDB.GetUser("user456")
//     userServiceWithDB.DeleteUser("user456")
// }

//
//
//
//
//
//
//
//
//
//
//
//
//
//
// // 3. 测试友好（Testing Friendly）
// // 作用：便于编写单元测试，可以使用mock实现
// package main

// import (
// 	"fmt"
// 	"testing"
// )

// // EmailSender 接口
// type EmailSender interface {
// 	SendEmail(to, subject, body string) error
// }

// // 真实的邮件发送服务
// type SMTPEmailSender struct {
// 	Server   string
// 	Port     int
// 	Username string
// 	Password string
// }

// func (s SMTPEmailSender) SendEmail(to, subject, body string) error {
// 	fmt.Printf("Sending email via SMTP to: %s, Subject: %s\n", to, subject)
// 	// 实际的SMTP发送逻辑
// 	return nil
// }

// // Mock邮件发送器 - 用于测试
// type MockEmailSender struct {
// 	SentEmails []Email
// }

// type Email struct {
// 	To      string
// 	Subject string
// 	Body    string
// }

// func (m *MockEmailSender) SendEmail(to, subject, body string) error {
// 	email := Email{To: to, Subject: subject, Body: body}
// 	m.SentEmails = append(m.SentEmails, email)
// 	fmt.Printf("Mock: Would send email to: %s, Subject: %s\n", to, subject)
// 	return nil
// }

// // 通知服务
// type NotificationService struct {
// 	emailSender EmailSender
// }

// func NewNotificationService(emailSender EmailSender) *NotificationService {
// 	return &NotificationService{emailSender: emailSender}
// }

// func (n *NotificationService) SendWelcomeEmail(userEmail, userName string) error {
// 	subject := "Welcome to our service!"
// 	body := fmt.Sprintf("Hello %s, welcome to our platform!", userName)
// 	return n.emailSender.SendEmail(userEmail, subject, body)
// }

// func (n *NotificationService) SendPasswordResetEmail(userEmail, resetToken string) error {
// 	subject := "Password Reset Request"
// 	body := fmt.Sprintf("Use this token to reset your password: %s", resetToken)
// 	return n.emailSender.SendEmail(userEmail, subject, body)
// }

// // 单元测试
// func TestNotificationService(t *testing.T) {
// 	// 使用mock进行测试，不依赖真实的SMTP服务
// 	mockSender := &MockEmailSender{}
// 	service := NewNotificationService(mockSender)

// 	// 测试发送欢迎邮件
// 	err := service.SendWelcomeEmail("test@example.com", "Test User")
// 	if err != nil {
// 		t.Errorf("Failed to send welcome email: %v", err)
// 	}

// 	// 验证邮件是否"发送"
// 	if len(mockSender.SentEmails) != 1 {
// 		t.Errorf("Expected 1 email, got %d", len(mockSender.SentEmails))
// 	}

// 	if mockSender.SentEmails[0].To != "test@example.com" {
// 		t.Errorf("Expected recipient test@example.com, got %s", mockSender.SentEmails[0].To)
// 	}
// }

// func main() {
// 	// 生产环境使用真实的SMTP发送器
// 	smtpSender := SMTPEmailSender{
// 		Server:   "smtp.example.com",
// 		Port:     587,
// 		Username: "user",
// 		Password: "pass",
// 	}

// 	notificationService := NewNotificationService(smtpSender)
// 	notificationService.SendWelcomeEmail("user@example.com", "John Doe")

// 	// 运行测试
// 	fmt.Println("\n=== Running Tests ===")
// 	testing.Main(func(pat, str string) (bool, error) { return true, nil },
// 		[]testing.InternalTest{
// 			{"TestNotificationService", TestNotificationService},
// 		},
// 		nil, nil)
// }

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// // 4. 扩展性（Extensibility）
// // 作用：易于扩展新功能，无需修改现有代码
// package main

// import "fmt"

// // Shape 接口定义几何形状
// type Shape interface {
// 	Area() float64
// 	Perimeter() float64
// 	Name() string
// }

// // 现有的形状实现
// type Rectangle struct {
// 	Width, Height float64
// }

// func (r Rectangle) Area() float64 {
// 	return r.Width * r.Height
// }

// func (r Rectangle) Perimeter() float64 {
// 	return 2 * (r.Width + r.Height)
// }

// func (r Rectangle) Name() string {
// 	return "Rectangle"
// }

// type Circle struct {
// 	Radius float64
// }

// func (c Circle) Area() float64 {
// 	return 3.14159 * c.Radius * c.Radius
// }

// func (c Circle) Perimeter() float64 {
// 	return 2 * 3.14159 * c.Radius
// }

// func (c Circle) Name() string {
// 	return "Circle"
// }

// // 新增的形状 - 无需修改现有代码
// type Triangle struct {
// 	Base, Height, SideA, SideB float64
// }

// func (t Triangle) Area() float64 {
// 	return 0.5 * t.Base * t.Height
// }

// func (t Triangle) Perimeter() float64 {
// 	return t.Base + t.SideA + t.SideB
// }

// func (t Triangle) Name() string {
// 	return "Triangle"
// }

// // 形状处理器 - 可以处理任何实现了Shape接口的类型
// type ShapeProcessor struct{}

// func (sp ShapeProcessor) ProcessShapes(shapes []Shape) {
// 	totalArea := 0.0
// 	totalPerimeter := 0.0

// 	for _, shape := range shapes {
// 		area := shape.Area()
// 		perimeter := shape.Perimeter()

// 		fmt.Printf("%s - Area: %.2f, Perimeter: %.2f\n",
// 			shape.Name(), area, perimeter)

// 		totalArea += area
// 		totalPerimeter += perimeter
// 	}

// 	fmt.Printf("Total - Area: %.2f, Perimeter: %.2f\n", totalArea, totalPerimeter)
// }

// // 新增的功能 - 形状渲染器
// type ShapeRenderer interface {
// 	Render() string
// }

// // 为现有形状添加新功能
// func (r Rectangle) Render() string {
// 	return fmt.Sprintf("📐 Rectangle %vx%v", r.Width, r.Height)
// }

// func (c Circle) Render() string {
// 	return fmt.Sprintf("⭕ Circle radius %v", c.Radius)
// }

// func (t Triangle) Render() string {
// 	return fmt.Sprintf("🔺 Triangle base %v height %v", t.Base, t.Height)
// }

// func RenderShapes(shapes []ShapeRenderer) {
// 	fmt.Println("\nRendering shapes:")
// 	for _, shape := range shapes {
// 		fmt.Println(shape.Render())
// 	}
// }

// func main() {
// 	processor := ShapeProcessor{}

// 	// 初始的形状集合
// 	shapes := []Shape{
// 		Rectangle{Width: 5, Height: 3},
// 		Circle{Radius: 4},
// 	}

// 	fmt.Println("Initial shapes:")
// 	processor.ProcessShapes(shapes)

// 	// 添加新形状 - 无需修改processor
// 	fmt.Println("\nAfter adding new shapes:")
// 	shapes = append(shapes, Triangle{Base: 6, Height: 4, SideA: 5, SideB: 5})
// 	processor.ProcessShapes(shapes)

// 	// 使用新的渲染功能
// 	renderableShapes := []ShapeRenderer{
// 		Rectangle{Width: 5, Height: 3},
// 		Circle{Radius: 4},
// 		Triangle{Base: 6, Height: 4, SideA: 5, SideB: 5},
// 	}
// 	RenderShapes(renderableShapes)
// }

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
// // 5. 插件架构（Plugin Architecture）
// // 作用：支持动态加载和替换组件
// package main

// import "fmt"

// // Processor 插件接口
// type Processor interface {
// 	Process(data string) string
// 	Name() string
// }

// // 文本大写处理器
// type UpperCaseProcessor struct{}

// func (u UpperCaseProcessor) Process(data string) string {
// 	// 模拟处理逻辑
// 	return "UPPERCASE: " + data
// }

// func (u UpperCaseProcessor) Name() string {
// 	return "UpperCaseProcessor"
// }

// // 文本反转处理器
// type ReverseProcessor struct{}

// func (r ReverseProcessor) Process(data string) string {
// 	// 反转字符串
// 	runes := []rune(data)
// 	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
// 		runes[i], runes[j] = runes[j], runes[i]
// 	}
// 	return "REVERSED: " + string(runes)
// }

// func (r ReverseProcessor) Name() string {
// 	return "ReverseProcessor"
// }

// // 加密处理器
// type EncryptProcessor struct{}

// func (e EncryptProcessor) Process(data string) string {
// 	// 简单加密示例
// 	encrypted := ""
// 	for _, char := range data {
// 		encrypted += string(char + 3) // 简单的字符偏移加密
// 	}
// 	return "ENCRYPTED: " + encrypted
// }

// func (e EncryptProcessor) Name() string {
// 	return "EncryptProcessor"
// }

// // 插件管理器
// type PluginManager struct {
// 	processors map[string]Processor
// }

// func NewPluginManager() *PluginManager {
// 	return &PluginManager{
// 		processors: make(map[string]Processor),
// 	}
// }

// func (pm *PluginManager) Register(processor Processor) {
// 	pm.processors[processor.Name()] = processor
// 	fmt.Printf("Registered plugin: %s\n", processor.Name())
// }

// func (pm *PluginManager) Unregister(name string) {
// 	delete(pm.processors, name)
// 	fmt.Printf("Unregistered plugin: %s\n", name)
// }

// func (pm *PluginManager) ProcessData(processorName, data string) (string, error) {
// 	processor, exists := pm.processors[processorName]
// 	if !exists {
// 		return "", fmt.Errorf("processor %s not found", processorName)
// 	}
// 	return processor.Process(data), nil
// }

// func (pm *PluginManager) ListPlugins() []string {
// 	var names []string
// 	for name := range pm.processors {
// 		names = append(names, name)
// 	}
// 	return names
// }

// func main() {
// 	pluginManager := NewPluginManager()

// 	// 注册插件
// 	pluginManager.Register(UpperCaseProcessor{})
// 	pluginManager.Register(ReverseProcessor{})
// 	pluginManager.Register(EncryptProcessor{})

// 	fmt.Printf("Available plugins: %v\n\n", pluginManager.ListPlugins())

// 	// 使用不同的处理器处理数据
// 	testData := "Hello, World!"

// 	result, _ := pluginManager.ProcessData("UpperCaseProcessor", testData)
// 	fmt.Println(result)

// 	result, _ = pluginManager.ProcessData("ReverseProcessor", testData)
// 	fmt.Println(result)

// 	result, _ = pluginManager.ProcessData("EncryptProcessor", testData)
// 	fmt.Println(result)

// 	// 动态卸载和加载插件
// 	fmt.Println("\n--- Dynamic plugin management ---")
// 	pluginManager.Unregister("ReverseProcessor")

// 	fmt.Printf("Available plugins after removal: %v\n", pluginManager.ListPlugins())

// 	// 尝试使用已卸载的插件
// 	_, err := pluginManager.ProcessData("ReverseProcessor", testData)
// 	if err != nil {
// 		fmt.Printf("Error: %v\n", err)
// 	}
// }
