package main

import (
	"fmt"
	"poker/handler"
	"time"
)

var Use7Cards bool = false

func main() {

	//start:=time.Now()
	//data := handler.ReadDataToModel("./match_result.json")
	//
	//for i := range data{
	//	turn := handler.CreateTurn(&data[i])
	//	handler.AnalyseFeatures(&turn)
	//	handler.AnalyseLevel(&turn)
	//	handler.JudgeWinner(&turn)
	//
	//	if turn.Winner+1 != data[i].Result{
	//		fmt.Println("-----Error-----")
	//		fmt.Println(data[i])
	//		fmt.Println(turn)
	//		panic("result wrong")
	//	}
	//}
	//
	//cost:=time.Since(start)
	//fmt.Println("成功！ 耗时：",cost)

	start := time.Now()
	data := handler.ReadDataToModel("./seven_cards_with_ghost.result.json")
	for i := range data {
		turn := handler.CreateTurn(&data[i])
		handler.JudgeWinner(&turn)
		if turn.Winner+1 != data[i].Result {
			fmt.Println("-----Error-----")
			fmt.Println(data[i])
			fmt.Println(turn)
			panic("result wrong")

			fmt.Println(turn)
		}
	}
	cost:=time.Since(start)
	fmt.Println("成功！ 耗时：",cost)
}
