package model

type Poker struct {
	Face int
	Color string
}


type Data struct {
	Alice string
	Bob string
	Result int
}

type InputData struct {
	Matches []Data
}

type Turn struct {
	AliceHandCard []Poker
	AliceFeature Feature
	AliceLevel int
	BobHandCard []Poker
	BobFeature Feature
	BobLevel int
	Winner int
}

type Feature struct {
	Continue int
	SameCards int
	Flush bool
}
