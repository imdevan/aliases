# Kicked

> ideas and features kicked down stream

- [ ] 2 bookmark sync command 
  - syncs bookmarks file based on:
    - **these should probably just be vars at the top of the bookmarks file**
    - config.shell
    - prioritize first shell in list as source of truth
    - confirm with user before updating out of sync alternate shells
    - config.home_char
    - config.navigation_tool
    - config.editor


## Feature 4: config
- [ ] 4.1 bookmark config; first call behavior
- check if initiated
  - if no: prompt user if they would like to init
  - if yes: continue to open config

## Feature 4: config
- [ ] 4.1 config init should check that bookmarks_location is imported in the users shell
  - if no: prompt user if they would like to add it

## feature 6: nu, fish, bash, aliases functions

## Feature 5: bookmark sync
- [ ] 5.1 sync to _
  - [ ] 5.1.0 sync from the current config to shell 
  - [ ] 5.1.1 sync strategies
  - [ ] 5.1.2 add sync_strategy = "add and skip" | "add and replace" | "overwrite" to config
- [ ] 5.2 sync from _

## Feature 5: bookmark sync
- [ ] 5.1 sync to _
  - [ ] 5.1.0 sync from the current config to shell 
  - [ ] 5.1.1 sync strategies
  - [ ] 5.1.2 add sync_strategy = "add and skip" | "add and replace" | "overwrite" to config
- [ ] 5.2 sync from _

