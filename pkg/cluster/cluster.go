/*
SPDX-License-Identifier: Apache-2.0

Copyright Contributors to the Submariner project.

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

package cluster

import (
	"context"

	"github.com/pkg/errors"
	"github.com/submariner-io/submariner-operator/api/submariner/v1alpha1"
	"github.com/submariner-io/submariner-operator/internal/constants"
	"github.com/submariner-io/submariner-operator/pkg/client"
	submarinerv1 "github.com/submariner-io/submariner/pkg/apis/submariner.io/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Info struct {
	Name            string
	ClientProducer  client.Producer
	Submariner      *v1alpha1.Submariner
}

func New(clusterName string, clientProducer client.Producer) (*Info, error) {
	cluster := &Info{
		Name:           clusterName,
		ClientProducer: clientProducer,
	}

	var err error

	cluster.Submariner, err = cluster.GetSubmariner()
	if err != nil {
		return nil, errors.Wrap(err, "Error retrieving Submariner")
	}

	return cluster, nil
}

func (c *Info) GetSubmariner() (*v1alpha1.Submariner, error) {
	submariner, err := c.ClientProducer.ForOperator().SubmarinerV1alpha1().Submariners(constants.SubmarinerNamespace).
		Get(context.TODO(), constants.SubmarinerName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.New( "Submariner not found")
		}
		return nil, err
	}
	return submariner, nil
}

func (c *Info) GetGateways() ([]submarinerv1.Gateway, error) {
	gateways, err := c.ClientProducer.ForSubmariner().SubmarinerV1().
		Gateways(constants.OperatorNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, errors.New("No gateways found")
		}

		return nil, err
	}

	return gateways.Items, nil
}