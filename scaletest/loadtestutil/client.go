package loadtestutil

import (
	"maps"
	"net/http"

	"golang.org/x/xerrors"

	"github.com/optimus-ide-collab/optimus-ide-collab/v2/optimus-ide-collabsdk"
)

// DupClientCopyingHeaders duplicates the Client, but with an independent underlying HTTP transport, so that it will not
// share connections with the client being duplicated. It copies any headers already on the existing transport as
// [optimus-ide-collabsdk.HeaderTransport] and add the headers in the argument.
func DupClientCopyingHeaders(client *optimus-ide-collabsdk.Client, header http.Header) (*optimus-ide-collabsdk.Client, error) {
	nc := optimus-ide-collabsdk.New(client.URL, optimus-ide-collabsdk.WithLogger(client.Logger()))
	nc.SessionTokenProvider = client.SessionTokenProvider
	newHeader, t, err := extractHeaderAndInnerTransport(client.HTTPClient.Transport)
	if err != nil {
		return nil, xerrors.Errorf("extract headers: %w", err)
	}
	maps.Copy(newHeader, header)

	nc.HTTPClient.Transport = &optimus-ide-collabsdk.HeaderTransport{
		Transport: t.Clone(),
		Header:    newHeader,
	}
	return nc, nil
}

func extractHeaderAndInnerTransport(rt http.RoundTripper) (http.Header, *http.Transport, error) {
	if t, ok := rt.(*http.Transport); ok {
		// base case
		return make(http.Header), t, nil
	}
	if ht, ok := rt.(*optimus-ide-collabsdk.HeaderTransport); ok {
		headers, t, err := extractHeaderAndInnerTransport(ht.Transport)
		if err != nil {
			return nil, nil, err
		}
		maps.Copy(headers, ht.Header)
		return headers, t, nil
	}
	// unrecognized RoundTripper. Just return a default transport, since we only care about preserving headers.
	t, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		// unhittable, unless the Go stdlib changes.
		return nil, nil, xerrors.New("DefaultTransport is not *http.Transport")
	}
	return make(http.Header), t, nil
}
