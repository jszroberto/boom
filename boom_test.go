package main

import (
	"github.com/geofffranks/simpleyaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	completeManifestPath             string = "examples/manifest.yml"
	manifestWithoutResourcePoolsPath string = "examples/manifest-without-resource-pools.yml"
)

var _ = Describe("Boom", func() {
	var (
		boom *Boom
	)

	BeforeEach(func() {
	})

	Context("SetInstances", func() {
		BeforeEach(func() {
			boom = New(completeManifestPath)
		})
		Context("when the job is found", func() {
			Context("and the resource_pool is specified", func() {

				Context("when resource_pools does not exist", func() {
					It("updates the value in the job", func() {
						boom = New(manifestWithoutResourcePoolsPath)
						err := boom.SetInstances("cell", 2)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(2))
					})
				})
				Context("when resource_pools exist", func() {
					It("updates the value in the job", func() {
						err := boom.SetInstances("cell", 2)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(2))
					})
					Context("when the resource pool is found", func() {
						Context("when the given value is the same", func() {
							It("does not modify resource pool", func() {
								err := boom.SetInstances("cell", 20)
								Expect(err).NotTo(HaveOccurred())
								result, err := simpleyaml.NewYaml([]byte(boom.String()))
								Expect(err).NotTo(HaveOccurred())
								Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(20))
								Expect(result.Get("resource_pools").GetIndex(0).Get("size").Int()).To(Equal(20))
							})
						})
						Context("when the given value is lower", func() {
							It("decreases the number of instances in resource pool", func() {

								err := boom.SetInstances("cell", 15)
								Expect(err).NotTo(HaveOccurred())
								result, err := simpleyaml.NewYaml([]byte(boom.String()))
								Expect(err).NotTo(HaveOccurred())
								Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(15))
								Expect(result.Get("resource_pools").GetIndex(0).Get("size").Int()).To(Equal(15))
							})
						})

						Context("when the given value is greater", func() {
							It("increases the number of instances in resource pool", func() {
								err := boom.SetInstances("cell", 25)
								Expect(err).NotTo(HaveOccurred())
								result, err := simpleyaml.NewYaml([]byte(boom.String()))
								Expect(err).NotTo(HaveOccurred())
								Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(25))
								Expect(result.Get("resource_pools").GetIndex(0).Get("size").Int()).To(Equal(25))
							})
						})
					})
				})
			})
			// TODO - Here is specified
			Context("and the resource_pool is not specified", func() {
				It("updates the value in the job", func() {
					err := boom.SetInstances("cell", 2)
					Expect(err).NotTo(HaveOccurred())
					result, err := simpleyaml.NewYaml([]byte(boom.String()))
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(2))
				})
			})
		})
		Context("when the job is not found", func() {
			It("returns an error", func() {
				err := boom.SetInstances("not-existing", 2)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("element `not-existing` not found"))

			})
		})
	})

	Context("ScaleInstances", func() {
		BeforeEach(func() {
			boom = New(completeManifestPath)
		})
		Context("when the job is found", func() {
			Context("when the value is not round", func() {
				It("don't update the value", func() {
					err := boom.ScaleInstances("cell", 1)
					Expect(err).NotTo(HaveOccurred())
					result, err := simpleyaml.NewYaml([]byte(boom.String()))
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(20))
				})
			})
			It("decreases the value", func() {
				err := boom.ScaleInstances("cell", 0.5)
				Expect(err).NotTo(HaveOccurred())
				result, err := simpleyaml.NewYaml([]byte(boom.String()))
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(10))
			})
			It("increases the value", func() {
				err := boom.ScaleInstances("cell", 2)
				Expect(err).NotTo(HaveOccurred())
				result, err := simpleyaml.NewYaml([]byte(boom.String()))
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(40))
			})
			Context("when the value is not round", func() {
				It("updates the value", func() {
					err := boom.ScaleInstances("cell", 1.5)
					Expect(err).NotTo(HaveOccurred())
					result, err := simpleyaml.NewYaml([]byte(boom.String()))
					Expect(err).NotTo(HaveOccurred())
					Expect(result.Get("jobs").GetIndex(0).Get("instances").Int()).To(Equal(30))
				})
			})
			Context("when the factor is 0", func() {
				It("returns an error", func() {
					err := boom.ScaleInstances("cell", 0)
					Expect(err).To(MatchError("factor 0 is not permitted"))
				})
			})

			Context("when the factor is too low to modify a unit", func() {
				Context("when force mode", func() {
					It("decreases the value", func() {
						err := boom.ScaleInstances("brain", 0.8)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("jobs").GetIndex(1).Get("instances").Int()).To(Equal(1))
					})
					It("increases the value", func() {
						boom.Force = true
						err := boom.ScaleInstances("brain", 1.2)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("jobs").GetIndex(1).Get("instances").Int()).To(Equal(3))
					})
				})
				Context("when force mode isn't used", func() {
					It("don't increase the value", func() {
						err := boom.ScaleInstances("brain", 1.2)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("resource_pools").GetIndex(1).Get("size").Int()).To(Equal(5))
					})
				})
			})
			Context("when the resource_pool is found", func() {
				Context("when a dedicated resource pool is used", func() {
					It("modifies the resource pool size", func() {
						err := boom.ScaleInstances("cell", 0.5)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("resource_pools").GetIndex(0).Get("size").Int()).To(Equal(10))
					})
				})
				Context("when a shared resource pool is used", func() {
					It("modifies the resource pool size", func() {
						err := boom.ScaleInstances("brain", 2)
						Expect(err).NotTo(HaveOccurred())
						result, err := simpleyaml.NewYaml([]byte(boom.String()))
						Expect(err).NotTo(HaveOccurred())
						Expect(result.Get("resource_pools").GetIndex(1).Get("size").Int()).To(Equal(7))
					})
				})
			})
		})

		Context("when the job is not found", func() {
			It("returns an error", func() {
				err := boom.SetInstances("not-existing", 2)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("element `not-existing` not found"))

			})
		})
	})
})
