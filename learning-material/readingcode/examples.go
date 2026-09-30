// Code generated from reading Markdown; DO NOT EDIT.

package readingcode

import "sort"

// Source: reading/01-go.md

type Movie struct {
    Title string
    Minutes int
}

func LookupMovie(movies map[string]Movie,
    id string) (Movie, bool) {
    movie, found := movies[id]
    return movie, found
}

// Source: reading/02-patterns.md

func LongestTwo(titles []string) int {
    count := make(map[string]int)
    left, best := 0, 0
    for right, title := range titles {
        count[title]++
        for len(count) > 2 {
            old := titles[left]
            count[old]--
            if count[old] == 0 {
                delete(count, old)
            }
            left++
        }
        length := right - left + 1
        if length > best {
            best = length
        }
    }
    return best
}

// Source: reading/03-order.md

type Version struct {
    At int64
    Value string
}

func VersionAt(v []Version, at int64) (string, bool) {
    lo, hi := 0, len(v)
    for lo < hi {
        mid := lo + (hi-lo)/2
        if v[mid].At <= at {
            lo = mid + 1
        } else {
            hi = mid
        }
    }
    if lo == 0 {
        return "", false
    }
    return v[lo-1].Value, true
}

// Source: reading/04-graphs.md

func HopDistances(edges map[string][]string,
    start string) map[string]int {
    distance := map[string]int{start: 0}
    queue := []string{start}
    for head := 0; head < len(queue); head++ {
        from := queue[head]
        for _, to := range edges[from] {
            if _, seen := distance[to]; seen {
                continue
            }
            distance[to] = distance[from] + 1
            queue = append(queue, to)
        }
    }
    return distance
}

// Source: reading/05-states.md

func BestNonAdjacent(reward []int64) int64 {
    var next, afterNext int64
    for i := len(reward) - 1; i >= 0; i-- {
        take := reward[i] + afterNext
        current := next // Skip position i.
        if take > current {
            current = take
        }
        afterNext = next
        next = current
    }
    return next
}

// Source: reading/06-cache.md

type cacheNode struct {
    key string
    value int
    prev, next *cacheNode
}

type ReadCache struct {
    byKey map[string]*cacheNode
    head, tail *cacheNode
}

func detach(n *cacheNode) {
    n.prev.next = n.next
    n.next.prev = n.prev
    n.prev, n.next = nil, nil
}

func insertFront(head, n *cacheNode) {
    first := head.next
    n.prev, n.next = head, first
    head.next = n
    first.prev = n
}

func (c *ReadCache) Get(key string) (int, bool) {
    n, found := c.byKey[key]
    if !found {
        return 0, false
    }
    detach(n)
    insertFront(c.head, n)
    return n.value, true
}

// Source: reading/07-time.md

type TimedValue struct {
    Value string
    Deadline int64
    Generation uint64
}

type Expiry struct {
    Key string
    Deadline int64
    Generation uint64
}

func ApplyExpiry(current map[string]TimedValue,
    record Expiry, now int64) bool {
    if record.Deadline > now {
        return false
    }
    entry, found := current[record.Key]
    if !found || entry.Generation != record.Generation {
        return false
    }
    delete(current, record.Key)
    return true
}

// Source: reading/08-limits.md

type WindowLimiter struct {
    Window int64
    Limit int
    accepted []int64
}

func (l *WindowLimiter) Allow(now int64) bool {
    cutoff := now - l.Window
    expired := 0
    for expired < len(l.accepted) &&
        l.accepted[expired] <= cutoff {
        expired++
    }
    l.accepted = l.accepted[expired:]
    if len(l.accepted) >= l.Limit {
        return false
    }
    l.accepted = append(l.accepted, now)
    return true
}

// Source: reading/09-events.md

type JobRecord struct {
    ID string
    Due int64
    Generation uint64
}

func ClaimJob(current map[string]uint64,
    record JobRecord, now int64) bool {
    if record.Due > now {
        return false
    }
    generation, found := current[record.ID]
    if !found || generation != record.Generation {
        return false
    }
    delete(current, record.ID)
    return true
}

// Source: reading/10-recommend.md

type Candidate struct {
    ID string
    Score int64
}

func TopTitles(candidates []Candidate, k int) []Candidate {
    if k <= 0 {
        return nil
    }
    ranked := append([]Candidate(nil), candidates...)
    sort.Slice(ranked, func(i, j int) bool {
        a, b := ranked[i], ranked[j]
        if a.Score != b.Score {
            return a.Score > b.Score
        }
        return a.ID < b.ID
    })
    if k < len(ranked) {
        ranked = ranked[:k]
    }
    return ranked
}

// Source: reading/11-practical.md

type Edit struct {
    Value string
    Deleted bool
}

func ReadLayered(base map[string]string,
    layers []map[string]Edit, key string) (string, bool) {
    for i := len(layers) - 1; i >= 0; i-- {
        edit, found := layers[i][key]
        if !found {
            continue
        }
        if edit.Deleted {
            return "", false
        }
        return edit.Value, true
    }
    value, found := base[key]
    return value, found
}

// Source: reading/12-case-studies.md

type WatchEvent struct {
    ID, Title string
    At, Minutes int64
}

func RecentTotals(events []WatchEvent,
    now, window int64) map[string]int64 {
    seen := make(map[string]bool)
    totals := make(map[string]int64)
    for _, event := range events {
        if event.At <= now-window || event.At > now {
            continue
        }
        if seen[event.ID] {
            continue
        }
        seen[event.ID] = true
        totals[event.Title] += event.Minutes
    }
    return totals
}
