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
		if len(turn.Alice.Pokers) == 7{
			if turn.Alice.Pokers[6].Face == 0 {
				turn.Alice = handleCardZero(turn.Alice.Pokers)
			}else {
				turn.Alice = internalCompetition(turn.Alice)
			}
			if turn.Bob.Pokers[6].Face ==0{
				turn.Bob = handleCardZero(turn.Bob.Pokers)
			}else {
				turn.Bob = internalCompetition(turn.Bob)
			}
		}

	return turn
}

// AnalyseFeature 分析这局游戏中双方的手牌特征
func AnalyseFeatures(turn *model.Turn)*model.Turn{
	analyseFeature(&turn.Alice)
	analyseFeature(&turn.Bob)

	return turn
}

// AnalyseLevel 根据特征值判断level
func AnalyseLevel(turn *model.Turn)*model.Turn{
	turn.Alice.Level = tools.GetLevel(turn.Alice.Feature)
	turn.Bob.Level = tools.GetLevel(turn.Bob.Feature)

	return turn
}

// JudgeWinner 根据level等级判断赢家
func JudgeWinner(turn *model.Turn)*model.Turn{
	if turn.Alice.Level > turn.Bob.Level{
		turn.Winner = 0
	}else if turn.Bob.Level > turn.Alice.Level{
		turn.Winner = 1
	}else{
		turn.Winner = advancedJudgement(&turn.Alice,&turn.Bob)
	}

	return turn
}

// advancedJudgement 如果等级相同则进一步判断
func advancedJudgement(Alice,Bob *model.HandCards)int{
	// 牌等级为9，只比较最大一张
	if Alice.Level == 9{
		if Alice.Pokers[0].Face > Bob.Pokers[0].Face{
			return 0
		}
		if Alice.Pokers[0].Face < Bob.Pokers[0].Face{
			return 1
		}
		return -1
	}

	//牌等级为8，比较中间一张,若相同，比较第一张，若还是相同，比较最后一张
	if Alice.Level ==8{
		if Alice.Pokers[2].Face > Bob.Pokers[2].Face{return 0}
		if Alice.Pokers[2].Face < Bob.Pokers[2].Face{return 1}
		if Alice.Pokers[2].Face == Bob.Pokers[2].Face{
			if Alice.Pokers[0].Face > Bob.Pokers[0].Face{return 0}
			if Alice.Pokers[0].Face < Bob.Pokers[0].Face{return 1}
			if Alice.Pokers[0].Face == Bob.Pokers[0].Face{
				if Alice.Pokers[4].Face > Bob.Pokers[4].Face{return 0}
				if Alice.Pokers[4].Face < Bob.Pokers[4].Face{return 1}
				if Alice.Pokers[4].Face == Bob.Pokers[4].Face{return -1}
			}
		}
	}

	//牌等级为7，先比较中间一张，若相同，再根据相同牌的特征值比较第2张或第4张
	if Alice.Level ==7{
		if Alice.Pokers[2].Face > Bob.Pokers[2].Face{
			return 0
		}
		if Alice.Pokers[2].Face < Bob.Pokers[2].Face{
			return 1
		}
		if Alice.Pokers[2].Face == Bob.Pokers[2].Face{
			var a,b int
			if Alice.Feature.SameCards ==12 {
				a = Alice.Pokers[4].Face
			}else if Alice.Feature.SameCards ==21 {
				a = Alice.Pokers[1].Face
			}else{
				fmt.Println(Alice.Feature.SameCards)
				panic("Error ! Feature.SameCard ")
			}
			if Bob.Feature.SameCards ==12 {
				a = Bob.Pokers[4].Face
			}else if Bob.Feature.SameCards ==21 {
				a = Bob.Pokers[1].Face
			}else{
				fmt.Println(Bob.Feature.SameCards)
				panic("Error ! Feature.SameCard ")
			}

			if a>b{
				return 0
			}else if a<b{
				return 1
			}else {
				return -1
			}
		}
	}

	//若牌等级为6，比依次比较
	if Alice.Level ==6{
		return tools.CompareEachCard(Alice.Pokers,Bob.Pokers)
	}

	//若牌等级为5，比较第一张
	if Alice.Level ==5{
		if Alice.Pokers[0].Face > Bob.Pokers[0].Face{
			return 0
		}
		if Alice.Pokers[0].Face < Bob.Pokers[0].Face{
			return 1
		}
	}

	//此处可能需要优化！！
	//若牌等级为4，先比较中间一张，再比较第一张或第五张
	if Alice.Level ==4{
		if Alice.Pokers[2].Face > Bob.Pokers[2].Face{
			return 0
		}
		if Alice.Pokers[2].Face < Bob.Pokers[2].Face{
			return 1
		}

		if Alice.Pokers[2].Face == Bob.Pokers[2].Face{
			//一组牌中除了三张相同的牌外，比较另外两张单牌的大小
			//cardsA中大的给a1，小的给a2，cardsB同理
			var a1,a2,b1,b2 int
			if Alice.Pokers[0].Face != Alice.Pokers[2].Face{
				a1=Alice.Pokers[0].Face
				if Alice.Pokers[1].Face == Alice.Pokers[2].Face{
					a2 = Alice.Pokers[4].Face
				}else {
					a2 = Alice.Pokers[1].Face
				}
			}else {
				a1 = Alice.Pokers[3].Face
				a2 = Alice.Pokers[4].Face
			}

			if Bob.Pokers[0].Face != Bob.Pokers[2].Face{
				b1=Bob.Pokers[0].Face
				if Bob.Pokers[1].Face == Bob.Pokers[2].Face{
					b2 = Bob.Pokers[3].Face
				}else {
					b2 = Bob.Pokers[1].Face
				}
			}else {
				b1 = Bob.Pokers[3].Face
				b2 = Bob.Pokers[4].Face
			}

			//比较单牌大小
			if a1 > b1{
				return 0
			}else if a1<b1{
				return 1
			}else {
				if a2>b2{
					return 0
				}else if a2<b2{
					return 1
				}else {
					return -1
				}
			}

		}

	}

	//若牌等级为3，需要先判断两个对子的大小，再判断单牌大小
	if Alice.Level ==3{
		return tools.AdvancedCompareTwoPair(Alice.Feature.SameCards,Bob.Feature.SameCards,Alice.Pokers,Bob.Pokers)
	}

	//若牌等级为2，需要先判断一个对子大小，再判断单牌大小
	if Alice.Level ==2{
		return tools.AdvancedCompareOnePair(Alice.Pokers,Bob.Pokers)
	}

	//若牌等级为1，直接逐个比较牌大小
	if Alice.Level ==1{
		return tools.CompareEachCard(Alice.Pokers,Bob.Pokers)
	}

	//出现错误
	return -2
}
// analyseFeature 分析这局游戏中双方的手牌特征
func analyseFeature(cards *model.HandCards){
	cards.Feature.Continue =tools.GetContinueLength(cards.Pokers)
	cards.Feature.Continue = tools.GetContinueLength(cards.Pokers)
	cards.Feature.SameCards = tools.GetSameCards(cards.Pokers)
	cards.Feature.SameCards = tools.GetSameCards(cards.Pokers)
	cards.Feature.Flush = tools.CheckFlush(cards.Pokers)
	cards.Feature.Flush = tools.CheckFlush(cards.Pokers)

	// 处理black Jack
	if len(cards.Pokers) == 5{
		if cards.Feature.Continue == 4 && cards.Pokers[0].Face == 14 && cards.Pokers[1].Face  == 5{
			cards.Pokers[0].Face  = 1
			cards.Feature.Continue = 5
			cards.Pokers = tools.AdjustCards(cards.Pokers)
		}
	}else if len(cards.Pokers) == 4{
		if cards.Feature.SameCards == 0 && cards.Pokers[0].Face == 14 && cards.Pokers[1].Face == 5{
			cards.Pokers[0].Face = 1
			cards.Pokers = tools.AdjustCards(cards.Pokers)
		}
	}else {
		panic("Error wrong pokers number!")
	}


}


