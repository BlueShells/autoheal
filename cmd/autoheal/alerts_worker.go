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
	"regexp"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/util/runtime"

	alertmanager "github.com/openshift/autoheal/pkg/alertmanager"
	monitoring "github.com/openshift/autoheal/pkg/apis/monitoring/v1alpha1"
)

func (h *Healer) runAlertsWorker() {
	for h.pickAlert() {
		// Nothing.
	}
}

func (h *Healer) pickAlert() bool {
	// Get the next item and end the work loop if asked to stop:
	item, stop := h.alertsQueue.Get()
	if stop {
		return false
	}

	// Process the item and make sure to always tell the queue that we are done with this item:
	err := func(item interface{}) error {
		h.alertsQueue.Done(item)

		// Check that the item we got from the queue is really an alert, and discard it otherwise:
		alert, ok := item.(*alertmanager.Alert)
		if !ok {
			h.alertsQueue.Forget(item)
		}

		// Process and then forget the alert:
		err := h.processAlert(alert)
		if err != nil {
			return err
		}
		h.alertsQueue.Forget(alert)

		return nil
	}(item)
	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

func (h *Healer) processAlert(alert *alertmanager.Alert) error {
	switch alert.Status {
	case alertmanager.AlertStatusFiring:
		return h.startHealing(alert)
	case alertmanager.AlertStatusResolved:
		return h.cancelHealing(alert)
	default:
		glog.Warningf(
			"Unknnown status '%s' reported by alert manager, will ignore it",
			alert.Status,
		)
		return nil
	}
}

// startHealing starts the healing process for the given alert.
//
func (h *Healer) startHealing(alert *alertmanager.Alert) error {
	// Find the rules that are activated for the alert:
	activated := make([]*monitoring.HealingRule, 0)
	h.rulesCache.Range(func(_, value interface{}) bool {
		rule := value.(*monitoring.HealingRule)
		matches, err := h.checkRule(rule, alert)
		if err != nil {
			glog.Errorf(
				"Error while checking if rule '%s' matches alert '%s': %s",
				rule.ObjectMeta.Name,
				alert.Name(),
				err,
			)
		} else if matches {
			glog.Infof(
				"Rule '%s' matches alert '%s'",
				rule.ObjectMeta.Name,
				alert.Name(),
			)
			activated = append(activated, rule)
		}
		return true
	})
	if len(activated) == 0 {
		glog.Infof("No rule matches alert '%s'", alert.Name())
		return nil
	}

	// Execute the activated rules:
	for _, rule := range activated {
		err := h.runRule(rule, alert)
		if err != nil {
			return err
		}
	}

	return nil
}

// cancelHealing cancels the healing process for the given alert.
//
func (h *Healer) cancelHealing(alert *alertmanager.Alert) error {
	return nil
}

func (h *Healer) checkRule(rule *monitoring.HealingRule, alert *alertmanager.Alert) (matches bool, err error) {
	glog.Infof(
		"Checking rule '%s' for alert '%s'",
		rule.ObjectMeta.Name,
		alert.Name(),
	)
	matches, err = h.checkMap(alert.Labels, rule.Labels)
	if !matches || err != nil {
		return
	}
	matches, err = h.checkMap(alert.Annotations, rule.Annotations)
	if !matches || err != nil {
		return
	}
	return
}

func (h *Healer) checkMap(values, patterns map[string]string) (result bool, err error) {
	if len(patterns) > 0 {
		if len(values) == 0 {
			return
		}
		for key, pattern := range patterns {
			value, present := values[key]
			if !present {
				return
			}
			var matches bool
			matches, err = regexp.MatchString(pattern, value)
			if !matches || err != nil {
				return
			}
		}
	}
	result = true
	return
}

func (h *Healer) runRule(rule *monitoring.HealingRule, alert *alertmanager.Alert) error {
	// Send the name of the rule to the log:
	glog.Infof(
		"Running rule '%s' for alert '%s'",
		rule.ObjectMeta.Name,
		alert.Name(),
	)

	// Prepare the template to process the action:
	template, err := NewObjectTemplateBuilder().
		Variable("alert", ".").
		Variable("labels", ".Labels").
		Variable("annotations", ".Annotations").
		Build()
	if err != nil {
		return err
	}

	// Decide which kind of action to run, and run it:
	if rule.AWXJob != nil {
		action := rule.AWXJob.DeepCopy()
		err = template.Process(action, alert)
		if err != nil {
			return err
		}
		return h.runAWXJob(rule, action, alert)
	} else if rule.BatchJob != nil {
		action := rule.BatchJob.DeepCopy()
		err = template.Process(action, alert)
		if err != nil {
			return err
		}
		return h.runBatchJob(rule, action, alert)
	} else {
		glog.Warningf(
			"There are no action details, rule '%s' will have no effect on alert '%s'",
			rule.ObjectMeta.Name,
			alert.Name(),
		)
	}
	return nil
}
