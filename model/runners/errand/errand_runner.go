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

package errand

import (
	. "github.com/ankyra/escape-client/model/interfaces"
	"github.com/ankyra/escape-client/model/runners"
)

func NewErrandRunner(errand Errand) Runner {
	return runners.NewRunner(func(ctx RunnerContext) error {
		step := runners.NewScriptStep(ctx, "deploy", errand.GetName(), true)
		step.ScriptPath = errand.GetScript()
		return step.Run(ctx)
	})
}