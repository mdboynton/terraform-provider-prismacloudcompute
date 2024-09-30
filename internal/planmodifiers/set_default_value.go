package planmodifiers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type setDefaultValue struct {
	value    []attr.Value
	elemType attr.Type
}

func SetDefaultValue(elemType attr.Type, value []attr.Value) planmodifier.Set {
	return &setDefaultValue{value: value, elemType: elemType}
}

var _ planmodifier.Set = (*setDefaultValue)(nil)

func (m *setDefaultValue) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m *setDefaultValue) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Sets the default value %v if the attribute is not set", m.value)
}

func (m *setDefaultValue) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if !req.ConfigValue.IsNull() {
		return
	}

	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue, resp.Diagnostics = types.SetValue(m.elemType, m.value)
}
