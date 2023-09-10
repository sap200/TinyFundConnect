package pangealogger

import (
	"context"
	"fmt"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/audit"
	"github.com/sap200/TinyFundConnect/secret"
)

func Log(event *CustomSchemaEvent) {
	token := secret.PANGEA_AUTHN_TOKEN

	auditcli, err := audit.New(
		&pangea.Config{
			Token:    token,
			Domain:   secret.PANGEA_DOMAIN,
			ConfigID: secret.PANGEA_AUDIT_SCHEMA_CONFIG_ID,
		},
		audit.WithCustomSchema(CustomSchemaEvent{}),
	)
	if err != nil {
		fmt.Println("Failed to create audit client")
		return
	}

	ctx := context.Background()

	fmt.Printf("Logging: %s\n", event)

	lr, err := auditcli.Log(ctx, event, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	e := (lr.Result.EventEnvelope.Event).(*CustomSchemaEvent)
	fmt.Printf("Logged event: %s", pangea.Stringify(e))
}
