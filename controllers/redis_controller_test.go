package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	simplev1 "github.com/spazzy757/simple-redis/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("redis controller", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		redisName           = "redis-test"
		redisNamespace      = "default"
		resourceMasterName  = "redis-test-master"
		resourceReplicaName = "redis-test-replica"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("when creating a redis instance", func() {

		masterLookup := types.NamespacedName{
			Name:      resourceMasterName,
			Namespace: redisNamespace,
		}
		replicaLookup := types.NamespacedName{
			Name:      resourceReplicaName,
			Namespace: redisNamespace,
		}

		It("should create all the required resources", func() {

			By("creating a redis resource")
			ctx := context.Background()
			redis := &simplev1.Redis{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "simple.simple.redis/v1",
					Kind:       "Redis",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      redisName,
					Namespace: redisNamespace,
				},
			}
			redisLookup := types.NamespacedName{Name: redisName, Namespace: redisNamespace}
			createdRedis := &simplev1.Redis{}
			Expect(k8sClient.Create(ctx, redis)).Should(Succeed())
			Eventually(func() bool {
				err := k8sClient.Get(ctx, redisLookup, createdRedis)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(createdRedis.Status.Status).Should(Equal(simplev1.StatusPending))

			By("creating a master redis deployment")
			masterdeploy := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, masterLookup, masterdeploy)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(masterdeploy.ObjectMeta.Name).Should(Equal(resourceMasterName))
			Expect(*masterdeploy.Spec.Replicas).Should(Equal(int32(1)))

			By("creating a master redis service")
			mastersvc := &v1.Service{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, masterLookup, mastersvc)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(mastersvc.ObjectMeta.Name).Should(Equal(resourceMasterName))

			By("creating a replica redis deployment")
			replicadeploy := &appsv1.Deployment{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, replicaLookup, replicadeploy)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(replicadeploy.ObjectMeta.Name).Should(Equal(resourceReplicaName))
			Expect(*replicadeploy.Spec.Replicas).Should(Equal(int32(0)))
		})
	})
})
