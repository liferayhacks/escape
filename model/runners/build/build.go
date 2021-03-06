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

package build

import (
	. "github.com/ankyra/escape/model/runners"
	"github.com/ankyra/escape/model/runners/deploy"
	"github.com/ankyra/escape-core/state"
)

var Stage = "build"

func NewPreBuildRunner() Runner {
	return NewPreScriptStepRunner(Stage, "pre_build", state.RunningPreStep, state.Failure)
}
func NewMainBuildRunner() Runner {
	return NewMainStepRunner(Stage, "build", state.RunningMainStep, state.Failure)
}
func NewPostBuildRunner() Runner {
	return NewPostScriptStepRunner(Stage, "post_build", state.RunningPostStep, state.Failure)
}
func NewTestRunner() Runner {
	return NewScriptRunner(Stage, "test", state.OK, state.TestFailure)
}

func NewBuildRunner() Runner {
	return NewCompoundRunner(
		NewDependencyRunner(deploy.Stage, Stage, deploy.NewDeployRunner, state.Failure),
		NewPreBuildRunner(),
		NewMainBuildRunner(),
		NewPostBuildRunner(),
		NewStatusCodeRunner(Stage, state.OK),
	)
}
