package function

import (
	"context"
	"fmt"
	"function/entity"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/sirikothe/gotextfsm"
)

// Handle an HTTP Request.
func Handle(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, error) {
	var request entity.CloudEventRequest
	if err := event.DataAs(&request); err != nil {
		return nil, fmt.Errorf("failed to get event data: %v", err)
	}
	if request.StdinCmdIsConfig() {
		networkInfo, err := parseIfConfigOutput(request.Data.Stdout)
		if err != nil {
			return nil, fmt.Errorf("error while parsing ifconfig output: %v", err)
		}

		responseEvent, err := convertNetworkInfoToCloudEvent(networkInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to convert network info to cloud event: %v", err)
		}
		return responseEvent, nil
	}

	return &event, nil
}

func parseIfConfigOutput(output string) ([]entity.Interface, error) {
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

func convertNetworkInfoToCloudEvent(networkInfo []entity.Interface) (*cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetSource("ifconfig-cmd")
	event.SetID(uuid.New().String())
	event.SetType("ifconfig")
	if err := event.SetData(cloudevents.ApplicationJSON, networkInfo); err != nil {
		return nil, err
	}
	if err := event.Validate(); err != nil {
		return nil, err
	}

	return &event, nil
}
