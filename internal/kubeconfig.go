package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

func newKubeconfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubeconfig",
		Short: "Get the kubeconfig for a cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doKubeconfig(cmd.Context())
		},
	}

	return cmd
}

func doKubeconfig(ctx context.Context) error {
	provider := cluster.NewProvider(
		// cluster.ProviderWithLogger(logger),
		// runtime.GetDefault(logger),
		getDefault(),
	)
	clusters, err := provider.List()
	if err != nil {
		return err
	}
	if len(clusters) == 0 {
		slog.Info("No kind clusters found.")
		return nil
	}

	clusterOptions := []huh.Option[string]{}
	for _, cluster := range clusters {
		clusterOptions = append(clusterOptions, huh.NewOption(cluster, cluster))
	}

	selectedCluster := ""
	kubeconfigFile := "test.kubeconfig"
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Cluster").
				Options(clusterOptions...).
				Value(&selectedCluster),

			huh.NewInput().
				Title("kubeconfig file").
				Value(&kubeconfigFile).
				Validate(func(s string) error {
					if s == "" {
						return errors.New("you must supply a kubeconfig file value")
					}
					return nil
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("failed running form: %w", err)
	}

	cfg, err := provider.KubeConfig(selectedCluster, false)
	if err != nil {
		return err
	}
	os.WriteFile(kubeconfigFile, []byte(cfg), os.ModePerm)

	slog.Info("Saved kubeconfig", "cluster", selectedCluster, "path", kubeconfigFile)

	return nil
}

// GetDefault selected the default runtime from the environment override
func getDefault() cluster.ProviderOption {
	switch p := os.Getenv("KIND_EXPERIMENTAL_PROVIDER"); p {
	case "":
		return nil
	case "podman":
		slog.Warn("using podman due to KIND_EXPERIMENTAL_PROVIDER")
		return cluster.ProviderWithPodman()
	case "docker":
		slog.Warn("using docker due to KIND_EXPERIMENTAL_PROVIDER")
		return cluster.ProviderWithDocker()
	case "nerdctl", "finch", "nerdctl.lima":
		slog.Warn(fmt.Sprintf("using %s due to KIND_EXPERIMENTAL_PROVIDER", p))
		return cluster.ProviderWithNerdctl(p)
	default:
		slog.Warn(fmt.Sprintf("ignoring unknown value %q for KIND_EXPERIMENTAL_PROVIDER", p))
		return nil
	}
}
