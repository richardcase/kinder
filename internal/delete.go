package internal

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"sigs.k8s.io/kind/pkg/cluster"
)

func newDeleteCommand() *cobra.Command {
	deleteAll := false
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete kind cluster(s)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doDelete(cmd.Context(), deleteAll)
		},
	}

	cmd.Flags().BoolVarP(&deleteAll, "all", "a", false, "delete all clusters")

	return cmd
}

func doDelete(ctx context.Context, deleteAll bool) error {
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

	selectedClusters := []string{}

	if !deleteAll {
		clusterOptions := []huh.Option[string]{}
		for _, cluster := range clusters {
			clusterOptions = append(clusterOptions, huh.NewOption(cluster, cluster))
		}

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewMultiSelect[string]().
					Title("Clusters to delete").
					Options(clusterOptions...).
					Value(&selectedClusters),
			),
		)

		if err := form.Run(); err != nil {
			return fmt.Errorf("failed running form: %w", err)
		}
	} else {
		slog.Info("Deleting ALL clusters")
		selectedClusters = clusters
	}

	for _, clusterToDelete := range selectedClusters {
		slog.Info("Deleting cluster", "name", clusterToDelete)
		if err := provider.Delete(clusterToDelete, ""); err != nil {
			// TODO: log error properly
			fmt.Printf("Failed to delete cluster %s\n", clusterToDelete)
		}
	}

	return nil
}
