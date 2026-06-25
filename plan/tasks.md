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
- [ ] 1.4 remove the confirmation messages for add, edit, and delete in interactive mode. 
- when listview is open
- when called on their own ineractively is okay
- "Updated:", "Cancelled", etc

## Feature 2: fixes
- [ ] 2.1 add error color to config use error color for * input cannot be empty on add

## Feature 3: improve bookmark add created 
- [ ] 3.1 border
- [ ] 3.2 some color

## Feature 4: improve bookmark edit confirmation
- [ ] 4.1 should use same confirmation wrapper as bookmark created (share border and color styles)
- [ ] 4.2 if name changed show from -> to

## Feature 5: config
- [ ] 5.1 bookmark config; first call behavior
- check if initiated
  - if no: prompt user if they would like to init
  - if yes: continue to open config
- [ ] 5.2 config init should check that bookmarks_location is imported in the users shell
  - if no: prompt user if they would like to add it

## Feature 6: nu, fish, bash, aliases functions
- [ ] 6.1 add alias functions like those in the zsh implementation. 
  - the bellow should work the the nu and fish implementations (but account for those shell syntax differences)
~/.bookmarks/bookmarks.sh:4
```bash
# Function wrapper to auto-source bookmarks after running bookmark commands
# This ensures new/updated bookmarks are immediately available in your shell
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
