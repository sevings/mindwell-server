// Code generated by go-swagger; DO NOT EDIT.

package entries

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"errors"
	"net/url"
	golangswaggerpaths "path"

	"github.com/go-openapi/swag"
)

// GetEntriesBestURL generates an URL for the get entries best operation
type GetEntriesBestURL struct {
	After       *string
	Before      *string
	Limit       *int64
	LongerThan  *int64
	MinRating   *int64
	ShorterThan *int64
	Tag         *string

	_basePath string
	// avoid unkeyed usage
	_ struct{}
}

// WithBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetEntriesBestURL) WithBasePath(bp string) *GetEntriesBestURL {
	o.SetBasePath(bp)
	return o
}

// SetBasePath sets the base path for this url builder, only required when it's different from the
// base path specified in the swagger spec.
// When the value of the base path is an empty string
func (o *GetEntriesBestURL) SetBasePath(bp string) {
	o._basePath = bp
}

// Build a url path and query string
func (o *GetEntriesBestURL) Build() (*url.URL, error) {
	var result url.URL

	var _path = "/entries/best"

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

	var longerThan string
	if o.LongerThan != nil {
		longerThan = swag.FormatInt64(*o.LongerThan)
	}
	if longerThan != "" {
		qs.Set("longer_than", longerThan)
	}

	var minRating string
	if o.MinRating != nil {
		minRating = swag.FormatInt64(*o.MinRating)
	}
	if minRating != "" {
		qs.Set("min_rating", minRating)
	}

	var shorterThan string
	if o.ShorterThan != nil {
		shorterThan = swag.FormatInt64(*o.ShorterThan)
	}
	if shorterThan != "" {
		qs.Set("shorter_than", shorterThan)
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
func (o *GetEntriesBestURL) Must(u *url.URL, err error) *url.URL {
	if err != nil {
		panic(err)
	}
	if u == nil {
		panic("url can't be nil")
	}
	return u
}

// String returns the string representation of the path with query string
func (o *GetEntriesBestURL) String() string {
	return o.Must(o.Build()).String()
}

// BuildFull builds a full url with scheme, host, path and query string
func (o *GetEntriesBestURL) BuildFull(scheme, host string) (*url.URL, error) {
	if scheme == "" {
		return nil, errors.New("scheme is required for a full url on GetEntriesBestURL")
	}
	if host == "" {
		return nil, errors.New("host is required for a full url on GetEntriesBestURL")
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
func (o *GetEntriesBestURL) StringFull(scheme, host string) string {
	return o.Must(o.BuildFull(scheme, host)).String()
}
