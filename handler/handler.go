package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"poker/model"
	"poker/tools"
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

	turn.Alice.Pokers = tools.AdjustCards(turn.Alice.Pokers)
	turn.Bob.Pokers = tools.AdjustCards(turn.Bob.Pokers)

	return turn
}

// AnalyseFeature 分析这局游戏中双方的手牌特征
func AnalyseFeature(turn *model.Turn)*model.Turn{
	turn.Alice.Feature.Continue = tools.GetContinueLength(turn.Alice.Pokers)
	turn.Bob.Feature.Continue = tools.GetContinueLength(turn.Bob.Pokers)
	turn.Alice.Feature.SameCards = tools.GetSameCards(turn.Alice.Pokers)
	turn.Bob.Feature.SameCards = tools.GetSameCards(turn.Bob.Pokers)
	turn.Alice.Feature.Flush = tools.CheckFlush(turn.Alice.Pokers)
	turn.Bob.Feature.Flush = tools.CheckFlush(turn.Bob.Pokers)

	// 处理black Jack
	if turn.Alice.Feature.Continue == 4 && turn.Alice.Pokers[0].Face == 14 && turn.Alice.Pokers[1].Face  == 5{
		turn.Alice.Pokers[0].Face  = 1
		turn.Alice.Feature.Continue = 5
		turn.Alice.Pokers = tools.AdjustCards(turn.Alice.Pokers)
	}

	if turn.Bob.Feature.Continue ==4 && turn.Bob.Pokers[0].Face == 14 && turn.Bob.Pokers[1].Face  == 5{
		turn.Bob.Pokers[0].Face  = 1
		turn.Bob.Feature.Continue = 5
		turn.Bob.Pokers = tools.AdjustCards(turn.Bob.Pokers)
	}

	return turn
}

// AnalyseLevel 根据特征值判断level
func AnalyseLevel(turn *model.Turn)*model.Turn{
	turn.Alice.Level = tools.GetLevel(&turn.Alice.Feature)
	turn.Bob.Level = tools.GetLevel(&turn.Bob.Feature)

	return turn
}

// JudgeWinner 根据level等级判断赢家
func JudgeWinner(turn *model.Turn)*model.Turn{
	if turn.Alice.Level > turn.Bob.Level{
		turn.Winner = 0
	}else if turn.Bob.Level > turn.Alice.Level{
		turn.Winner = 1
	}else{
		advancedJudgement(turn)
	}

	return turn
}

