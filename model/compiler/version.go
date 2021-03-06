/*
Copyright 2017 Ankyra

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

package compiler

import (
	"fmt"
	"github.com/ankyra/escape-core/parsers"
	"github.com/ankyra/escape-core/script"
	"strings"
)

func compileVersion(ctx *CompilerContext) error {
	version := ctx.Plan.Version
	_, err := script.ParseScript(version)
	if err != nil {
		return fmt.Errorf("Couldn't parse expression '%s' in version field: %s", version, err.Error())
	}
	str, err := ctx.RunScriptForCompileStep(version)
	if err != nil {
		return fmt.Errorf("Couldn't evaluate expression '%s' in version field: %s", version, err.Error())
	}
	version = strings.TrimSpace(str)
	if version == "auto" { // backwards compatibility
		version = "@"
	}
	if err := parsers.ValidateVersion(version); err != nil {
		return err
	}
	if strings.HasSuffix(version, "@") {
		prefix := version[:len(version)-1]
		nextVersion, err := ctx.Inventory.QueryNextVersion(ctx.Metadata.Project, ctx.Metadata.Name, prefix)
		if err != nil {
			return err
		}
		version = nextVersion
	}
	ctx.Metadata.Version = version
	return nil
}
