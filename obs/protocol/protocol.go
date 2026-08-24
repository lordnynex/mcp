package protocol

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/andreykaipov/goobs/api/events/subscriptions"
)

const (
	EventsResourceURI = "obs://events"
	EventURIPrefix    = "obs://events/"
	EventURITemplate  = "obs://events/{eventType}"
	EventBufferSize   = 256
)

// AlwaysToolNames are registered before an OBS connection.
var AlwaysToolNames = []string{"Connect", "ConnectionStatus"}

// ConnectedSessionToolNames are registered only after a successful OBS connection.
var ConnectedSessionToolNames = []string{
	"Disconnect",
	"RequestBatch",
	"SubscribeEvents",
	"UnsubscribeEvents",
}

// RequestNames is the official obs-websocket request list, in protocol TOC order.
var RequestNames = []string{
	"GetVersion", "GetStats", "BroadcastCustomEvent", "CallVendorRequest",
	"GetHotkeyList", "TriggerHotkeyByName", "TriggerHotkeyByKeySequence", "Sleep",
	"GetPersistentData", "SetPersistentData", "GetSceneCollectionList",
	"SetCurrentSceneCollection", "CreateSceneCollection", "GetProfileList",
	"SetCurrentProfile", "CreateProfile", "RemoveProfile", "GetProfileParameter",
	"SetProfileParameter", "GetVideoSettings", "SetVideoSettings",
	"GetStreamServiceSettings", "SetStreamServiceSettings", "GetRecordDirectory",
	"SetRecordDirectory",
	"GetSourceActive", "GetSourceScreenshot", "SaveSourceScreenshot",
	"GetCanvasList",
	"GetSceneList", "GetGroupList", "GetCurrentProgramScene", "SetCurrentProgramScene",
	"GetCurrentPreviewScene", "SetCurrentPreviewScene", "CreateScene", "RemoveScene",
	"SetSceneName", "GetSceneSceneTransitionOverride", "SetSceneSceneTransitionOverride",
	"GetInputList", "GetInputKindList", "GetSpecialInputs", "CreateInput", "RemoveInput",
	"SetInputName", "GetInputDefaultSettings", "GetInputSettings", "SetInputSettings",
	"GetInputMute", "SetInputMute", "ToggleInputMute", "GetInputVolume", "SetInputVolume",
	"GetInputAudioBalance", "SetInputAudioBalance", "GetInputAudioSyncOffset",
	"SetInputAudioSyncOffset", "GetInputAudioMonitorType", "SetInputAudioMonitorType",
	"GetInputAudioTracks", "SetInputAudioTracks", "GetInputDeinterlaceMode",
	"SetInputDeinterlaceMode", "GetInputDeinterlaceFieldOrder",
	"SetInputDeinterlaceFieldOrder", "GetInputPropertiesListPropertyItems",
	"PressInputPropertiesButton",
	"GetTransitionKindList", "GetSceneTransitionList", "GetCurrentSceneTransition",
	"SetCurrentSceneTransition", "SetCurrentSceneTransitionDuration",
	"SetCurrentSceneTransitionSettings", "GetCurrentSceneTransitionCursor",
	"TriggerStudioModeTransition", "SetTBarPosition",
	"GetSourceFilterKindList", "GetSourceFilterList", "GetSourceFilterDefaultSettings",
	"CreateSourceFilter", "RemoveSourceFilter", "SetSourceFilterName", "GetSourceFilter",
	"SetSourceFilterIndex", "SetSourceFilterSettings", "SetSourceFilterEnabled",
	"GetSceneItemList", "GetGroupSceneItemList", "GetSceneItemId", "GetSceneItemSource",
	"CreateSceneItem", "RemoveSceneItem", "DuplicateSceneItem", "GetSceneItemTransform",
	"SetSceneItemTransform", "GetSceneItemEnabled", "SetSceneItemEnabled",
	"GetSceneItemLocked", "SetSceneItemLocked", "GetSceneItemIndex", "SetSceneItemIndex",
	"GetSceneItemBlendMode", "SetSceneItemBlendMode",
	"GetVirtualCamStatus", "ToggleVirtualCam", "StartVirtualCam", "StopVirtualCam",
	"GetReplayBufferStatus", "ToggleReplayBuffer", "StartReplayBuffer", "StopReplayBuffer",
	"SaveReplayBuffer", "GetLastReplayBufferReplay", "GetOutputList", "GetOutputStatus",
	"ToggleOutput", "StartOutput", "StopOutput", "GetOutputSettings", "SetOutputSettings",
	"GetStreamStatus", "ToggleStream", "StartStream", "StopStream", "SendStreamCaption",
	"GetRecordStatus", "ToggleRecord", "StartRecord", "StopRecord", "ToggleRecordPause",
	"PauseRecord", "ResumeRecord", "SplitRecordFile", "CreateRecordChapter",
	"GetMediaInputStatus", "SetMediaInputCursor", "OffsetMediaInputCursor",
	"TriggerMediaInputAction",
	"GetStudioModeEnabled", "SetStudioModeEnabled", "OpenInputPropertiesDialog",
	"OpenInputFiltersDialog", "OpenInputInteractDialog", "GetMonitorList",
	"OpenVideoMixProjector", "OpenSourceProjector",
}

