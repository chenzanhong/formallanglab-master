package service

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type QACache struct {
	data map[string]string // map[问题]答案
}

func NewQACache() *QACache {
	return &QACache{
		data: make(map[string]string),
	}
}

// aliasPattern 用于匹配 Markdown 文件中的 aliases 元数据注释行，
// 格式为 <!-- aliases: ["问题1","问题2"] -->
var aliasPattern = regexp.MustCompile(`<!--\s*aliases:\s*(\[.*?\])\s*-->`)

// LoadCache 从指定目录递归加载所有 .md 文件，解析其中的 aliases 元数据，并将“问题→答案”映射写入缓存。
// 每个文件的答案内容默认去掉首行的 aliases 注释行。
func (c *QACache) LoadCache(dirPath string) error {
	// 使用 filepath.WalkDir 遍历目录，仅处理 .md 文件
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// 跳过子目录与非 Markdown 文件
		if d.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		// 读取文件全部内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %w", path, err)
		}

		text := string(content)
		// 提取 aliases 元数据
		matches := aliasPattern.FindStringSubmatch(text)
		// fmt.Printf("%+v\n", matches)
		if len(matches) < 2 {
			// 没有找到 aliases，跳过或记录警告
			fmt.Printf("警告：文件 %s 缺少 aliases 元数据\n", path)
			return nil
		}

		// 将 JSON 数组字符串解析为字符串切片
		var aliases []string
		if err := json.Unmarshal([]byte(matches[1]), &aliases); err != nil {
			return fmt.Errorf("解析 aliases 失败 in %s: %w", path, err)
		}

		// 提取 answer：去掉第一行（即注释行）
		scanner := bufio.NewScanner(strings.NewReader(text))
		var answerLines []string
		firstLine := true
		for scanner.Scan() {
			line := scanner.Text()
			if firstLine && aliasPattern.MatchString(line) {
				firstLine = false
				continue // 跳过元数据行
			}
			answerLines = append(answerLines, line)
		}
		answer := strings.Join(answerLines, "\n")

		// 存入缓存
		for _, alias := range aliases {
			alias = strings.TrimSpace(alias)
			if alias != "" {
				c.data[alias] = answer
			}
		}
		return nil
	})
}

func (c *QACache) Get(key string) string {
	key = c.normalize(key)
	fmt.Println("规格化后的提问：", key)
	return c.data[key]
}

var spaceRegex = regexp.MustCompile(`\s+`)

// 标准化提问key
// 1. 去除前后空白
// 2. 去除行内所有空格
// 3. 去除末尾中英文句子常见的结尾标点符号
// 考虑到一些名称，不改变大小写
func (c *QACache) normalize(key string) string {
	// 去除前后空白
	key = strings.TrimSpace(key)
	// 压缩连续空格
	key = spaceRegex.ReplaceAllString(key, "")
	// 只去除真正可能出现在句尾的结束标点
	key = strings.TrimRightFunc(key, func(r rune) bool {
		return r == '。' || r == '！' || r == '？' ||
			r == '.' || r == '!' || r == '?'
	})
	return key
}

func (c *QACache) Set(key, value string) {
	c.data[key] = value
}

func (c *QACache) Delete(key string) {
	delete(c.data, key)
}
