package prepare

import (
	"github.com/spf13/cobra"
)

type Row struct {
	Id            int       `csv:"id"`
	Url           string    `csv:"url"`
	Title         string    `csv:"title"`
	Text          string    `csv:"text"`
	TitleVector   []float64 `csv:"title_vector"`
	ContentVector []float64 `csv:"content_vector"`
	VectorId      int       `csv:"vector_id"`
}

var cmd = &cobra.Command{
	Use:   "prepare <action>",
	Short: "Prepare vector data",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	cmdCsv.Flags().StringP("file", "f", "./data.csv", "csv file path.")

	cmd.AddCommand(cmdCsv)
}

func Register(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(cmd)

	return nil
}
