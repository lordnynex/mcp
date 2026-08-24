package tools

import (
	"testing"

	"github.com/lordnynex/gest"
)

func TestRequestBatch(t *testing.T) {
	gest.Run(t, "RequestBatch", func(s *gest.S) {
		s.It("rejects unknown request types without sending", func(t *gest.T) {
			_, err := runRequestBatch(t.Context(), nil, nil, batchInput{
				Requests: []batchRequest{{RequestType: "InventedRequest"}},
			})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.MatchRegexp("UnknownRequestType"))
		})

		s.It("rejects SerialFrame", func(t *gest.T) {
			exec := batchSerialFrame
			_, err := runRequestBatch(t.Context(), nil, nil, batchInput{
				ExecutionType: &exec,
				Requests:      []batchRequest{{RequestType: "GetVersion"}},
			})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.MatchRegexp("SerialFrame"))
		})

		s.It("rejects Sleep outside SerialRealtime", func(t *gest.T) {
			exec := batchParallel
			_, err := runRequestBatch(t.Context(), nil, nil, batchInput{
				ExecutionType: &exec,
				Requests:      []batchRequest{{RequestType: "Sleep"}},
			})
			t.Expect(err).NotTo(gest.BeNil())
			t.Expect(err.Error()).To(gest.MatchRegexp("Sleep"))
		})

		s.It("rejects an empty request list", func(t *gest.T) {
			_, err := runRequestBatch(t.Context(), nil, nil, batchInput{})
			t.Expect(err).NotTo(gest.BeNil())
		})
	})
}
