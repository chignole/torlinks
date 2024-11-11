/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"chignole/torlinks/internal/files"
	"chignole/torlinks/internal/utils"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt" // Careful - v0.48 introduces some breaking changes
)

var inboxFiles = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "torlinksInboxFiles",
		Help: "Torrents currently in Torlinks inbox folder",
	},
)

var failedFiles = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "torlinksFailedFiles",
		Help: "Failed torrents currently in Torlinks inbox folder",
	},
)

var multiFiles = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "torlinksMultiFiles",
		Help: "Torrent files getting multiple matches and requiring manual treatment",
	},
)

var inboxSize = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "torlinksInboxSize",
		Help: "Size of files currently in inbox folder",
	},
)

// statsCmd represents the inbox command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Provides some useful stats about your inbox folder.",
	Long:  `Provides some useful stats about your inbox folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		source := viper.GetString("general.torrentsInbox")
		metricsFile := viper.GetString("metrics.file")
		metricsPing := viper.GetString("metrics.ping")

		// Creating new Prometheus registry
		reg := prometheus.NewRegistry()
		reg.MustRegister(inboxFiles)
		reg.MustRegister(failedFiles)
		reg.MustRegister(inboxSize)
		reg.MustRegister(multiFiles)

		inboxFiles.Set(float64(getInboxFiles(source)))
		inboxSize.Set(float64(getTotalSize(source)))
		failedFiles.Set(float64(getFailedFiles(source)))
		multiFiles.Set(float64(getMultiFiles(source)))

		// Collecting data
		mfs, err := reg.Gather()
		if err != nil {
			log.Println("[ERROR] Could not gather metrics:", err)
			return
		}

		// Opening output file
		file, err := os.Create(metricsFile)
		if err != nil {
			log.Println("[ERROR] Could not create metrics file:", err)
			return
		}
		defer file.Close()

		// Writing output file
		encoder := expfmt.NewEncoder(file, expfmt.Format(expfmt.FmtOpenMetrics_0_0_1))
		for _, mf := range mfs {
			err := encoder.Encode(mf)
			if err != nil {
				log.Println("[ERROR] Could not encode metric family:", err)
				return
			}
		}

		log.Println("[INFO] Metrics file saved at:", metricsFile)

		// Sending ping to healtcheck URL
		if metricsPing != "" {
			utils.PingHealthCheck(metricsPing)
		}

		// log.Printf("[STATS] Total files size....: %dGb\n", )
		// log.Printf("[STATS] Total failed files..: %d\n", totalFailedFiles)
		// log.Printf("[STATS] Total inbox files...: %d\n", totalInboxFiles)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func getTotalSize(source string) int64 {
	var totalSize int64
	torrents := files.Find(source, ".torrent")
	for t := range torrents {
		files := files.Parse(torrents[t].File)
		for _, f := range files {
			totalSize = totalSize + f.Size
		}
	}
	totalSize = totalSize / 1073741824
	return totalSize
}

func getInboxFiles(source string) int64 {
	inboxFiles := files.Find(source, ".torrent")
	totalInboxFiles := len(inboxFiles)
	return int64(totalInboxFiles)
}

func getFailedFiles(source string) int64 {
	failedFiles := files.Find(source, ".delete")
	totalFailedFiles := len(failedFiles)
	return int64(totalFailedFiles)
}

func getMultiFiles(source string) int64 {
	multiFiles := files.Find(source, ".multi")
	totalMultiFiles := len(multiFiles)
	return int64(totalMultiFiles)
}

func getBiggestFiles(source string) []string {
	torrents := files.Find(source, ".torrent")
	for t := range torrents {
		files := files.Parse(torrents[t].File)
		for _, f := range files {
			log.Println(f)
		}
	}
	return nil
}
