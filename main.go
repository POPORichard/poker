package main

import (
	"fmt"
	"poker/handler"
	"time"
)

func main() {

	data := handler.ReadDataToModel("./match_result.json")
	//fmt.Println(data)

	start:=time.Now()

	for i := range data{
		turn := handler.CreateTurn(&data[i])
		handler.AnalyseFeature(&turn)
		handler.AnalyseLevel(&turn)
		handler.JudgeWinner(&turn)

		if turn.Winner+1 != data[i].Result{
			fmt.Println("-----Error-----")
			fmt.Println(data[i])
			fmt.Println(turn)
			panic("result wrong")
		}
	}

	cost:=time.Since(start)
	fmt.Println("成功！ 耗时：",cost)


}