// EventNames is the official obs-websocket event list, in protocol TOC order.
var EventNames = []string{
	"ExitStarted", "VendorEvent", "CustomEvent",
	"CurrentSceneCollectionChanging", "CurrentSceneCollectionChanged", "SceneCollectionListChanged",
	"CurrentProfileChanging", "CurrentProfileChanged", "ProfileListChanged",
	"CanvasCreated", "CanvasRemoved", "CanvasNameChanged",
	"SceneCreated", "SceneRemoved", "SceneNameChanged", "CurrentProgramSceneChanged",
	"CurrentPreviewSceneChanged", "SceneListChanged",
	"InputCreated", "InputRemoved", "InputNameChanged", "InputSettingsChanged",
	"InputActiveStateChanged", "InputShowStateChanged", "InputMuteStateChanged",
	"InputVolumeChanged", "InputAudioBalanceChanged", "InputAudioSyncOffsetChanged",
	"InputAudioTracksChanged", "InputAudioMonitorTypeChanged", "InputVolumeMeters",
	"CurrentSceneTransitionChanged", "CurrentSceneTransitionDurationChanged",
	"SceneTransitionStarted", "SceneTransitionEnded", "SceneTransitionVideoEnded",
	"SourceFilterListReindexed", "SourceFilterCreated", "SourceFilterRemoved",
	"SourceFilterNameChanged", "SourceFilterSettingsChanged", "SourceFilterEnableStateChanged",
	"SceneItemCreated", "SceneItemRemoved", "SceneItemListReindexed",
	"SceneItemEnableStateChanged", "SceneItemLockStateChanged", "SceneItemSelected",
	"SceneItemTransformChanged",
	"StreamStateChanged", "RecordStateChanged", "RecordFileChanged",
	"ReplayBufferStateChanged", "VirtualcamStateChanged", "ReplayBufferSaved",
	"MediaInputPlaybackStarted", "MediaInputPlaybackEnded", "MediaInputActionTriggered",
	"StudioModeStateChanged", "ScreenshotSaved",
}

// HighVolumeEvents maps official high-volume event names to Identify bits.
var HighVolumeEvents = map[string]int{
	"InputVolumeMeters":         subscriptions.InputVolumeMeters,
	"InputActiveStateChanged":   subscriptions.InputActiveStateChanged,
	"InputShowStateChanged":     subscriptions.InputShowStateChanged,
	"SceneItemTransformChanged": subscriptions.SceneItemTransformChanged,
}

