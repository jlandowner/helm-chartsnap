package gomega

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ObjectSnapshot(obj client.Object) client.Object {
	t := obj.DeepCopyObject()
	o := t.(client.Object)
	RemoveDynamicFields(o)
	return o
}

func RemoveDynamicFields(o client.Object) {
	o.SetCreationTimestamp(metav1.Time{})
	o.SetResourceVersion("")
	o.SetGeneration(0)
	o.SetUID(types.UID(""))
	o.SetManagedFields(nil)

	ownerRefs := make([]metav1.OwnerReference, len(o.GetOwnerReferences()))
	for i, v := range o.GetOwnerReferences() {
		v.UID = ""
		ownerRefs[i] = v
	}
	o.SetOwnerReferences(ownerRefs)
}
