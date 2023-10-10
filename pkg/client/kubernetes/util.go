package kubernetes

import (
	"bytes"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/client-go/kubernetes/scheme"
)

func ObjectToYaml(object runtime.Object) string {
	buffer := new(bytes.Buffer)
	printer := printers.YAMLPrinter{}
	if err := printer.PrintObj(object, buffer); err != nil {
		log.Fatalf("ObjectToYaml: %#v\n", err)
	}
	return buffer.String()
}

// https://github.com/kubernetes/client-go/issues/193
// https://medium.com/@harshjniitr/reading-and-writing-k8s-resource-as-yaml-in-golang-81dc8c7ea800
func YamlToDeployment(data string) *appsv1.Deployment {
	decoder := scheme.Codecs.UniversalDeserializer().Decode
	object, _, err := decoder([]byte(data), nil, nil)
	if err != nil {
		log.Fatalf("YamlToDeployment: %#v\n", err)
	}
	return object.(*appsv1.Deployment)
}