// advancedJudgement 如果等级相同则进一步判断
func advancedJudgement(turn *model.Turn)*model.Turn{
	// 牌等级为9，只比较最大一张
	if turn.Alice.Level == 9{
		if turn.Alice.Pokers[0].Face > turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[0].Face < turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn
	}

	//牌等级为8，比较中间一张,若相同，比较第一张，若还是相同，比较最后一张
	if turn.Alice.Level ==8{
		if turn.Alice.Pokers[2].Face > turn.Bob.Pokers[2].Face{turn.Winner = 0}
		if turn.Alice.Pokers[2].Face < turn.Bob.Pokers[2].Face{turn.Winner = 1}
		if turn.Alice.Pokers[2].Face == turn.Bob.Pokers[2].Face{
			if turn.Alice.Pokers[0].Face > turn.Bob.Pokers[0].Face{turn.Winner = 0}
			if turn.Alice.Pokers[0].Face < turn.Bob.Pokers[0].Face{turn.Winner = 1}
			if turn.Alice.Pokers[0].Face == turn.Bob.Pokers[0].Face{
				if turn.Alice.Pokers[4].Face > turn.Bob.Pokers[4].Face{turn.Winner = 0}
				if turn.Alice.Pokers[4].Face < turn.Bob.Pokers[4].Face{turn.Winner = 1}
				if turn.Alice.Pokers[4].Face == turn.Bob.Pokers[4].Face{turn.Winner = -1}
			}
		}
		return turn

	}

	//牌等级为7，先比较中间一张，若相同，再根据相同牌的特征值比较第2张或第4张
	if turn.Alice.Level ==7{
		if turn.Alice.Pokers[2].Face > turn.Bob.Pokers[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[2].Face < turn.Bob.Pokers[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[2].Face == turn.Bob.Pokers[2].Face{
			var a,b int
			if turn.Alice.Feature.SameCards ==12 {
				a = turn.Alice.Pokers[4].Face
			}else if turn.Alice.Feature.SameCards ==21 {
				a = turn.Alice.Pokers[1].Face
			}else{
				fmt.Println(turn.Alice.Feature.SameCards)
				panic("Error ! Feature.SameCard ")
			}
			if turn.Bob.Feature.SameCards ==12 {
				a = turn.Bob.Pokers[4].Face
			}else if turn.Bob.Feature.SameCards ==21 {
				a = turn.Bob.Pokers[1].Face
			}else{
				fmt.Println(turn.Bob.Feature.SameCards)
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
	if turn.Alice.Level ==6{
		if turn.Alice.Pokers[0].Face > turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[0].Face < turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn
	}

	//若牌等级为5，比较第一张
	if turn.Alice.Level ==5{
		if turn.Alice.Pokers[0].Face > turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[0].Face < turn.Bob.Pokers[0].Face{
			turn.Winner = 0
			return turn
		}
		return turn

	}

	//此处可能需要优化！！
	//若牌等级为4，先比较中间一张，再比较第一张或第五张
	if turn.Alice.Level ==4{
		if turn.Alice.Pokers[2].Face > turn.Bob.Pokers[2].Face{
			turn.Winner = 0
			return turn
		}
		if turn.Alice.Pokers[2].Face < turn.Bob.Pokers[2].Face{
			turn.Winner = 0
			return turn
		}

		if turn.Alice.Pokers[2].Face == turn.Bob.Pokers[2].Face{
			//一组牌中除了三张相同的牌外，比较另外两张单牌的大小
			//cardsA中大的给a1，小的给a2，cardsB同理
			var a1,a2,b1,b2 int
			if turn.Alice.Pokers[0].Face != turn.Alice.Pokers[2].Face{
				a1=turn.Alice.Pokers[0].Face
				if turn.Alice.Pokers[1].Face == turn.Alice.Pokers[2].Face{
					a2 = turn.Alice.Pokers[3].Face
				}else {
					a2 = turn.Alice.Pokers[1].Face
				}
			}else {
				a1 = turn.Alice.Pokers[3].Face
				a2 = turn.Alice.Pokers[4].Face
			}

			if turn.Bob.Pokers[0].Face != turn.Bob.Pokers[2].Face{
				a1=turn.Bob.Pokers[0].Face
				if turn.Bob.Pokers[1].Face == turn.Bob.Pokers[2].Face{
					a2 = turn.Bob.Pokers[3].Face
				}else {
					a2 = turn.Bob.Pokers[1].Face
				}
			}else {
				a1 = turn.Bob.Pokers[3].Face
				a2 = turn.Bob.Pokers[4].Face
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
	if turn.Alice.Level ==3{
		turn.Winner = tools.AdvancedCompareTwoPair(turn.Alice.Feature.SameCards,turn.Bob.Feature.SameCards,turn.Alice.Pokers,turn.Bob.Pokers)
	}

	//若牌等级为2，需要先判断一个对子大小，再判断单牌大小
	if turn.Alice.Level ==2{
		turn.Winner = tools.AdvancedCompareOnePair(turn.Alice.Pokers,turn.Bob.Pokers)
	}

	//若牌等级为一，直接逐个比较牌大小
	if turn.Alice.Level ==1{
		turn.Winner = tools.CompareEachCard(turn.Alice.Pokers,turn.Bob.Pokers)
	}

	return turn
}

//use7cards

//func InternalCompetition(turn *model.Turn)*model.Turn{
//	Combinations := tools.CalculateAllPossibilities(turn.Alice.Pokers)
//	turn.Alice.Pokers = GetBestCombination(Combinations).Pokers
//
//	return nil
//}



//func GetBestCombination(Combinations []model.HandCards)(best model.HandCards){
//return nil
//}






