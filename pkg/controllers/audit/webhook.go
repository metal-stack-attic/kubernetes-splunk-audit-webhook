package audit

import (
	// "encoding/base64"
	"io/ioutil"
	"net/http"

	hec "github.com/fuyufjh/splunk-hec-go"

	"go.uber.org/zap"
)

// Controller that retrieves audit webhooks and reports event to the resource-sink
type Controller struct {
	logger *zap.SugaredLogger
	client hec.HEC
}

// NewController returns a new accounting controller
func NewController(logger *zap.SugaredLogger, client hec.HEC) *Controller {
	controller := &Controller{
		logger: logger,
		client: client,
	}
	return controller
}

// AuditEvent handles an audit event
func (c *Controller) AuditEvent(response http.ResponseWriter, request *http.Request) {
	BodyStringBase64, _ := ioutil.ReadAll(request.Body)
	c.logger.Infow("received audit event", "request", BodyStringBase64)
	// BodyString, err := base64.URLEncoding.DecodeString(string(BodyStringBase64))
	// if err != nil {
	// 	c.logger.Errorw("error base64 decoding the body", "error", err)
	// 	response.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// c.logger.Infow("received audit event", "request base64decoded", BodyString)

	event := hec.NewEvent(BodyStringBase64)
	// event.SetTime(time.Now())

	c.logger.Infow("HEC Event", event.Event)

	err := c.client.WriteEvent(event)
	if err != nil {
		c.logger.Errorw("error sending event to splunk", "error", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}
