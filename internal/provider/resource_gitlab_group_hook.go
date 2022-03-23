package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_group_hook` + "`" + ` resource allows to manage the lifecycle of a group hook.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#hooks)`,

		CreateContext: resourceGitlabGroupHookCreate,
		ReadContext:   resourceGitlabGroupHookRead,
		UpdateContext: resourceGitlabGroupHookUpdate,
		DeleteContext: resourceGitlabGroupHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabGroupHookStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The name or id of the group to add the hook to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "The url of the hook to invoke.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Description: "A token to present when invoking the hook. The token is not available for imported resources.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"push_events": {
				Description: "Invoke the hook for push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"push_events_branch_filter": {
				Description: "Invoke the hook for push events on matching branches only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"issues_events": {
				Description: "Invoke the hook for issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"confidential_issues_events": {
				Description: "Invoke the hook for confidential issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"merge_requests_events": {
				Description: "Invoke the hook for merge requests.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tag_push_events": {
				Description: "Invoke the hook for tag push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"note_events": {
				Description: "Invoke the hook for notes events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"confidential_note_events": {
				Description: "Invoke the hook for confidential notes events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"job_events": {
				Description: "Invoke the hook for job events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"pipeline_events": {
				Description: "Invoke the hook for pipeline events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"wiki_page_events": {
				Description: "Invoke the hook for wiki page events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"deployment_events": {
				Description: "Invoke the hook for deployment events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"releases_events": {
				Description: "Invoke the hook for releases events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"subgroup_events": {
				Description: "Invoke the hook when a subgroup is created or removed.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"enable_ssl_verification": {
				Description: "Enable ssl verification when invoking the hook.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
})

func resourceGitlabGroupHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	options := &gitlab.AddGroupHookOptions{
		URL:                      gitlab.String(d.Get("url").(string)),
		PushEvents:               gitlab.Bool(d.Get("push_events").(bool)),
		PushEventsBranchFilter:   gitlab.String(d.Get("push_events_branch_filter").(string)),
		IssuesEvents:             gitlab.Bool(d.Get("issues_events").(bool)),
		ConfidentialIssuesEvents: gitlab.Bool(d.Get("confidential_issues_events").(bool)),
		MergeRequestsEvents:      gitlab.Bool(d.Get("merge_requests_events").(bool)),
		TagPushEvents:            gitlab.Bool(d.Get("tag_push_events").(bool)),
		NoteEvents:               gitlab.Bool(d.Get("note_events").(bool)),
		ConfidentialNoteEvents:   gitlab.Bool(d.Get("confidential_note_events").(bool)),
		JobEvents:                gitlab.Bool(d.Get("job_events").(bool)),
		PipelineEvents:           gitlab.Bool(d.Get("pipeline_events").(bool)),
		WikiPageEvents:           gitlab.Bool(d.Get("wiki_page_events").(bool)),
		DeploymentEvents:         gitlab.Bool(d.Get("deployment_events").(bool)),
		ReleasesEvents:           gitlab.Bool(d.Get("releases_events").(bool)),
		SubGroupEvents:           gitlab.Bool(d.Get("subgroup_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if v, ok := d.GetOk("token"); ok {
		options.Token = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab group hook %q", *options.URL)

	hook, _, err := client.Groups.AddGroupHook(group, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", hook.ID))
	d.Set("token", options.Token)

	return resourceGitlabGroupHookRead(ctx, d, meta)
}

func resourceGitlabGroupHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab group hook %s/%d", group, hookId)

	hook, _, err := client.Groups.GetGroupHook(group, hookId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group hook not found %s/%d", group, hookId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("url", hook.URL)
	d.Set("push_events", hook.PushEvents)
	d.Set("push_events_branch_filter", hook.PushEventsBranchFilter)
	d.Set("issues_events", hook.IssuesEvents)
	d.Set("confidential_issues_events", hook.ConfidentialIssuesEvents)
	d.Set("merge_requests_events", hook.MergeRequestsEvents)
	d.Set("tag_push_events", hook.TagPushEvents)
	d.Set("note_events", hook.NoteEvents)
	d.Set("confidential_note_events", hook.ConfidentialNoteEvents)
	d.Set("job_events", hook.JobEvents)
	d.Set("pipeline_events", hook.PipelineEvents)
	d.Set("wiki_page_events", hook.WikiPageEvents)
	d.Set("deployment_events", hook.DeploymentEvents)
	d.Set("releases_events", hook.ReleasesEvents)
	d.Set("subgroup_events", hook.SubGroupEvents)
	d.Set("enable_ssl_verification", hook.EnableSSLVerification)
	return nil
}

func resourceGitlabGroupHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	options := &gitlab.EditGroupHookOptions{
		URL:                      gitlab.String(d.Get("url").(string)),
		PushEvents:               gitlab.Bool(d.Get("push_events").(bool)),
		PushEventsBranchFilter:   gitlab.String(d.Get("push_events_branch_filter").(string)),
		IssuesEvents:             gitlab.Bool(d.Get("issues_events").(bool)),
		ConfidentialIssuesEvents: gitlab.Bool(d.Get("confidential_issues_events").(bool)),
		MergeRequestsEvents:      gitlab.Bool(d.Get("merge_requests_events").(bool)),
		TagPushEvents:            gitlab.Bool(d.Get("tag_push_events").(bool)),
		NoteEvents:               gitlab.Bool(d.Get("note_events").(bool)),
		ConfidentialNoteEvents:   gitlab.Bool(d.Get("confidential_note_events").(bool)),
		JobEvents:                gitlab.Bool(d.Get("job_events").(bool)),
		PipelineEvents:           gitlab.Bool(d.Get("pipeline_events").(bool)),
		WikiPageEvents:           gitlab.Bool(d.Get("wiki_page_events").(bool)),
		DeploymentEvents:         gitlab.Bool(d.Get("deployment_events").(bool)),
		ReleasesEvents:           gitlab.Bool(d.Get("releases_events").(bool)),
		SubGroupEvents:           gitlab.Bool(d.Get("subgroup_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if d.HasChange("token") {
		options.Token = gitlab.String(d.Get("token").(string))
	}

	log.Printf("[DEBUG] update gitlab group hook %s", d.Id())

	_, _, err = client.Groups.EditGroupHook(group, hookId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabGroupHookRead(ctx, d, meta)
}

func resourceGitlabGroupHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab group hook %s", d.Id())

	_, err = client.Groups.DeleteGroupHook(group, hookId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabGroupHookStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Group Hook import format; expected '{group_id}:{hook_id}'")
	}
	group, id := s[0], s[1]

	d.SetId(id)
	d.Set("group", group)

	return []*schema.ResourceData{d}, nil
}
