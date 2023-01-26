package v1

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("redis webhook", func() {
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
		It("should create validate the log level", func() {
			By("creating a redis resource")
			ctx := context.Background()
			redis := &Redis{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "simple.simple.redis/v1",
					Kind:       "Redis",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      redisName,
					Namespace: redisNamespace,
				},
				Spec: RedisSpec{
					LogLevel: "test",
				},
			}
			Expect(k8sClient.Create(ctx, redis)).ShouldNot(Succeed())
		})
		It("should create validate the cluster size", func() {
			By("creating a redis resource")
			ctx := context.Background()
			redis := &Redis{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "simple.simple.redis/v1",
					Kind:       "Redis",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      redisName,
					Namespace: redisNamespace,
				},
				Spec: RedisSpec{
					ClusterSize: -1,
				},
			}
			Expect(k8sClient.Create(ctx, redis)).ShouldNot(Succeed())
		})
		It("should set default values", func() {
			By("creating a redis resource")
			ctx := context.Background()
			redis := &Redis{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "simple.simple.redis/v1",
					Kind:       "Redis",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      redisName,
					Namespace: redisNamespace,
				},
			}
			Expect(k8sClient.Create(ctx, redis)).Should(Succeed())
			redisLookup := types.NamespacedName{Name: redisName, Namespace: redisNamespace}
			createdRedis := &Redis{}
			Eventually(func() bool {
				err := k8sClient.Get(ctx, redisLookup, createdRedis)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
			Expect(createdRedis.Spec.LogLevel).Should(Equal(RLogLevelNotice))
			Expect(createdRedis.Spec.ClusterSize).Should(Equal(1))
		})
	})
})
