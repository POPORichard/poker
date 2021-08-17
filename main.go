package main

import (
	"fmt"
	"poker/handler"
)

func main() {

	data := handler.ReadDataToModel("./match.json")
	//fmt.Println(data)

	for i := range data{
		turn := handler.CreateTurn(&data[i])
		handler.AnalyseFeature(&turn)
		handler.AnalyseLevel(&turn)
		handler.JudgeWinner(&turn)

		fmt.Printf("%v ",turn.Winner)
	}


}
