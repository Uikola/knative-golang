package parser

import (
	"github.com/Uikola/knative-golang/internal/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/sirikothe/gotextfsm"
	"strconv"
)

func ParseKubectlGetPodsOutput(output string) ([]entity.Pod, error) {
	pods := make([]entity.Pod, 0)

	template := `Value NAME (\S+)
Value READY (\S+)
Value STATUS (\S+)
Value RESTARTS (\d+)
Value AGE (\S+)

Start
 ^${NAME}\s+${READY}\s+${STATUS}\s+${RESTARTS}\s+${AGE} -> Record`
	fsm := gotextfsm.TextFSM{}
	err := fsm.ParseString(template)
	if err != nil {
		return nil, err
	}

	parser := gotextfsm.ParserOutput{}
	err = parser.ParseTextString(output, fsm, true)
	if err != nil {
		return nil, err
	}

	for _, record := range parser.Dict {
		if record["NAME"] == "NAME" {
			continue
		}

		restarts, _ := strconv.Atoi(record["RESTARTS"].(string))
		pod := entity.Pod{
			Name:     record["NAME"].(string),
			Ready:    record["READY"].(string),
			Status:   record["STATUS"].(string),
			Restarts: restarts,
			Age:      record["AGE"].(string),
		}
		pods = append(pods, pod)
	}
	return pods, nil
}

func ConvertPodsToCloudEvent(pods []entity.Pod) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetSource("kubectlgetpods-cmd")
	event.SetID(uuid.New().String())
	event.SetType("kubectlgetpods")
	if err := event.SetData(cloudevents.ApplicationJSON, pods); err != nil {
		return cloudevents.Event{}, err
	}
	if err := event.Validate(); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}
