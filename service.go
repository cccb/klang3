package main

import (
	"github.com/cameliot/alpaca"
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
		}
	}

}
