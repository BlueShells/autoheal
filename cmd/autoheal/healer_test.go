/*
Copyright (c) 2018 Red Hat, Inc.

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
	"path/filepath"
	"testing"

	alertmanager "github.com/openshift/autoheal/pkg/alertmanager"
	monitoring "github.com/openshift/autoheal/pkg/apis/monitoring/v1alpha1"
)

func TestRuleWithExactLabel(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithExactAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithMatchingLabel(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithMatchingAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithNonMatchingLabel(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithNonMatchingAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithTwoMatchingLabels(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel":   "my.*",
			"yourlabel": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel":   "myvalue",
			"yourlabel": "yourvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithTwoMatchingAnnotations(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation":   "my.*",
			"yourannotation": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation":   "myvalue",
			"yourannotation": "yourvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithMatchingAndNotMatchingLabels(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel":   "my.*",
			"yourlabel": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel":   "myvalue",
			"yourlabel": "ugly",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithMatchingAndNotMatchingAnnotations(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation":   "my.*",
			"yourannotation": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation":   "myvalue",
			"yourannotation": "ugly",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithMatchingLabelAndAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "my.*",
		},
		Annotations: map[string]string{
			"myannotation": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithMatchingLabelAndNonMatchingAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "my.*",
		},
		Annotations: map[string]string{
			"myannotation": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
		Annotations: map[string]string{
			"myannotation": "ugly",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithNonMatchingLabelAndMatchingAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "my.*",
		},
		Annotations: map[string]string{
			"myannotation": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "ugly",
		},
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithNonMatchingAndIgnoredLabels(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel":   "myvalue",
			"yourlabel": "yourvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithNonMatchingAndIgnoredAnnotations(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation": "my.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation":   "myvalue",
			"yourannotation": "yourvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestRuleWithMatchingAndMissingLabels(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Labels: map[string]string{
			"mylabel":   "my.*",
			"yourlabel": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestRuleWithMatchingAndMissingAnnotations(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{
		Annotations: map[string]string{
			"myannotation":   "my.*",
			"yourannotation": "your.*",
		},
	}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if matches {
		t.Fail()
	}
}

func TestEmptyRuleMatchesEmptyAlert(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{}
	alert := &alertmanager.Alert{}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestEmptyRuleMatchesAlertWithLabel(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{}
	alert := &alertmanager.Alert{
		Labels: map[string]string{
			"mylabel": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func TestEmptyRuleMatchesAlertWithAnnotation(t *testing.T) {
	healer := makeHealer(t, "empty")
	rule := &monitoring.HealingRule{}
	alert := &alertmanager.Alert{
		Annotations: map[string]string{
			"myannotation": "myvalue",
		},
	}
	matches, err := healer.checkRule(rule, alert)
	if err != nil {
		t.Error(err)
	}
	if !matches {
		t.Fail()
	}
}

func makeHealer(t *testing.T, name string) *Healer {
	file := filepath.Join("..", "..", "testdata", name+"-config.yml")
	healer, err := NewHealerBuilder().
		ConfigFile(file).
		Build()
	if err != nil {
		t.Error(err)
		return nil
	}
	return healer
}
