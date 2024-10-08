package prismatypes 

import (
    "context"
    "fmt"
    "reflect"
    
    "github.com/hashicorp/terraform-plugin-go/tftypes"

    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    //"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func CollectionObjectType() types.ObjectType {
    return types.ObjectType{
        AttrTypes: CollectionObjectAttrTypeMap(),
    }
}

func CollectionObjectAttrTypeMap() map[string]attr.Type {
    return map[string]attr.Type{
        "account_ids":  types.SetType{ElemType: types.StringType},
        "app_ids":  types.SetType{ElemType: types.StringType},
        "clusters":  types.SetType{ElemType: types.StringType},
        "color":        types.StringType,
        "containers":  types.SetType{ElemType: types.StringType},
        "description":  types.StringType,
        "functions":  types.SetType{ElemType: types.StringType},
        "hosts":  types.SetType{ElemType: types.StringType},
        "images":  types.SetType{ElemType: types.StringType},
        "labels":  types.SetType{ElemType: types.StringType},
        "modified": types.StringType,
        "name":         types.StringType,
        "namespaces":  types.SetType{ElemType: types.StringType},
        "owner":        types.StringType,
        "prisma":       types.BoolType,
        "system":       types.BoolType,
    }
}

//func CollectionSetType() types.SetType {
//    return CollectionType 
//}

// Ensure the implementation satisifies the expected interfaces
//var _ basetypes.SetTypable = CollectionType{}
var _ basetypes.SetTypable = (*CollectionType)(nil)
//var _ basetypes.ObjectTypable = (*CollectionType)(nil)
//var _ xattr.TypeWithValidate = CollectionType{}

type CollectionType struct {
    basetypes.SetType
    //basetypes.ObjectType
}

//var (
//    CollectionType = CollectionType{basetypes.SetType{ElemType: CollectionObjectType()}}
//    //CollectionType = CollectionType{basetypes.ObjectType{AttrTypes: CollectionObjectAttrTypeMap()}}
//)

//func (t CollectionType) ToTerraformValue(ctx context.Context, in tfsdk.Value) 

func (t CollectionType) Equal(o attr.Type) bool {
    other, ok := o.(CollectionType)

    if !ok {
        return false
    }

    //return t.ObjectType.Equal(other.ObjectType)
    return t.SetType.Equal(other.SetType)
}

func (t CollectionType) String() string {
    return "prismacloudcompute.CollectionType"
}

// func(t CollectionType) ValueFromString(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics)

//func (t CollectionType) ValueFromObject(ctx context.Context, in basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
func (t CollectionType) ValueFromSet(ctx context.Context, in basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("in ValueFromSet")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    
    var diags diag.Diagnostics

    if in.IsNull() {
        //return collectionValue{basetypes.NewSetNull(CollectionObjectType())}
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        fmt.Println("returning null")
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        return NewCollectionNull(), diags
    }

    if in.IsUnknown() {
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        fmt.Println("returning unknown")
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        //return collectionValue{basetypes.NewSetUnknown(CollectionObjectType())}
        return NewCollectionUnknown(), diags
    }

    //v, d := basetypes.NewSetValue(CollectionObjectType(), in.Elements())
    //_, d := basetypes.NewObjectValue(CollectionObjectAttrTypeMap(), in.Attributes())
    v, d := basetypes.NewSetValue(CollectionObjectType(), in.Elements())
    diags.Append(d...)
    if diags.HasError() {
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        fmt.Println("returning unknown from NewObjectValue error")
        fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
        return NewCollectionUnknown(), diags
        //return CollectionValue{basetypes.NewSetUnknown(CollectionObjectType())}, diags
    }

    //return collectionValue{basetypes.SetValue(CollectionObjectType, 
    //return CollectionValue{basetypes.NewSetUnknown(CollectionObjectType())}, diags
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("returning CollectionValue")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    //return CollectionValue{basetypes.NewObjectUnknown(CollectionObjectAttrTypeMap())}, diags
    //return NewCollectionUnknown(), diags
    //return CollectionValue{v}, diags

    value := CollectionValue{
        SetValue: v,
    }

    return value, diags
}

func(t CollectionType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("in ValueFromTerraform")
    //fmt.Println(in)
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    
    //attrValue, err := t.ObjectType.ValueFromTerraform(ctx, in)
    attrValue, err := t.SetType.ValueFromTerraform(ctx, in)

    if err != nil {
        return nil, err
    }

    //objectValue, ok := attrValue.(basetypes.ObjectValue)
    setValue, ok := attrValue.(basetypes.SetValue)

    if !ok {
        return nil, fmt.Errorf("unexpected value type of %T", attrValue) 
    }

    //objectValuable, diags := t.ValueFromObject(ctx, objectValue)
    setValuable, diags := t.ValueFromSet(ctx, setValue)

    if diags.HasError() {
        //return nil, fmt.Errorf("unexpected error converting ObjectValue to ObjectValuable: %v", diags)
        return nil, fmt.Errorf("unexpected error converting SetValue to SetValuable: %v", diags)
    }

    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    //fmt.Println("returning objectValuable")
    fmt.Println("returning setValuable from ValueFromTerraform")
    fmt.Println(reflect.TypeOf(setValuable))
    fmt.Println(setValuable)
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")

    //return objectValuable, nil
    return setValuable, nil
    //return CollectionType
}

func (t CollectionType) ValueType(ctx context.Context) attr.Value {
    return CollectionValue{}
}

func (t CollectionType) Validate(ctx context.Context, in tftypes.Value, path path.Path) diag.Diagnostics {
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("in Validate")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    
    var diags diag.Diagnostics

    if in.Type() == nil {
        return diags
    }

    //if !in.Type().Is(tftypes.Set) {
    //if !in.Type().Is(CollectionValue.Type(ctx).TerraformType(ctx)) {
    if !in.Type().Is(t.SetType.TerraformType(ctx)) {
		err := fmt.Errorf("expected Set value, received %T with value: %v", in, in)
		diags.AddAttributeError(
			path,
			"Prisma Collection Type Validation Error",
			"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
				"Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return diags
	}

    //var valueSet jj //t.SetType.TerraformType(ctx)

    //if err := in.ElementsAs(&valueSet); err != nil {
	//	diags.AddAttributeError(
	//		path,
	//		"Prisma Collection Type Validation Error",
	//		"An unexpected error was encountered trying to validate an attribute value. This is always an error in the provider. "+
	//			"Please report the following to the provider developer:\n\n"+err.Error(),
	//	)

	//	return diags
	//}
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")
    fmt.Println("exiting Validate")
    fmt.Println("^^^^^^^^^^^^^^^^^^^^^")

    return diags
}
