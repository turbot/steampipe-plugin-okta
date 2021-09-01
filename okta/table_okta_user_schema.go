package okta

import (
	"context"
	// "fmt"
	"strings"
	// "time"
	// "gjson"
	"encoding/json"

	"github.com/okta/okta-sdk-golang/v2/okta"
	// "github.com/okta/okta-sdk-golang/v2/okta/query"
	// "github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	// "github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOktaUserSchema() *plugin.Table {
	return &plugin.Table{
		Name:        "okta_user_schema",
		Description: "Okta User Schema",
		List: &plugin.ListConfig{
			ParentHydrate: listOktaUserTypes,
			Hydrate: listOktaUserSchemas,
		},
		Columns: []*plugin.Column{
			// Top Columns
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the Policy."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the Policy."},
			// {Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the Policy."},
			// {Name: "created", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was created."},

			// // Other Columns
			// {Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the Policy was last modified."},
			// {Name: "priority", Type: proto.ColumnType_INT, Description: "Priority of the Policy."},
			// {Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the Policy: ACTIVE or INACTIVE."},
			// {Name: "system", Type: proto.ColumnType_BOOL, Description: "This is set to true on system policies, which cannot be deleted."},

			// // JSON Columns
			// {Name: "conditions", Type: proto.ColumnType_JSON, Description: "Conditions for Policy."},
			// {Name: "rules", Type: proto.ColumnType_JSON, Transform: transform.FromField("Embedded.rules"), Description: "Each Policy may contain one or more Rules. Rules, like Policies, contain conditions that must be satisfied for the Rule to be applied."},
			// {Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings of the password policy."},

			// // Steampipe Columns
			// {Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("Name"), Description: titleDescription},
		},
	}
}

type SchemaDetails struct {
	Href string
}

func listOktaUserSchemas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	client, err := Connect(ctx, d)

	// input := &query.Params{}
	if err != nil {
		logger.Error("listOktaUserSchemas", "connect", err)
		return nil, err
	}

	var schemaId string
	if h.Item != nil {
		userTypeLinks := h.Item.(*okta.UserType).Links

		a := &SchemaDetails{}
		json.Unmarshal([]byte(userTypeLinks.(string)), &a)
		// var linkDetails interface{}
		// parseErr := json.Unmarshal([]byte(userTypeLinks.(string)), &linkDetails)
		// if parseErr != nil {
		// 	return nil, parseErr
		// }

		// mappedDetails := linkDetails.(map[string]interface{})
		// schemaMap := mappedDetails["schema"]
		// hrefDetails := schemaMap.(map[string]interface{})
		// schemaHref := hrefDetails["href"].(string)
		schemaId = strings.Split(a.Href, "/")[8]
	}

	if schemaId == "" {
		return nil, nil
	}

	// url := "/api/v1/meta/schemas/user/" + userTypeId
	policies, resp, err := listUserSchemas(ctx, *client, schemaId)
	if err != nil {
		logger.Error("listPolicies", "list policies", err)
		return nil, err
	}
	// requestExecutor := client.GetRequestExecutor()
	// req, err := requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	// if err != nil {
	// 	logger.Error("listOktaUserSchemas", "list schemas request", err)
	// 	return nil, err
	// }

	// var policies []*okta.UserSchema

	// resp, err := requestExecutor.Do(ctx, req, &policies)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "Not found") {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	for _, policy := range policies {
		d.StreamListItem(ctx, policy)
	}

	// paging
	for resp.HasNextPage() {
		var nextFactorSet []*okta.UserSchema
		resp, err = resp.Next(ctx, &nextFactorSet)
		if err != nil {
			logger.Error("listOktaUserSchemas", "list schemas paging", err)
			return nil, err
		}
		for _, factor := range nextFactorSet {
			d.StreamListItem(ctx, factor)
		}
	}

	return nil, err
}

func listUserSchemas(ctx context.Context, client okta.Client, schemaId string) ([]*okta.UserSchema, *okta.Response, error) {
	url := "/api/v1/meta/schemas/user/" + schemaId
	// if qp != nil {
	// 	url = url + qp.String()
	// }

	requestExecutor := client.GetRequestExecutor()
	req, err := requestExecutor.WithAccept("application/json").WithContentType("application/json").NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	var policies []*okta.UserSchema

	resp, err := requestExecutor.Do(ctx, req, &policies)
	if err != nil {
		return nil, resp, err
	}

	return policies, resp, nil
}

func (a *SchemaDetails) UnmarshalJSON(b []byte) error {

	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})

	foomap := m["schema"]
	v := foomap.(map[string]interface{})

	a.Href = v["href"].(string)

	return nil
}