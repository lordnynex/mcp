package tools

import (
	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/record"
	"github.com/lordnynex/mcp/obs/session"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerRecordTools(s *mcp.Server, h *session.Host) {
	addRequest(s, h, "GetRecordStatus", "Gets the status of the record output.",
		func(c *goobs.Client, in *record.GetRecordStatusParams) (*record.GetRecordStatusResponse, error) {
			return c.Record.GetRecordStatus(in)
		})
	addRequest(s, h, "ToggleRecord", "Toggles the status of the record output.",
		func(c *goobs.Client, in *record.ToggleRecordParams) (*record.ToggleRecordResponse, error) {
			return c.Record.ToggleRecord(in)
		})
	addRequest(s, h, "StartRecord", "Starts the record output.",
		func(c *goobs.Client, in *record.StartRecordParams) (*record.StartRecordResponse, error) {
			return c.Record.StartRecord(in)
		})
	addRequest(s, h, "StopRecord", "Stops the record output.",
		func(c *goobs.Client, in *record.StopRecordParams) (*record.StopRecordResponse, error) {
			return c.Record.StopRecord(in)
		})
	addRequest(s, h, "ToggleRecordPause", "Toggles pause on the record output.",
		func(c *goobs.Client, in *record.ToggleRecordPauseParams) (*record.ToggleRecordPauseResponse, error) {
			return c.Record.ToggleRecordPause(in)
		})
	addRequest(s, h, "PauseRecord", "Pauses the record output.",
		func(c *goobs.Client, in *record.PauseRecordParams) (*record.PauseRecordResponse, error) {
			return c.Record.PauseRecord(in)
		})
	addRequest(s, h, "ResumeRecord", "Resumes the record output.",
		func(c *goobs.Client, in *record.ResumeRecordParams) (*record.ResumeRecordResponse, error) {
			return c.Record.ResumeRecord(in)
		})
	addRequest(s, h, "SplitRecordFile", "Splits the current file being recorded into a new file.",
		func(c *goobs.Client, in *record.SplitRecordFileParams) (*record.SplitRecordFileResponse, error) {
			return c.Record.SplitRecordFile(in)
		})
	addRequest(s, h, "CreateRecordChapter", "Adds a new chapter marker to the file currently being recorded.",
		func(c *goobs.Client, in *record.CreateRecordChapterParams) (*record.CreateRecordChapterResponse, error) {
			return c.Record.CreateRecordChapter(in)
		})
}
