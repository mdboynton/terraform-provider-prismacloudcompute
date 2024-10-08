package prismatypes 

import (
    "context"
    "fmt" 
    
    "github.com/hashicorp/terraform-plugin-go/tftypes"
    
    "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

//func CollectionObjectValue() map[string]attr.Value {
//
//}

func DefaultCollectionObjectValue() map[string]attr.Value {
    return map[string]attr.Value{
        "account_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "app_ids": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "clusters": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "color": types.StringValue("#3FA2F7"),
        "containers": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "description": types.StringValue("System - all resources collection"),
        "functions": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "hosts": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "images": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "labels": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "modified": basetypes.NewStringUnknown(),
        //"modified": types.StringValue(""),
        "name": types.StringValue("All"),
        "namespaces": types.SetValueMust(types.StringType, []attr.Value{ types.StringValue("*") }),
        "owner": types.StringValue("admin"),
        "prisma": types.BoolValue(false),
        "system": types.BoolValue(true),
    }
}


// Ensure the implementation satisfies the expected interfaces
//var _ basetypes.ObjectValuable = CollectionValue{}
var _ basetypes.SetValuable = CollectionValue{}
//var _ basetypes.SetValuableWithSemanticEquals = CollectionValue{}

//type collectionValue struct {
type CollectionValue struct {
    //basetypes.ObjectValue
    basetypes.SetValue
}

func (v CollectionValue) Type(ctx context.Context) attr.Type {
    //return collectionType{}
    //return CollectionType
    return CollectionType{basetypes.SetType{ElemType: CollectionObjectType()}}
}

func (v CollectionValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("in ToTerraformValue")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")

    //return CollectionValue{
    //    SetValue: v.ToTerraformValue(ctx),
    //}
    return v.ToTerraformValue(ctx)

    //return v.
    //setType := tftypes.Set{ElementType: CollectionObjectType().TerraformType(ctx)}

    //switch s
    //return tftypes.NewValue(setType, nil), nil
}

//var (
//    CollectionValue = collectionValue{basetypes.SetType{ElemType: CollectionObjectType(), }}
//)

func NewCollectionNull() CollectionValue {
    //return CollectionValue{basetypes.NewObjectNull(CollectionObjectAttrTypeMap())}
    return CollectionValue{basetypes.NewSetNull(CollectionObjectType())}
}

func NewCollectionUnknown() CollectionValue {
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("in NewCollectionUnknown")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    //return CollectionValue{basetypes.NewObjectUnknown(CollectionObjectAttrTypeMap())}
    return CollectionValue{basetypes.NewSetUnknown(CollectionObjectType())}
}

//func NewCollectionValue(


