package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "service-operator/api/v1"
)

var _ = Describe("Service Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-service"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		service := &appsv1alpha1.Service{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Service")
			err := k8sClient.Get(ctx, typeNamespacedName, service)
			if err != nil && errors.IsNotFound(err) {
				resource := &appsv1alpha1.Service{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: appsv1alpha1.ServiceSpec{
						Image: "nginx:1.21",
						Port:  80,
						ConfigData: map[string]string{
							"config.yaml": "test: value",
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &appsv1alpha1.Service{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Service")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ServiceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking if Deployment was created")
			deployment := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, typeNamespacedName, deployment)
			}, time.Minute, time.Second).Should(Succeed())

			By("Checking if Service was created")
			k8sService := &corev1.Service{}
			Eventually(func() error {
				return k8sClient.Get(ctx, typeNamespacedName, k8sService)
			}, time.Minute, time.Second).Should(Succeed())

			By("Checking if ConfigMap was created")
			configMap := &corev1.ConfigMap{}
			configMapName := types.NamespacedName{
				Name:      resourceName + "-config",
				Namespace: "default",
			}
			Eventually(func() error {
				return k8sClient.Get(ctx, configMapName, configMap)
			}, time.Minute, time.Second).Should(Succeed())
		})

		It("should create ingress when enabled", func() {
			By("Creating service with ingress enabled")
			service := &appsv1alpha1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-service-ingress",
					Namespace: "default",
				},
				Spec: appsv1alpha1.ServiceSpec{
					Image: "nginx:1.21",
					Port:  80,
					Ingress: &appsv1alpha1.IngressSpec{
						Enabled: true,
						Host:    "test.example.com",
						Path:    "/",
					},
				},
			}
			Expect(k8sClient.Create(ctx, service)).To(Succeed())

			controllerReconciler := &ServiceReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			ingressTypeNamespacedName := types.NamespacedName{
				Name:      "test-service-ingress",
				Namespace: "default",
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: ingressTypeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())

			By("Checking if Ingress was created")
			ingress := &networkingv1.Ingress{}
			Eventually(func() error {
				return k8sClient.Get(ctx, ingressTypeNamespacedName, ingress)
			}, time.Minute, time.Second).Should(Succeed())

			// Cleanup
			Expect(k8sClient.Delete(ctx, service)).To(Succeed())
		})
	})
})
