package scriptModeling

type Histogram []int
type Distribution []float64

func newHistogram(topics int) Histogram {
  return make([]int, topics)
}