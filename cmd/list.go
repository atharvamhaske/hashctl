package cmd

import (
	"fmt"

	"github.com/atharvamhaske/hashctl/internal/hasher"
	"github.com/atharvamhaske/hashctl/internal/tui"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available algorithms",
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	byCategory := hasher.GetAlgorithmsByCategory()
	categories := []hasher.Category{
		hasher.CategoryChecksum,
		hasher.CategoryFastHash,
		hasher.CategoryPasswordHash,
	}

	fmt.Println()
	fmt.Println(tui.LogoStyle.Render("hashctl") + tui.LogoAccent.Render(" algorithms"))
	fmt.Println()

	for _, cat := range categories {
		algs, ok := byCategory[cat]
		if !ok || len(algs) == 0 {
			continue
		}

		catStyle := tui.CategoryStyle
		if cat == hasher.CategoryPasswordHash {
			catStyle = tui.WarningCategoryStyle
		}
		fmt.Println(catStyle.Render(cat.String()))

		for _, alg := range algs {
			key := getKeyForAlg(alg.Name)
			name := tui.LabelStyle.Render(fmt.Sprintf("%-14s", key))
			desc := tui.MutedStyle.Render(alg.Description)
			fmt.Printf("  %s %s\n", name, desc)
		}
		fmt.Println()
	}
}

func getKeyForAlg(name string) string {
	for key, alg := range hasher.Registry {
		if alg.Name == name {
			return key
		}
	}
	return name
}