// EventCategories maps EventSubscription category names to official event types.
// High-volume events are omitted and must be named explicitly.
var EventCategories = map[string][]string{
	"General": {
		"ExitStarted", "VendorEvent", "CustomEvent",
	},
	"Config": {
		"CurrentSceneCollectionChanging", "CurrentSceneCollectionChanged", "SceneCollectionListChanged",
		"CurrentProfileChanging", "CurrentProfileChanged", "ProfileListChanged",
	},
	"Canvases": {
		"CanvasCreated", "CanvasRemoved", "CanvasNameChanged",
	},
	"Scenes": {
		"SceneCreated", "SceneRemoved", "SceneNameChanged", "CurrentProgramSceneChanged",
		"CurrentPreviewSceneChanged", "SceneListChanged",
	},
	"Inputs": {
		"InputCreated", "InputRemoved", "InputNameChanged", "InputSettingsChanged",
		"InputMuteStateChanged", "InputVolumeChanged", "InputAudioBalanceChanged",
		"InputAudioSyncOffsetChanged", "InputAudioTracksChanged", "InputAudioMonitorTypeChanged",
	},
	"Transitions": {
		"CurrentSceneTransitionChanged", "CurrentSceneTransitionDurationChanged",
		"SceneTransitionStarted", "SceneTransitionEnded", "SceneTransitionVideoEnded",
	},
	"Filters": {
		"SourceFilterListReindexed", "SourceFilterCreated", "SourceFilterRemoved",
		"SourceFilterNameChanged", "SourceFilterSettingsChanged", "SourceFilterEnableStateChanged",
	},
	"Outputs": {
		"StreamStateChanged", "RecordStateChanged", "RecordFileChanged",
		"ReplayBufferStateChanged", "VirtualcamStateChanged", "ReplayBufferSaved",
	},
	"SceneItems": {
		"SceneItemCreated", "SceneItemRemoved", "SceneItemListReindexed",
		"SceneItemEnableStateChanged", "SceneItemLockStateChanged", "SceneItemSelected",
	},
	"MediaInputs": {
		"MediaInputPlaybackStarted", "MediaInputPlaybackEnded", "MediaInputActionTriggered",
	},
	"Vendors": {
		"VendorEvent",
	},
	"Ui": {
		"StudioModeStateChanged", "ScreenshotSaved",
	},
}

var knownEventNames = func() map[string]struct{} {
	out := make(map[string]struct{}, len(EventNames))
	for _, n := range EventNames {
		out[n] = struct{}{}
	}
	return out
}()

var knownRequestNames = func() map[string]struct{} {
	out := make(map[string]struct{}, len(RequestNames))
	for _, n := range RequestNames {
		out[n] = struct{}{}
	}
	return out
}()

// ConnectedToolNames are removed on disconnect.
func ConnectedToolNames() []string {
	return append(append([]string{}, ConnectedSessionToolNames...), RequestNames...)
}

// IsKnownRequest reports whether name is an official obs-websocket request.
func IsKnownRequest(name string) bool {
	_, ok := knownRequestNames[name]
	return ok
}

// IsKnownEvent reports whether name is an official obs-websocket event.
func IsKnownEvent(name string) bool {
	_, ok := knownEventNames[name]
	return ok
}

// EventResourceURI returns obs://events/{eventType}.
func EventResourceURI(eventType string) string {
	return EventURIPrefix + eventType
}

// ExpandEventTypes expands official event names and EventSubscription category aliases.
func ExpandEventTypes(names []string) ([]string, error) {
	if len(names) == 0 {
		return nil, fmt.Errorf("eventTypes must not be empty")
	}
	out := make([]string, 0, len(names))
	seen := make(map[string]struct{}, len(names))
	for _, n := range names {
		if cat, ok := EventCategories[n]; ok {
			for _, e := range cat {
				if _, dup := seen[e]; dup {
					continue
				}
				seen[e] = struct{}{}
				out = append(out, e)
			}
			continue
		}
		if !IsKnownEvent(n) {
			return nil, fmt.Errorf("unknown event type %q", n)
		}
		if _, dup := seen[n]; dup {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out, nil
}

// ParseTOCNames extracts bullet link labels from a markdown TOC section.
func ParseTOCNames(md, heading string) []string {
	sc := bufio.NewScanner(strings.NewReader(md))
	in := false
	var names []string
	for sc.Scan() {
		line := sc.Text()
		if !in {
			if strings.HasPrefix(line, heading) {
				in = true
			}
			continue
		}
		if strings.HasPrefix(line, "## ") && line != heading {
			break
		}
		if !strings.Contains(line, "[") || !strings.Contains(line, "](") {
			continue
		}
		start := strings.Index(line, "[")
		end := strings.Index(line[start:], "]")
		if start < 0 || end < 0 {
			continue
		}
		name := line[start+1 : start+end]
		if name == "" || strings.Contains(name, " ") || strings.Contains(name, "::") {
			continue
		}
		names = append(names, name)
	}
	return names
}
