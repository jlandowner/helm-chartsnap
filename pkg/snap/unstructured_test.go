package snap

import (
	"io"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/jlandowner/helm-chartsnap/pkg/unstructured"
)

var _ = Describe("Unstructured Snapshot", func() {
	f := func(m OmegaMatcher, filePath string) (success bool, err error) {
		f, err := os.Open(filePath)
		Expect(err).NotTo(HaveOccurred())
		defer f.Close()

		buf, err := io.ReadAll(f)
		Expect(err).NotTo(HaveOccurred())

		manifests, errs := unstructured.Decode(string(buf))
		Expect(len(errs)).To(BeZero())

		return UnstructuredMatch(m, manifests)
	}

	It("should match", func() {
		m := UnstructuredSnapShotMatcher("testdata/unstructured.snap", "default")
		success, err := f(m, "testdata/unstructured.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(success).To(BeTrue())
	})

	It("should not match and output diff N=1", func() {
		m := UnstructuredSnapShotMatcher("testdata/unstructured.snap", "default", WithDiffContextLineN(1))
		success, err := f(m, "testdata/unstructured_diff.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(success).To(BeFalse())

		Expect(m.FailureMessage(nil)).Should(Equal(`Expected to match
--- kind=Deployment name=chartsnap-app1 line=8
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0

--- kind=Deployment name=chartsnap-app1 line=24
                      app.kubernetes.io/name: app1
-                     app.kubernetes.io/version: 1.16.0
+                     app.kubernetes.io/version: 1.21.0
                      helm.sh/chart: app1-0.1.0

--- kind=Deployment name=chartsnap-app1 line=29
                  containers:
-                     - image: nginx:1.16.0
+                     - image: nginx:1.21.0
                        imagePullPolicy: IfNotPresent

--- kind=Deployment name=chartsnap-app1 line=32
                        imagePullPolicy: IfNotPresent
-                       livenessProbe:
-                         httpGet:
-                             path: /
-                             port: http
                        name: app1

--- kind=Pod name=chartsnap-app1-test-connection line=59
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0

--- kind=Service name=chartsnap-app1 line=80
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0

--- kind=Service name=chartsnap-app1 line=90
                targetPort: http
+             - name: https
+               port: 443
+               protocol: TCP
+               targetPort: https
          selector:

--- kind=ServiceAccount name= line=107
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0


`))
	})

	It("should not match and output full diff", func() {
		m := UnstructuredSnapShotMatcher("testdata/unstructured.snap", "default")
		success, err := f(m, "testdata/unstructured_diff.yaml")
		Expect(err).NotTo(HaveOccurred())
		Expect(success).To(BeFalse())

		Expect(m.FailureMessage(nil)).Should(Equal(`Expected to match
  - object:
      apiVersion: apps/v1
      kind: Deployment
      metadata:
          labels:
              app.kubernetes.io/instance: chartsnap
              app.kubernetes.io/managed-by: Helm
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0
          name: chartsnap-app1
      spec:
          replicas: 1
          selector:
              matchLabels:
                  app.kubernetes.io/instance: chartsnap
                  app.kubernetes.io/name: app1
          template:
              metadata:
                  labels:
                      app.kubernetes.io/instance: chartsnap
                      app.kubernetes.io/managed-by: Helm
                      app.kubernetes.io/name: app1
-                     app.kubernetes.io/version: 1.16.0
+                     app.kubernetes.io/version: 1.21.0
                      helm.sh/chart: app1-0.1.0
              spec:
                  containers:
-                     - image: nginx:1.16.0
+                     - image: nginx:1.21.0
                        imagePullPolicy: IfNotPresent
-                       livenessProbe:
-                         httpGet:
-                             path: /
-                             port: http
                        name: app1
                        ports:
                          - containerPort: 80
                            name: http
                            protocol: TCP
                        readinessProbe:
                          httpGet:
                              path: /
                              port: http
                        resources: {}
                        securityContext: {}
                  securityContext: {}
                  serviceAccountName: chartsnap-app1
  - object:
      apiVersion: v1
      kind: Pod
      metadata:
          annotations:
              helm.sh/hook: test
          labels:
              app.kubernetes.io/instance: chartsnap
              app.kubernetes.io/managed-by: Helm
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0
          name: chartsnap-app1-test-connection
      spec:
          containers:
              - args:
                  - chartsnap-app1:80
                command:
                  - wget
                image: busybox
                name: wget
          restartPolicy: Never
  - object:
      apiVersion: v1
      kind: Service
      metadata:
          labels:
              app.kubernetes.io/instance: chartsnap
              app.kubernetes.io/managed-by: Helm
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0
          name: chartsnap-app1
      spec:
          ports:
              - name: http
                port: 80
                protocol: TCP
                targetPort: http
+             - name: https
+               port: 443
+               protocol: TCP
+               targetPort: https
          selector:
              app.kubernetes.io/instance: chartsnap
              app.kubernetes.io/name: app1
          type: ClusterIP
  - object:
      apiVersion: v1
      automountServiceAccountToken: true
      kind: ServiceAccount
      metadata:
          labels:
              app.kubernetes.io/instance: chartsnap
              app.kubernetes.io/managed-by: Helm
              app.kubernetes.io/name: app1
-             app.kubernetes.io/version: 1.16.0
+             app.kubernetes.io/version: 1.21.0
              helm.sh/chart: app1-0.1.0
          name: chartsnap-app1
  

`))
	})
})
