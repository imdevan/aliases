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
- [x] 2.1 bookmark add should use the same form as interactive add bookmark (and bookmark edit)
- [x] 2.3 add error color to config use error color for form validation errors 
  - e.g. * input cannot be empty on add / edit
- [x] 2.4 validate form fields on blur
- [x] 2.5 bookmark add should accept arguments and flags, same as root command
  - -xfistTy, functionally same as calling bookmark [alias] [flags]
  - open interactively prefilled unless -y passed or 
- [ ] 2.5 ensure all quotes are escaped correctly
  - [ ] 2.5.1 add explanation to docs

## Feature 3: import l scripts
- [x] 3.1 import l scripts

## Feature 4: internal/docs
- [ ] 4.1 rewrite scripts/doc_generate.sh as  scripts/docs/gen.go 
  - [ ] 4.1.1 use scripts/docs/templates/ for storing template content over in-lining

## Feature 5: parse commands rewrite
- [ ] 5.1 move parse_commands to scripts/docs/parse_commands
- [ ] 5.2 search for @docs-commands:root, grep, or something not line by line.
  - [ ] 5.2.1 cache location and check there first on subsequent updates?
- [ ] 5.3 update comment parsing
  - see cmd/bookmark/root.go:55 for an example of what command comment docs will look like
- [ ] 5.4 update flag parsing
  - [ ] 5.4.1 flag groups will be defined by comment blocks above a group of flags
    - see cmd/bookmark/root.go:114 for example
  - [ ] 5.4.2 flag group will list flags to include and in which order from the flag prop
  - [ ] 5.4.3 flag shorthands and descriptions should also be included in the docs
- [ ] 5.5 commands and subcommands should be found by following the add command tree from the root command
  - example: cmd/bookmark/root.go:169
  - subcommand example: cmd/bookmark/config.go:41

    

## Feature 6: improve bookmark confirmation messages
- [x] 6.1 add config  plain_text = "false" | "true"
  - [x] 6.1.1 add `PlainText bool` field to `domain.Config` with `toml:"plain_text"`
  - [x] 6.1.2 default `false` in `DefaultConfig()`
  - [x] 6.1.3 add `PlainText *bool` to `partialConfig` and `applyPartial` in `config/manager.go`
  - [x] 6.1.4 add commented entry to `renderConfigTemplate` in `config_init.go`
- when true: output confirmation messages, errors, 
- when false (default) show pretty outputs keep as is
- all of the following changes apply when plain_text is false:
- [x] 6.2 border (config border color) + padding(1,1)
- [x] 6.3 use home_icon in place of listing full directory
- [x] 6.4 add success color to config
- [x] 6.5 use success color for successful confirmation messages. title "Bookmark deleted:" is green the bookmark value is default text color
- bookmark added, edited, deleted
- [x] 6.6 add and edit header get their own line. Break after title. 
- [x] 6.7 remove "'" around bookmark value on add and edit
- [x] 6.8 use nerd font icon in header  for success messages
- [x] 6.9 improve Cancel exit_message
  - [x] 6.9.1 should have error color
  - [x] 6.9.2 should have border when not "plain_text"
  - [x] 6.9.3 should have  icon when not "plain_text"
- [x] 6.10 improve delete bookmark confirmation message title should use error color
 
## feature 4: nu, fish, bash, aliases functions
- [x] 7.1 add alias functions like those in the zsh implementation. 
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

## Feature 7: additional config
- [x] 7.2 confirm_delete = true | false; when false do not ask before deleting a bookmark
  - [x] 7.2.1 add config 
  - [x] 7.2.2 add to config init
  - [x] 7.2.3 add to config docs
