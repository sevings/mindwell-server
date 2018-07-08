// Code generated by go-swagger; DO NOT EDIT.

package me

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GetMeFavoritesURL generates an URL for the get me favorites operation
type GetMeFavoritesURL struct {
	After  *string
	Before *string
	Limit  *int64

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetMeFavoritesURL) WithBasePath(bp string) *GetMeFavoritesURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetMeFavoritesURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetMeFavoritesURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/me/favorites"

	_basePath := o._basePath
	if _basePath == "" {
		_basePath = "/api/v1"
	}
	result.Path = golangswaggerpaths.Join(_basePath, _path)

	qs := make(url.Values)

	var after string
	if o.After != nil {
		after = *o.After
	}
	if after != "" {
		qs.Set("after", after)
	}

	var before string
	if o.Before != nil {
		before = *o.Before
	}
	if before != "" {
		qs.Set("before", before)
	}

	var limit string
	if o.Limit != nil {
		limit = swag.FormatInt64(*o.Limit)
	}
	if limit != "" {
		qs.Set("limit", limit)
	}

	result.RawQuery = qs.Encode()

	return &result, nil
}

// Must is a helper function to panic when the url builder returns an error
func (o *GetMeFavoritesURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetMeFavoritesURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetMeFavoritesURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetMeFavoritesURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetMeFavoritesURL")
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
func (o *GetMeFavoritesURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}