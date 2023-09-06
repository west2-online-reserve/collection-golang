package main

import (
	"fmt"
	"math/rand"
)

// 生成一段随机的整型序列
func generateIntSequence(min, max int, seq []int, size int) {
	rand.Seed(0)
	for i := 0; i < size; i++ {
		seq[i] = rand.Intn(max-min) + min
	}
}

func main() {
	value_map := make(map[int]([]int))
	seq := make([]int, 25)
	go generateIntSequence(0, 25, seq, cap(seq))
	except_num := rand.Intn(25)
	fmt.Printf("期望和值:%v\n已生成序列:%v", except_num, seq)
	for i := 0; i < len(seq); i++ {
		value_map[seq[i]] = append(value_map[seq[i]], i+1)
	}
	//fmt.Print(value_map)	//test
	for i := 0; i <= except_num; i-- {
		if except_num == i*2 && len(value_map[i]) > 1 {
			fmt.Print("\n[", value_map[i][0], ",", value_map[i][1], "]")
			fmt.Printf("\nReason:seq[%v]=%v,seq[%v]=%v\n", value_map[i][0], i, value_map[i][1], i)
			return
		} else if len(value_map[except_num-i]) > 0 && len(value_map[i]) > 0 {
			fmt.Print("\n[", value_map[except_num-i][0], ",", value_map[i][0], "]")
			fmt.Printf("\nReason:seq[%v]=%v,seq[%v]=%v\n", value_map[except_num-i][0], except_num-i, value_map[i][0], i)
			return
		}
	}
}

//以上即为复杂度为O(n)的算法
//思路如下：
//对于任意的序列seq，不妨令任意元素的下标为index，值为value
//则index与value的关系为seq[index]=value，即通过序号查找其值
//在物理结构上为 index(n)=value(m)，即index映射到value
//如果有一个物理结构可以使得f(value)=index，那么便可以通过f(except_number-value)直接读取另一个序号
//简单来说就是 哈希表(HashMap)
//不过序列里的元素是可以重复的，那么数据存储的结构应当为Item为Slice(切片)的Map(集合)
//这样，一个value就可以映射到多个index，其形式为Slice(切片)
//对于except_number==value*2，仅需输出HashMap[value][0],HashMap[value][1]
//对于except_number!=value*2，仅需输出HashMap[value][0],HashMap[except_number-value][0]
//对于其他情况，则不输出
//变化后的数据存储方式如下：
//						HashMap[0]=nil
//						HashMap[1]=[0,1]
//seq=[1,1,3,2,5]  <=>	HashMap[2]=[3]
//						HashMap[3]=[2]
//						HashMap[5]=[4]
