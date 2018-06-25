package lib

type LimitParams struct {
	// 随机区间（生成记录随机的单次距离区间）
	RandDistance Float64Range
	// 限制区间（目标系统限制的单次距离区间）
	LimitSingleDistance Float64Range
	// 限制区间（目标系统限制的总距离区间）
	LimitTotalDistance Float64Range
	// 每条记录的时间区间
	MinuteDuration IntRange
}

type Float64Range struct {
	Min float64
	Max float64
}
type IntRange struct {
	Min int
	Max int
}
