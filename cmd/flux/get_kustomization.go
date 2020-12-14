/*
Copyright 2020 The Flux authors

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

package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/fluxcd/flux2/internal/utils"
	"github.com/fluxcd/pkg/apis/meta"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta1"
	"github.com/spf13/cobra"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var getKsCmd = &cobra.Command{
	Use:     "kustomizations",
	Aliases: []string{"ks"},
	Short:   "Get Kustomization statuses",
	Long:    "The get kustomizations command prints the statuses of the resources.",
	Example: `  # List all kustomizations and their status
  flux get kustomizations
`,
	RunE: getKsCmdRun,
}

func init() {
	getCmd.AddCommand(getKsCmd)
}

func getKsCmdRun(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	kubeClient, err := utils.KubeClient(kubeconfig, kubecontext)
	if err != nil {
		return err
	}

	var listOpts []client.ListOption
	if !allNamespaces {
		listOpts = append(listOpts, client.InNamespace(namespace))
	}
	var list kustomizev1.KustomizationList
	err = kubeClient.List(ctx, &list, listOpts...)
	if err != nil {
		return err
	}

	if len(list.Items) == 0 {
		logger.Failuref("no kustomizations found in %s namespace", namespace)
		return nil
	}

	header := []string{"Name", "Ready", "Message", "Revision", "Suspended"}
	if allNamespaces {
		header = append([]string{"Namespace"}, header...)
	}
	var rows [][]string
	for _, kustomization := range list.Items {
		row := []string{}
		if c := apimeta.FindStatusCondition(kustomization.Status.Conditions, meta.ReadyCondition); c != nil {
			row = []string{
				kustomization.GetName(),
				string(c.Status),
				c.Message,
				kustomization.Status.LastAppliedRevision,
				strings.Title(strconv.FormatBool(kustomization.Spec.Suspend)),
			}
		} else {
			row = []string{
				kustomization.GetName(),
				string(metav1.ConditionFalse),
				"waiting to be reconciled",
				kustomization.Status.LastAppliedRevision,
				strings.Title(strconv.FormatBool(kustomization.Spec.Suspend)),
			}
		}
		if allNamespaces {
			row = append([]string{kustomization.Namespace}, row...)
		}
		rows = append(rows, row)
	}
	utils.PrintTable(os.Stdout, header, rows)
	return nil
}