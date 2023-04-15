package prepare

import (
	"lark-vkm/internal/initialization"
	"lark-vkm/pkg/qdrantkit"
	"log"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "prepare <action>",
	Short: "Prepare vector data",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	var cmdRun = &cobra.Command{
		Use:   "csv",
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
			last := 0

			//数据向量化
			batchSize := 100
			points := make([]qdrantkit.Point, 0, batchSize)

			err = gocsv.UnmarshalToCallback(fp, func(row Row) {
				count++
				points = append(points, qdrantkit.Point{
					ID:      uuid.New().String(),
					Payload: row,
					Vector:  row.ContentVector,
				})
				last = row.Id

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
					return
				}
			})

			log.Printf("loaded %d items, last: %d, err: %v", count, last, err)
		},
	}

	cmdRun.Flags().StringP("file", "f", "./data.csv", "csv file path.")

	cmd.AddCommand(cmdRun)
}

func Register(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(cmd)

	return nil
}

type Row struct {
	Id            int       `csv:"id"`
	Url           string    `csv:"url"`
	Title         string    `csv:"title"`
	Text          string    `csv:"text"`
	TitleVector   []float64 `csv:"title_vector"`
	ContentVector []float64 `csv:"content_vector"`
	VectorId      int       `csv:"vector_id"`
}
