package prepare

import (
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"lark-vkm/internal/initialization"
	"lark-vkm/pkg/openai"
	"log"
	"os"
	"strings"
)

var tempFile = "__temp.csv"

var cmdAnalyze = &cobra.Command{
	Use:   "analyze",
	Short: "analyze vector data from csv file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cmd.Flags().GetString("config")
		config := initialization.LoadConfig(cfg)
		gpt := openai.NewChatGPT(*config)

		inputFile, err := cmd.Flags().GetString("file")
		if err != nil {
			panic(err)
		}

		f, err := os.OpenFile(inputFile, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer closeFile(f)

		reader := csv.NewReader(f)
		reader.FieldsPerRecord = -1

		header, err := readHeader(reader)
		if err != nil {
			log.Fatal(err)
		}

		if !headerContainsFields(header, "title", "content") {
			log.Fatal("header must contain title and content")
		}

		newHeader := addMissingFields(header, "title_vector", "content_vector")

		if len(newHeader) > len(header) {
			err = updateHeader(reader, newHeader)
			if err != nil {
				panic(err)
			}
		}

		fTemp, err := os.OpenFile(tempFile, os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}
		defer closeFile(fTemp)
		defer func() {
			err := os.Remove(tempFile)
			if err != nil {
				panic(err)
			}
		}()
		reader = csv.NewReader(fTemp)
		reader.FieldsPerRecord = -1

		// 读取所有记录
		records, err := reader.ReadAll()
		if err != nil {
			panic(err)
		}
		// 获取 title 和 content 的字段索引
		titleIndex, contentIndex := getFieldIndexes(newHeader, "title", "content")
		// 获取 title_vector 和 content_vector 的字段索引
		titleVectorIndex, contentVectorIndex := getFieldIndexes(newHeader, "title_vector", "content_vector")
		// 遍历所有记录
		finalRecords := make([][]string, 0)
		for i := 1; i < len(records); i++ {
			//最少四位元素
			record := make([]string, len(newHeader))

			copy(record, records[i])
			// 如果 title_vector 字段为空，调用 OpenAI API 查询
			if isEmpty(record[titleVectorIndex]) {
				titleVector, err := getEmbedding(
					record[titleIndex], gpt)
				if err != nil {
					panic(err)
				}
				record[titleVectorIndex] = titleVector
			}

			// 如果 content_vector 字段为空，调用 OpenAI API 查询
			if isEmpty(record[contentVectorIndex]) {
				contentVector, err := getEmbedding(
					record[contentIndex], gpt)
				if err != nil {
					panic(err)
				}
				record[contentVectorIndex] = contentVector
			}
			finalRecords = append(finalRecords, record)
		}
		// 写入新的 CSV 文件
		outputFile, err := cmd.Flags().GetString("output")
		if err != nil {
			panic(err)
		}

		output, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		defer closeFile(output)

		writer := csv.NewWriter(output)
		err = writer.Write(newHeader)
		if err != nil {
			panic(err)
		}

		err = writer.WriteAll(finalRecords)
		if err != nil {
			panic(err)
		}
		fmt.Println("Successfully analyzed data and wrote"+
			" results to", outputFile)
	},
}

func closeFile(f *os.File) {
	err := f.Close()
	if err != nil {

	}
}

func readHeader(reader *csv.Reader) ([]string, error) {
	header, err := reader.Read()
	if err != nil {
		return nil, err
	}
	for i, field := range header {
		header[i] = clearField(field)
	}
	return header, nil
}

func headerContainsFields(header []string, fields ...string) bool {
	for _, field := range fields {
		if !contains(header, field) {
			return false
		}
	}
	return true
}

func addMissingFields(header []string, fields ...string) []string {
	newHeader := make([]string, len(header))
	copy(newHeader, header)
	for _, field := range fields {
		if !contains(header, field) {
			newHeader = append(newHeader, field)
		}
	}
	fmt.Println("new header:", newHeader)
	return newHeader
}

func updateHeader(reader *csv.Reader, newHeader []string) error {
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}
	records = append([][]string{newHeader}, records...)
	outputFile, err := os.Create(tempFile)
	defer closeFile(outputFile)
	if err != nil {
		return err
	}
	//fmt.Println("record:", records)
	//fmt.Println("outputFile:", tempFile)
	writer := csv.NewWriter(outputFile)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}
	return nil
}

func contains(slice []string, s string) bool {
	for _, str := range slice {
		str = clearField(str)
		if str == s {
			return true
		}
	}
	return false
}

func clearField(str string) string {
	str = strings.TrimSpace(str)
	str = strings.Trim(str, "\xef\xbb\xbf")
	return str
}

func getFieldIndexes(header []string, fields ...string) (int, int) {
	var titleIndex, contentIndex int
	for i, field := range header {
		switch clearField(field) {
		case fields[0]:
			titleIndex = i
		case fields[1]:
			contentIndex = i
		}
	}
	return titleIndex, contentIndex
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
func getEmbedding(text string, gpt *openai.ChatGPT) (string,
	error) {
	response, err := gpt.Embeddings(text)
	if err != nil {
		return "", err
	}
	fmt.Println("正在获取向量:" + text)
	return toString(response.Data[0].Embedding), nil
	//return "xxx", nil
}

func toString(vector []float64) string {
	//[2,3]
	return strings.Trim(strings.Replace(fmt.Sprint(vector), " ", ",", -1), "[]")
}
