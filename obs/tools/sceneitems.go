package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sceneitems"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerSceneItemsTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetSceneItemList", "Gets a list of all scene items in a scene.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemListParams) (*sceneitems.GetSceneItemListResponse, error) {
			return c.SceneItems.GetSceneItemList(in)
		})
	addRequest(s, h, "GetGroupSceneItemList", "Basically GetSceneItemList, but for groups.",
		func(c *goobs.Client, in *sceneitems.GetGroupSceneItemListParams) (*sceneitems.GetGroupSceneItemListResponse, error) {
			return c.SceneItems.GetGroupSceneItemList(in)
		})
	addRequest(s, h, "GetSceneItemId", "Searches a scene for a source, and returns its id.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemIdParams) (*sceneitems.GetSceneItemIdResponse, error) {
			return c.SceneItems.GetSceneItemId(in)
		})
	addRequest(s, h, "GetSceneItemSource", "Gets the source associated with a scene item.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemSourceParams) (*sceneitems.GetSceneItemSourceResponse, error) {
			return c.SceneItems.GetSceneItemSource(in)
		})
	addRequest(s, h, "CreateSceneItem", "Creates a new scene item using a source.",
		func(c *goobs.Client, in *sceneitems.CreateSceneItemParams) (*sceneitems.CreateSceneItemResponse, error) {
			return c.SceneItems.CreateSceneItem(in)
		})
	addRequest(s, h, "RemoveSceneItem", "Removes a scene item from a scene.",
		func(c *goobs.Client, in *sceneitems.RemoveSceneItemParams) (*sceneitems.RemoveSceneItemResponse, error) {
			return c.SceneItems.RemoveSceneItem(in)
		})
	addRequest(s, h, "DuplicateSceneItem", "Duplicates a scene item, copying all transform and crop info.",
		func(c *goobs.Client, in *sceneitems.DuplicateSceneItemParams) (*sceneitems.DuplicateSceneItemResponse, error) {
			return c.SceneItems.DuplicateSceneItem(in)
		})
	addRequest(s, h, "GetSceneItemTransform", "Gets the transform and crop info of a scene item.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemTransformParams) (*sceneitems.GetSceneItemTransformResponse, error) {
			return c.SceneItems.GetSceneItemTransform(in)
		})
	addRequest(s, h, "SetSceneItemTransform", "Sets the transform and crop info of a scene item.",
		func(c *goobs.Client, in *sceneitems.SetSceneItemTransformParams) (*sceneitems.SetSceneItemTransformResponse, error) {
			return c.SceneItems.SetSceneItemTransform(in)
		})
	addRequest(s, h, "GetSceneItemEnabled", "Gets the enable state of a scene item.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemEnabledParams) (*sceneitems.GetSceneItemEnabledResponse, error) {
			return c.SceneItems.GetSceneItemEnabled(in)
		})
	addRequest(s, h, "SetSceneItemEnabled", "Sets the enable state of a scene item.",
		func(c *goobs.Client, in *sceneitems.SetSceneItemEnabledParams) (*sceneitems.SetSceneItemEnabledResponse, error) {
			return c.SceneItems.SetSceneItemEnabled(in)
		})
	addRequest(s, h, "GetSceneItemLocked", "Gets the lock state of a scene item.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemLockedParams) (*sceneitems.GetSceneItemLockedResponse, error) {
			return c.SceneItems.GetSceneItemLocked(in)
		})
	addRequest(s, h, "SetSceneItemLocked", "Sets the lock state of a scene item.",
		func(c *goobs.Client, in *sceneitems.SetSceneItemLockedParams) (*sceneitems.SetSceneItemLockedResponse, error) {
			return c.SceneItems.SetSceneItemLocked(in)
		})
	addRequest(s, h, "GetSceneItemIndex", "Gets the index position of a scene item in a scene.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemIndexParams) (*sceneitems.GetSceneItemIndexResponse, error) {
			return c.SceneItems.GetSceneItemIndex(in)
		})
	addRequest(s, h, "SetSceneItemIndex", "Sets the index position of a scene item in a scene.",
		func(c *goobs.Client, in *sceneitems.SetSceneItemIndexParams) (*sceneitems.SetSceneItemIndexResponse, error) {
			return c.SceneItems.SetSceneItemIndex(in)
		})
	addRequest(s, h, "GetSceneItemBlendMode", "Gets the blend mode of a scene item.",
		func(c *goobs.Client, in *sceneitems.GetSceneItemBlendModeParams) (*sceneitems.GetSceneItemBlendModeResponse, error) {
			return c.SceneItems.GetSceneItemBlendMode(in)
		})
	addRequest(s, h, "SetSceneItemBlendMode", "Sets the blend mode of a scene item.",
		func(c *goobs.Client, in *sceneitems.SetSceneItemBlendModeParams) (*sceneitems.SetSceneItemBlendModeResponse, error) {
			return c.SceneItems.SetSceneItemBlendMode(in)
		})
}
