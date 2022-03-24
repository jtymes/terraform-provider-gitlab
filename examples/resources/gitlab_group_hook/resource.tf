resource "gitlab_group_hook" "example" {
  group                 = "example/hooked"
  url                   = "https://example.com/hook/example"
  merge_requests_events = true
}
