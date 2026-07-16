
## Feature 4: SQLite index
  - [ ] 4.1 Create and migrate sqlite db schema (name, value, description, source, mtime)
  - [ ] 4.2 Index aliases from alias file + all `index_folders` glob matches
  - [ ] 4.3 Async background refresh on shell load, debounced by `cache_interval`
    - notes: must not block shell startup; skip re-index if mtime unchanged

    ## Feature 8: `am search`
  - [ ] 8.1 Non-interactive search against sqlite index
  - [ ] 8.2 Output matching aliases (name, value, description)

# v0.2.0

## Feature 10: `script_icons` display
  - [ ] 10.1 Detect tokens in alias value and map to devicons
  - [ ] 10.2 Show icon in list view next to alias value

## Feature 11: Index folders
  - [ ] 11.1 Support glob patterns in `index_folders` config
  - [ ] 11.2 Scan matched files, merge aliases into db with source attribution
  - [ ] 11.3 Edit/delete aliases in indexed (non-default) files

## Feature 12: View modes
  - [ ] 12.1 Config option to select list style: stacked, tight, table, grouped
  - [ ] 12.2 Implement grouped view (group by common prefix or tag)
