// Code generated by go-swagger; DO NOT EDIT.

package users

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"
	"strings"

	"github.com/go-openapi/swag"
)

// GetUsersIDEntriesURL generates an URL for the get users ID entries operation
type GetUsersIDEntriesURL struct {
	ID int64

	Limit *int64
	Skip  *int64
	Sort  *string
	Tag   *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetUsersIDEntriesURL) WithBasePath(bp string) *GetUsersIDEntriesURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetUsersIDEntriesURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetUsersIDEntriesURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/users/{id}/entries"

	id := swag.FormatInt64(o.ID)
	if id != "" {
		_path = strings.Replace(_path, "{id}", id, -1)
	} else {
		return nil, errors.New("ID is required on GetUsersIDEntriesURL")
	}
	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api/v1"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var limit string
	if o.Limit != nil {
		limit = swag.FormatInt64(*o.Limit)
	}
	if limit != "" {
		qs.Set("limit", limit)
	}

	var skip string
	if o.Skip != nil {
		skip = swag.FormatInt64(*o.Skip)
	}
	if skip != "" {
		qs.Set("skip", skip)
	}

	var sort string
	if o.Sort != nil {
		sort = *o.Sort
	}
	if sort != "" {
		qs.Set("sort", sort)
	}

	var tag string
	if o.Tag != nil {
		tag = *o.Tag
	}
	if tag != "" {
		qs.Set("tag", tag)
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetUsersIDEntriesURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetUsersIDEntriesURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetUsersIDEntriesURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetUsersIDEntriesURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetUsersIDEntriesURL")
	}

	base, err := o.Build()
	if err != nil {
		return nil, err
	}

	base.Scheme = scheme
	base.Host = host
	return base, nil
}

// StringFull returns the string representation of a complete url
func (o *GetUsersIDEntriesURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
