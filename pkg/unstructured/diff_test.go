package unstructured

import (
	"testing"

	"github.com/aryann/difflib"
)

var (
	testdata = []difflib.DiffRecord{
		{Delta: difflib.LeftOnly, Payload: "apiVersion: v2"},
		{Delta: difflib.RightOnly, Payload: "apiVersion: v1"},
		{Delta: difflib.Common, Payload: "automountServiceAccountToken: true"},
		{Delta: difflib.Common, Payload: "kind: ServiceAccount"},
		{Delta: difflib.Common, Payload: "  labels:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "  name: chartsnap-app1"},
		{Delta: difflib.Common, Payload: "---"},
		{Delta: difflib.LeftOnly, Payload: "apiVersion: v1"},
		{Delta: difflib.LeftOnly, Payload: "kind: Namespace"},
		{Delta: difflib.LeftOnly, Payload: "  annotations:"},
		{Delta: difflib.LeftOnly, Payload: "    helm.sh/hook: test"},
		{Delta: difflib.LeftOnly, Payload: "  labels:"},
		{Delta: difflib.LeftOnly, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.LeftOnly, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.LeftOnly, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.LeftOnly, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.LeftOnly, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.LeftOnly, Payload: "  name: chartsnap-app1-namespace"},
		{Delta: difflib.Common, Payload: "---"},
		{Delta: difflib.Common, Payload: "apiVersion: v1"},
		{Delta: difflib.LeftOnly, Payload: "  ca.crt: '###DYNAMIC_FIELD###'"},
		{Delta: difflib.LeftOnly, Payload: "  tls.crt: '###DYNAMIC_FIELD###'"},
		{Delta: difflib.LeftOnly, Payload: "  tls.key: '###DYNAMIC_FIELD###'"},
		{Delta: difflib.RightOnly, Payload: "  ca.crt: IyMjRFlOQU1JQ19GSUVMRCMjIw=="},
		{Delta: difflib.RightOnly, Payload: "  tls.crt: IyMjRFlOQU1JQ19GSUVMRCMjIw=="},
		{Delta: difflib.RightOnly, Payload: "  tls.key: IyMjRFlOQU1JQ19GSUVMRCMjIw=="},
		{Delta: difflib.Common, Payload: "kind: Secret"},
		{Delta: difflib.Common, Payload: "  labels:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "  name: app1-cert"},
		{Delta: difflib.Common, Payload: "  namespace: default"},
		{Delta: difflib.Common, Payload: "type: kubernetes.io/tls"},
		{Delta: difflib.Common, Payload: "---"},
		{Delta: difflib.Common, Payload: "apiVersion: v1"},
		{Delta: difflib.Common, Payload: "kind: Service"},
		{Delta: difflib.Common, Payload: "  labels:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "  name: chartsnap-app1"},
		{Delta: difflib.Common, Payload: "  ports:"},
		{Delta: difflib.Common, Payload: "    - name: http"},
		{Delta: difflib.Common, Payload: "      port: 80"},
		{Delta: difflib.Common, Payload: "      protocol: TCP"},
		{Delta: difflib.Common, Payload: "      targetPort: http"},
		{Delta: difflib.Common, Payload: "  selector:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.LeftOnly, Payload: "  type: LoadBalancer"},
		{Delta: difflib.RightOnly, Payload: "  type: ClusterIP"},
		{Delta: difflib.Common, Payload: "---"},
		{Delta: difflib.Common, Payload: "apiVersion: apps/v1"},
		{Delta: difflib.Common, Payload: "kind: Deployment"},
		{Delta: difflib.Common, Payload: "  labels:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.LeftOnly, Payload: "    app.kubernetes.io/version: 1.15.0"},
		{Delta: difflib.RightOnly, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "  name: chartsnap-app1"},
		{Delta: difflib.Common, Payload: "  replicas: 1"},
		{Delta: difflib.Common, Payload: "  selector:"},
		{Delta: difflib.Common, Payload: "    matchLabels:"},
		{Delta: difflib.Common, Payload: "      app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "      app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "  template:"},
		{Delta: difflib.Common, Payload: "    metadata:"},
		{Delta: difflib.Common, Payload: "      labels:"},
		{Delta: difflib.Common, Payload: "        app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "        app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "        app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "        app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "        helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "    spec:"},
		{Delta: difflib.Common, Payload: "      containers:"},
		{Delta: difflib.Common, Payload: "        - image: nginx:1.16.0"},
		{Delta: difflib.Common, Payload: "          imagePullPolicy: IfNotPresent"},
		{Delta: difflib.Common, Payload: "          livenessProbe:"},
		{Delta: difflib.Common, Payload: "            httpGet:"},
		{Delta: difflib.Common, Payload: "              path: /"},
		{Delta: difflib.Common, Payload: "              port: http"},
		{Delta: difflib.Common, Payload: "          name: app1"},
		{Delta: difflib.Common, Payload: "          ports:"},
		{Delta: difflib.Common, Payload: "            - containerPort: 80"},
		{Delta: difflib.Common, Payload: "              name: http"},
		{Delta: difflib.Common, Payload: "              protocol: TCP"},
		{Delta: difflib.Common, Payload: "          readinessProbe:"},
		{Delta: difflib.Common, Payload: "            httpGet:"},
		{Delta: difflib.Common, Payload: "              path: /"},
		{Delta: difflib.Common, Payload: "              port: http"},
		{Delta: difflib.Common, Payload: "          resources: {}"},
		{Delta: difflib.Common, Payload: "          securityContext: {}"},
		{Delta: difflib.Common, Payload: "      securityContext: {}"},
		{Delta: difflib.Common, Payload: "      serviceAccountName: chartsnap-app1"},
		{Delta: difflib.RightOnly, Payload: "---"},
		{Delta: difflib.RightOnly, Payload: "apiVersion: networking.k8s.io/v1"},
		{Delta: difflib.RightOnly, Payload: "kind: Ingress"},
		{Delta: difflib.RightOnly, Payload: "  annotations:"},
		{Delta: difflib.RightOnly, Payload: "    cert-manager.io/cluster-issuer: nameOfClusterIssuer"},
		{Delta: difflib.RightOnly, Payload: "  labels:"},
		{Delta: difflib.RightOnly, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.RightOnly, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.RightOnly, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.RightOnly, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.RightOnly, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.RightOnly, Payload: "  name: chartsnap-app1"},
		{Delta: difflib.RightOnly, Payload: "  ingressClassName: nginx"},
		{Delta: difflib.RightOnly, Payload: "  rules:"},
		{Delta: difflib.RightOnly, Payload: "    - host: chart-example.local"},
		{Delta: difflib.RightOnly, Payload: "      http:"},
		{Delta: difflib.RightOnly, Payload: "        paths:"},
		{Delta: difflib.RightOnly, Payload: "          - backend:"},
		{Delta: difflib.RightOnly, Payload: "              service:"},
		{Delta: difflib.RightOnly, Payload: "                name: chartsnap-app1"},
		{Delta: difflib.RightOnly, Payload: "                port:"},
		{Delta: difflib.RightOnly, Payload: "                  number: 80"},
		{Delta: difflib.RightOnly, Payload: "            path: /"},
		{Delta: difflib.RightOnly, Payload: "            pathType: ImplementationSpecific"},
		{Delta: difflib.RightOnly, Payload: "  tls:"},
		{Delta: difflib.RightOnly, Payload: "    - hosts:"},
		{Delta: difflib.RightOnly, Payload: "        - chart-example.local"},
		{Delta: difflib.RightOnly, Payload: "      secretName: chart-example-tls"},
		{Delta: difflib.Common, Payload: "---"},
		{Delta: difflib.Common, Payload: "apiVersion: v1"},
		{Delta: difflib.Common, Payload: "kind: Pod"},
		{Delta: difflib.Common, Payload: "  annotations:"},
		{Delta: difflib.Common, Payload: "    helm.sh/hook: test"},
		{Delta: difflib.Common, Payload: "  labels:"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/instance: chartsnap"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/managed-by: Helm"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/name: app1"},
		{Delta: difflib.Common, Payload: "    app.kubernetes.io/version: 1.16.0"},
		{Delta: difflib.Common, Payload: "    helm.sh/chart: app1-0.1.0"},
		{Delta: difflib.Common, Payload: "  name: chartsnap-app1-test-connection"},
		{Delta: difflib.Common, Payload: "  containers:"},
		{Delta: difflib.Common, Payload: "    - args:"},
		{Delta: difflib.Common, Payload: "        - chartsnap-app1:80"},
		{Delta: difflib.Common, Payload: "      command:"},
		{Delta: difflib.Common, Payload: "        - wget"},
		{Delta: difflib.Common, Payload: "      image: busybox"},
		{Delta: difflib.Common, Payload: "      name: wget"},
		{Delta: difflib.Common, Payload: "  restartPolicy: Never"},
		{Delta: difflib.LeftOnly, Payload: "a"},
	}
)

func TestMergeDiffOptions(t *testing.T) {
	opts := []DiffOptions{
		{ContextLineN: 5},
		{ContextLineN: 3},
		{ContextLineN: 7},
	}

	merged := MergeDiffOptions(opts)

	if merged.ContextLineN != 7 {
		t.Errorf("Expected DiffContextLineN to be 7, got %d", merged.ContextLineN)
	}
}

func Test_printDiff(t *testing.T) {
	d := difflib.DiffRecord{
		Delta:   difflib.LeftOnly,
		Payload: "abc",
	}

	want := "- abc\n"

	if got := printDiff(d); got != want {
		t.Errorf("printDiff() = %v, want %v", got, want)
	}
}

func Test_printHeader(t *testing.T) {
	kind := "TestKind"
	name := "TestName"
	lineN := 5

	want := "@@ KIND=TestKind NAME=TestName LINE=5\n"

	if got := printHeader(kind, name, lineN); got != want {
		t.Errorf("printHeader() = %v, want %v", got, want)
	}
}

func Test_findNextKind(t *testing.T) {
	diffs := []difflib.DiffRecord{
		{Delta: difflib.Common, Payload: "abc"},
		{Delta: difflib.Common, Payload: "kind: Pod"},
		{Delta: difflib.Common, Payload: "def"},
	}

	want := "Pod"

	if got := findNextKind(diffs); got != want {
		t.Errorf("findNextKind() = %v, want %v", got, want)
	}
}

func Test_findNextName(t *testing.T) {
	diffs := []difflib.DiffRecord{
		{Delta: difflib.Common, Payload: "abc"},
		{Delta: difflib.Common, Payload: "metadata:"},
		{Delta: difflib.Common, Payload: "  name: TestName"},
		{Delta: difflib.Common, Payload: "def"},
	}

	want := "TestName"

	if got := findNextName(diffs); got != want {
		t.Errorf("findNextName() = %v, want %v", got, want)
	}
}
