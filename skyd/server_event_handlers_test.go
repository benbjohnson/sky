package skyd

import (
  "testing"
)

// Ensure that we can put an event on the server.
func TestServerUpdateEvents(t *testing.T) {
  runTestServer(func(s *Server) {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    setupTestProperty("foo", "baz", "action", "integer")
    
    // Send two new events.
    resp, _ := sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"myValue"}, "action":{"baz":12}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T03:00:00Z", "application/json", `{"data":{"bar":"myValue2"}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    
    // Replace the first one.
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"myValue3"}, "action":{"baz":1000}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")

    // Merge the second one.
    resp, _ = sendTestHttpRequest("PATCH", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T03:00:00Z", "application/json", `{"data":{"bar":"myValue2"}, "action":{"baz":20}}`)
    assertResponse(t, resp, 200, "", "PATCH /tables/:name/objects/:objectId/events failed.")

    // Check our work.
    resp, _ = sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/objects/xyz/events", "application/json", "")
    assertResponse(t, resp, 200, `[{"action":{"baz":1000},"data":{"bar":"myValue3"},"timestamp":"2012-01-01T02:00:00Z"},{"action":{"baz":20},"data":{"bar":"myValue2"},"timestamp":"2012-01-01T03:00:00Z"}]`+"\n", "GET /tables/:name/objects/:objectId/events failed.")

    // Grab a single event.
    resp, _ = sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T03:00:00Z", "application/json", "")
    assertResponse(t, resp, 200, `{"action":{"baz":20},"data":{"bar":"myValue2"},"timestamp":"2012-01-01T03:00:00Z"}`+"\n", "GET /tables/:name/objects/:objectId/events/:timestamp failed.")
  })
}

// Ensure that we can delete all events for an object.
func TestServerDeleteEvent(t *testing.T) {
  runTestServer(func(s *Server) {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    
    // Send two new events.
    resp, _ := sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"myValue"}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T03:00:00Z", "application/json", `{"data":{"bar":"myValue2"}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    
    // Delete one of the events.
    resp, _ = sendTestHttpRequest("DELETE", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T02:00:00Z", "application/json", "")
    assertResponse(t, resp, 200, "", "DELETE /tables/:name/objects/:objectId/events failed.")

    // Check our work.
    resp, _ = sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/objects/xyz/events", "application/json", "")
    assertResponse(t, resp, 200, `[{"data":{"bar":"myValue2"},"timestamp":"2012-01-01T03:00:00Z"}]`+"\n", "GET /tables/:name/objects/:objectId/events failed.")
  })
}

// Ensure that we can delete all events for an object.
func TestServerDeleteEvents(t *testing.T) {
  runTestServer(func(s *Server) {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    
    // Send two new events.
    resp, _ := sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"myValue"}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events/2012-01-01T03:00:00Z", "application/json", `{"data":{"bar":"myValue2"}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
    
    // Delete the events.
    resp, _ = sendTestHttpRequest("DELETE", "http://localhost:8585/tables/foo/objects/xyz/events", "application/json", "")
    assertResponse(t, resp, 200, "", "DELETE /tables/:name/objects/:objectId/events failed.")

    // Check our work.
    resp, _ = sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/objects/xyz/events", "application/json", "")
    assertResponse(t, resp, 200, "[]\n", "GET /tables/:name/objects/:objectId/events failed.")
  })
}