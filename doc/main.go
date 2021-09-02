package main

import (
	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/apps"
	"github.com/blushft/go-diagrams/nodes/aws"
	"github.com/blushft/go-diagrams/nodes/generic"
	"log"
)

func main() {
	d, err := diagram.New(diagram.Filename("diagram"))
	if err != nil {
		log.Fatal(err)
	}

	server := generic.Compute.Rack().Label("acmevault server")
	vault := apps.Security.Vault().Label("Vault")
	iam := aws.Security.IdentityAndAccessManagementIam().Label("IAM")
	route53 := aws.Network.Route53().Label("Route53")

	dc := diagram.NewGroup("clients")
	dc.NewGroup("clients").
		Label("acmevault clients").
		Add(
			apps.Client.User().Label("Client"),
			apps.Client.User().Label("Client"),
		).
		ConnectAllTo(vault.ID(), diagram.Forward(), func(opt *diagram.EdgeOptions) {
			opt.Label = "Read Cert Bundle"
		})

	d.Connect(server, vault, func(options *diagram.EdgeOptions) {
		options.Label = "Read IAM role"
	})

	d.Connect(vault, server, func(options *diagram.EdgeOptions) {
		options.Label = "Write Certs"
	})

	d.Group(dc)

	d.Connect(vault, iam, func(opt *diagram.EdgeOptions) {
		opt.Label = "Create temp IAM credentials"
	})

	d.Connect(server, route53, func(opt *diagram.EdgeOptions) {
		opt.Label = "Solve DNS01 Challenge"
	})

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
