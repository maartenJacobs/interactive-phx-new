package main

import (
	"reflect"
	"testing"
)

func TestBuildCommandWithoutFeatures(t *testing.T) {
	project := project{name: "foo", features: []string{}, database: "postgres", adapter: "bandit"}
	expectedCommand := []string{
		"phx.new", "--install",
		"--database", "postgres", "--adapter", "bandit",
		"--no-esbuild", "--no-ecto", "--no-gettext", "--no-html",
		"--no-dashboard", "--no-live", "--no-mailer", "--no-tailwind",
		"foo",
	}
	if command := project.buildCommand(); !reflect.DeepEqual(command, expectedCommand) {
		t.Errorf("Expected %v, got %v", expectedCommand, command)
	}
}

func TestBuildCommandWithFeatures(t *testing.T) {
	project := project{name: "foo", features: []string{"Tailwind", "ESBuild"}, database: "postgres", adapter: "bandit"}
	expectedCommand := []string{
		"phx.new", "--install",
		"--database", "postgres", "--adapter", "bandit",
		"--no-ecto", "--no-gettext", "--no-html",
		"--no-dashboard", "--no-live", "--no-mailer",
		"foo",
	}
	if command := project.buildCommand(); !reflect.DeepEqual(command, expectedCommand) {
		t.Errorf("Expected %v, got %v", expectedCommand, command)
	}
}
