package analyze

import (
	"fmt"
	"github.com/spf13/cobra"
	xlog "github.com/txchat/im-util/pkg/log"
	"github.com/txchat/im-util/pkg/util"
	"github.com/txchat/im-util/pressure/pkg/analyze"
	"os"
)

var Cmd = &cobra.Command{
	Use:     "ana",
	Short:   "",
	Long:    "",
	Example: "",
	Run:     do,
}

var (
	pressureOutputPath string
	outputPath         string
)

func init() {
	Cmd.Flags().StringVarP(&outputPath, "out", "o", "./analyze_output.txt", "")
	Cmd.Flags().StringVarP(&pressureOutputPath, "in", "i", "./pressure_output.txt", "")
}

func do(cmd *cobra.Command, args []string) {
	// 打开文件
	fd, closer, err := util.WriteFile(outputPath)
	if err != nil {
		panic(err)
	}
	defer closer()
	outLog := xlog.NewLogger(fd)
	log := xlog.NewLogger(os.Stdout)

	//load users
	log.Info().Msg("start analyze")
	lines, err := util.ReadAllLines(pressureOutputPath)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}
	log.Info().Msg(fmt.Sprintf("source lines: %d", len(lines)))
	aze := analyze.NewAnalyzeStore(lines)
	err = aze.LoadAll()
	if err != nil {
		panic(err)
	}
	err = aze.Start()
	if err != nil {
		panic(err)
	}

	log.Info().Msg("start print out")
	log.Info().Msg(fmt.Sprintf("transmit msg count: %d", analyze.GetTransmitMsgStatic().GetAllTransmitMsgCount()))
	pt := analyze.NewPrinter(analyze.GetTransmitMsgStatic(), outLog)
	pt.PrintAllLevel1()

	log.Info().Msg(fmt.Sprintf("failed count: %d", aze.FailedCount()))
	log.Info().Msg(fmt.Sprintf("message tranport success count: %d -- failed count: %d", pt.GetSuccessCount(), pt.GetFailedCount()))
	log.Info().Msg("done")
}
