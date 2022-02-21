/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

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
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/submariner-io/submariner-operator/internal/cli"
	"github.com/submariner-io/submariner-operator/internal/constants"
	"github.com/submariner-io/submariner-operator/internal/exit"
	"github.com/submariner-io/submariner-operator/pkg/client"
	"github.com/submariner-io/submariner-operator/pkg/uninstall"
)

var noPrompt bool

var uninstallCmd = &cobra.Command{
	Use:     "uninstall",
	Short:   "Uninstall Submariner and its components",
	Long:    "This command uninstalls Submariner and its components",
	PreRunE: restConfigProducer.CheckVersionMismatch,
	Run: func(cmd *cobra.Command, args []string) {
		status := cli.NewStatus()

		config, err := restConfigProducer.ForCluster()
		exit.OnError(status.Error(err, "Error creating REST config"))

		clientProducer, err := client.NewProducerFromRestConfig(config)
		exit.OnError(status.Error(err, "Error creating client producer"))

		if !noPrompt {
			result := false
			prompt := &survey.Confirm{
				Message: "This will completely uninstall Submariner from the cluster. Are you sure you want to continue?",
			}

			_ = survey.AskOne(prompt, &result)

			if !result {
				return
			}
		}

		exit.OnError(uninstall.All(clientProducer, namespace, status))
	},
}

func init() {
	uninstallCmd.Flags().StringVarP(&namespace, "namespace", "n", constants.SubmarinerNamespace,
		"namespace in which Submariner is installed")
	uninstallCmd.Flags().BoolVarP(&noPrompt, "yes", "y", false, "automatically answer yes to confirmation prompt")

	restConfigProducer.AddKubeConfigFlag(uninstallCmd)
	rootCmd.AddCommand(uninstallCmd)
}