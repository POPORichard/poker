package Tools

import (
	"poker/model"
	"strconv"
)

func ChangeFaceToNumber(f string)int{
	switch f {
	case "T":
		return 10
	case "J":
		return 11
	case "Q":
		return 12
	case "K":
		return 13
	case "A":
		return 14
	default:
		int,_:=strconv.Atoi(f)
		return int

	}
}

func PutCardIntoHand(data *model.Data) (turn model.Turn) {
	a := make([]model.Poker,5,5)
	b := make([]model.Poker,5,5)


	for i:=0;i<5;i++{
			a[i].Face = ChangeFaceToNumber(string(data.Alice[i*2]))
			a[i].Color = string(data.Alice[i*2+1])
			b[i].Face = ChangeFaceToNumber(string(data.Bob[i*2]))
			b[i].Color = string(data.Bob[i*2+1])
	}

	turn.AliceHandCard = a
	turn.BobHandCard = b


	return
}

func AdjustCards(turn *model.Turn){
	var tmp model.Poker
	var target int
	var biggest int

	for i :=range turn.AliceHandCard{
		biggest = turn.AliceHandCard[i].Face
		target=i
		for t:=i+1;t<len(turn.AliceHandCard);t++{
			if biggest < turn.AliceHandCard[t].Face{
				biggest = turn.AliceHandCard[t].Face
				target = t
			}
		}
		if target !=i {
			tmp = turn.AliceHandCard[i]
			turn.AliceHandCard[i] = turn.AliceHandCard[target]
			turn.AliceHandCard[target] = tmp
		}
	}
	// 下面是重复代码务必后面改进
	for i :=range turn.BobHandCard{
		biggest = turn.BobHandCard[i].Face
		target=i
		for t:=i+1;t<len(turn.BobHandCard);t++{
			if biggest < turn.BobHandCard[t].Face{
				biggest = turn.BobHandCard[t].Face
				target = t
			}
		}
		if target !=i {
			tmp = turn.BobHandCard[i]
			turn.BobHandCard[i] = turn.BobHandCard[target]
			turn.BobHandCard[target] = tmp
		}
	}
}

func GetContinueLength(pokers []model.Poker) int {
	length := len(pokers)
	result := 1
	t:=1
	for i:=0;i< length-1;i++{
		if pokers[i].Face-1 == pokers[i+1].Face {
			t++
		}else if pokers[i].Face == pokers[i+1].Face{
			continue
		}else {
			t = 1
		}

		if t>result{
			result = t
		}
	}
	return result
}

// GetSameCards 获取这组牌中相同face的牌的数量
// sameCards个位表示第一个相同组的牌的数量、十位即第二个相同组
// 每个位上数字加一即为这一相同组的牌的数量
func GetSameCards(pokers []model.Poker)(sameCards int){
	sameCards = 0
	t:=1
	length := len(pokers)

	for i:=0;i<length-1;i++{
		if pokers[i].Face == pokers[i+1].Face {
			sameCards = sameCards + t
		} else{
			t=t*10
		}
	}

	return
}

func CheckFlush(pokers []model.Poker)bool{
	length := len(pokers)
	for i:=0;i<length-1;i++{
		if pokers[i].Color != pokers[i+1].Color{
			return false
		}
	}
	return true
}

func GetLevel(feature *model.Feature)int{
	if feature.Flush && feature.Continue ==5{
		return 9
	}
	if feature.SameCards == 3{
		return 8
	}
	if feature.SameCards ==21 || feature.SameCards ==12{
		return 7
	}
	if feature.Flush{
		return 6
	}
	if feature.Continue ==5 {
		return 5
	}
	if feature.SameCards ==2{
		return 4
	}
	if feature.SameCards ==11 || feature.SameCards ==101 || feature.SameCards ==110{
		return 3
	}
	if feature.SameCards ==1{
		return 2
	}
	return 1


}

func CompareEachCard(cardsA,cardsB []model.Poker)int{
	for i := range cardsA{
		if cardsA[i].Face > cardsB[i].Face{
			return 0
		}
		if cardsA[i].Face < cardsB[i].Face{
			return 1
		}
	}
	return -1
}

func AdvancedCompare(sameCardFeatureA,sameCardFeatureB int,cardsA,cardsB []model.Poker)int{
	lenth := len(cardsA)

	if sameCardFeatureA ==1{
		a := 0
		b := 0
		for i:=0;i<lenth-1;i++{
			if cardsA[i].Face == cardsA[i+1].Face{
				a = cardsA[i].Face
			}
		}
		if a==0{
			panic("error")
		}
		for i:=0;i<lenth-1;i++{
			if cardsB[i].Face == cardsB[i+1].Face{
				b = cardsB[i].Face
			}
		}
		if b==0{
			panic("error")
		}

		if a>b{return 0}
		if a<b{return 1}
		if a==b{
			return CompareEachCard(cardsA,cardsB)
		}
	}
	var B1,B2,B3 int
	if sameCardFeatureB == 11{
		B1 = 1
		B2 = 3
		B3 = 0
	}else if sameCardFeatureB ==101{
		B1 = 0
		B2 = 3
		B3 = 2
	}else if sameCardFeatureB == 110{
		B1 = 0
		B2 = 2
		B3 = 4
	}else {
		panic("Error")
	}

	if sameCardFeatureA == 11{
		if cardsA[1].Face > cardsB[B1].Face{return 0}
		if cardsA[1].Face < cardsB[B1].Face{return 1}
		if cardsA[1].Face == cardsB[B1].Face{
			if cardsA[3].Face > cardsB[B2].Face{return 0}
			if cardsA[3].Face < cardsB[B2].Face{return 1}
			if cardsA[3].Face ==  cardsB[B2].Face{
				if cardsA[0].Face > cardsB[B3].Face{return 1}
				if cardsA[0].Face < cardsB[B3].Face{return 0}
				if cardsA[0].Face == cardsB[B3].Face{return -1}
			}
		}
	}

	if sameCardFeatureA ==101{
		if cardsA[0].Face > cardsB[B1].Face{return 0}
		if cardsA[0].Face < cardsB[B1].Face{return 1}
		if cardsA[0].Face == cardsB[B1].Face{
			if cardsA[3].Face > cardsB[B2].Face{return 0}
			if cardsA[3].Face < cardsB[B2].Face{return 1}
			if cardsA[3].Face ==  cardsB[B2].Face{
				if cardsA[2].Face > cardsB[B3].Face{return 1}
				if cardsA[2].Face < cardsB[B3].Face{return 0}
				if cardsA[2].Face == cardsB[B3].Face{return -1}
			}
		}
	}

	if sameCardFeatureA ==110{
		if cardsA[0].Face > cardsB[B1].Face{return 0}
		if cardsA[0].Face < cardsB[B1].Face{return 1}
		if cardsA[0].Face == cardsB[B1].Face{
			if cardsA[2].Face > cardsB[B2].Face{return 0}
			if cardsA[2].Face < cardsB[B2].Face{return 1}
			if cardsA[2].Face ==  cardsB[B2].Face{
				if cardsA[4].Face > cardsB[B3].Face{return 1}
				if cardsA[4].Face < cardsB[B3].Face{return 0}
				if cardsA[4].Face == cardsB[B3].Face{return -1}
			}
		}
	}




	return -2
}





