package tools

import (
	"poker/model"
	"strconv"
)
// ChangeFaceToNumber 把牌面从string转换为int
// 牌面为花牌则按 T:10;J:11;Q:12,K:13;A:14来进行转换
// 返回牌面的数值int型
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

// PutCardIntoHand 将字符串信息转换写入model.poker并放到当前回合里
// 返回model.Turn
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

// AdjustCards 把牌按Face从大到小排列
// 输入model.turn的指针，无返回
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

// GetContinueLength 计算一手牌中连续牌的长度
// 输入model.Poker的切片，返回连续牌长度的int
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
// 输入model.Poker的切片,返回相同牌的特征值
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

// CheckFlush 判断一手牌是否同花色
// 输入model.Poker的切片,返回是否同花色的bool值
func CheckFlush(pokers []model.Poker)bool{
	length := len(pokers)
	for i:=0;i<length-1;i++{
		if pokers[i].Color != pokers[i+1].Color{
			return false
		}
	}
	return true
}

// GetLevel 根据特征值判断一手牌的等级
// 输入一组牌的特征值model.Feature,返回其等级Level的int值
func GetLevel(feature *model.Feature)int{

	//straight flush or royal flush
	if feature.Flush && feature.Continue ==5{
		return 9
	}

	//four of a kind
	if feature.SameCards == 3{
		return 8
	}

	//full house
	if feature.SameCards ==21 || feature.SameCards ==12{
		return 7
	}

	//flush
	if feature.Flush{
		return 6
	}

	//straight
	if feature.Continue ==5 {
		return 5
	}

	//three of a kind
	if feature.SameCards ==2{
		return 4
	}

	//two pairs
	if feature.SameCards ==11 || feature.SameCards ==101 || feature.SameCards ==110{
		return 3
	}

	//one pair
	if feature.SameCards ==1{
		return 2
	}

	//high card
	return 1


}

// CompareEachCard 从两手牌第一张开始逐一相互判断牌面的大小
// 输入两个[]model.Poker,若第一组大则返回int 0 若第二组大则返回int 1 若完全相同则返回int -1
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

// AdvancedCompare 当两手牌中同时出现一对或两对时使用
// 输入两手牌的相同牌特征值，并输入两手牌[]model.Poker
// 若第一组手牌大则返回0,若第二组手牌大则返回1，若两手牌同样大则返回-1
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




