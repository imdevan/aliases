# v0.2.0
## Feature 1: overlay interactive inputs
- [ ] 1.1 use [bubble-overlay](https://github.com/floatpane/bubble-overlay) for interactive crud operations and confirmation dialogs
  - keep list view visible and render inputs as overlays
  - [ ] 1.1.1 add
  - [ ] 1.1.2 edit
  - [ ] 1.1.3 delete
  - [ ] 1.1.4 overwrite


## Feature 2: fixes
- [ ] 2.1 add error color to config use error color for * input cannot be empty on add

## Feature 3: improve bookmark created 
- [ ] 3.1 border
- [ ] 3.2 some color

## Feature 4: improve bookmark edit
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
