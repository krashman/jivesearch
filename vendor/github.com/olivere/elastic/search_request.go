// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

import (
	"encoding/json"
	"strings"
)

// SearchRequest combines a search request and its
// query details (see SearchSource).
// It is used in combination with MultiSearch.
type SearchRequest struct {
	searchType        string
	indices           []string
	types             []string
	routing           *string
	preference        *string
	requestCache      *bool
	ignoreUnavailable *bool
	allowNoIndices    *bool
	expandWildcards   string
	scroll            string
	source            interface{}
}

// NewSearchRequest creates a new search request.
func NewSearchRequest() *SearchRequest {
	return &SearchRequest{}
}

// SearchRequest must be one of "dfs_query_then_fetch" or
// "query_then_fetch".
func (r *SearchRequest) SearchType(searchType string) *SearchRequest {
	r.searchType = searchType
	return r
}

// SearchTypeDfsQueryThenFetch sets search type to dfs_query_then_fetch.
func (r *SearchRequest) SearchTypeDfsQueryThenFetch() *SearchRequest {
	return r.SearchType("dfs_query_then_fetch")
}

// SearchTypeQueryThenFetch sets search type to query_then_fetch.
func (r *SearchRequest) SearchTypeQueryThenFetch() *SearchRequest {
	return r.SearchType("query_then_fetch")
}

func (r *SearchRequest) Index(indices ...string) *SearchRequest {
	r.indices = append(r.indices, indices...)
	return r
}

func (r *SearchRequest) HasIndices() bool {
	return len(r.indices) > 0
}

func (r *SearchRequest) Type(types ...string) *SearchRequest {
	r.types = append(r.types, types...)
	return r
}

func (r *SearchRequest) Routing(routing string) *SearchRequest {
	r.routing = &routing
	return r
}

func (r *SearchRequest) Routings(routings ...string) *SearchRequest {
	if routings != nil {
		routings := strings.Join(routings, ",")
		r.routing = &routings
	} else {
		r.routing = nil
	}
	return r
}

func (r *SearchRequest) Preference(preference string) *SearchRequest {
	r.preference = &preference
	return r
}

func (r *SearchRequest) RequestCache(requestCache bool) *SearchRequest {
	r.requestCache = &requestCache
	return r
}

// IgnoreUnavailable indicates whether specified concrete indices should be
// ignored when unavailable (missing or closed).
func (s *SearchRequest) IgnoreUnavailable(ignoreUnavailable bool) *SearchRequest {
	s.ignoreUnavailable = &ignoreUnavailable
	return s
}

// AllowNoIndices indicates whether to ignore if a wildcard indices
// expression resolves into no concrete indices. (This includes `_all` string or when no indices have been specified).
func (s *SearchRequest) AllowNoIndices(allowNoIndices bool) *SearchRequest {
	s.allowNoIndices = &allowNoIndices
	return s
}

// ExpandWildcards indicates whether to expand wildcard expression to
// concrete indices that are open, closed or both.
func (s *SearchRequest) ExpandWildcards(expandWildcards string) *SearchRequest {
	s.expandWildcards = expandWildcards
	return s
}

func (r *SearchRequest) Scroll(scroll string) *SearchRequest {
	r.scroll = scroll
	return r
}

func (r *SearchRequest) SearchSource(searchSource *SearchSource) *SearchRequest {
	return r.Source(searchSource)
}

func (r *SearchRequest) Source(source interface{}) *SearchRequest {
	r.source = source
	return r
}

// header is used e.g. by MultiSearch to get information about the search header
// of one SearchRequest.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.0/search-multi-search.html
func (r *SearchRequest) header() interface{} {
	h := make(map[string]interface{})
	if r.searchType != "" {
		h["search_type"] = r.searchType
	}

	switch len(r.indices) {
	case 0:
	case 1:
		h["index"] = r.indices[0]
	default:
		h["indices"] = r.indices
	}

	switch len(r.types) {
	case 0:
	case 1:
		h["type"] = r.types[0]
	default:
		h["types"] = r.types
	}

	if r.routing != nil && *r.routing != "" {
		h["routing"] = *r.routing
	}
	if r.preference != nil && *r.preference != "" {
		h["preference"] = *r.preference
	}
	if r.requestCache != nil {
		h["request_cache"] = *r.requestCache
	}
	if r.ignoreUnavailable != nil {
		h["ignore_unavailable"] = *r.ignoreUnavailable
	}
	if r.allowNoIndices != nil {
		h["allow_no_indices"] = *r.allowNoIndices
	}
	if r.expandWildcards != "" {
		h["expand_wildcards"] = r.expandWildcards
	}
	if r.scroll != "" {
		h["scroll"] = r.scroll
	}

	return h
}

// Body allows to access the search body of the request, as generated by the DSL.
// Notice that Body is read-only. You must not change the request body.
//
// Body is used e.g. by MultiSearch to get information about the search body
// of one SearchRequest.
// See https://www.elastic.co/guide/en/elasticsearch/reference/6.0/search-multi-search.html
func (r *SearchRequest) Body() (string, error) {
	switch t := r.source.(type) {
	default:
		body, err := json.Marshal(r.source)
		if err != nil {
			return "", err
		}
		return string(body), nil
	case *SearchSource:
		src, err := t.Source()
		if err != nil {
			return "", err
		}
		body, err := json.Marshal(src)
		if err != nil {
			return "", err
		}
		return string(body), nil
	case json.RawMessage:
		return string(t), nil
	case *json.RawMessage:
		return string(*t), nil
	case string:
		return t, nil
	case *string:
		if t != nil {
			return *t, nil
		}
		return "{}", nil
	}
}
