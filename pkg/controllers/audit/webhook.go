package audit

import (
	"io/ioutil"
	"net/http"
	"time"

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
	BodyString, err := ioutil.ReadAll(request.Body)
	c.logger.Infow("received audit event", "request", BodyString)

	event := hec.NewEvent(request.Body)
	event.SetTime(time.Now())

	err := c.client.WriteEvent(event)
	if err != nil {
		c.logger.Errorw("error sending event to splunk", "error", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.WriteHeader(http.StatusOK)
}
