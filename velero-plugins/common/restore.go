package common

import (
	"github.com/sirupsen/logrus"
	"github.com/vmware-tanzu/velero/pkg/plugin/velero"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// RestorePlugin is a restore item action plugin for Heptio Ark.
type RestorePlugin struct {
	Log logrus.FieldLogger
}

// AppliesTo returns a velero.ResourceSelector that applies to the listed resources in the slice.
func (p *RestorePlugin) AppliesTo() (velero.ResourceSelector, error) {
	return velero.ResourceSelector{
		IncludedResources: []string{"*"},
	}, nil
}

// Execute sets a custom annotation on the item being restored.
func (p *RestorePlugin) Execute(input *velero.RestoreItemActionExecuteInput) (*velero.RestoreItemActionExecuteOutput, error) {
	metadata, annotations, err := getMetadataAndAnnotations(input.Item)
	if err != nil {
		return nil, err
	}
	name := metadata.GetName()
	p.Log.Infof("[common-restore] common restore plugin for %s", name)

	annotations[MigrationRegistry] = input.Restore.Annotations[MigrationRegistry]

	// Set migmigration and migplan labels on all resources, except ServiceAccounts
	switch input.Item.DeepCopyObject().(type) {
	case *corev1.ServiceAccount:
		break
	default:
		migMigrationLabel, exist := input.Restore.Labels[MigMigrationLabelKey]
		if !exist {
			p.Log.Info("migmigration label was not found on restore")
		}
		migPlanLabel, exist := input.Restore.Labels[MigPlanLabelKey]
		if !exist {
			p.Log.Info("migplan label was not found on restore")
		}
		labels := metadata.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels[MigMigrationLabelKey] = migMigrationLabel
		labels[MigPlanLabelKey] = migPlanLabel

		metadata.SetLabels(labels)
	}

	metadata.SetAnnotations(annotations)

	return velero.NewRestoreItemActionExecuteOutput(input.Item), nil
}

func getMetadataAndAnnotations(item runtime.Unstructured) (metav1.Object, map[string]string, error) {
	metadata, err := meta.Accessor(item)
	if err != nil {
		return nil, nil, err
	}

	annotations := metadata.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	return metadata, annotations, nil
}
