package unstructured

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5"
	yamlv3 "gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func Encode(arr []unstructured.Unstructured) ([]byte, error) {
	sort.SliceStable(arr, func(i, j int) bool {
		if arr[i].GetAPIVersion() != arr[j].GetAPIVersion() {
			return arr[i].GetAPIVersion() < arr[j].GetAPIVersion()
		}
		if arr[i].GetKind() != arr[j].GetKind() {
			return arr[i].GetKind() < arr[j].GetKind()
		}
		return arr[i].GetName() < arr[j].GetName()
	})
	return yamlv3.Marshal(arr)
}

func Decode(source string) ([]unstructured.Unstructured, error) {
	splitString := regexp.MustCompile(`(?m)^---$`).Split(source, -1)
	resources := make([]unstructured.Unstructured, 0, len(splitString))
	for _, v := range splitString {
		if strings.TrimSpace(v) == "" {
			continue
		}
		// split lines and if all lines start with '#', it is skipped
		empty := true
		for _, line := range strings.Split(v, "\n") {
			l := strings.TrimPrefix(line, " ")
			l = strings.TrimPrefix(l, "\t")
			if l != "" && !strings.HasPrefix(l, "#") {
				empty = false
				break
			}
		}
		if empty {
			continue
		}
		_, obj, err := StringToUnstructured(v)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %w input='%s'", err, v)
		}
		resources = append(resources, *obj)
	}

	return resources, nil
}

func Replace(obj unstructured.Unstructured, key, value string) (*unstructured.Unstructured, error) {
	str_patch := fmt.Sprintf(`[{"op": "replace", "path": "%s", "value": "%s"}]`, key, value)
	bytes_obj, err := UnstructuredToJSONBytes(obj.DeepCopy())
	if err != nil {
		return nil, fmt.Errorf("failed to encode unstructured to JSON: %w", err)
	}

	// patch JSON6902
	patch, err := jsonpatch.DecodePatch([]byte(str_patch))
	if err != nil {
		return nil, fmt.Errorf("failed to decode patch: %w: %s", err, str_patch)
	}

	patched, err := patch.Apply(bytes_obj)
	if err != nil {
		return nil, fmt.Errorf("failed to patch JSON6902: %w: patch=%s: buf=%s", err, str_patch, bytes_obj)
	}

	_, patchedObj, err := BytesToUnstructured(patched)

	return patchedObj, err
}

func StringToUnstructured(data string) (*schema.GroupVersionKind, *unstructured.Unstructured, error) {
	return BytesToUnstructured([]byte(data))
}

func BytesToUnstructured(data []byte) (*schema.GroupVersionKind, *unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode([]byte(data), nil, obj)
	if err != nil {
		return nil, nil, err
	}
	return gvk, obj, nil
}

func UnstructuredToJSONBytes(obj *unstructured.Unstructured) ([]byte, error) {
	return json.Marshal(obj)
}