//use7cards

func internalCompetition(cards model.HandCards)model.HandCards{
	Combinations := tools.Choose5From7(cards.Pokers)
	return getBestCombination(Combinations)
}

func getBestCombination(Combinations []model.HandCards)(best model.HandCards){
	t := 0
	var bests []model.HandCards
	for i := range Combinations{
		Combinations[i].Pokers = tools.AdjustCards(Combinations[i].Pokers)
		analyseFeature(&Combinations[i])
		tmp := tools.GetLevel(Combinations[i].Feature)
		Combinations[i].Level = tmp
		if tmp > t{
			t = tmp
		}
	}

	for i := range Combinations{
		if Combinations[i].Level == t{
			bests = append(bests,Combinations[i])
		}
	}

	length := len(bests)

	if length == 0 {
		panic("Error")
	}

	if length == 1{
		return bests[0]
	}

	best = bests[0]
	for i:=1;i<length;i++{
		t = advancedJudgement(&best,&bests[i])
		if t == 1{
			best = bests[i]
		}
	}

return
}

func handleCardZero(pokers []model.Poker)model.HandCards{
	Combinations := tools.CalculateAllPossibilitiesWithCardZero(pokers)
	return getBestCombinationWithCardZero(Combinations)

}

func getBestCombinationWithCardZero(Combinations []model.HandCards)(best model.HandCards){
	t := 0
	var bests []model.HandCards
	for i := range Combinations{
		Combinations[i].Pokers = tools.AdjustCards(Combinations[i].Pokers)
		analyseFeature(&Combinations[i])
		tools.AssemblyZeroCard(&Combinations[i])
		tmp := Combinations[i].Level
		Combinations[i].Level = tmp
		if tmp > t{
			t = tmp
		}
	}

	for i := range Combinations{
		if Combinations[i].Level == t{
			bests = append(bests,Combinations[i])
		}
	}

	length := len(bests)

	if length == 0 {
		panic("Error")
	}

	if length == 1{
		return bests[0]
	}

	best = bests[0]
	tools.AdjustCards(best.Pokers)
	for i:=1;i<length;i++{
		tools.AdjustCards(bests[i].Pokers)
		t = advancedJudgement(&best,&bests[i])
		if t == 1{
			best = bests[i]
		}
	}

	return

}






