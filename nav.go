package foundation

type navEntry struct {
	viewID string
	data   any
}

type navStack struct {
	entries []navEntry
}

func (s *navStack) push(viewID string, data any) {
	s.entries = append(s.entries, navEntry{viewID: viewID, data: data})
}

func (s *navStack) pop() {
	if len(s.entries) == 0 {
		return
	}
	s.entries = s.entries[:len(s.entries)-1]
}

func (s *navStack) replace(viewID string, data any) {
	if len(s.entries) == 0 {
		return
	}
	s.entries[len(s.entries)-1] = navEntry{viewID: viewID, data: data}
}

func (s *navStack) current() *navEntry {
	if len(s.entries) == 0 {
		return nil
	}
	return &s.entries[len(s.entries)-1]
}

func (s *navStack) len() int {
	return len(s.entries)
}

func (s *navStack) breadcrumbs() []string {
	ids := make([]string, len(s.entries))
	for i, e := range s.entries {
		ids[i] = e.viewID
	}
	return ids
}
