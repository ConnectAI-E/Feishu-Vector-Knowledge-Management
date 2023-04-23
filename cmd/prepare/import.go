package prepare

import (
	"github.com/k0kubun/pp/v3"
	"lark-vkm/internal/initialization"
	"lark-vkm/pkg/qdrantkit"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type Client struct { // Our example struct with a custom type (DateTime)
	Title         string `csv:"title"`
	Content       string `csv:"content"`
	TitleVector   string `csv:"title_vector"`
	ContentVector string `csv:"content_vector"`
}

var cmdCsv = &cobra.Command{
	Use:   "import",
	Short: "load vector data from csv file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Println(err)
			return
		}
		config := initialization.LoadConfig(cfg)

		client := qdrantkit.New(config.QdrantHost, config.QdrantCollection)

		file, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Println(err)
			return
		}
		_, err = os.Stat(file)
		if err != nil {
			log.Println(err)
			return
		}
		fp, err := os.Open(file)
		if err != nil {
			log.Println(err)
			return
		}

		count := 0

		//数据向量化
		batchSize := 3
		points := make([]qdrantkit.Point, 0, batchSize)
		var clients []*Client
		if err := gocsv.UnmarshalFile(fp,
			&clients); err != nil { // Load clients from file
			panic(err)
		}

		for _, row := range clients {
			pp.Println(clearField(row.Title), clearField(row.Content))
			//pp.Println(count)
			count++

			newPayload := map[string]interface{}{
				"Title":   clearField(row.Title),
				"Content": clearField(row.Content),
			}
			points = append(points, qdrantkit.Point{
				ID:      uuid.New().String(),
				Payload: newPayload,
				Vector:  stringToFloat64(row.ContentVector),
			})
			if count%batchSize == 0 {
				pr := qdrantkit.PointRequest{
					Points: points,
				}
				//存储
				err := client.CreatePoints(config.QdrantCollection, pr)
				if err != nil {
					log.Println(err)
					return
				}
				points = make([]qdrantkit.Point, 0, batchSize)
			}

		}

		log.Printf("loaded %d items, err: %v", count, err)
	},
}

// [2,3,4] -> [2.0, 3.0, 4.0]
func stringToFloat64(s string) []float64 {
	s = strings.Trim(s, "[]")
	s = strings.ReplaceAll(s, " ", "")
	split := strings.Split(s, ",")
	result := make([]float64, len(split))
	for i, v := range split {
		result[i], _ = strconv.ParseFloat(v, 64)
	}
	return result
}
