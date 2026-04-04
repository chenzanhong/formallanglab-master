package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/chenzanhong/zlog"

	"github.com/chenzanhong/formallanglab-master/internal/domain/dto"
	"github.com/chenzanhong/formallanglab-master/internal/domain/model"
)

// PerformanceTestConfig 性能测试配置
type PerformanceTestConfig struct {
	BaseURL         string `json:"base_url"`
	Concurrency     int    `json:"concurrency"`
	RequestsPerUser int    `json:"requests_per_user"`
	Timeout         int    `json:"timeout"`
}

// AutomatonTestSuite 自动机性能测试套件
type AutomatonTestSuite struct {
	config *PerformanceTestConfig
	client *http.Client
}

// NewAutomatonTestSuite 创建新的测试套件
func NewAutomatonTestSuite(config *PerformanceTestConfig) *AutomatonTestSuite {
	return &AutomatonTestSuite{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
	}
}

// TestAutomatonValidate 测试自动机验证接口
func (s *AutomatonTestSuite) TestAutomatonValidate() error {
	zlog.Info("开始测试自动机验证接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				automaton := s.createRandomAutomaton()
				err := s.callAutomatonValidate(automaton)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("自动机验证请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestAutomatonCleanup 测试自动机清理接口
func (s *AutomatonTestSuite) TestAutomatonCleanup() error {
	zlog.Info("开始测试自动机清理接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				automaton := s.createRandomAutomaton()
				err := s.callAutomatonCleanup(automaton)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("自动机清理请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestDFAMinimize 测试 DFA 最小化接口
func (s *AutomatonTestSuite) TestDFAMinimize() error {
	zlog.Info("开始测试 DFA 最小化接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				dfa := s.createRandomDFA()
				err := s.callDFAMinimize(dfa)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("DFA 最小化请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestAutomatonStringRecognize 测试自动机字符串识别接口
func (s *AutomatonTestSuite) TestAutomatonStringRecognize() error {
	zlog.Info("开始测试自动机字符串识别接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				automaton := s.createRandomAutomaton()
				str := s.generateRandomString(automaton.Alphabet, 5)
				err := s.callAutomatonStringRecognize(automaton, str)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("自动机字符串识别请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestNFADeterminization 测试 NFA 确定化接口
func (s *AutomatonTestSuite) TestNFADeterminization() error {
	zlog.Info("开始测试 NFA 确定化接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				nfa := s.createRandomNFA()
				err := s.callNFADeterminization(nfa)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("NFA 确定化请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestAutomatonEquivalenceCheck 测试自动机等价性检查接口
func (s *AutomatonTestSuite) TestAutomatonEquivalenceCheck() error {
	zlog.Info("开始测试自动机等价性检查接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				automaton1 := s.createRandomAutomaton()
				automaton2 := s.createRandomAutomaton()
				err := s.callAutomatonEquivalenceCheck(automaton1, automaton2)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("自动机等价性检查请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// TestAutomatonGenerateExampleString 测试自动机生成示例字符串接口
func (s *AutomatonTestSuite) TestAutomatonGenerateExampleString() error {
	zlog.Info("开始测试自动机生成示例字符串接口")

	var wg sync.WaitGroup
	var successful, failed int64

	startTime := time.Now()

	for i := 0; i < s.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < s.config.RequestsPerUser; j++ {
				automaton := s.createRandomAutomaton()
				err := s.callAutomatonGenerateExampleString(automaton)
				if err != nil {
					atomic.AddInt64(&failed, 1)
					zlog.Errorw("自动机生成示例字符串请求失败", "error", err.Error())

					continue
				}

				atomic.AddInt64(&successful, 1)
			}
		}()
	}

	wg.Wait()
	totalTime := time.Since(startTime)
	totalRequests := int64(s.config.Concurrency * s.config.RequestsPerUser)

	fmt.Printf("测试完成: 总请求数=%d, 成功=%d, 失败=%d, 总耗时=%v\n",
		totalRequests, successful, failed, totalTime)

	return nil
}

// Helper Functions

func (s *AutomatonTestSuite) callAutomatonValidate(automaton model.Automaton) error {
	reqBody, err := json.Marshal(dto.AutomatonValidateRequest{Automaton: automaton})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/validate", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callAutomatonCleanup(automaton model.Automaton) error {
	reqBody, err := json.Marshal(dto.AutomatonCleanupRequest{Automaton: automaton})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/cleanup", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callDFAMinimize(automaton model.Automaton) error {
	reqBody, err := json.Marshal(dto.DFAMinimizeRequest{Automaton: automaton})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/minimize", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callAutomatonStringRecognize(automaton model.Automaton, str string) error {
	reqBody, err := json.Marshal(dto.AutomatonStringRecognizeRequest{
		Automaton: automaton,
		Str:       str,
	})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/recognize", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callNFADeterminization(nfa model.Automaton) error {
	reqBody, err := json.Marshal(dto.NFADeterminizationRequest{Automaton: nfa})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/nfatodfa", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callAutomatonEquivalenceCheck(automaton1, automaton2 model.Automaton) error {
	reqBody, err := json.Marshal(dto.AutomatonEquivalenceCheckRequest{
		Automaton1: automaton1,
		Automaton2: automaton2,
	})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/equivalence", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) callAutomatonGenerateExampleString(automaton model.Automaton) error {
	reqBody, err := json.Marshal(dto.AutomatonGenerateExampleStringRequest{Automaton: automaton})
	if err != nil {
		return err
	}

	resp, err := s.client.Post(
		fmt.Sprintf("%s/gdesign/master/automaton/generate", s.config.BaseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *AutomatonTestSuite) createRandomAutomaton() model.Automaton {
	states := []model.State{"q0", "q1", "q2"}
	alphabet := []model.Symbol{"a", "b"}

	transitions := []model.Transition{
		{
			FromState: "q0",
			Input:     "a",
			ToStates:  []model.State{"q1"},
		},
		{
			FromState: "q0",
			Input:     "b",
			ToStates:  []model.State{"q0"},
		},
		{
			FromState: "q1",
			Input:     "a",
			ToStates:  []model.State{"q2"},
		},
		{
			FromState: "q1",
			Input:     "b",
			ToStates:  []model.State{"q1"},
		},
		{
			FromState: "q2",
			Input:     "a",
			ToStates:  []model.State{"q2"},
		},
		{
			FromState: "q2",
			Input:     "b",
			ToStates:  []model.State{"q2"},
		},
	}

	return model.Automaton{
		States:          states,
		Alphabet:        alphabet,
		Transitions:     transitions,
		InitialState:    "q0",
		AcceptingStates: []model.State{"q2"},
		Type:            model.DFA,
	}
}

func (s *AutomatonTestSuite) createRandomDFA() model.Automaton {
	return s.createRandomAutomaton()
}

func (s *AutomatonTestSuite) createRandomNFA() model.Automaton {
	automaton := s.createRandomAutomaton()
	automaton.Type = model.NFA

	automaton.Transitions = append(automaton.Transitions, model.Transition{
		FromState: "q0",
		Input:     "a",
		ToStates:  []model.State{"q0", "q1"},
	})

	return automaton
}

func (s *AutomatonTestSuite) generateRandomString(alphabet []model.Symbol, length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += string(alphabet[i%len(alphabet)])
	}

	return result
}

// RunAllTests 运行所有测试
func (s *AutomatonTestSuite) RunAllTests() error {
	tests := []struct {
		name string
		fn   func() error
	}{
		{"AutomatonValidate", s.TestAutomatonValidate},
		{"AutomatonCleanup", s.TestAutomatonCleanup},
		{"DFAMinimize", s.TestDFAMinimize},
		{"AutomatonStringRecognize", s.TestAutomatonStringRecognize},
		{"NFADeterminization", s.TestNFADeterminization},
		{"AutomatonEquivalenceCheck", s.TestAutomatonEquivalenceCheck},
		{"AutomatonGenerateExampleString", s.TestAutomatonGenerateExampleString},
	}

	fmt.Println("\n" + repeatString("=", 80))
	fmt.Println("自动机性能测试")
	fmt.Println(repeatString("=", 80))

	for _, test := range tests {
		fmt.Printf("\n测试: %s\n", test.name)
		fmt.Println(repeatString("-", 40))

		err := test.fn()
		if err != nil {
			fmt.Printf("测试失败: %v\n", err)
			continue
		}
	}

	return nil
}

// RunTest 运行单个测试
func (s *AutomatonTestSuite) RunTest(testName string) error {
	tests := map[string]func() error{
		"AutomatonValidate":              s.TestAutomatonValidate,
		"AutomatonCleanup":               s.TestAutomatonCleanup,
		"DFAMinimize":                    s.TestDFAMinimize,
		"AutomatonStringRecognize":       s.TestAutomatonStringRecognize,
		"NFADeterminization":             s.TestNFADeterminization,
		"AutomatonEquivalenceCheck":      s.TestAutomatonEquivalenceCheck,
		"AutomatonGenerateExampleString": s.TestAutomatonGenerateExampleString,
	}

	fn, ok := tests[testName]
	if !ok {
		return fmt.Errorf("未知测试: %s", testName)
	}

	return fn()
}

// repeatString 重复字符串
func repeatString(s string, n int) string {
	result := ""
	for i := 0; i < n; i++ {
		result += s
	}

	return result
}
