// Copyright 2023 StreamNative, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package notifications

import (
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/streamnative/oxia/cmd/client/common"
)

var Cmd = &cobra.Command{
	Use:   "notifications",
	Short: "Get notifications stream",
	Long:  `Follow the change notifications stream`,
	Args:  cobra.NoArgs,
	RunE:  exec,
}

func exec(_ *cobra.Command, _ []string) error {
	client, err := common.Config.NewClient()
	if err != nil {
		return err
	}

	defer client.Close()

	notifications, err := client.GetNotifications()
	if err != nil {
		return err
	}

	defer notifications.Close()

	for notification := range notifications.Ch() {
		slog.Info(
			"",
			slog.Any("type", notification.Type),
			slog.String("key", notification.Key),
			slog.Int64("version-id", notification.VersionId),
		)
	}

	return nil
}
