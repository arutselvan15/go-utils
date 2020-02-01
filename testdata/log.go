package testdata

import (
	"github.com/arutselvan15/go-utils/log"
	lc "github.com/arutselvan15/go-utils/logconstants"
)

type SampleObj struct {
	Name, Description string
	Price             int
	Stock             struct {
		Total, Returned, Sold int
	}
	Category []string
}

var l = log.NewLogger().SetFormatterType(log.TextFormatterType)

var cObj = SampleObj{
	Name:        "iphone",
	Description: "testdata obj",
	Price:       100,
	Stock: struct {
		Total, Returned, Sold int
	}{Total: 50, Returned: 0, Sold: 20},
	Category: []string{"cellphones"},
}

var nObj = SampleObj{
	Name:        "iphone",
	Description: "testdata obj",
	Price:       150,
	Stock: struct {
		Total, Returned, Sold int
	}{Total: 50, Returned: 0, Sold: 30},
	Category: []string{"cellphones", "electronics"},
}

func mutate() {
	l.SetStep("mutate").SetStepState(lc.Start)
	l.Debug("mutate prep")
	l.Debug("mutate apply")
	l.SetStepState(lc.Complete).Debug("mutate complete")
}

func validate() {
	l.SetStep("validate").SetStepState(lc.Start)
	l.Debug("validate prep")
	l.Debug("validate check")
	l.SetStepState(lc.Complete).Debug("validate complete")
}

func webhook() {
	l.SetComponent("webhook")
	mutate()
	validate()
}

func controller() {
	l.SetComponent("controller")
	process()
	notify()
}

func process() {
	l.SetStep("process").SetStepState(lc.Start)
	l.Debug("process prep")
	l.Debug("process item")
	l.LogAuditAPI("GET", "/orders/iphone", "", "", 404)
	l.LogAuditAPI("POST", "/orders", "{name: iphone}", "{name: iphone}", 201)
	l.LogAuditAPI("GET", "/orders/iphone", "", "{name: iphone}", 200)

	l.SetStepState(lc.Complete).Debug("process complete")
}

func notify() {
	l.SetStep("notify").SetStepState(lc.Start)
	l.Debug("notify prep")
	l.Debug("notify email")
	l.Debug("notify msg")
	l.SetStepState(lc.Complete).Debug("notify complete")
}

func SampleLogging() {
	l.SetCluster("minikube").SetApplication("estore").SetResource("product").SetLevel(log.DebugLevel)
	l.SetObjectState(lc.Received)
	l.SetObjectName("iphone").SetUser("johnny")

	l.LogAuditObject(cObj, nObj)

	l.SetOperation(lc.Create).SetObjectState(lc.Processing)

	webhook()
	l.LogAuditEvent("webhook complete")

	controller()
	l.LogAuditEvent("controller complete")

	l.SetObjectState(lc.Successful).Debug("product complete")
}
