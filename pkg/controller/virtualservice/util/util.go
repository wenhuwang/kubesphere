package util

import (
	"github.com/knative/pkg/apis/istio/v1alpha3"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	AppLabel                = "app"
	VersionLabel            = "version"
	ApplicationNameLabel    = "app.kubernetes.io/name"
	ApplicationVersionLabel = "app.kubernetes.io/version"
	ServiceMeshEnabledLabel = "servicemesh.kubesphere.io/enabled"
)

// resource with these following labels considered as part of servicemesh
var ApplicationLabels = [...]string{
	ApplicationNameLabel,
	ServiceMeshEnabledLabel,
	AppLabel,
}

var TrimChars = [...]string{".", "_", "-"}

// normalize version names
// strip [_.-]
func NormalizeVersionName(version string) string {
	for _, char := range TrimChars {
		version = strings.ReplaceAll(version, char, "")
	}
	return version
}

func GetComponentName(meta *metav1.ObjectMeta) string {
	if len(meta.Labels[AppLabel]) > 0 {
		return meta.Labels[AppLabel]
	}
	return ""
}

func GetComponentVersion(meta *metav1.ObjectMeta) string {
	if len(meta.Labels[VersionLabel]) > 0 {
		return meta.Labels[VersionLabel]
	}
	return ""
}

func ExtractApplicationLabels(meta *metav1.ObjectMeta) map[string]string {

	labels := make(map[string]string, len(ApplicationLabels))
	for _, label := range ApplicationLabels {
		if _, ok := meta.Labels[label]; !ok {
			return nil
		} else {
			labels[label] = meta.Labels[label]
		}
	}

	return labels
}

func IsApplicationComponent(lbs map[string]string) bool {

	for _, label := range ApplicationLabels {
		if _, ok := lbs[label]; !ok {
			return false
		}
	}

	return true
}

// if virtualservice not specified with port number, then fill with service first port
func FillDestinationPort(vs *v1alpha3.VirtualService, service *v1.Service) {
	// fill http port
	for i := range vs.Spec.Http {
		for j := range vs.Spec.Http[i].Route {
			if vs.Spec.Http[i].Route[j].Destination.Port.Number == 0 {
				vs.Spec.Http[i].Route[j].Destination.Port.Number = uint32(service.Spec.Ports[0].Port)
			}
		}

		if vs.Spec.Http[i].Mirror != nil && vs.Spec.Http[i].Mirror.Port.Number == 0 {
			vs.Spec.Http[i].Mirror.Port.Number = uint32(service.Spec.Ports[0].Port)
		}
	}

	// fill tcp port
	for i := range vs.Spec.Tcp {
		for j := range vs.Spec.Tcp[i].Route {
			if vs.Spec.Tcp[i].Route[j].Destination.Port.Number == 0 {
				vs.Spec.Tcp[i].Route[j].Destination.Port.Number = uint32(service.Spec.Ports[0].Port)
			}
		}
	}
}