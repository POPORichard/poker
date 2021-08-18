package model

// Poker 每张牌
type Poker struct {
	Face int
	Color string
}

// Data 输入的Data结构
type Data struct {
	Alice string
	Bob string
	Result int
}

// InputData 从json获取的数据结构
type InputData struct {
	Matches []Data
}

// Turn 每一局的所以数据
type Turn struct {
	AliceHandCard []Poker
	AliceFeature Feature
	AliceLevel int
	BobHandCard []Poker
	BobFeature Feature
	BobLevel int
	Winner int
}
// Feature 每一组牌的特征值
type Feature struct {
	Continue int
	SameCards int
	Flush bool
}
