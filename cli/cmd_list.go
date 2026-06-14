package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) listCmd() *cobra.Command {
	var sort, mangaType string
	var pages int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List manga from the MangaFire catalog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			limit := a.effectiveLimit(0)
			if pages == 0 {
				pages = 1
			}
			items, err := a.client.ListMangas(cmd.Context(), sort, mangaType, pages, limit)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(items, len(items))
		},
	}
	cmd.Flags().StringVarP(&sort, "sort", "s", "most_viewed", "sort order: most_viewed|latest_updated|new_release|title_az")
	cmd.Flags().StringVarP(&mangaType, "type", "t", "", "filter by type: manga|manhwa|manhua|novel|one-shot|doujinshi")
	cmd.Flags().IntVarP(&pages, "pages", "p", 1, "number of pages to fetch (30 items per page)")
	return cmd
}
