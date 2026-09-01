package codes

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"strings"
	"time"
)

type Namespace struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Authority  string `json:"authority"`
	EntityType string `json:"entityType"`
}
type Code struct {
	Namespace    string   `json:"namespace"`
	Value        string   `json:"code"`
	Name         string   `json:"name"`
	EntityID     string   `json:"entityId"`
	EntityType   string   `json:"entityType"`
	DistrictType string   `json:"districtType,omitempty"`
	RegionCode   string   `json:"regionCode,omitempty"`
	ValidFrom    string   `json:"validFrom"`
	ValidTo      *string  `json:"validTo"`
	Status       string   `json:"status"`
	SourceID     string   `json:"sourceId"`
	Aliases      []string `json:"aliases,omitempty"`
}
type Dataset struct {
	Version             string      `json:"version"`
	EffectiveAt         string      `json:"effectiveAt"`
	GeneratedFromSHA256 string      `json:"generatedFromSha256"`
	Namespaces          []Namespace `json:"namespaces"`
	Codes               []Code      `json:"codes"`
}
type Resolver struct {
	dataset  Dataset
	byKey    map[string][]Code
	byEntity map[string][]Code
}
type Resolution struct {
	Status     string            `json:"status"`
	Query      map[string]string `json:"query"`
	Candidates []Code            `json:"candidates"`
	Version    string            `json:"version"`
}

func key(namespace, value string) string {
	return strings.ToLower(namespace) + "\x00" + strings.ToUpper(strings.TrimSpace(value))
}
func Load(path string) (*Resolver, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var data Dataset
	if err = json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}
	return New(data)
}
func New(data Dataset) (*Resolver, error) {
	namespaces := map[string]bool{}
	for _, namespace := range data.Namespaces {
		namespaces[namespace.ID] = true
	}
	byKey := map[string][]Code{}
	byEntity := map[string][]Code{}
	for _, code := range data.Codes {
		if !namespaces[code.Namespace] {
			return nil, errors.New("code references unknown namespace")
		}
		byKey[key(code.Namespace, code.Value)] = append(byKey[key(code.Namespace, code.Value)], code)
		byEntity[code.EntityID] = append(byEntity[code.EntityID], code)
	}
	return &Resolver{data, byKey, byEntity}, nil
}
func activeAt(code Code, at string) bool {
	if at == "" {
		at = time.Now().UTC().Format("2006-01-02")
	}
	if at < code.ValidFrom {
		return false
	}
	return code.ValidTo == nil || at <= *code.ValidTo
}
func (r *Resolver) Dataset() Dataset { return r.dataset }
func (r *Resolver) Resolve(namespace, value, at string) Resolution {
	candidates := []Code{}
	for _, code := range r.byKey[key(namespace, value)] {
		if activeAt(code, at) {
			candidates = append(candidates, code)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].EntityID < candidates[j].EntityID })
	status := "not_found"
	if len(candidates) == 1 {
		status = "resolved"
	} else if len(candidates) > 1 {
		status = "ambiguous"
	}
	return Resolution{status, map[string]string{"namespace": namespace, "code": value, "at": at}, candidates, r.dataset.Version}
}
func (r *Resolver) Crosswalk(fromNamespace, value, toNamespace, at string) Resolution {
	source := r.Resolve(fromNamespace, value, at)
	if source.Status != "resolved" {
		return source
	}
	candidates := []Code{}
	for _, code := range r.byEntity[source.Candidates[0].EntityID] {
		if code.Namespace == toNamespace && activeAt(code, at) {
			candidates = append(candidates, code)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Value < candidates[j].Value })
	status := "not_found"
	if len(candidates) == 1 {
		status = "resolved"
	} else if len(candidates) > 1 {
		status = "ambiguous"
	}
	return Resolution{status, map[string]string{"fromNamespace": fromNamespace, "code": value, "toNamespace": toNamespace, "at": at}, candidates, r.dataset.Version}
}
