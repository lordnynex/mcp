package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/filters"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerFiltersTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetSourceFilterKindList", "Gets an array of all available source filter kinds.",
		func(c *goobs.Client, in *filters.GetSourceFilterKindListParams) (*filters.GetSourceFilterKindListResponse, error) {
			return c.Filters.GetSourceFilterKindList(in)
		})
	addRequest(s, h, "GetSourceFilterList", "Gets an array of all of a source's filters.",
		func(c *goobs.Client, in *filters.GetSourceFilterListParams) (*filters.GetSourceFilterListResponse, error) {
			return c.Filters.GetSourceFilterList(in)
		})
	addRequest(s, h, "GetSourceFilterDefaultSettings", "Gets the default settings for a filter kind.",
		func(c *goobs.Client, in *filters.GetSourceFilterDefaultSettingsParams) (*filters.GetSourceFilterDefaultSettingsResponse, error) {
			return c.Filters.GetSourceFilterDefaultSettings(in)
		})
	addRequest(s, h, "CreateSourceFilter", "Creates a new filter, adding it to the specified source.",
		func(c *goobs.Client, in *filters.CreateSourceFilterParams) (*filters.CreateSourceFilterResponse, error) {
			return c.Filters.CreateSourceFilter(in)
		})
	addRequest(s, h, "RemoveSourceFilter", "Removes a filter from a source.",
		func(c *goobs.Client, in *filters.RemoveSourceFilterParams) (*filters.RemoveSourceFilterResponse, error) {
			return c.Filters.RemoveSourceFilter(in)
		})
	addRequest(s, h, "SetSourceFilterName", "Sets the name of a source filter (rename).",
		func(c *goobs.Client, in *filters.SetSourceFilterNameParams) (*filters.SetSourceFilterNameResponse, error) {
			return c.Filters.SetSourceFilterName(in)
		})
	addRequest(s, h, "GetSourceFilter", "Gets the info for a specific source filter.",
		func(c *goobs.Client, in *filters.GetSourceFilterParams) (*filters.GetSourceFilterResponse, error) {
			return c.Filters.GetSourceFilter(in)
		})
	addRequest(s, h, "SetSourceFilterIndex", "Sets the index position of a filter on a source.",
		func(c *goobs.Client, in *filters.SetSourceFilterIndexParams) (*filters.SetSourceFilterIndexResponse, error) {
			return c.Filters.SetSourceFilterIndex(in)
		})
	addRequest(s, h, "SetSourceFilterSettings", "Sets the settings of a source filter.",
		func(c *goobs.Client, in *filters.SetSourceFilterSettingsParams) (*filters.SetSourceFilterSettingsResponse, error) {
			return c.Filters.SetSourceFilterSettings(in)
		})
	addRequest(s, h, "SetSourceFilterEnabled", "Sets the enable state of a source filter.",
		func(c *goobs.Client, in *filters.SetSourceFilterEnabledParams) (*filters.SetSourceFilterEnabledResponse, error) {
			return c.Filters.SetSourceFilterEnabled(in)
		})
}
