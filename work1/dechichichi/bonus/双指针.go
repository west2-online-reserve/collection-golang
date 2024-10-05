package main

func main() {
	li := make([]int, 0)
	var target int
	//一系列输入
	front := 0
	rear := len(li) - 1
	for {
		if li[front]+li[rear] == target {
			//输出
			break
		}
		if li[front]+li[rear] < target {
			front++
		}
		if li[front]+li[rear] > target {
			rear--
		}
	}
}
