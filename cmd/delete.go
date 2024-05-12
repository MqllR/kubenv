package cmd

import (
	"fmt"
	"slices"

	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"k8s.io/klog"

	"github.com/mqllr/kubenv/cmd/helpers"
	"github.com/mqllr/kubenv/pkg/k8s"
)

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "delete a context",
		Aliases: []string{"rm"},
		Run: func(cmd *cobra.Command, args []string) {
			deleteContext()
		},
	}
}

func deleteContext() {
	k, err := helpers.NewKubeConfig()
	if err != nil {
		klog.Fatalf("cannot read the kubeconfig file: %s", err)
	}

	context, err := k.GetContextByContextName(k.CurrentContext)
	if err != nil {
		klog.Errorf("cannot retrieve the context %s: %s", k.CurrentContext, err.Error())
	}

	choiceContext := false
	err = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Do you really want to delete context %s ?", k.CurrentContext),
	}, &choiceContext)
	if err != nil {
		klog.Fatal("failed to get context answer")
	}

	if !choiceContext {
		return
	}

	k.Contexts = slices.DeleteFunc(k.Contexts, func(n *k8s.ContextWithName) bool {
		return n.Name == k.CurrentContext
	})

	choiceUsers := false
	err = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Do you want to delete user %s ?", context.User),
	}, &choiceUsers)
	if err != nil {
		klog.Fatal("failed to get user answer")
	}

	if choiceUsers {
		k.Users = slices.DeleteFunc(k.Users, func(n *k8s.UserWithName) bool {
			return n.Name == context.User
		})
	}

	choiceClusters := false
	err = survey.AskOne(&survey.Confirm{
		Message: fmt.Sprintf("Do you want to delete cluster %s ?", context.Cluster),
	}, &choiceClusters)
	if err != nil {
		klog.Fatal("failed to get cluster answer")
	}

	if choiceClusters {
		k.Clusters = slices.DeleteFunc(k.Clusters, func(n *k8s.ClusterWithName) bool {
			return n.Name == context.Cluster
		})
	}

	if err := helpers.SaveKubeConfig(k); err != nil {
		klog.Fatalf("error when saving config file: %s", err.Error())
	}

	fmt.Printf("%v Config saved!", promptui.IconGood)
}
