package model


// =================== 自动机识别字符串的过程记录 ===================
type RecognitionStep struct {
	Step      int    `json:"step"`
	State     State  `json:"state"`
	Input     Symbol `json:"input"`
	NextState State  `json:"nextState"`
}

type RecognitionResult struct {
	IsAccepted bool              `json:"isAccepted"`
	Steps      []RecognitionStep `json:"steps"`
}


// ===================	  DFA最小化的过程记录    ===================
// MinimizationStep 表示DFA最小化过程中的一个步骤
type MinimizationStep struct {
	Step         int                  `json:"step"`      // 步骤编号
	PartitionMap []map[State]struct{} `json:"-"`         // 计算时使用，响应时转为Partition
	Partition    [][]State            `json:"partition"` // 当前的状态划分
	Actions      []string             `json:"actions"`   // 执行的操作描述
}

// 把PartitionMap数据复制到Partition
func (m *MinimizationStep) PMap2P() {
	m.Partition = make([][]State, len(m.PartitionMap))
	for i, p := range m.PartitionMap {
		m.Partition[i] = make([]State, 0, len(p))
		for s := range p {
			m.Partition[i] = append(m.Partition[i], s)
		}
	}
}

// MinimizationProcess 表示DFA最小化的完整过程记录
type MinimizationProcess struct {
	Steps []MinimizationStep `json:"steps"` // 最小化步骤列表
}