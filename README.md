# DevFlow
__NOTE__: This project is in beta

This is mainly for my fellow friends at ApoEx.

## Why?
I'm lazy. Like, super lazy. But I also really like the workflow at my current
workplace.

We plan our sprint in TargetProcess and then we name our branches to match the
UserStory in TP:  
`:tp_id-some-descriptive-title`.  

My shortterm memory is really bad, so I always end up going back to TP to see
find the ID for the story that I'm working on.  

This gem aims to solve that.

## How?
1. Install the binary via homebrew: `brew install micke/devflow/devflow`
2. The bin accepts settings from flags or it can also read them from config file a file in `~/.devflow.yaml`
  * `accesstoken` Your TP access token
  * `baseurl` Your organizations TP URL, ex: `https://project.tpondemand.com`
  * `userid` Your TP user id
3. run `devflow checkout` when you start work on a new story
4. Profit?

# Bonus

Install [hub](https://github.com/github/hub) and configure hub to work with
devflow. This will use the `pr` command from devflow to populate your template
with story id, url and title.

### To open the PR in your `$EDITOR`
```
[alias]
  p-r = "!devflow pr | hub -c core.commentChar=';' pull-request -oe -F -"
```

### To open the PR in your browser directly
```
[alias]
  p-r = "!devflow pr | hub -c core.commentChar=';' pull-request -o --no-edit -F -"
```
