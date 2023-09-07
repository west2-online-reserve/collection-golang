package main

import (
	"fmt"
)

func generate(ch chan int) { //����ԭʼ�������������͸�ͨ��ch
	for i := 2; ; i++ {
		ch <- i
	}
}

func filter(in chan int, out chan int, prime int) { //����ɸ��ͨ������inͨ�������ݣ������ܷ����������ж��Ƿ�Ϊ���������������͸�ͨ��out
	for {
		num := <-in
		if num%prime != 0 {
			out <- num
		}
	}
}

func main() {
	ch := make(chan int)     //����ͨ��ch
	go generate(ch)          //����Э�̲���ԭʼ����
	for i := 0; i < 6; i++ { //����ѭ������������ɸѡ����������
		prime := <-ch //��������
		fmt.Printf("prime:%d\n", prime)
		out := make(chan int)
		go filter(ch, out, prime) //����Э�̽�������ɸ
		ch = out                  //filter��ɸ��������ֵ��ͨ��ch���Խ�����һ�׶�����ɸ
	}
}

//1.����ʵ��������ɸ��ͨ������main������forѭ������������������������
//2.����������goԭ��֧�ֲ��������ԣ�����channelʵ�ֲ�����������ɸ
//3.ͨ�������������������ٶ�
