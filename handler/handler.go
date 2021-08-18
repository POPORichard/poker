package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"poker/tools"
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
	turn := tools.PutCardIntoHand(match)

	turn.AliceHandCard = tools.AdjustCards(turn.AliceHandCard)
	turn.BobHandCard = tools.AdjustCards(turn.BobHandCard)

	return turn
}

// AnalyseFeature 分析这局游戏中双方的手牌特征
func AnalyseFeature(turn *model.Turn)*model.Turn{
	turn.AliceFeature.Continue = tools.GetContinueLength(turn.AliceHandCard)
	turn.BobFeature.Continue = tools.GetContinueLength(turn.BobHandCard)
	turn.AliceFeature.SameCards = tools.GetSameCards(turn.AliceHandCard)
	turn.BobFeature.SameCards = tools.GetSameCards(turn.BobHandCard)
	turn.AliceFeature.Flush = tools.CheckFlush(turn.AliceHandCard)
	turn.BobFeature.Flush = tools.CheckFlush(turn.BobHandCard)

	// 处理black Jack
	if turn.AliceFeature.Continue == 4 && turn.AliceHandCard[0].Face == 14 && turn.AliceHandCard[1].Face  == 5{
		turn.AliceHandCard[0].Face  = 1
		turn.AliceFeature.Continue = 5
		turn.AliceHandCard = tools.AdjustCards(turn.AliceHandCard)
	}

	if turn.BobFeature.Continue ==4 && turn.BobHandCard[0].Face == 14 && turn.BobHandCard[1].Face  == 5{
		turn.BobHandCard[0].Face  = 1
		turn.BobFeature.Continue = 5
		turn.BobHandCard = tools.AdjustCards(turn.BobHandCard)
	}

	return turn
}

// AnalyseLevel 根据特征值判断level
func AnalyseLevel(turn *model.Turn)*model.Turn{
	turn.AliceLevel = tools.GetLevel(&turn.AliceFeature)
	turn.BobLevel = tools.GetLevel(&turn.BobFeature)

	return turn
}

// JudgeWinner 根据level等级判断赢家
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

// advancedJudgement 如果等级相同则进一步判断
func advancedJudgement(turn *model.Turn)*model.Turn{
	// 牌等级为9，只比较最大一张
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

	//牌等级为8，比较中间一张,若相同，比较第一张，若还是相同，比较最后一张
	if turn.AliceLevel ==8{
		if turn.AliceHandCard[2].Face > turn.BobHandCard[2].Face{turn.Winner = 0}
		if turn.AliceHandCard[2].Face < turn.BobHandCard[2].Face{turn.Winner = 1}
		if turn.AliceHandCard[2].Face == turn.BobHandCard[2].Face{
			if turn.AliceHandCard[0].Face > turn.BobHandCard[0].Face{turn.Winner = 0}
			if turn.AliceHandCard[0].Face < turn.BobHandCard[0].Face{turn.Winner = 1}
			if turn.AliceHandCard[0].Face == turn.BobHandCard[0].Face{
				if turn.AliceHandCard[4].Face > turn.BobHandCard[4].Face{turn.Winner = 0}
				if turn.AliceHandCard[4].Face < turn.BobHandCard[4].Face{turn.Winner = 1}
				if turn.AliceHandCard[4].Face == turn.BobHandCard[4].Face{turn.Winner = -1}
			}
		}
		return turn

	}

	//牌等级为7，先比较中间一张，若相同，再根据相同牌的特征值比较第2张或第4张
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

	//若牌等级为7，比较第一张
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

	//若牌等级为5，比较第一张
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

	//此处可能需要优化！！
	//若牌等级为4，先比较中间一张，再比较第一张或第五张
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
			//一组牌中除了三张相同的牌外，比较另外两张单牌的大小
			//cardsA中大的给a1，小的给a2，cardsB同理
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

			//比较单牌大小
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

	//若牌等级为3，需要先判断两个对子的大小，再判断单牌大小
	if turn.AliceLevel ==3{
		turn.Winner = tools.AdvancedCompareTwoPair(turn.AliceFeature.SameCards,turn.BobFeature.SameCards,turn.AliceHandCard,turn.BobHandCard)
	}

	//若牌等级为2，需要先判断一个对子大小，再判断单牌大小
	if turn.AliceLevel ==2{
		turn.Winner = tools.AdvancedCompareOnePair(turn.AliceHandCard,turn.BobHandCard)
	}

	//若牌等级为一，直接逐个比较牌大小
	if turn.AliceLevel ==1{
		turn.Winner = tools.CompareEachCard(turn.AliceHandCard,turn.BobHandCard)
	}

	return turn
}






