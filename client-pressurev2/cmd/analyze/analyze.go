package analyze

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/txchat/im-util/client-pressurev2/pkg/analyze"
	"github.com/txchat/im-util/client-pressurev2/pkg/filehelper"
	"github.com/txchat/im-util/client-pressurev2/pkg/logger"
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
	sysLogPath         string
)

func init() {
	Cmd.Flags().StringVarP(&outputPath, "out", "o", "./analyze_output.txt", "")
	Cmd.Flags().StringVarP(&pressureOutputPath, "in", "i", "./pressure_output.txt", "")
	Cmd.Flags().StringVarP(&sysLogPath, "syslog", "", "./analyze_sys_log.txt", "")
}

func do(cmd *cobra.Command, args []string) {
	sysFd, sysCloser, err := filehelper.WriteFile(sysLogPath)
	if err != nil {
		panic(err)
	}
	defer sysCloser()
	msgFd, msgCloser, err := filehelper.WriteFile(outputPath)
	if err != nil {
		panic(err)
	}
	defer msgCloser()

	//load users
	log := logger.NewSysLog(sysFd)

	log.Info().Msg("start analyze")
	fr := filehelper.NewFileReader()
	err = fr.ReadFile(pressureOutputPath)
	if err != nil {
		log.Error().Err(err).Msg("ReadFile error")
		return
	}
	lines := fr.GetAllLines()
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
	pt := analyze.NewPrinter(analyze.GetTransmitMsgStatic(), logger.NewMsgLog(msgFd))
	pt.PrintAllLevel1()

	log.Info().Msg(fmt.Sprintf("failed count: %d", aze.FailedCount()))
	log.Info().Msg(fmt.Sprintf("message tranport success count: %d -- failed count: %d", pt.GetSuccessCount(), pt.GetFailedCount()))
	log.Info().Msg("done")
}
