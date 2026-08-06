package feedback

import (
	"fmt"
	"sort"
	"strings"
)

func GenerateReport(entries []Entry) string {
	if len(entries) == 0 {
		return "No feedback recorded yet. Use `virgil feedback gap|friction|win|session` to start collecting."
	}

	var b strings.Builder
	b.WriteString("## Virgil Feedback Report\n")

	gaps := filterByType(entries, "gap")
	frictions := filterByType(entries, "friction")
	wins := filterByType(entries, "win")
	sessions := filterByType(entries, "session")

	writeSection(&b, "Top Gaps (by frequency x impact)", aggregateWeighted(gaps))
	writeSectionSimple(&b, "Top Friction Points", aggregateSimple(frictions))
	writeSectionSimple(&b, "Wins (preserve these)", aggregateSimple(wins))
	writeSessionSummary(&b, sessions)
	writeDocCoverage(&b, entries)

	return b.String()
}

func filterByType(entries []Entry, t string) []Entry {
	var result []Entry
	for _, e := range entries {
		if e.Type == t {
			result = append(result, e)
		}
	}
	return result
}

func impactWeight(impact string) int {
	switch impact {
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 1
	}
}

type aggregated struct {
	Doc      string
	Scenario string
	Count    int
	Impact   string
	Score    int
}

func aggregateWeighted(entries []Entry) []aggregated {
	type key struct{ Doc, Scenario string }
	m := make(map[key]*aggregated)

	for _, e := range entries {
		k := key{e.Doc, e.Scenario}
		if a, ok := m[k]; ok {
			a.Count++
			w := impactWeight(e.Impact)
			if w > impactWeight(a.Impact) {
				a.Impact = e.Impact
			}
			a.Score += w
		} else {
			w := impactWeight(e.Impact)
			m[k] = &aggregated{
				Doc:      e.Doc,
				Scenario: e.Scenario,
				Count:    1,
				Impact:   e.Impact,
				Score:    w,
			}
		}
	}

	var result []aggregated
	for _, a := range m {
		result = append(result, *a)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Score > result[j].Score
	})
	return result
}

func aggregateSimple(entries []Entry) []aggregated {
	type key struct{ Doc, Scenario string }
	m := make(map[key]*aggregated)

	for _, e := range entries {
		k := key{e.Doc, e.Scenario}
		if a, ok := m[k]; ok {
			a.Count++
		} else {
			m[k] = &aggregated{
				Doc:      e.Doc,
				Scenario: e.Scenario,
				Count:    1,
			}
		}
	}

	var result []aggregated
	for _, a := range m {
		result = append(result, *a)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})
	return result
}

func writeSection(b *strings.Builder, title string, items []aggregated) {
	b.WriteString(fmt.Sprintf("\n### %s\n", title))
	if len(items) == 0 {
		b.WriteString("None recorded.\n")
		return
	}
	for i, item := range items {
		b.WriteString(fmt.Sprintf("%d. %s — %s (%d hits, %s impact)\n", i+1, item.Doc, item.Scenario, item.Count, item.Impact))
	}
}

func writeSectionSimple(b *strings.Builder, title string, items []aggregated) {
	b.WriteString(fmt.Sprintf("\n### %s\n", title))
	if len(items) == 0 {
		b.WriteString("None recorded.\n")
		return
	}
	for i, item := range items {
		b.WriteString(fmt.Sprintf("%d. %s — %s (%d reports)\n", i+1, item.Doc, item.Scenario, item.Count))
	}
}

func writeSessionSummary(b *strings.Builder, sessions []Entry) {
	b.WriteString("\n### Session Summary\n")
	if len(sessions) == 0 {
		b.WriteString("No sessions recorded.\n")
		return
	}

	gradeValues := map[string]int{"A": 4, "B": 3, "C": 2, "D": 1, "F": 0}
	gradeNames := []string{"F", "D", "C", "B", "A"}

	totalGrade := 0
	gradeCount := 0
	docsUsedCount := make(map[string]int)
	docsMissingCount := make(map[string]int)

	for _, s := range sessions {
		if v, ok := gradeValues[s.MethodologyGrade]; ok {
			totalGrade += v
			gradeCount++
		}
		for _, d := range splitCSV(s.DocsUsed) {
			docsUsedCount[d]++
		}
		for _, d := range splitCSV(s.DocsMissing) {
			docsMissingCount[d]++
		}
	}

	avgGrade := "N/A"
	if gradeCount > 0 {
		avgIdx := totalGrade / gradeCount
		if avgIdx >= 0 && avgIdx < len(gradeNames) {
			avgGrade = gradeNames[avgIdx]
		}
	}

	b.WriteString(fmt.Sprintf("- Sessions recorded: %d\n", len(sessions)))
	b.WriteString(fmt.Sprintf("- Average methodology grade: %s\n", avgGrade))
	b.WriteString(fmt.Sprintf("- Most used docs: %s\n", topN(docsUsedCount, 5)))
	b.WriteString(fmt.Sprintf("- Most missed docs: %s\n", topN(docsMissingCount, 5)))
}

func writeDocCoverage(b *strings.Builder, entries []Entry) {
	b.WriteString("\n### Doc Coverage\n")

	type docStats struct {
		SessionsUsed    int
		GapsReported    int
		FrictionReported int
		WinsReported    int
	}

	m := make(map[string]*docStats)

	for _, e := range entries {
		doc := e.Doc
		if doc == "" {
			continue
		}
		if _, ok := m[doc]; !ok {
			m[doc] = &docStats{}
		}
		switch e.Type {
		case "gap":
			m[doc].GapsReported++
		case "friction":
			m[doc].FrictionReported++
		case "win":
			m[doc].WinsReported++
		case "session":
			for _, d := range splitCSV(e.DocsUsed) {
				if _, ok := m[d]; !ok {
					m[d] = &docStats{}
				}
				m[d].SessionsUsed++
			}
		}
	}

	if len(m) == 0 {
		b.WriteString("No doc data available.\n")
		return
	}

	var docs []string
	for doc := range m {
		docs = append(docs, doc)
	}
	sort.Strings(docs)

	b.WriteString("| Doc | Sessions used | Gaps reported | Friction reported | Grade |\n")
	b.WriteString("|-----|--------------|---------------|-------------------|-------|\n")

	for _, doc := range docs {
		s := m[doc]
		grade := docGrade(s.GapsReported, s.WinsReported)
		b.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %s |\n", doc, s.SessionsUsed, s.GapsReported, s.FrictionReported, grade))
	}
}

func docGrade(gaps, wins int) string {
	if gaps > 0 && wins == 0 {
		return "F"
	}
	switch {
	case gaps == 0:
		return "A"
	case gaps == 1:
		return "B"
	case gaps <= 3:
		return "C"
	default:
		return "D"
	}
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func topN(m map[string]int, n int) string {
	if len(m) == 0 {
		return "none"
	}

	type kv struct {
		Key   string
		Value int
	}
	var pairs []kv
	for k, v := range m {
		pairs = append(pairs, kv{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	if len(pairs) > n {
		pairs = pairs[:n]
	}

	var names []string
	for _, p := range pairs {
		names = append(names, p.Key)
	}
	return strings.Join(names, ", ")
}
