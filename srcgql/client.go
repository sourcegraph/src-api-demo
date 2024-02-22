package srcgql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type tokenAuthTransport struct {
	token         string
	csrfCookie    *http.Cookie
	sessionCookie *http.Cookie

	wrapped http.RoundTripper
}

func (t *tokenAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", t.token))
	} else {
		// NOTE: This header is required to authenticate our session with a session cookie, see:
		// https://docs.sourcegraph.com/dev/security/csrf_security_model#authentication-in-api-endpoints
		req.Header.Set("X-Requested-With", "Sourcegraph")
		req.AddCookie(t.sessionCookie)

		// Older versions of Sourcegraph require a CSRF cookie.
		if t.csrfCookie != nil {
			req.AddCookie(t.csrfCookie)
		}
	}
	return t.wrapped.RoundTrip(req)
}

// NewGraphQLClient creates a new GraphQL client for the given sourcegraph endpoint
// and use bearer token for authentication.
func NewGraphQLClient(endpoint string, token string) graphql.Client {
	return newTracedClient(endpoint, &http.Client{
		Transport: &tokenAuthTransport{
			token:   token,
			wrapped: http.DefaultTransport,
		},
	})
}

// NewGraphQLClient creates a new GraphQL client for the given sourcegraph endpoint
// and use cookies for authentication.
func NewGraphQLCookieClient(endpoint string, csrfCookie *http.Cookie, sessionCookie *http.Cookie) graphql.Client {
	return newTracedClient(endpoint, &http.Client{
		Transport: &tokenAuthTransport{
			csrfCookie:    csrfCookie,
			sessionCookie: sessionCookie,
			wrapped:       http.DefaultTransport,
		},
	})
}

func newTracedClient(endpoint string, httpClient graphql.Doer) graphql.Client {
	client := graphql.NewClient(endpoint, httpClient)
	return &tracedClient{endpoint: endpoint, client: client}
}

type tracedClient struct {
	endpoint string
	client   graphql.Client
}

var tracer = otel.Tracer("internal/srcgql")

func (tc *tracedClient) MakeRequest(
	ctx context.Context,
	req *graphql.Request,
	resp *graphql.Response,
) error {
	// Start a span
	ctx, span := tracer.Start(ctx, fmt.Sprintf("GraphQL: %s", req.OpName),
		trace.WithAttributes(
			attribute.String("endpoint", tc.endpoint),
			attribute.String("query", req.Query),
		))

	// Do the request
	err := tc.client.MakeRequest(ctx, req, resp)

	// Assess the result
	if err != nil {
		span.RecordError(err)
	}
	if len(resp.Errors) > 0 {
		span.RecordError(resp.Errors)
	}
	span.End()

	return err
}
