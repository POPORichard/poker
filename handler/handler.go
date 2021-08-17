package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"poker/Tools"
	"poker/model"
)
// ReadDataToModel 读取json文件数据
func ReadDataToModel(path string)[]model.Data{
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取json文件失败", err)
		return nil
	}

	var inputData model.InputData
	err = json.Unmarshal(bytes, &inputData)
	if err != nil {
		fmt.Println("解析数据失败", err)
		return nil
	}

	//length := len(inputData.Matches)

	data := inputData.Matches
	return data
}

// CreateTurn 创建一局游戏
func CreateTurn(match *model.Data) model.Turn {
	turn := Tools.PutCardIntoHand(match)
	Tools.AdjustCards(&turn)

	return turn
}
// AnalyseFeature 分析这局游戏中双方的手牌特征
func AnalyseFeature(turn *model.Turn)*model.Turn{
	turn.AliceFeature.Continue = Tools.GetContinueLength(turn.AliceHandCard)
	turn.BobFeature.Continue = Tools.GetContinueLength(turn.BobHandCard)
	turn.AliceFeature.SameCards = Tools.GetSameCards(turn.AliceHandCard)
	turn.BobFeature.SameCards = Tools.GetSameCards(turn.BobHandCard)
	turn.AliceFeature.Flush = Tools.CheckFlush(turn.AliceHandCard)
	turn.BobFeature.Flush = Tools.CheckFlush(turn.BobHandCard)

	return turn
}
// AnalyseLevel 根据特征值判断level
func AnalyseLevel(turn *model.Turn)*model.Turn{
	turn.AliceLevel = Tools.GetLevel(&turn.AliceFeature)
	turn.BobLevel = Tools.GetLevel(&turn.BobFeature)

	return turn
}

func JudgeWinner(turn *model.Turn)*model.Turn{
	if turn.AliceLevel > turn.BobLevel{
		turn.Winner = 0
	}else if turn.BobLevel > turn.AliceLevel{
		turn.Winner = 1
	}else{
		advancedJudgement(turn)
	}

	return turn
}

func advancedJudgement(turn *model.Turn)*model.Turn{
	if turn.AliceLevel == 9{
		if turn.AliceHandCard[0].Face > turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[0].Face < turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn
	}

	if turn.AliceLevel ==8{
		if turn.AliceHandCard[2].Face > turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[2].Face < turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		return turn

	}

	if turn.AliceLevel ==7{
		if turn.AliceHandCard[2].Face > turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[2].Face < turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[2].Face == turn.BobHandCard[2].Face{
			var a,b int
			if turn.AliceFeature.SameCards ==12 {
				a = turn.AliceHandCard[4].Face
			}else if turn.AliceFeature.SameCards ==21 {
				a = turn.AliceHandCard[1].Face
			}else{
				fmt.Println(turn.AliceFeature.SameCards)
				panic("Error ! Feature.SameCard ")
			}
			if turn.BobFeature.SameCards ==12 {
				a = turn.BobHandCard[4].Face
			}else if turn.BobFeature.SameCards ==21 {
				a = turn.BobHandCard[1].Face
			}else{
				fmt.Println(turn.BobFeature.SameCards)
				panic("Error ! Feature.SameCard ")
			}

			if a>b{
				turn.Winner = 0
			}else if a<b{
				turn.Winner = 1
			}else {
				turn.Winner = -1
			}
		}
		return turn
	}

	if turn.AliceLevel ==6{
		if turn.AliceHandCard[0].Face > turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[0].Face < turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn
	}

	if turn.AliceLevel ==5{
		if turn.AliceHandCard[0].Face > turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[0].Face < turn.BobHandCard[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn

	}

	if turn.AliceLevel ==4{
		if turn.AliceHandCard[2].Face > turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[2].Face < turn.BobHandCard[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.AliceHandCard[2].Face == turn.BobHandCard[2].Face{
			var a1,a2,b1,b2 int
			if turn.AliceHandCard[0].Face != turn.AliceHandCard[2].Face{
				a1=turn.AliceHandCard[0].Face
				if turn.AliceHandCard[1].Face == turn.AliceHandCard[2].Face{
					a2 = turn.AliceHandCard[3].Face
				}else {
					a2 = turn.AliceHandCard[1].Face
				}
			}else {
				a1 = turn.AliceHandCard[3].Face
				a2 = turn.AliceHandCard[4].Face
			}

			if turn.BobHandCard[0].Face != turn.BobHandCard[2].Face{
				a1=turn.BobHandCard[0].Face
				if turn.BobHandCard[1].Face == turn.BobHandCard[2].Face{
					a2 = turn.BobHandCard[3].Face
				}else {
					a2 = turn.BobHandCard[1].Face
				}
			}else {
				a1 = turn.BobHandCard[3].Face
				a2 = turn.BobHandCard[4].Face
			}

			if a1 > b1{
				turn.Winner = 0
			}else if a1<b1{
				turn.Winner = 1
			}else {
				if a2>b2{
					turn.Winner = 0
				}else if a2<b2{
					turn.Winner = 1
				}else {
					turn.Winner = -1
				}
			}

			return turn
		}

	}

	if turn.AliceLevel ==3{
		turn.Winner = Tools.AdvancedCompare(turn.AliceFeature.SameCards,turn.BobFeature.SameCards,turn.AliceHandCard,turn.BobHandCard)

	}
	if turn.AliceLevel ==2{
		turn.Winner = Tools.AdvancedCompare(turn.AliceFeature.SameCards,turn.BobFeature.SameCards,turn.AliceHandCard,turn.BobHandCard)
	}

	if turn.AliceLevel ==1{
		turn.Winner = Tools.CompareEachCard(turn.AliceHandCard,turn.BobHandCard)
	}

	return turn
}





