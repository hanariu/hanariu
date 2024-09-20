package hanariu

import (
	"bytes"
	"context"
	"io"
	"reflect"

	"github.com/a-h/templ"
)

type Boson map[string]string
type Attrs map[string]string

func CreateClasses(prop string, val string, b interface{}, boson Boson) string {
	bosonKey := GetTagDefault(prop, val, b)
	if len(boson[bosonKey]) != 0 {
		return boson[val]
	} else {
		return ""
	}
}

func GetAttrs(attrs Attrs) string {
	attr := ""
	for k, v := range attrs {
		if len(v) != 0 {
			if v == "true" {
				attr += " " + k + "=true"
			} else {
				attr += " " + k + "=\"" + v + "\""
			}
		} else {
			attr += " " + k
		}

	}
	return attr
}

func GetTagDefault(prop string, val string, b interface{}) string {
	if len(val) != 0 {
		return val
	} else {
		st := reflect.TypeOf(b)
		field, found := st.FieldByName(prop)
		if found {
			defaultTag := field.Tag.Get("default")
			if len(defaultTag) > 0 {
				return defaultTag
			}
		}
		return ""
	}
}

func CreateComponent(tag string, classes string, attrs Attrs, short bool) templ.ComponentFunc {
	footer := ""
	header := "<" + tag
	if len(classes) > 0 {
		header += " class=\"" + classes + "\""
	}
	if len(attrs) > 0 {
		header += GetAttrs(attrs)
	}
	if short {
		header += ">\n"
	} else {
		header += ">\n"
		footer += "</" + tag + ">"
	}

	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (Err error) {
		Buffer, IsBuffer := w.(*bytes.Buffer)
		if !IsBuffer {
			Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(Buffer)
		}
		_, Err = Buffer.WriteString(header)
		if Err != nil {
			return Err
		}
		if !short {
			ctx = templ.InitializeContext(ctx)
			children := templ.GetChildren(ctx)
			if children == nil {
				children = templ.NopComponent
			}
			ctx = templ.ClearChildren(ctx)

			Err = children.Render(ctx, Buffer)
			if Err != nil {
				return Err
			}
			_, Err = Buffer.WriteString(footer)
			if Err != nil {
				return Err
			}
		}

		if !IsBuffer {
			_, Err = Buffer.WriteTo(w)
		}
		return Err
	})
}
