/*
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

package resource

import (
	"context"
	"fmt"

	"github.com/prateekgogia/kit/pkg/apis/infrastructure/v1alpha1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AutoScalingGroup struct {
	KubeClient client.Client
}

func (l *AutoScalingGroup) Create(ctx context.Context, controlPlane *v1alpha1.ControlPlane) error {
	for _, component := range v1alpha1.ComponentsSupported {
		if err := l.exists(ctx, controlPlane.Namespace, ObjectName(controlPlane, component)); err != nil {
			if errors.IsNotFound(err) {
				if err := l.create(ctx, component, controlPlane); err != nil {
					return fmt.Errorf("creating auto scaling group kube object, %w", err)
				}
				continue
			}
			return fmt.Errorf("getting auto scaling group object, %w", err)
		}
	}
	// TODO verify existing object matches the desired else update
	return nil
}

func (l *AutoScalingGroup) create(ctx context.Context, component string, controlPlane *v1alpha1.ControlPlane) error {
	if err := l.KubeClient.Create(ctx, &v1alpha1.AutoScalingGroup{
		ObjectMeta: ObjectMeta(controlPlane, component),
		Spec: v1alpha1.AutoScalingGroupSpec{
			ClusterName:   controlPlane.Name,
			InstanceCount: 3,
		},
	}); err != nil {
		return fmt.Errorf("creating auto scaling group kube object, %w", err)
	}
	zap.S().Debugf("Successfully created auto scaling group object %v for cluster %v",
		ObjectMeta(controlPlane, component).Name, controlPlane.Name)
	return nil
}

func (l *AutoScalingGroup) exists(ctx context.Context, ns, objName string) error {
	result := &v1alpha1.AutoScalingGroup{}
	if err := l.KubeClient.Get(ctx, NamespacedName(ns, objName), result); err != nil {
		return err
	}
	return nil
}