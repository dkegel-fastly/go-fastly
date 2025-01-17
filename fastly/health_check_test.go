package fastly

import (
	"testing"
)

func TestClient_HealthChecks(t *testing.T) {
	t.Parallel()

	var err error
	var tv *Version
	record(t, "health_checks/version", func(c *Client) {
		tv = testVersion(t, c)
	})

	// Create
	var hc *HealthCheck
	record(t, "health_checks/create", func(c *Client) {
		hc, err = c.CreateHealthCheck(&CreateHealthCheckInput{
			ServiceID:        testServiceID,
			ServiceVersion:   tv.Number,
			Name:             "test-healthcheck",
			Method:           "HEAD",
			Host:             "example.com",
			Path:             "/foo",
			HTTPVersion:      "1.1",
			Timeout:          Uint(1500),
			CheckInterval:    Uint(2500),
			ExpectedResponse: Uint(200),
			Window:           Uint(5000),
			Threshold:        Uint(10),
			Initial:          Uint(10),
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	// Ensure deleted
	defer func() {
		record(t, "health_checks/cleanup", func(c *Client) {
			c.DeleteHealthCheck(&DeleteHealthCheckInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "test-healthcheck",
			})

			c.DeleteHealthCheck(&DeleteHealthCheckInput{
				ServiceID:      testServiceID,
				ServiceVersion: tv.Number,
				Name:           "new-test-healthcheck",
			})
		})
	}()

	if hc.Name != "test-healthcheck" {
		t.Errorf("bad name: %q", hc.Name)
	}
	if hc.Method != "HEAD" {
		t.Errorf("bad address: %q", hc.Method)
	}
	if hc.Host != "example.com" {
		t.Errorf("bad host: %q", hc.Host)
	}
	if hc.Path != "/foo" {
		t.Errorf("bad path: %q", hc.Path)
	}
	if hc.HTTPVersion != "1.1" {
		t.Errorf("bad http_version: %q", hc.HTTPVersion)
	}
	if hc.Timeout != 1500 {
		t.Errorf("bad timeout: %q", hc.Timeout)
	}
	if hc.CheckInterval != 2500 {
		t.Errorf("bad check_interval: %q", hc.CheckInterval)
	}
	if hc.ExpectedResponse != 200 {
		t.Errorf("bad timeout: %q", hc.ExpectedResponse)
	}
	if hc.Window != 5000 {
		t.Errorf("bad window: %q", hc.Window)
	}
	if hc.Threshold != 10 {
		t.Errorf("bad threshold: %q", hc.Threshold)
	}
	if hc.Initial != 10 {
		t.Errorf("bad initial: %q", hc.Initial)
	}

	// List
	var hcs []*HealthCheck
	record(t, "health_checks/list", func(c *Client) {
		hcs, err = c.ListHealthChecks(&ListHealthChecksInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hcs) < 1 {
		t.Errorf("bad health checks: %v", hcs)
	}

	// Get
	var nhc *HealthCheck
	record(t, "health_checks/get", func(c *Client) {
		nhc, err = c.GetHealthCheck(&GetHealthCheckInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-healthcheck",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if hc.Name != nhc.Name {
		t.Errorf("bad name: %q (%q)", hc.Name, nhc.Name)
	}
	if hc.Method != nhc.Method {
		t.Errorf("bad address: %q", hc.Method)
	}
	if hc.Host != nhc.Host {
		t.Errorf("bad host: %q", hc.Host)
	}
	if hc.Path != nhc.Path {
		t.Errorf("bad path: %q", hc.Path)
	}
	if hc.HTTPVersion != nhc.HTTPVersion {
		t.Errorf("bad http_version: %q", hc.HTTPVersion)
	}
	if hc.Timeout != nhc.Timeout {
		t.Errorf("bad timeout: %q", hc.Timeout)
	}
	if hc.CheckInterval != nhc.CheckInterval {
		t.Errorf("bad check_interval: %q", hc.CheckInterval)
	}
	if hc.ExpectedResponse != nhc.ExpectedResponse {
		t.Errorf("bad timeout: %q", hc.ExpectedResponse)
	}
	if hc.Window != nhc.Window {
		t.Errorf("bad window: %q", hc.Window)
	}
	if hc.Threshold != nhc.Threshold {
		t.Errorf("bad threshold: %q", hc.Threshold)
	}
	if hc.Initial != nhc.Initial {
		t.Errorf("bad initial: %q", hc.Initial)
	}

	// Update
	var uhc *HealthCheck
	record(t, "health_checks/update", func(c *Client) {
		uhc, err = c.UpdateHealthCheck(&UpdateHealthCheckInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "test-healthcheck",
			NewName:        String("new-test-healthcheck"),
		})
	})
	if err != nil {
		t.Fatal(err)
	}
	if uhc.Name != "new-test-healthcheck" {
		t.Errorf("bad name: %q", uhc.Name)
	}

	// Delete
	record(t, "health_checks/delete", func(c *Client) {
		err = c.DeleteHealthCheck(&DeleteHealthCheckInput{
			ServiceID:      testServiceID,
			ServiceVersion: tv.Number,
			Name:           "new-test-healthcheck",
		})
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_ListHealthChecks_validation(t *testing.T) {
	var err error
	_, err = testClient.ListHealthChecks(&ListHealthChecksInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.ListHealthChecks(&ListHealthChecksInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_CreateHealthCheck_validation(t *testing.T) {
	var err error
	_, err = testClient.CreateHealthCheck(&CreateHealthCheckInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.CreateHealthCheck(&CreateHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_GetHealthCheck_validation(t *testing.T) {
	var err error
	_, err = testClient.GetHealthCheck(&GetHealthCheckInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHealthCheck(&GetHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.GetHealthCheck(&GetHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_UpdateHealthCheck_validation(t *testing.T) {
	var err error
	_, err = testClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	_, err = testClient.UpdateHealthCheck(&UpdateHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}

func TestClient_DeleteHealthCheck_validation(t *testing.T) {
	var err error
	err = testClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		ServiceID: "",
	})
	if err != ErrMissingServiceID {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 0,
	})
	if err != ErrMissingServiceVersion {
		t.Errorf("bad error: %s", err)
	}

	err = testClient.DeleteHealthCheck(&DeleteHealthCheckInput{
		ServiceID:      "foo",
		ServiceVersion: 1,
		Name:           "",
	})
	if err != ErrMissingName {
		t.Errorf("bad error: %s", err)
	}
}
