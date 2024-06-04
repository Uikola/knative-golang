package parser

import (
	"github.com/Uikola/knative-golang/internal/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/sirikothe/gotextfsm"
)

func ParseIfConfigOutput(output string) ([]entity.Interface, error) {
	networkInfo := make([]entity.Interface, 0)

	template := `Value INTERFACE (\S+)
Value MTU (\d+)
Value RX_PACKETS (\d+)
Value RX_BYTES (\d+)
Value TX_PACKETS (\d+)
Value TX_BYTES (\d+)

Start
  ^${INTERFACE}: flags=\S+  mtu ${MTU}
  ^\s+RX packets ${RX_PACKETS}  bytes ${RX_BYTES} \(\d+\.?\d* \w+\)
  ^\s+TX packets ${TX_PACKETS}  bytes ${TX_BYTES} \(\d+\.?\d* \w+\) -> Record
`
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
		inface := entity.NewInterface(
			record["INTERFACE"].(string),
			record["MTU"].(string),
			record["RX_PACKETS"].(string),
			record["RX_BYTES"].(string),
			record["TX_PACKETS"].(string),
			record["TX_BYTES"].(string),
		)
		networkInfo = append(networkInfo, inface)
	}

	return networkInfo, nil
}

func ConvertNetworkInfoToCloudEvent(networkInfo []entity.Interface) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetSource("ifconfig-cmd")
	event.SetID(uuid.New().String())
	event.SetType("ifconfig")
	if err := event.SetData(cloudevents.ApplicationJSON, networkInfo); err != nil {
		return cloudevents.Event{}, err
	}
	if err := event.Validate(); err != nil {
		return cloudevents.Event{}, err
	}

	return event, nil
}
