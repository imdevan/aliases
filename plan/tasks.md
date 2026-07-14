<!-- # v0.2.0 -->
## Feature 1: overlay interactive inputs
- [x] 1.1 use [bubble-overlay](https://github.com/floatpane/bubble-overlay) for interactive crud operations and confirmation dialogs
  - keep list view visible and render inputs as overlays
  - [x] 1.1.1 add
  - [x] 1.1.2 edit
  - [x] 1.1.3 delete
  - [x] 1.1.4 overwrite
  - [x] 1.1.5 all content under the overlay should be styled with Faint. cmd/bookmark/root.go:978 
    - only when the modal is open
- [x] 1.2 move overlay to internal/ui module 
  - >or possibly extend the functionality of list view so that less ui logic lives in root.go 
- [x] 1.3 Update edit / add modal title
  - [x] 1.3.1 edit: Edit bookmark
  - [x] 1.3.2 edit, but no bookmark found: '' Not Found, Add Bookmark
  - [x] 1.3.3 add: Add bookmark
- [x] 1.4 remove the confirmation messages for add, edit, and delete in interactive mode. 
- when listview is open
- when called on their own ineractively is okay
- "Updated:", "Cancelled", etc

## Feature 2: fixes
- [x] 2.1 bookmark add should use the same form as interactive add bookmark ( prand bookmark edit)
- [x] 2.3 add error color to config use error color for form validation errors 
  - e.g. * input cannot be empty on add / edit
- [x] 2.4 validate form fields on blur
- [x] 2.5 bookmark add should accept arguments and flags, same as root command
  - -xfistTy, functionally same as calling bookmark [alias] [flags]
  - open interactively prefilled unless -y passed or 
- [ ] 2.5 ensure all args are escaped correctly. files paths, scripts with quotes, descriptions
  - [ ] 2.5.1 add explanation to docs

## Feature 3: import l scripts
- [x] 3.1 import l scripts

## Feature 4: improve bookmark confirmation messages
- [x] 4.1 add config  plain_text = "false" | "true"
  - [x] 4.1.1 add `PlainText bool` field to `domain.Config` with `toml:"plain_text"`
  - [x] 4.1.2 default `false` in `DefaultConfig()`
  - [x] 4.1.3 add `PlainText *bool` to `partialConfig` and `applyPartial` in `config/manager.go`
  - [x] 4.1.4 add commented entry to `renderConfigTemplate` in `config_init.go`
- when true: output confirmation messages, errors, 
- when false (default) show pretty outputs keep as is
- all of the following changes apply when plain_text is false:
- [x] 4.2 border (config border color) + padding(1,1)
- [x] 4.3 use home_icon in place of listing full directory
- [x] 4.4 add success color to config
- [x] 4.5 use success color for successful confirmation messages. title "Bookmark deleted:" is green the bookmark value is default text color
- bookmark added, edited, deleted
- [x] 4.6 add and edit header get their own line. Break after title. 
- [x] 4.7 remove "'" around bookmark value on add and edit
- [x] 4.8 use nerd font icon in header  for success messages
- [x] 4.9 improve Cancel exit_message
  - [x] 4.9.1 should have error color
  - [x] 4.9.2 should have border when not "plain_text"
  - [x] 4.9.3 should have  icon when not "plain_text"
- [x] 4.10 improve delete bookmark confirmation message title should use error color
 
## feature 4: nu, fish, bash, aliases functions
- [x] 5.1 add alias functions like those in the zsh implementation. 
  - the bellow should work the the nu and fish implementations (but account for those shell syntax differences)
~/.bookmarks/bookmarks.sh:4
```bash
# function wrapper to auto-source bookmarks after running bookmark commands
# this ensures new/updated bookmarks are immediately available in your shell
bookmark() {
	command bookmark "$@" && source /home/devy/.bookmarks/bookmarks.sh
}

# Interactive bookmark navigation function
# Displays the TUI and executes the selected bookmark command
bm() {
	local cmd=$(CLICOLOR_FORCE=1 bookmark -i)
	if [[ -n "$cmd" ]]; then
		eval "$cmd"
	fi
	source /home/devy/.bookmarks/bookmarks.sh
}

```

## Feature 5: additional config
- [x] 5.2 confirm_delete = true | false; when false do not ask before deleting a bookmark
  - [x] 5.2.1 add config 
  - [x] 5.2.2 add to config init
  - [x] 5.2.3 add to config docs
