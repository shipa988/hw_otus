/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/shipa988/hw_otus/hw12_13_14_15_calendar/cmd/calendar_scheduler/internal/data/app"
)

// runCmd represents the run command.
var runCmd = &cobra.Command{
	Use: "run",
	Short: "a `scheduler` process that periodically scans the main database,	choosing events to remind",
	Long: `
    -at startup, the process should connect to RabbitMQ and create all the necessary structures(topics, etc.) in it;
    - the process must select the events for which the notification should be sent (the event has a corresponding field), create for each Notification. serialize it (for example, to JSON) and add it to the queue;
    - the process should clean up old (more than 1 year ago) events`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.NewApp()
		if err := a.Run(cfg, debug); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
