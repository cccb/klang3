package main

import (
	"github.com/cameliot/alpaca"
)

/* Actions */
const GROUPS_LIST_REQUEST = "@sampler/GROUPS_LIST_REQUEST"
const GROUPS_LIST_SUCCESS = "@sampler/GROUPS_LIST_SUCCESS"
const GROUPS_LIST_ERROR = "@sampler/GROUPS_LIST_ERROR"

const SAMPLES_LIST_REQUEST = "@sampler/SAMPLES_LIST_REQUEST"
const SAMPLES_LIST_SUCCESS = "@sampler/SAMPLES_LIST_SUCCESS"
const SAMPLES_LIST_ERROR = "@sampler/SAMPLES_LIST_ERROR"

const SAMPLE_START_REQUEST = "@sampler/SAMPLE_START_REQUEST"
const SAMPLE_START_SUCCESS = "@sampler/SAMPLE_START_SUCCESS"
const SAMPLE_START_ERROR = "@sampler/SAMPLE_START_ERROR"

const SAMPLE_STOP_REQUEST = "@sampler/SAMPLE_STOP_REQUEST"
const SAMPLE_STOP_SUCCESS = "@sampler/SAMPLE_STOP_SUCCESS"
const SAMPLE_STOP_ERROR = "@sampler/SAMPLE_STOP_ERROR"

/* Payloads */
type GroupsPayload struct {
	Groups []string `json:"groups"`
}

type SamplesPayload struct {
	Group   string    `json:"group"`
	Samples []*Sample `json:"samples"`
}

type SampleIdPayload struct {
	SampleId int `json:"sample_id"`
}

type SamplePayload struct {
	Sample *Sample `json:"sample"`
}

type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SampleIdErrorPayload struct {
	SampleId int    `json:"sample_id"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
}

/* Action Creators */
func GroupsListSuccess(groups []string) aplaca.Action {
	return alpaca.Action{
		Type: GROUPS_LIST_SUCCESS,
		Payload: GroupsPayload{
			Groups: groups,
		},
	}
}

func GroupsListError(code int, err error) alpaca.Action {
	return alpaca.Action{
		Type: GROUPS_LIST_ERROR,
		Payload: ErrorPayload{
			Code:    code,
			Message: err.Error(),
		},
	}
}

func SamplesListSuccess(group string, samples []*Sample) alpaca.Action {
	return alpaca.Action{
		Type: SAMPLES_LIST_SUCCESS,
		Payload: SamplesPayload{
			Group:   group,
			Samples: samples,
		},
	}
}

func SamplesListError(code int, err error) alpaca.Action {
	return alpaca.Action{
		Type: SAMPLES_LIST_ERROR,
		Payload: ErrorPayload{
			Code:    code,
			Message: err.Error(),
		},
	}
}

func SampleStartSuccess(sampleId int) alpaca.Action {
	return alpaca.Action{
		Type: SAMPLE_START_SUCCESS,
		Payload: SampleIdPayload{
			SampleId: sampleId,
		},
	}
}

func SampleStartError(sampleId int, code int, err error) alpaca.Action {
	return alpaca.Action{
		Type: SAMPLE_START_SUCCESS,
		Payload: SampleIdErrorPayload{
			SampleId: sampleId,
			Code:     code,
			Message:  err.Error(),
		},
	}
}

func handle(actions alpaca.Actions, dispatch alpaca.Dispatch) {

}
