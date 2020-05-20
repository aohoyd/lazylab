# GitLab merge requests viewer

This is a viewer for GitLab's merge requests. Filters are currently configured right in [main function](main.go).

## Control

- arrows + `hjkl`: move cursor
- `Enter` in bottom menu: choose section
- `Enter` in MRs list: show current MR on fullscreen
- `Esc` in MRs fullscreen mode: exit from fullscreen
- `Esc` in usual mode: exit from MRs section

## Configuration

This tool requires `~/.lazylab` file with next json:
```json
{"server": "https://your_gitlab_server", "user": "your_username", "user_id": userid_int, "token": "API_access_token"}
```

[How to get personal API access token](https://docs.gitlab.com/ce/user/profile/personal_access_tokens.html)