
## Feature 4: SQLite index
  - [x] 4.1 Create and migrate sqlite db schema (name, value, description, source, mtime, global (boolean), suffix (boolean))
  - [x] 4.2 add `index_folders` to config 1, or array of folders, accept glob patterns, "!" exclude, "*." etc.
  - [x] 4.3 Index aliases from alias file + all `index_folders` glob matches
  - [x] 4.4 Async background refresh on shell load, debounced by `cache_interval`
    - notes: must not block shell startup; skip re-index if mtime unchanged

    ## Feature 8: `al search` command
  - [x] 8.1 Non-interactive search against sqlite index
    - [x] search across name, value, and description
  - [x] 8.2 searches ending with "=" search for exact alias name
  - [x] 8.3 Output matching aliases (name, value, description,location [grouped by file])
  - [x] 8.4 update interactive list to populate and search against sql index as well

## Feature 9: `script_icons` display
  - [ ] 9.1 Detect tokens in alias value and map to devicons
  - [ ] 9.2 Show icon in list view next to alias value


## Feature 11: View modes
  - [ ] 11.1 Config option to select list style: stacked, tight, table, grouped
  - [ ] 11.2 Implement grouped view (group by common prefix or tag)
