package main

import (
	"github.com/cameliot/alpaca"

	"fmt"
)

type SamplerSvc struct {
	repo *Repository
}

func NewSamplerSvc(repository *Repository) *SamplerSvc {
	svc := &SamplerSvc{
		repo: repository,
	}

	return svc
}

func (self *SamplerSvc) handleGroupsListRequest() alpaca.Action {
	groups := self.repo.Groups()

	return GroupsListSuccess(groups)
}

func (self *SamplerSvc) handleSamplesListRequest(request *SamplesListPayload) alpaca.Action {
	samples := self.repo.Samples(request.Group)

	return SamplesListSuccess(request.Group, samples)
}

func (self *SamplerSvc) handleSampleStartRequest(dispatch alpaca.Dispatch, request *SampleIdPayload) alpaca.Action {
	sample := self.repo.GetSampleById(request.SampleId)
	if sample == nil {
		return SampleStartError(request.SampleId, 404, fmt.Errorf("Sample not found"))
	}

	// Start playback
	done, err := sample.Start()
	if err != nil {
		return SampleStartError(request.SampleId, 501, err)
	}

	go func() {
		// Wait until done, dispatch finished callback
		<-done

		dispatch(SampleStopSuccess(sample.Id))
	}()

	return SampleStartSuccess(sample.Id)
}

func (self *SamplerSvc) handleSampleStopRequest(request *SampleIdPayload) alpaca.Action {
	if request.SampleId == 0 {
		for _, sample := range self.repo.AllSamples() {
			sample.Stop() // Just stop do not care
		}

		return SampleStopSuccess(0)
	}

	sample := self.repo.GetSampleById(request.SampleId)
	if sample == nil {
		return SampleStopError(0, 404, fmt.Errorf("Sample not found"))
	}

	sample.Stop()

	return SampleStopSuccess(sample.Id)
}

func (self *SamplerSvc) Handle(actions alpaca.Actions, dispatch alpaca.Dispatch) {

	for action := range actions {
		switch action.Type {
		case GROUPS_LIST_REQUEST:
			dispatch(self.handleGroupsListRequest())
			break

		case SAMPLES_LIST_REQUEST:
			var request SamplesListPayload
			action.DecodePayload(&request)
			dispatch(self.handleSamplesListRequest(&request))
			break

		case SAMPLE_START_REQUEST:
			var request SampleIdPayload
			action.DecodePayload(&request)
			dispatch(self.handleSampleStartRequest(dispatch, &request))
			break

		case SAMPLE_STOP_REQUEST:
			var request SampleIdPayload
			action.DecodePayload(&request)
			dispatch(self.handleSampleStopRequest(&request))
			break
		}
	}

}
