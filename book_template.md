# {{.Book.Title}}

**Author:** {{.Book.Author}}  
**Category:** {{.Book.Category}}  
**Source:** {{.Book.Source}}  
**Highlights:** {{.Book.NumHighlights}}  
**Last Highlight:** {{.Book.LastHighlight.Format "January 2, 2006"}}  
**Synced:** {{.SyncDate}}

---

## Overview

{{if .Book.CoverImageURL}}
![Book Cover]({{.Book.CoverImageURL}})
{{end}}

## Highlights & Notes

{{range $index, $highlight := .Highlights}}
### Highlight {{add $index 1}}

> {{$highlight.Text}}

{{if $highlight.Note}}
**My Note:** {{$highlight.Note}}
{{end}}

**Location:** {{$highlight.Location}} ({{$highlight.LocationType}})  
**Highlighted:** {{$highlight.HighlightedAt.Format "January 2, 2006 15:04"}}  
{{if $highlight.Color}}**Color:** {{$highlight.Color}}{{end}}

---
{{end}}

## Summary

*Add your own summary and thoughts here...*

## Action Items

- [ ] Review key concepts
- [ ] Apply insights to current projects
- [ ] Share interesting quotes with team

## Related

*Link to related books, articles, or notes...*