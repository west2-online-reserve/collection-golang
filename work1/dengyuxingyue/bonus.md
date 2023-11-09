##### Bonus

1. ###### 打印乘法表

    

   ```go
   package main
   
   import (
   	"fmt"
   	"os"
   )
   
   func main() {
   	//采取双重循环打印乘法表
   	file, err := os.Create("ninenine.txt")
   	if err != nil {
   		fmt.Println("无法创建文件:", err)
   		return
   	}
   
   	defer file.Close()
   
   	fmt.Fprintln(file, "九九乘法表：")
   
   	for i := 1; i <= 9; i++ {
   		for j := 1; j <= i; j++ {
   			fmt.Fprintf(file, "%d*%d=%2d ", j, i, (i * j))
   		}
   		fmt.Fprintln(file)
   	}
   
   }
   
   ```

      （1）如果已经创建文件，creat函数会将文件截断为空文件，删除文件原先的所有内容。

      （2）Fprintf与Printf函数类似c语言，可以方便的控制格式输出。

   

2. Go语言中的切片和数组的区别有哪些？

   

   1. 长度固定性：数组的长度是固定的，声明时需要指定长度；而切片是动态的，长度可以根据实际情况自动扩展或缩减。

   2. 内存分配：数组是一个连续的内存块，元素的存储位置是连续的；而切片可以理解为数组的引用，包括指向底层数组的指针，容量和长度。切片在进行扩容时，底层数组的内存会发生改变但依然保持连续性，而切片本身的地址不会发生改变。多个切片可能会共用一个底层数组，修改其中一个切片的值会引起其它切片值得改变，验证如下。

      ```go
      package main
      
      import "fmt"
      
      func main() {
      	arr := [5]int{1, 2, 3, 4, 5}
      	slice1 := arr[:]
      	slice2 := arr[1:3]
          //引用arr得到了slice1,slice2
          
      	fmt.Println(slice1)//1，2，3，4，5
      	fmt.Println(slice2)//2，3
      
      	//对切片2进行修改
      	slice2[0] = 6
      	fmt.Println(slice1)//1，6，3，4，5
      	fmt.Println(slice2)//6，3
      	fmt.Println(arr)//1，6，3，4，5
      
      }
      ```

      

   3. 传递方式：数组在函数之间的传递是值拷贝，会复制整个数组；切片在函数之间传递时是引用传递，只会复制切片的指针、长度和容量。所以如果向函数里直接以slice []int的形式或者arr [5]int的形式传递时，函数会改变原本切片的值而数组就不会改变

       
   
      ```go
      package main
      
      import "fmt"
      
      func main() {
      
      	a := [5]int{1, 2, 3, 4, 5}
      	b := []int{1, 2, 3, 4, 5}
      
      	fmt.Println(a) //1,2,3,4,5
      	fmt.Println(b) //1,2,3,4,5
      
      	modifyArray(a)
      	modifySlice(b)
      
      	fmt.Println(a) //1,2,3,4,5
      	fmt.Println(b) //2,4,6,8,10
      
      }
      
      func modifySlice(slice []int) {
      	// 修改切片的元素
      	for i := 0; i < len(slice); i++ {
      		slice[i] = slice[i] * 2
      	}
      }
      func modifyArray(arr [5]int) {
      	// 修改数组的元素
      	for i := 0; i < len(arr); i++ {
      		arr[i] = arr[i] * 2
      	}
      }
      
      ```

      

   4. 动态性：切片支持动态增加和删除元素的操作，可以通过append()函数来添加元素；数组的长度是固定的，无法动态改变。

      

   

   Go中创建切片有几种方式？

      Go 语言中包含三种初始化切片的方式：
   
   1. 通过下标的方式获得数组或者切片的一部分；

      ```go
          arr := [5]int{1, 2, 3, 4, 5}
      	slice1 := slice1[:]
      	slice2 := arr[1:3]
      ```

      
   
   2. 使用字面量初始化新的切片；
   
      ```go
      s := []int{1, 2, 3, 4, 5}
      ```
   
      
   
   3. 使用关键字 make 创建切片：
   
      ```
      s := make([]int, 5, 10)//指定长度与容量
      s:=make([]int,5)	   //长度与容量相等
      ```
   
      
   
   创建map呢？
   
   1. 通过字面量
   
   ```go
   m := map[int]int{
   		1: 10,
   		2: 20,
   		3: 30,
   	}
   	fmt.Println(m[1])
   ```
   
   ​	2.通过关键字
   
   ```go
   // 声明一个map变量
   var m map[keyType]valueType
   m = make(map[keyType]valueType)
   
   
   // 声明并初始化
   m := make(map[keyType]valueType)
   ```
   
   
   
   ##### 数组元素之妈妈找儿子
   
   这个代码长度有点超出我的预期，因为我不知道怎样在go里简洁优雅地输入一个不知道长度的数组，实在是不想用事先知道个数的方法或者输入一个奇怪的数字来终止输入，便写了段代码来处理数组的输入。
   
   如果哈希查找的复杂度确实为1，那么该代码的算法复杂度为O(n)
   
   ```go
   package main
   
   import (
   	"bufio"
   	"fmt"
   	"os"
   	"strconv"
   )
   
   
   
   //使用哈希表来实现查找
   func twoSum(nums []int, target int) []int {
   	numMap := make(map[int]int)
   
   	for i, num := range nums {
   		diff := target - num
   		if j, ok := numMap[diff]; ok && j != i {
   			return []int{j, i}
   		}
   		numMap[num] = i
   	}
   
   	return []int{}
   }
   
   func main() {
       
       
       //以下代码均为处理数据输入的代码，实现输入数字，空格区分，换行终止，得到数组的功能
       
       var numint []int //我们的目标数组
   	var str string = ""//过度字符串
   
   	fmt.Println("请输入你想要的整数数组,换行代表结束")
   
   	reader := bufio.NewReader(os.Stdin)
    
   	nums, _, err := reader.ReadLine()
   
   	if nil != err {
   		fmt.Println("reader.ReadLine() error:", err)
   	}
   
   	for i := 0; i < len(nums); i++ {
   
   		if int(nums[i]) != 32 {
   			str += string(nums[i])
   		}
   		if int(nums[i]) == 32 || i == len(nums)-1 {
   
   			num, err := strconv.Atoi(str)
   			if err == nil {
   				fmt.Println("正在加载数据：", num)
   			} else {
   				fmt.Println("error")
   			}
   			numint = append(numint, num)
   			str = ""
   		}
   	}
       //输入功能结束
       
       
   	fmt.Println("加载完毕，请输入你想要的目标值")
   	var target int
   	fmt.Scan(&target)
   
   	fmt.Println(twoSum(numint, target))
   
   }
   
   ```
   
   ##### 通过通道并发筛选素数
   
   1）代码实现输出素数的功能
   
   2）代码利用了 Go 语言的并发特性，通过使用通道（channel）实现了协程间的通信和并发计算。
   
   3）相较于普通的实现方式，这段代码确实有性能上的提升。通过并发的方式同时进行素数的生成和筛选，能够加快求解速度。在筛选非常大的素数集合时，这样的并发计算会更加高效。
   
   ```go
   package main
   
   import (
   	"fmt"
   )
   
   
   //向管道ch里入队待筛选的数字
   func generate(ch chan int) { 
   	for i := 2; ; i++ {
   		ch <- i
   	}
   }
   
   
   //利用素数筛选法，in是上一次筛选出来的通道，out是新的筛选后的管道，prime是当前用来筛选的素数
   func filter(in chan int, out chan int, prime int) {
   
   	for {
   
   		num := <-in
   		if num%prime != 0 {
   			out <- num
   		}
   	}
   }
   
   func main() {
   	ch := make(chan int)
       
       //新建协程调用generate产生数字
   	go generate(ch)
       
       //for循环执行次数即为输出的素数个数
   	for i := 0; i < 100; i++ {
   		prime := <-ch
   		fmt.Printf("prime:%d\n", prime)
   		out := make(chan int)
   		go filter(ch, out, prime)
   		ch = out
   	}
   }
   ```
   
   