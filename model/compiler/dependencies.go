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
	"github.com/ankyra/escape-core"
)

func compileDependencies(ctx *CompilerContext) error {
	result := []string{}
	for _, depend := range ctx.Plan.Depends {
		dep, err := core.NewDependencyFromString(depend)
		if err != nil {
			return err
		}
		metadata, err := resolveVersion(ctx, dep)
		if err != nil {
			return err
		}
		for _, consume := range metadata.GetConsumes() {
			ctx.Metadata.AddConsumes(consume)
		}
		for _, input := range metadata.GetInputs() {
			if !input.HasDefault() {
				ctx.Metadata.AddInputVariable(input)
			}
		}
		ctx.VariableCtx[dep.Name] = metadata
		ctx.Metadata.SetVariableInContext(dep.Name, metadata.GetQualifiedReleaseId())

		if dep.VariableName != "" {
			ctx.VariableCtx[dep.VariableName] = metadata
			ctx.Metadata.SetVariableInContext(dep.VariableName, metadata.GetQualifiedReleaseId())
		}

		result = append(result, dep.GetQualifiedReleaseId())
	}
	ctx.Metadata.SetDependencies(result)
	return nil
}

func resolveVersion(ctx *CompilerContext, d *core.Dependency) (*core.ReleaseMetadata, error) {
	if d.NeedsResolving() {
		if ctx.ReleaseQuery == nil {
			return nil, fmt.Errorf("Missing release query function")
		}
		metadata, err := ctx.ReleaseQuery(d)
		if err != nil {
			return nil, err
		}
		d.Version = metadata.Version
	}
	if ctx.DependencyFetcher == nil {
		return nil, fmt.Errorf("Missing dependency fetcher")
	}
	metadata, err := ctx.DependencyFetcher(d.GetQualifiedReleaseId())
	if err != nil {
		return nil, err
	}
	return metadata, nil
}
