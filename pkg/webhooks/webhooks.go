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

package webhooks

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	knativeinjection "knative.dev/pkg/injection"
	"knative.dev/pkg/webhook/resourcesemantics"
	"knative.dev/pkg/webhook/resourcesemantics/conversion"
	"knative.dev/pkg/webhook/resourcesemantics/defaulting"
	"knative.dev/pkg/webhook/resourcesemantics/validation"

	"github.com/awslabs/operatorpkg/object"

	v1 "github.com/aws/karpenter-provider-aws/pkg/apis/v1"
	"github.com/aws/karpenter-provider-aws/pkg/apis/v1beta1"
)

var (
	Resources = map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
		object.GVK(&v1beta1.EC2NodeClass{}): &v1beta1.EC2NodeClass{},
	}
	ConversionResource = map[schema.GroupKind]conversion.GroupKindConversion{
		object.GVK(&v1.EC2NodeClass{}).GroupKind(): {
			DefinitionName: "ec2nodeclasses.karpenter.k8s.aws",
			HubVersion:     "v1",
			Zygotes: map[string]conversion.ConvertibleObject{
				"v1":      &v1.EC2NodeClass{},
				"v1beta1": &v1beta1.EC2NodeClass{},
			},
		},
	}
)

func NewWebhooks() []knativeinjection.ControllerConstructor {
	return []knativeinjection.ControllerConstructor{
		NewCRDDefaultingWebhook,
		NewCRDValidationWebhook,
	}
}

func NewCRDDefaultingWebhook(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	return defaulting.NewAdmissionController(ctx,
		"defaulting.webhook.karpenter.k8s.aws",
		"/default/karpenter.k8s.aws",
		Resources,
		func(ctx context.Context) context.Context { return ctx },
		true,
	)
}

func NewCRDValidationWebhook(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	return validation.NewAdmissionController(ctx,
		"validation.webhook.karpenter.k8s.aws",
		"/validate/karpenter.k8s.aws",
		Resources,
		func(ctx context.Context) context.Context { return ctx },
		true,
	)
}

func NewCRDConversionWebhook(ctx context.Context, _ configmap.Watcher) *controller.Impl {
	return conversion.NewConversionController(ctx,
		"/conversion/karpenter.sh",
		ConversionResource,
		func(ctx context.Context) context.Context { return ctx },
	)
}
