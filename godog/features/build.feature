Feature: Running the build phase 
    As a Infrastructure engineer
    In order to validate my changes
    I need to be able to build a project locally

    Scenario: Build the default Escape plan
      Given a new Escape plan called "my-release"
      When I build the application
      Then "_/my-release" version "0.0.0" is present in the build state

    Scenario: Build with default input variables
      Given a new Escape plan called "my-release"
        And input variable "input_variable" with default "test"
        And input variable "input_variable2" with default "test2"
      When I build the application
      Then "_/my-release" version "0.0.0" is present in the build state
       And its calculated input "input_variable" is set to "test"
       And its calculated input "input_variable2" is set to "test2"

    Scenario: Default input variables update on every build
      Given a new Escape plan called "my-release"
        And input variable "input_variable" with default "test"
       When I build the application
       Then "_/my-release" version "0.0.0" is present in the build state
        And its calculated input "input_variable" is set to "test"
        And input variable "input_variable" with default "new default baby"
       When I build the application
       Then "_/my-release" version "0.0.0" is present in the build state
        And its calculated input "input_variable" is set to "new default baby"

    Scenario: Build a template
      Given a new Escape plan called "my-release"
        And template "test.txt.tpl" containing "Hello World"
       When I build the application
       Then I should have a file "test.txt" with contents "Hello World"

    Scenario: Build a template that uses a variable
      Given a new Escape plan called "my-release"
        And input variable "input_variable" with default "test"
        And template "test.txt.tpl" containing "{{{version}}} {{input_variable}}"
       When I build the application
       Then I should have a file "test.txt" with contents "0.0.0 test"

    Scenario: Build skips template when not in build scope
      Given a new Escape plan called "my-release"
        And template "not_in_scope.txt.tpl" containing "Hello World" with "deploy" scope
       When I build the application
       Then I should not have a file "not_in_scope.txt"

    Scenario: Build does not skip template when in build scope
      Given a new Escape plan called "my-release"
        And template "inscope.txt.tpl" containing "Hello World" with "build" scope
       When I build the application
       Then I should have a file "inscope.txt" with contents "Hello World"

    Scenario: Default build with dependencies
      Given a new Escape plan called "my-dependency"
        And I release the application
        And a new Escape plan called "my-second-dependency"
        And it has "my-dependency-latest" as a dependency 
        And I release the application
        And a new Escape plan called "my-release"
        And it has "my-second-dependency-latest" as a dependency 
      When I build the application
      Then "_/my-release" version "0.0.0" is present in the build state
       And "_/my-second-dependency" version "0.0.0" is present in its deployment state

    Scenario: Default input variables update for dependencies on every build
      Given a new Escape plan called "my-input-dependency"
        And input variable "input_variable" with default "test"
        And I release the application
      Given a new Escape plan called "my-release"
        And it has "my-input-dependency-latest" as a dependency 
       When I build the application
       Then "_/my-release" version "0.0.0" is present in the build state
        And "_/my-input-dependency" version "0.0.0" is present in its deployment state
        And its calculated input "input_variable" is set to "test"
        And input variable "input_variable" with default "new default baby"
       When I build the application again
       Then "_/my-release" version "0.0.0" is present in the build state
        And "_/my-input-dependency" version "0.0.0" is present in its deployment state
        And its calculated input "PREVIOUS_input_variable" is set to "test"
        And its calculated input "input_variable" is set to "new default baby"
       
    Scenario: Build with provider (simplest case)
      Given a new Escape plan called "my-provider"
        And it provides "provider"
        And I release the application
      Given a new Escape plan called "my-consumer"
        And it consumes "provider"
       When I deploy "my-provider-v0.0.0"
        And I build the application
       Then "_/my-consumer" version "0.0.0" is present in the build state
        And "_/my-provider" is the provider for "provider"


    Scenario: Build with provider (output -> default linking)
      Given a new Escape plan called "my-provider2"
        And it provides "provider"
        And output variable "output_variable" with default "test"
        And I release the application
       When I deploy "_/my-provider2-v0.0.0"
       Then "_/my-provider2" version "0.0.0" is present in the deploy state
        And its calculated output "output_variable" is set to "test"

      Given a new Escape plan called "my-consumer"
        And it consumes "provider"
        And input variable "input_variable" with default "$provider.outputs.output_variable"
       When I build the application
       Then "_/my-consumer" version "0.0.0" is present in the build state
        And "_/my-provider2" is the provider for "provider"
        And its calculated input "input_variable" is set to "test"


    Scenario: Build with provider (output -> items linking)
      Given a new Escape plan called "my-provider3"
        And it provides "provider"
        And output variable "default_output" with default "a"
        And list output variable "output_variable" with default '["a", "b", "c"]'
        And I release the application
       When I deploy "_/my-provider3-v0.0.0"
       Then "_/my-provider3" version "0.0.0" is present in the deploy state
        And its calculated output "default_output" is set to "a"

      Given a new Escape plan called "my-consumer"
        And it consumes "provider"
        And input variable "input_variable" with default "$provider.outputs.default_output" and items "$provider.outputs.output_variable"
       When I build the application
       Then "_/my-consumer" version "0.0.0" is present in the build state
        And "_/my-provider3" is the provider for "provider"
        And its calculated input "input_variable" is set to "a"

    Scenario: If a consumer is not in scope no provider is necessary
      Given a new Escape plan called "my-consumer"
        And it consumes "provider" in the "deploy" scope
       When I build the application
       Then "_/my-consumer" version "0.0.0" is present in the build state

    Scenario: If a consumer is in scope it will be used
      Given a new Escape plan called "my-provider3"
        And it provides "provider"
        And I release the application
       When I deploy "_/my-provider3-v0.0.0"
       Then "_/my-provider3" version "0.0.0" is present in the deploy state
      Given a new Escape plan called "my-consumer"
        And it consumes "provider" in the "build" scope
       When I build the application
       Then "_/my-consumer" version "0.0.0" is present in the build state
        And "_/my-provider3" is the provider for "provider"
