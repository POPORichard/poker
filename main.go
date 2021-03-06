package main

import (
	"fmt"
	"poker/handler"
	"time"
)

var Use7Cards bool = true

func main() {

	start:=time.Now()
	data := handler.ReadDataToModel("./seven_cards_with_ghost.result.json")
	for i := range data{
		turn := handler.CreateTurn(&data[i])
		if !Use7Cards {
			handler.AnalyseFeatures(&turn)
			handler.AnalyseLevel(&turn)
		}
		handler.JudgeWinner(&turn)
		if turn.Winner+1 != data[i].Result{
			fmt.Println("-----Error-----")
			fmt.Println("position:",i)
			fmt.Println(data[i])
			fmt.Println(turn)
			panic("result wrong")
		}
	}

	cost:=time.Since(start)
	fmt.Println("成功！ 耗时：",cost)

}
