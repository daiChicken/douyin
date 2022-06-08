package tool

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
)

/**
 * @author  daijizai
 * @date  2022/6/8 18:39
 * @version  1.0
 * @description 敏感词过滤器
 */

const REPLACEMENT = "**"

type trieNode struct {
	isKeywordEnd bool               // 关键词结束标识
	subNodes     map[rune]*trieNode // 子节点(key是下级字符,value是下级节点)
}

//trieNode的构造函数
func newTrieNode() *trieNode {
	subNode := new(trieNode)
	subNode.isKeywordEnd = false
	subNode.subNodes = make(map[rune]*trieNode)
	return subNode
}

var rootNode *trieNode

func Init() error {
	rootNode = newTrieNode()

	//从文件中读取敏感词
	filepath := "./sensitive-words.txt"
	file, err := os.OpenFile(filepath, os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		//把敏感词添加到前缀树中
		addKeyWord(line)

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

//将一个敏感词添加到前缀树中
func addKeyWord(originalKeyword string) {
	dummyNode := rootNode

	keyword := []rune(originalKeyword)
	for i := 0; i < len(keyword); i++ {
		c := keyword[i]
		subNode := dummyNode.subNodes[c]

		if subNode == nil {
			//初始化子节点
			subNode = newTrieNode()
			dummyNode.subNodes[c] = subNode
		}

		//指向子节点，进入下一轮循环
		dummyNode = subNode

		//设置敏感词结束标识
		if i == len(keyword)-1 {
			dummyNode.isKeywordEnd = true
		}
	}
}

func Filter(originalText string) string {
	if originalText == "" {
		return ""
	}

	dummyNode := rootNode
	begin := 0
	position := 0
	var res bytes.Buffer

	text := []rune(originalText)
	for position < len(text) {
		c := text[position]

		//检查下级节点
		dummyNode = dummyNode.subNodes[c]
		if dummyNode == nil {
			// 以begin开头的字符串不是敏感词
			res.WriteString(string(text[begin]))
			// 进入下一个位置
			begin++
			position = begin
			// 重新指向trie根节点
			dummyNode = rootNode
		} else if dummyNode.isKeywordEnd {
			// 发现敏感词,将begin~position字符串替换掉
			res.WriteString(REPLACEMENT)
			// 进入下一个位置
			position++
			begin = position
			// 重新指向根节点
			dummyNode = rootNode
		} else {
			// 检查下一个字符
			position++
		}
	}
	res.WriteString(string(text[begin:]))

	return res.String()
}
