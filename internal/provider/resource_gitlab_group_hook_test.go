package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupHook_basic(t *testing.T) {
	var hook gitlab.GroupHook
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabGroupHookDestroy,
		Steps: []resource.TestStep{
			// Create a group and hook with default options
			{
				Config: testAccGitlabGroupHookConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupHookExists("gitlab_group_hook.foo", &hook),
					testAccCheckGitlabGroupHookAttributes(&hook, &testAccGitlabGroupHookExpectedAttributes{
						URL:                   fmt.Sprintf("https://example.com/group-hook-%d", rInt),
						PushEvents:            true,
						EnableSSLVerification: true,
					}),
				),
			},
			// Update the group hook to toggle all the values to their inverse
			{
				Config: testAccGitlabGroupHookUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupHookExists("gitlab_group_hook.foo", &hook),
					testAccCheckGitlabGroupHookAttributes(&hook, &testAccGitlabGroupHookExpectedAttributes{
						URL:                      fmt.Sprintf("https://example.com/group-hook-%d", rInt),
						PushEvents:               true,
						PushEventsBranchFilter:   "devel",
						IssuesEvents:             false,
						ConfidentialIssuesEvents: false,
						MergeRequestsEvents:      true,
						TagPushEvents:            true,
						NoteEvents:               true,
						ConfidentialNoteEvents:   true,
						JobEvents:                true,
						PipelineEvents:           true,
						WikiPageEvents:           true,
						DeploymentEvents:         true,
						ReleasesEvents:           true,
						SubGroupEvents:           true,
						EnableSSLVerification:    false,
					}),
				),
			},
			// Update the group hook to toggle the options back
			{
				Config: testAccGitlabGroupHookConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupHookExists("gitlab_group_hook.foo", &hook),
					testAccCheckGitlabGroupHookAttributes(&hook, &testAccGitlabGroupHookExpectedAttributes{
						URL:                   fmt.Sprintf("https://example.com/group-hook-%d", rInt),
						PushEvents:            true,
						EnableSSLVerification: true,
					}),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_group_hook.foo",
				ImportStateIdFunc:       getGroupHookImportID("gitlab_group_hook.foo"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccCheckGitlabGroupHookExists(n string, hook *gitlab.GroupHook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		hookID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupID := rs.Primary.Attributes["group"]
		if groupID == "" {
			return fmt.Errorf("No group ID is set")
		}

		gotHook, _, err := testGitlabClient.Groups.GetGroupHook(groupID, hookID)
		if err != nil {
			return err
		}
		*hook = *gotHook
		return nil
	}
}

type testAccGitlabGroupHookExpectedAttributes struct {
	URL                      string
	PushEvents               bool
	PushEventsBranchFilter   string
	IssuesEvents             bool
	ConfidentialIssuesEvents bool
	MergeRequestsEvents      bool
	TagPushEvents            bool
	NoteEvents               bool
	ConfidentialNoteEvents   bool
	JobEvents                bool
	PipelineEvents           bool
	WikiPageEvents           bool
	DeploymentEvents         bool
	ReleasesEvents           bool
	SubGroupEvents           bool
	EnableSSLVerification    bool
}

func testAccCheckGitlabGroupHookAttributes(hook *gitlab.GroupHook, want *testAccGitlabGroupHookExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if hook.URL != want.URL {
			return fmt.Errorf("got url %q; want %q", hook.URL, want.URL)
		}

		if hook.EnableSSLVerification != want.EnableSSLVerification {
			return fmt.Errorf("got enable_ssl_verification %t; want %t", hook.EnableSSLVerification, want.EnableSSLVerification)
		}

		if hook.PushEvents != want.PushEvents {
			return fmt.Errorf("got push_events %t; want %t", hook.PushEvents, want.PushEvents)
		}

		if hook.PushEventsBranchFilter != want.PushEventsBranchFilter {
			return fmt.Errorf("got push_events_branch_filter %q; want %q", hook.PushEventsBranchFilter, want.PushEventsBranchFilter)
		}

		if hook.IssuesEvents != want.IssuesEvents {
			return fmt.Errorf("got issues_events %t; want %t", hook.IssuesEvents, want.IssuesEvents)
		}

		if hook.ConfidentialIssuesEvents != want.ConfidentialIssuesEvents {
			return fmt.Errorf("got confidential_issues_events %t; want %t", hook.ConfidentialIssuesEvents, want.ConfidentialIssuesEvents)
		}

		if hook.MergeRequestsEvents != want.MergeRequestsEvents {
			return fmt.Errorf("got merge_requests_events %t; want %t", hook.MergeRequestsEvents, want.MergeRequestsEvents)
		}

		if hook.TagPushEvents != want.TagPushEvents {
			return fmt.Errorf("got tag_push_events %t; want %t", hook.TagPushEvents, want.TagPushEvents)
		}

		if hook.NoteEvents != want.NoteEvents {
			return fmt.Errorf("got note_events %t; want %t", hook.NoteEvents, want.NoteEvents)
		}

		if hook.ConfidentialNoteEvents != want.ConfidentialNoteEvents {
			return fmt.Errorf("got confidential_note_events %t; want %t", hook.ConfidentialNoteEvents, want.ConfidentialNoteEvents)
		}

		if hook.JobEvents != want.JobEvents {
			return fmt.Errorf("got job_events %t; want %t", hook.JobEvents, want.JobEvents)
		}

		if hook.PipelineEvents != want.PipelineEvents {
			return fmt.Errorf("got pipeline_events %t; want %t", hook.PipelineEvents, want.PipelineEvents)
		}

		if hook.WikiPageEvents != want.WikiPageEvents {
			return fmt.Errorf("got wiki_page_events %t; want %t", hook.WikiPageEvents, want.WikiPageEvents)
		}

		if hook.DeploymentEvents != want.DeploymentEvents {
			return fmt.Errorf("got deployment_events %t; want %t", hook.DeploymentEvents, want.DeploymentEvents)
		}

		if hook.ReleasesEvents != want.ReleasesEvents {
			return fmt.Errorf("got releases_events %t; want %t", hook.ReleasesEvents, want.ReleasesEvents)
		}

		if hook.SubGroupEvents != want.SubGroupEvents {
			return fmt.Errorf("got subgroup_events %t; want %t", hook.SubGroupEvents, want.SubGroupEvents)
		}

		return nil
	}
}

func testAccCheckGitlabGroupHookDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group" {
			continue
		}

		gotRepo, _, err := testGitlabClient.Groups.GetGroup(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func getGroupHookImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		hookID := rs.Primary.ID
		if hookID == "" {
			return "", fmt.Errorf("No hook ID is set")
		}
		groupID := rs.Primary.Attributes["group"]
		if groupID == "" {
			return "", fmt.Errorf("No group ID is set")
		}
		return fmt.Sprintf("%s:%s", groupID, hookID), nil
	}
}

func testAccGitlabGroupHookConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_group_hook" "foo" {
  group = "${gitlab_group.foo.id}"
  url = "https://example.com/hook-%d"
}
	`, rInt, rInt)
}

func testAccGitlabGroupHookUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_group_hook" "foo" {
  group = "${gitlab_group.foo.id}"
  url = "https://example.com/hook-%d"
  enable_ssl_verification = false
  push_events = true
  push_events_branch_filter = "devel"
  issues_events = false
  confidential_issues_events = false
  merge_requests_events = true
  tag_push_events = true
  note_events = true
  confidential_note_events = true
  job_events = true
  pipeline_events = true
  wiki_page_events = true
  deployment_events = true
  releases_events = true
  subgroup_events = true
}
	`, rInt, rInt)
}
