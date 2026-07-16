# Goal
 
create `a` an alias manager for "your favorite shell"

# Context

bookmark already exsists
see ../b for current

# Commands

am
am add
am edit
am delete
am search
am config
am config init

## search
Search for an alias. 

By default search all aliases in `~/.aliases` and `index_folders`.

By default search should print the aliases list response
### flags
- `-o,--only-am` only serach in `~/.aliases` 
- -i, --interactive: open search in interactive mode?
  - could be limitations if not able to run full command in interactive

# Architecture


Default config location ~/.config/alias-manager/config.toml already implemented.

New config opions: `index_folders`: a list of glob patterns for the user to defined directories in addition 
  to the default alias location which the user would like to index. e.g. ~/.bookmarks, or ~/dotfiles which 
  would would allow for alias in those folders to be found when the user uses the `alias search` command

Sqlite db for indexing aliases. To enable fast searching via `alias search`

db should be regularly updated so that if a users changes a file in an indexed folder it it will update the db, but not so much that it the user takes a peformance hit. if  it is a huge thing, possibly just update asynchronously on shell load but cache or debounce at some `cache_interval` so that it is not being called every load. I really want this app to be performant. 

Default alias storage location ~/.alias-manager/aliases.zsh (or appropriate shell). The place where aliases are added via `alias add` 

am likely needs some kind of wrapper function to re-source the alias file after add/edit/delete

# Config options

`script_icons` - show devicons for context of alias value 

```
dliv  - "list containers open in nvim"                               
   docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim
                                                                     
nds  - "run npm build and docker start"                              
   npm run build && docker compuse up -d   # include multiple icons?

```

`custom_icons` - custom icons to show in list view for given tokens (v2)


# Open Questions

## Figure out naming

a | am | alias-manager | alias ?

## edit and delete all indexed aliases?

Even those not in ~/.aliases/aliases.sh?

# View options
## Bookmark default
advantage: least to change, easiest to start impelmentation

 ╭──────────────────────────────────────────────────────────────────────────────────────────────╮
 │    Aliases (524)                                                                            │
 │                                                                                              │
 │ │ a  - "this is a possible description"                                                      │
 │ │  alias                                                                                    │
 │                                                                                              │
 │   ld - "this is a possible description"                                                      │
 │    lazydocker                                                                               │
 │                                                                                              │
 │   dl - "this is a possible description"                                                      │
 │     docker logs                                                                             │
 │                                                                                              │
 │   dli  - "list containers"                                                                   │
 │     docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'                              │
 │                                                                                              │
 │   dliv  - "list containers open in nvim"                                                     │
 │     docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim                       │
 │                                                                                              │
 │   nds  - "run npm build and docker start"                                                    │
 │      npm run build && docker compuse up -d   # include multiple icons?                     │
 │                                                                                              │
 │                                                                                              │
 │   ••••••••••                                                                                 │
 │                                                                                              │
 │   ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more    │
 │                                                                                              │
 ╰──────────────────────────────────────────────────────────────────────────────────────────────╯

 
## Tight optional desc. on bottom
 ╭──────────────────────────────────────────────────────────────────────────────────────────────╮
 │                                                                                              │
 │    Aliases (432)                                                                            │
 │                                                                                              │
 │ │ a       alias                                                                             │
 │ │ "this is a possible description"                                                           │
 │   ld     $  lazydocker  "possible description"                                               │
 │   dl       docker logs                                                                      │
 │   dli      docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'                      │
 │   "Docker - list containers"                                                                 │
 │   dliv     docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim               │
 │   "Docker - list containers - open in vim"                                                   │
 │                                                                                              │
 │                                                                                              │
 │   ••••••••••                                                                                 │
 │                                                                                              │
 │   ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more    │
 │                                                                                              │
 ╰──────────────────────────────────────────────────────────────────────────────────────────────╯

## Grouped

 ╭──────────────────────────────────────────────────────────────────────────────────────────────╮
 │                                                                                              │
 │    Aliases (65)                                                                             │
 │                                                                                              │
 │ │ a       alias                                                                             │
 │ │ "this is a possible description"                                                           │
 │   ld     $  lazydocker  "possible description"                                               │
 │                                                                                              │
 │     Docker                                                                                  │
 │   ────────────────────────────────────────────────────────────────────────────────────────   │
 │   dl       docker logs                                                                      │
 │   dli      docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'                      │
 │   dliv     docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim               │
 │                                                                                              │
 │     NPM                                                                                     │
 │   ────────────────────────────────────────────────────────────────────────────────────────   │
 │                                                                                              │
 │   ••••••••••                                                                                 │
 │                                                                                              │
 │   ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more    │
 │                                                                                              │
 ╰──────────────────────────────────────────────────────────────────────────────────────────────╯

## Space

 ╭──────────────────────────────────────────────────────────────────────────────────────────────╮
 │                                                                                              │
 │    Aliases (65)                                                                             │
 │                                                                                              │
 │ │ a       alias                                                                             │
 │ │ "this is a possible description"                                                           │
 │                                                                                              │
 │   ld     $  lazydocker  "possible description"                                               │
 │   "this is a possible description"                                                           │
 │                                                                                              │
 │   dl       docker logs                                                                      │
 │                                                                                              │
 │   dli      docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'                      │
 │   "this is a possible description"                                                           │
 │                                                                                              │
 │   dliv     docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim               │
 │   "this is a possible description"                                                           │
 │                                                                                              │
 │                                                                                              │
 │   ••••••••••                                                                                 │
 │                                                                                              │
 │   ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more    │
 │                                                                                              │
 ╰──────────────────────────────────────────────────────────────────────────────────────────────╯

## Table view

 ╭──────────────────────────────────────────────────────────────────────────────────────────────╮
 │                                                                                              │
 │    Aliases (1,234)                                                                          │
 │                                                                                              │
 │   a      │ do the thing              │ this is a possible description                       │
 │   ld     │ lazydocker                │ this is a possible description                       │
 │   dl     │ docker logs               │ this is a possible description                       │
 │   dli    │ docker ps --fo..          │ this is a possible description                       │
 │   dliv   │ docker ps --fo..          │ this is a possible description                       │
 │                                                                                              │
 │   ••••••••••                                                                                 │
 │                                                                                              │
 │   ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more    │
 │                                                                                              │
 ╰──────────────────────────────────────────────────────────────────────────────────────────────╯


## Full width table view
 Aliases (1,234)                                                                      
                                                                                       
a      │ do the thing                                                           │ this is a possible description  
ld     │ lazydocker                                                             │ this is a possible description  
dl     │ docker logs                                                            │ this is a possible description  
dli    │ docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'             │ this is a possible description  
dliv   │ docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim      │ this is a possible description  
                                                                                       
••••••••••                                                                             
                                                                                       
↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more


## Full width fat highlight
  Aliases (1,234)                                                                      
                                                                                       
 a      │ do the thing                                                           │ this is a possible description  
│ld     │ lazydocker                                                             
│this is a possible description  
 dl     │ docker logs                                                             │ this is a possible description  
 dli    │ docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}'             │ this is a possible description  
 dliv   │ docker ps --format 'table {{.ID}}	{{.Names}}	{{.Status}}' | nvim      │ this is a possible description  
                                                                                       
 ••••••••••                                                                             
                                                                                        
 ↑/k up • ↓/j down • / filter • enter copy cd command • a add bookmark • q quit • ? more

# Icon library







 │                                                                                              │
