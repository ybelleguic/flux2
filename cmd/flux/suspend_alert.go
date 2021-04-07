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
	notificationv1 "github.com/fluxcd/notification-controller/api/v1beta1"
	"github.com/spf13/cobra"
)

var suspendAlertCmd = &cobra.Command{
	Use:   "alert [name]",
	Short: "Suspend reconciliation of Alert",
	Long:  "The suspend command disables the reconciliation of a Alert resource.",
	Example: `  # Suspend reconciliation for an existing Alert
  flux suspend alert main`,
	RunE: suspendCommand{
		apiType: alertType,
		object:  &alertAdapter{&notificationv1.Alert{}},
	}.run,
}

func init() {
	suspendCmd.AddCommand(suspendAlertCmd)
}

func (obj alertAdapter) isSuspended() bool {
	return obj.Alert.Spec.Suspend
}

func (obj alertAdapter) setSuspended() {
	obj.Alert.Spec.Suspend = true
}