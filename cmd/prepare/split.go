package prepare

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var cmdSplit = &cobra.Command{
	Use:   "split",
	Short: "split text into question and answer and write to csv file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 从命令行参数获取输入文件名
		inputFile, err := cmd.Flags().GetString("file")
		outputFile, err := cmd.Flags().GetString("output")
		// 打开输入文件
		file, err := os.Open(inputFile)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		// 创建CSV文件
		fileCsv, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		defer func(fileCsv *os.File) {
			err := fileCsv.Close()
			if err != nil {

			}
		}(fileCsv)

		// utf-8
		writer := csv.NewWriter(transform.NewWriter(fileCsv,
			unicode.UTF8BOM.NewEncoder()))
		defer writer.Flush()

		// 写入表头 title, content
		err = writer.Write([]string{"title", "content"})

		// 处理输入文件中的每一行文本
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			text := scanner.Text()

			// 将文本分隔为问题和答案
			question, answer, err := splitText(text)
			if err != nil {
				log.Printf("Error splitting text %s: %v", text, err)
				continue
			}

			//删除空格
			question = strings.Replace(question, " ", "", -1)
			answer = strings.Replace(answer, " ", "", -1)

			// 将问题和答案写入CSV文件
			err = writer.Write([]string{question, answer})
			if err != nil {
				log.Printf("Error writing CSV data for text %s: %v", text, err)
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		fmt.Printf("输入文件 %s 中的文本已分割并写入CSV文件 %s 中。\n", inputFile, outputFile)
	},
}

// 将文本分隔为问题和答案
func splitText(text string) (question, answer string, err error) {
	// 中文问号 或者 英文问号
	index := strings.IndexAny(text, "? ？")

	if index == -1 {
		err = fmt.Errorf("Question mark not found in text: %s", text)

		return
	}

	// 将问题和答案分隔为两个字段
	question = strings.TrimSpace(text[:index+1])
	answer = strings.TrimSpace(text[index+1:])

	return
}
