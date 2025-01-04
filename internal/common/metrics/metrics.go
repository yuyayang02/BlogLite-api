package metrics

type MockMetrics struct{}

func NewMockMetrics() MockMetrics {
	return MockMetrics{}
}

func (t MockMetrics) Inc(_ string, _ int) {

}
