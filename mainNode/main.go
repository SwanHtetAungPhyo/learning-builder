package main

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/learning/mainNode/cmd"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/config"
	"github.com/SwanHtetAungPhyo/learning/mainNode/internal/handler"

	"github.com/sirupsen/logrus"
	"runtime"
)

var configuration *config.Config

func init() {
	configuration = config.NewConfig()
}
func main() {
	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filePlusLine := fmt.Sprintf("%s:%d", f.File, f.Line)
			return f.Function, filePlusLine
		},
		PrettyPrint: true,
	}
	logger.Println(configuration.Validators)
	server := cmd.NewServer(logger, handler.NewImpl(logger, configuration.Validators))

	server.Start()
	//mainNode := avl.NewNode(configuration.Validators[0])
	//for _, nodes := range configuration.Validators {
	//	logger.Println(nodes)
	//	mainNode.Insert(nodes)
	//}
	//
	//logger.Println(mainNode.CheckConsensus(), " Coseensus")
	//logger.Println(mainNode.GetHighestValidator())

}
