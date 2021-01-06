package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix/rolling"
	"time"
)

func main() {
	n := rolling.NewNumber()
	for i := 1; i <= 30; i++ {
		n.Increment(float64(i))
		fmt.Printf("sliding window sum value : %v , avg value : %v  , current :%v \n",
			n.Sum(time.Now()), n.Avg(time.Now()), i)
		time.Sleep(1 * time.Second)
	}
}
