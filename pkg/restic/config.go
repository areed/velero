/*
Copyright 2018 the Heptio Ark contributors.

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

package restic

import (
	"fmt"
	"strings"

	arkv1api "github.com/heptio/ark/pkg/apis/ark/v1"
)

// getRepoPrefix returns the prefix of the value of the --repo flag for
// restic commands, i.e. everything except the "/<repo-name>".
func getRepoPrefix(config arkv1api.ObjectStorageProviderConfig) string {
	if BackendType(config.Name) == AWSBackend {
		var url string
		switch {
		// non-AWS, S3-compatible object store
		case config.Config["s3Url"] != "":
			url = config.Config["s3Url"]
		default:
			url = "s3.amazonaws.com"
		}

		return fmt.Sprintf("s3:%s/%s", url, config.ResticLocation)
	}

	var (
		parts        = strings.SplitN(config.ResticLocation, "/", 2)
		bucket, path string
	)

	if len(parts) >= 1 {
		bucket = parts[0]
	}
	if len(parts) >= 2 {
		path = parts[1]
	}

	var prefix string
	switch BackendType(config.Name) {
	case AzureBackend:
		prefix = "azure"
	case GCPBackend:
		prefix = "gs"
	}

	return fmt.Sprintf("%s:%s:/%s", prefix, bucket, path)
}

// GetRepoIdentifier returns the string to be used as the value of the --repo flag in
// restic commands for the given repository.
func GetRepoIdentifier(config arkv1api.ObjectStorageProviderConfig, name string) string {
	prefix := getRepoPrefix(config)

	return fmt.Sprintf("%s/%s", strings.TrimSuffix(prefix, "/"), name)
}
