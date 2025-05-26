package cmd

import (
	"fmt"
	"log"

	"github.com/kp/pager/databases/kafka"
	"github.com/spf13/cobra"
)

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Manage Kafka topics and configurations",
}

var createTopicsCmd = &cobra.Command{
	Use:   "create-topics",
	Short: "Create default Kafka topics",
	Run: func(cmd *cobra.Command, args []string) {
		brokers, err := cmd.Flags().GetStringSlice("brokers")
		if err != nil {
			log.Fatalf("Error getting brokers: %v", err)
		}

		if len(brokers) == 0 {
			log.Fatal("At least one Kafka broker must be specified")
		}

		err = kafka.RunMigrations(brokers)
		if err != nil {
			log.Fatalf("Failed to create topics: %v", err)
		}

		fmt.Println("Successfully created Kafka topics")
	},
}

func init() {
	createTopicsCmd.Flags().StringSlice("brokers", []string{"host.docker.internal:9092"}, "Kafka broker addresses")
	kafkaCmd.AddCommand(createTopicsCmd)
	rootCmd.AddCommand(kafkaCmd)
}
