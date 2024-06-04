package parser

import (
	"github.com/Uikola/knative-golang/internal/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/sirikothe/gotextfsm"
)

func ParseKubectlGetNsOutput(output string) ([]entity.NameSpace, error) {
	namespaceInfo := make([]entity.NameSpace, 0)

	template := `Value NAME (\S+)
Value STATUS (\S+)
Value AGE (\S+)

Start
  ^${NAME}\s+${STATUS}\s+${AGE} -> Record`
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
		namespace := entity.NameSpace{
			Name:   record["NAME"].(string),
			Status: record["STATUS"].(string),
			Age:    record["AGE"].(string),
		}
		namespaceInfo = append(namespaceInfo, namespace)
	}
	return namespaceInfo, nil
}

func ConvertNamespaceInfoToCloudEvent(namespaceInfo []entity.NameSpace) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetSource("kubectlgetns-cmd")
	event.SetID(uuid.New().String())
	event.SetType("kubectlgetns")
	if err := event.SetData(cloudevents.ApplicationJSON, namespaceInfo); err != nil {
		return cloudevents.Event{}, err
	}
	if err := event.Validate(); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}
