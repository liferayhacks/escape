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

package core

import (
	. "gopkg.in/check.v1"
)

func (s *metadataSuite) Test_DependencyConfig_Mapping_is_set(c *C) {
	metadata := NewReleaseMetadata("name", "1.0")
	dep := NewDependencyConfig("my-dependency")
	dep.Mapping = nil
	c.Assert(dep.Mapping, IsNil)
	c.Assert(dep.Validate(metadata), IsNil)
	c.Assert(dep.Mapping, Not(IsNil))
	c.Assert(dep.Mapping, HasLen, 0)
	c.Assert(dep.Scopes, DeepEquals, []string{"build", "deploy"})
	c.Assert(dep.Consumes, DeepEquals, map[string]string{})
}

func (s *metadataSuite) Test_NewDependencyConfigFromMap(c *C) {
	dep, err := NewDependencyConfigFromMap(map[interface{}]interface{}{
		"release_id": "test-latest",
		"mapping": map[interface{}]interface{}{
			"input_variable1": "test",
		},
		"scopes": []interface{}{"build"},
		"consumes": map[interface{}]interface{}{
			"test": "whatver",
		},
	})
	c.Assert(err, IsNil)
	c.Assert(dep.ReleaseId, Equals, "test-latest")
	c.Assert(dep.Mapping, Not(IsNil))
	c.Assert(dep.Mapping, HasLen, 1)
	c.Assert(dep.Mapping["input_variable1"], Equals, "test")
	c.Assert(dep.Scopes, DeepEquals, []string{"build"})
	c.Assert(dep.Consumes, DeepEquals, map[string]string{"test": "whatver"})
}
