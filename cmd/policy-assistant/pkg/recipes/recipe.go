package recipes

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/network-policy-api/policy-assistant/pkg/connectivity/probe"
	"sigs.k8s.io/network-policy-api/policy-assistant/pkg/generator"
	"sigs.k8s.io/network-policy-api/policy-assistant/pkg/matcher"
	"sigs.k8s.io/network-policy-api/policy-assistant/pkg/utils"
)

type Recipe struct {
	PolicyYamls []string
	Resources   *probe.Resources
	Protocol    v1.Protocol
	Port        int
}

func (r *Recipe) Policies() []*networkingv1.NetworkPolicy {
	var policies []*networkingv1.NetworkPolicy
	for _, yamlString := range r.PolicyYamls {
		netpol, err := utils.ParseYaml[networkingv1.NetworkPolicy]([]byte(yamlString))
		utils.DoOrDie(err)
		policies = append(policies, netpol)
	}
	return policies
}

func (r *Recipe) RunProbe() *probe.Table {
	runner := probe.NewSimulatedRunner(matcher.BuildNetworkPolicies(true, r.Policies()), &probe.JobBuilder{TimeoutSeconds: 5})
	return runner.RunProbeForConfig(generator.NewProbeConfig(intstr.FromInt(r.Port), r.Protocol, generator.ProbeModeServiceName), r.Resources)
}

var AllRecipes = []*Recipe{
	{[]string{Recipe01}, Resources01, v1.ProtocolTCP, 80},
	{[]string{Recipe02}, Resources02, v1.ProtocolTCP, 80},
	{[]string{Recipe01, Recipe02A}, Resources02A, v1.ProtocolTCP, 80},
	{[]string{Recipe03}, Resources03, v1.ProtocolTCP, 80},
	{[]string{Recipe04}, Resources04, v1.ProtocolTCP, 80},
	{[]string{Recipe01, Recipe05}, Resources05, v1.ProtocolTCP, 80},
	{[]string{Recipe06}, Resources06, v1.ProtocolTCP, 80},
	{[]string{Recipe07}, Resources07, v1.ProtocolTCP, 80},
	{[]string{Recipe01, Recipe08}, Resources08, v1.ProtocolTCP, 80},
	{[]string{Recipe09}, Resources09, v1.ProtocolTCP, 5000},
	{[]string{Recipe10}, Resources10, v1.ProtocolTCP, 80},
	{[]string{Recipe11_1}, Resources11_1, v1.ProtocolTCP, 80},
	{[]string{Recipe11_2}, Resources11_2, v1.ProtocolTCP, 53},
	{[]string{Recipe12}, Resources12, v1.ProtocolTCP, 80},
	{[]string{Recipe14}, Resources14, v1.ProtocolTCP, 80},
}

func Run() {
	for _, recipe := range AllRecipes {
		table := recipe.RunProbe()

		fmt.Printf("Policies:\n%s\n", matcher.BuildNetworkPolicies(true, recipe.Policies()).ExplainTable())

		fmt.Printf("resources:\n%s\n", recipe.Resources.RenderTable())

		fmt.Printf("Results:\n%s\n", table.RenderTable())

		fmt.Printf("Ingress:\n%s\n", table.RenderIngress())

		fmt.Printf("Egress:\n%s\n", table.RenderEgress())

		fmt.Printf("\n\n\n")
	}
}
