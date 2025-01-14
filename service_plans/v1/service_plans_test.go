package service_plans

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/cf-performance-tests/helpers"
)

var _ = Describe("service plans", func() {
	Describe("GET /v3/service_plans", func() {
		Context("as admin", func() {
			Measure("list all", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, "/v3/service_plans")
				})
			}, testConfig.Samples)
			Measure(fmt.Sprintf("list with page size %d", testConfig.LargePageSize), func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf("/v3/service_plans?per_page=%d", testConfig.LargePageSize))
				})
			}, testConfig.Samples)
		})
		Context("as regular user", func() {
			Measure("list all", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, "/v3/service_plans")
				})
			}, testConfig.Samples)
		})
	})

	Describe("GET /v3/service_plans/:guid", func() {
		var servicePlanGUID string
		Context("as admin", func() {
			BeforeEach(func() {
				servicePlanGUIDs := helpers.GetGUIDs(testSetup.AdminUserContext(), testConfig, "/v3/service_plans")
				Expect(servicePlanGUIDs).NotTo(BeNil())
				servicePlanGUID = servicePlanGUIDs[rand.Intn(len(servicePlanGUIDs))]
			})
			Measure("show one", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/service_plans/%s", servicePlanGUID))
				})
			}, testConfig.Samples)

		})
		Context("as regular user", func() {
			BeforeEach(func() {
				servicePlanGUID = getRandomLimitedServicePlanGuid()
			})
			Measure("show one", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf("/v3/service_plans/%s", servicePlanGUID))
				})
			}, testConfig.Samples)
		})
	})

	Describe("GET /v3/service_plans/:guid/visibility", func() {
		var servicePlanGUID string
		BeforeEach(func() {
			servicePlanGUID = getRandomLimitedServicePlanGuid()
		})
		Context("as admin", func() {
			Measure("show visibility", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/service_plans/%s/visibility", servicePlanGUID))
				})
			}, testConfig.Samples)
		})
		Context("as regular user", func() {
			Measure("show visibility", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf("/v3/service_plans/%s/visibility", servicePlanGUID))
				})
			}, testConfig.Samples)
		})
	})

	Describe("GET /v3/service_plans?service_offering_guids=", func() {
		var serviceOfferingGuidsList []string
		BeforeEach(func() {
			serviceOfferingGuidsList = nil
			serviceOfferingGuids := helpers.ExecuteSelectStatement(ccdb, ctx,
				"SELECT guid FROM services ORDER BY random() LIMIT 50")
			for _, guid := range serviceOfferingGuids {
				serviceOfferingGuidsList = append(serviceOfferingGuidsList, helpers.ConvertToString(guid))
			}
		})
		Context("as admin", func() {
			Measure("filter for list of service_offerings", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf(
						"/v3/service_plans?service_offering_guids=%v", strings.Join(serviceOfferingGuidsList[:], ",")))
				})
			}, testConfig.Samples)
			Measure(fmt.Sprintf("filter for list of service_offerings with page size %d", testConfig.LargePageSize), func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf(
						"/v3/service_plans?service_offering_guids=%v&per_page=%d",
						strings.Join(serviceOfferingGuidsList[:], ","), testConfig.LargePageSize))
				})
			}, testConfig.Samples)
		})
		Context("as regular user", func() {
			Measure("filter for list of service_offerings", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf(
						"/v3/service_plans?service_offering_guids=%v", strings.Join(serviceOfferingGuidsList[:], ",")))
				})
			}, testConfig.Samples)
		})
	})

	Describe("GET /v3/service_plans?service_instance_guids=", func() {
		var serviceInstanceGuidsList []string
		BeforeEach(func() {
			serviceInstanceGuidsList = nil
			serviceInstanceGuids := helpers.ExecuteSelectStatement(ccdb, ctx,
				"SELECT guid FROM service_instances ORDER BY random() LIMIT 50")
			for _, guid := range serviceInstanceGuids {
				serviceInstanceGuidsList = append(serviceInstanceGuidsList, helpers.ConvertToString(guid))
			}
		})
		Context("as admin", func() {
			Measure("filter for list of service_instances", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf(
						"/v3/service_plans?service_instance_guids=%v", strings.Join(serviceInstanceGuidsList[:], ",")))
				})
			}, testConfig.Samples)
			Measure(fmt.Sprintf("filter for list of service_instances with page size %d", testConfig.LargePageSize), func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.AdminUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf(
						"/v3/service_plans?service_instance_guids=%v&per_page=%d",
						strings.Join(serviceInstanceGuidsList[:], ","), testConfig.LargePageSize))
				})
			}, testConfig.Samples)
		})
		Context("as regular user", func() {
			Measure("filter for list of service_instances", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf(
						"/v3/service_plans?service_instance_guids=%v", strings.Join(serviceInstanceGuidsList[:], ",")))
				})
			}, testConfig.Samples)
			Measure(fmt.Sprintf("filter for list of service_instances with page size %d", testConfig.LargePageSize), func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.LongTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.LongTimeout, fmt.Sprintf(
						"/v3/service_plans?service_instance_guids=%v&per_page=%d",
						strings.Join(serviceInstanceGuidsList[:], ","), testConfig.LargePageSize))
				})
			}, testConfig.Samples)
		})
	})

	Describe("GET /v3/service_plans?organization_guids=&space_guids=", func() {
		var orgGuidsList []string
		var spaceGuidsList []string
		BeforeEach(func() {
			orgGuidsList = nil
			selectOrgGuidsStatement := fmt.Sprintf("SELECT guid FROM organizations WHERE name LIKE '%s-org-%%' ORDER BY random() LIMIT 50", testConfig.GetNamePrefix())
			orgGuids := helpers.ExecuteSelectStatement(ccdb, ctx, selectOrgGuidsStatement)
			for _, guid := range orgGuids {
				orgGuidsList = append(orgGuidsList, helpers.ConvertToString(guid))
			}
			spaceGuidsList = nil
			selectSpaceGuidsStatement := fmt.Sprintf("SELECT guid FROM spaces WHERE name LIKE '%s-space-%%' ORDER BY random() LIMIT 50", testConfig.GetNamePrefix())
			spaceGuids := helpers.ExecuteSelectStatement(ccdb, ctx, selectSpaceGuidsStatement)
			for _, guid := range spaceGuids {
				spaceGuidsList = append(spaceGuidsList, helpers.ConvertToString(guid))
			}
		})
		Context("as regular user", func() {
			Measure("filter by org and space guids", func(b Benchmarker) {
				workflowhelpers.AsUser(testSetup.RegularUserContext(), testConfig.BasicTimeout, func() {
					helpers.TimeCFCurl(b, testConfig.BasicTimeout, fmt.Sprintf(
						"/v3/service_plans?organization_guids=%v&space_guids=%v", strings.Join(orgGuidsList[:], ","), strings.Join(spaceGuidsList[:], ",")))
				})
			}, testConfig.Samples)
		})
	})
})

func getRandomLimitedServicePlanGuid() string {
	servicePlanGUIDs := helpers.ExecuteSelectStatement(ccdb, ctx,
		"SELECT s_p.guid FROM service_plans s_p INNER JOIN service_plan_visibilities s_p_v ON s_p.id = s_p_v.service_plan_id ORDER BY random() LIMIT 1")
	return helpers.ConvertToString(servicePlanGUIDs[0])
}
