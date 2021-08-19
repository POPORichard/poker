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
	case "X":
		return 0
	default:
		int,_:=strconv.Atoi(f)
		return int

	}
}

// PutCardIntoHand 将字符串信息转换写入model.poker并放到当前回合里
// 返回model.Turn
func PutCardIntoHand(data *model.Data) (turn model.Turn) {
	length := len(data.Alice)/2
	a := make([]model.Poker, length, length)
	b := make([]model.Poker, length, length)

	for i:=0;i< length;i++{
			a[i].Face = ChangeFaceToNumber(string(data.Alice[i*2]))
			a[i].Color = string(data.Alice[i*2+1])
			b[i].Face = ChangeFaceToNumber(string(data.Bob[i*2]))
			b[i].Color = string(data.Bob[i*2+1])
	}

	turn.Alice.Pokers = a
	turn.Bob.Pokers = b

	return
}

// AdjustCards 把牌按Face从大到小排列
// 输入model.turn的指针，无返回
func AdjustCards(pokers []model.Poker)[]model.Poker{
	var tmp model.Poker
	length := len(pokers)

	var target int
	var biggest int

	for i :=range pokers{
		biggest = pokers[i].Face
		target=i
		for t:=i+1;t< length;t++{
			if biggest < pokers[t].Face{
				biggest = pokers[t].Face
				target = t
			}
		}
		if target !=i {
			tmp = pokers[i]
			pokers[i] = pokers[target]
			pokers[target] = tmp
		}
	}
	return pokers
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
func GetLevel(feature model.Feature)int{

	//straight flush or royal flush
	if feature.Flush && feature.Continue ==5{
		return 9
	}

	//four of a kind
	if feature.SameCards == 3 || feature.SameCards == 30{
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
	if feature.SameCards ==2 || feature.SameCards == 20 || feature.SameCards == 200{
		return 4
	}

	//two pairs
	if feature.SameCards ==11 || feature.SameCards ==101 || feature.SameCards ==110{
		return 3
	}

	//one pair
	if feature.SameCards ==1 || feature.SameCards ==10 || feature.SameCards ==100 || feature.SameCards ==1000{
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

// AdvancedCompareTwoPair 当两手牌中同时两对时使用
// 输入两手牌的相同牌特征值，并输入两手牌[]model.Poker
// 若第一组手牌大则返回0,若第二组手牌大则返回1，若两手牌同样大则返回-1
func AdvancedCompareTwoPair(sameCardFeatureA,sameCardFeatureB int,cardsA,cardsB []model.Poker)int{

	//根据sameCard特征值找到该先比较那一组卡
	//sameCard特征值为11: 前面4张牌两两成对
	//sameCard特征值为101: 前面2张和最后2张牌两两成对
	//sameCard特征值为110: 后面4张牌两两成对
	var B1,B2,B3 int
	if sameCardFeatureB == 11{
		B1 = 0
		B2 = 2
		B3 = 4
	}else if sameCardFeatureB ==101{
		B1 = 0
		B2 = 3
		B3 = 2
	}else if sameCardFeatureB == 110{
		B1 = 1
		B2 = 3
		B3 = 0
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
				if cardsA[0].Face > cardsB[B3].Face{return 0}
				if cardsA[0].Face < cardsB[B3].Face{return 1}
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
				if cardsA[2].Face > cardsB[B3].Face{return 0}
				if cardsA[2].Face < cardsB[B3].Face{return 1}
				if cardsA[2].Face == cardsB[B3].Face{return -1}
			}
		}
	}

	if sameCardFeatureA ==110{
		if cardsA[1].Face > cardsB[B1].Face{return 0}
		if cardsA[1].Face < cardsB[B1].Face{return 1}
		if cardsA[1].Face == cardsB[B1].Face{
			if cardsA[3].Face > cardsB[B2].Face{return 0}
			if cardsA[3].Face < cardsB[B2].Face{return 1}
			if cardsA[3].Face ==  cardsB[B2].Face{
				if cardsA[0].Face > cardsB[B3].Face{return 0}
				if cardsA[0].Face < cardsB[B3].Face{return 1}
				if cardsA[0].Face == cardsB[B3].Face{return -1}
			}
		}
	}




	return -2
}

// AdvancedCompareOnePair 当两手牌中同时两对时使用
// 输入两手牌的相同牌特征值，并输入两手牌[]model.Poker
// 若第一组手牌大则返回0,若第二组手牌大则返回1，若两手牌同样大则返回-1
func AdvancedCompareOnePair(cardsA,cardsB []model.Poker)int{
	length := len(cardsA)
	a := 0
	b := 0

	//找出cardsA中一对牌的Face
	for i:=0;i< length-1;i++{
		if cardsA[i].Face == cardsA[i+1].Face{
			a = cardsA[i].Face
		}
	}
	if a==0{
		panic("error")
	}

	//找出cardsB中一对牌的Face
	for i:=0;i< length-1;i++{
		if cardsB[i].Face == cardsB[i+1].Face{
			b = cardsB[i].Face
		}
	}
	if b==0{
		panic("error")
	}

	//对成对的牌比较，若相同，则依次比较每张牌
	if a>b{return 0}
	if a<b{return 1}
	if a==b{
		return CompareEachCard(cardsA,cardsB)
	}

	//错误情况下返回-2
	return -2
}

//use7cards

func chooseColor(pokers []model.Poker)string{
	a := [4]int{0}
	t := 0
	for _,poker := range pokers{
		if poker.Color == "s"{
			a[0]++
		}
		if poker.Color == "h"{
			a[1]++
		}
		if poker.Color == "d"{
			a[2]++
		}
		if poker.Color == "c"{
			a[3]++
		}
	}

	n:=a[0]
	for i:=1;i<4;i++{
		if n<a[i]{
			n = a[i]
			t = i
		}
	}
	if t == 0{return "s"}
	if t == 1{return "h"}
	if t == 2{return "d"}
	if t == 3{return "c"}
	return "n"

}

func choose4From6(pokers []model.Poker)[]model.HandCards{
	pokers = pokers[:6]
	pointer := 0
	results := make([]model.HandCards,15,15)

	for t:=0;t<5;t++{
		for i:=t+1;i<6;i++{
			result := make([]model.Poker,0,7)
			tmp := make([]model.Poker,0,7)
			for x:= range pokers{
				tmp = append(tmp,pokers[x])
			}
			result = append(tmp[:t],tmp[t+1:]...)
			result = append(result[:i-1],result[i:]...)
			results[pointer].Pokers = result
			pointer++
		}
	}

	//for i :=range results{
	//	results[i].Pokers = AdjustCards(results[i].Pokers)
	//	fmt.Println(results[i].Pokers)
	//}
	//
	//fmt.Println("==============")
	//
	//for i:= range results{
	//	for t:=i+1;t<15;t++{
	//		if reflect.DeepEqual(results[i],results[t]){
	//			fmt.Println(results[i],i)
	//			fmt.Println(results[t],t)
	//		}
	//	}
	//
	//}
	//fmt.Println("len:",len(results))
	//
	//panic("end")
	return results
}

func Choose5From7(pokers []model.Poker)[]model.HandCards{
	pointer := 0
	results := make([]model.HandCards,21,21)

	for t:=0;t<6;t++{
		for i:=t+1;i<7;i++{
			result := make([]model.Poker,0,7)
			tmp := make([]model.Poker,0,7)
			for x:= range pokers{
				tmp = append(tmp,pokers[x])
			}
				result = append(tmp[:t],tmp[t+1:]...)
				result = append(result[:i-1],result[i:]...)
				results[pointer].Pokers = result
				pointer++
		}
	}

	//for i :=range results{
	//	results[i].Pokers = AdjustCards(results[i].Pokers)
	//	fmt.Println(results[i].Pokers)
	//}
	//
	//fmt.Println("==============")
	//
	//for i:= range results{
	//	for t:=i+1;t<21;t++{
	//		if reflect.DeepEqual(results[i],results[t]){
	//			fmt.Println(results[i],i)
	//			fmt.Println(results[t],t)
	//		}
	//	}
	//
	//}
	//fmt.Println("len:",len(results))

	return results
}

func CalculateAllPossibilitiesWithCardZero(pokers []model.Poker)(AllPossibilities[]model.HandCards){
	AllPossibilities = choose4From6(pokers)
	return
}

func AssemblyZeroCard(card *model.HandCards){

	cardZero := model.Poker{
		Face:  0,
		Color: chooseColor(card.Pokers),
	}

	//Level 9
	if card.Feature.Continue == 4 && card.Feature.Flush{
		if card.Pokers[0].Face == 14{
			cardZero.Face = 10
		}else{
			cardZero.Face = card.Pokers[0].Face+1
		}
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 9

		return
	}
	if card.Feature.Continue ==3 && card.Feature.Flush{
		if card.Pokers[0].Face-1 == card.Pokers[1].Face{
			if card.Pokers[2].Face-2 == card.Pokers[3].Face{
				cardZero.Face = card.Pokers[2].Face-1
				card.Pokers = append(card.Pokers,cardZero)
				card.Level = 9
				return
			}
		}else {
			if card.Pokers[0].Face-2 == card.Pokers[1].Face{
				cardZero.Face = card.Pokers[0].Face-2
				card.Pokers = append(card.Pokers,cardZero)
				card.Level = 9
				return
			}
		}
	}
	if card.Feature.Continue == 2 && card.Feature.Flush{
		if card.Pokers[0].Face -1 == card.Pokers[1].Face && card.Pokers[2].Face-1 == card.Pokers[3].Face && card.Pokers[1].Face-2 == card.Pokers[2].Face{
			cardZero.Face = card.Pokers[1].Face-1
			card.Pokers = append(card.Pokers,cardZero)
			card.Level = 9
			return
		}
	}


	//Level 8
	if card.Feature.SameCards ==3{
		if card.Pokers[0].Face == 14{
			cardZero.Face = 13
		}else{
			cardZero.Face = card.Pokers[0].Face+1
		}
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 8
		return
	}
	if card.Feature.SameCards == 2 {
		cardZero.Face = card.Pokers[0].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 8
		return
	}
	if card.Feature.SameCards == 20{
		cardZero.Face = card.Pokers[1].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 8
		return
	}

	//Level 7

	if card.Feature.SameCards == 11{
		cardZero.Face = card.Pokers[0].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 7
		return
	}

	//Level 6
	if card.Feature.Flush{
		cardZero.Face = 14
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 6
		return
	}

	//Level 5
	if card.Feature.Continue == 4{
		if card.Pokers[0].Face == 14{
			cardZero.Face = 10
		}else{
			cardZero.Face = card.Pokers[0].Face+1
		}
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 5
		return
	}
	if card.Feature.Continue ==3{
		if card.Pokers[0].Face-1 == card.Pokers[1].Face{
			if card.Pokers[2].Face-2 == card.Pokers[3].Face{
				cardZero.Face = card.Pokers[2].Face-1
				card.Pokers = append(card.Pokers,cardZero)
				card.Level = 5
				return
			}
		}else {
			if card.Pokers[0].Face-2 == card.Pokers[1].Face{
				cardZero.Face = card.Pokers[0].Face-1
				card.Pokers = append(card.Pokers,cardZero)
				card.Level = 5
				return
			}
		}
	}
	if card.Feature.Continue == 2{
		if card.Pokers[0].Face -1 == card.Pokers[1].Face && card.Pokers[2].Face-1 == card.Pokers[3].Face && card.Pokers[1].Face-2 == card.Pokers[2].Face{
			cardZero.Face = card.Pokers[1].Face-1
			card.Pokers = append(card.Pokers,cardZero)
			card.Level = 5
			return
		}
	}

	//Level 4
	if card.Feature.SameCards == 1{
		cardZero.Face = card.Pokers[0].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 4
		return
	}
	if card.Feature.SameCards == 10{
		cardZero.Face = card.Pokers[1].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 4
		return
	}
	if card.Feature.SameCards == 100{
		cardZero.Face = card.Pokers[2].Face
		card.Pokers = append(card.Pokers,cardZero)
		card.Level = 4
		return
	}

	//Level3
	//Level2
	cardZero.Face = card.Pokers[0].Face
	card.Pokers = append(card.Pokers,cardZero)
	card.Level = 2
	return

}








