package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

func NewCategoriesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "categories",
		Short: "List all unique categories",
		Long:  `Display all unique categories from time entries in alphabetical order.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			categories, err := timeService.GetCategories(ctx)
			if err != nil {
				return fmt.Errorf("failed to get categories: %w", err)
			}

			if len(categories) == 0 {
				fmt.Println("No categories found.")
				return nil
			}

			sort.Strings(categories)
			for _, category := range categories {
				fmt.Println(category)
			}

			return nil
		},
	}

	return cmd
}
