package codes

import "testing"

func testResolver(t *testing.T) *Resolver {
	t.Helper()
	resolver, err := Load("../../data/codes.json")
	if err != nil {
		t.Fatal(err)
	}
	return resolver
}
func TestDatasetCompleteness(t *testing.T) {
	data := testResolver(t).Dataset()
	if len(data.Namespaces) != 3 || len(data.Codes) != 293 {
		t.Fatalf("got %d namespaces and %d codes", len(data.Namespaces), len(data.Codes))
	}
}
func TestResolveDistrict(t *testing.T) {
	result := testResolver(t).Resolve("gss-phc-2021-district", "0101", "2021-06-27")
	if result.Status != "resolved" || result.Candidates[0].Name != "Jomoro" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
func TestCrosswalkRegionCodeToAbbreviation(t *testing.T) {
	result := testResolver(t).Crosswalk("gss-phc-2021-region", "03", "gss-phc-2021-region-abbreviation", "2021-06-27")
	if result.Status != "resolved" || result.Candidates[0].Value != "GAR" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
func TestEffectiveDate(t *testing.T) {
	result := testResolver(t).Resolve("gss-phc-2021-district", "0101", "2020-01-01")
	if result.Status != "not_found" {
		t.Fatalf("historically premature code resolved")
	}
}
func TestAmbiguityReturnsCandidates(t *testing.T) {
	end := "2025-12-31"
	data := Dataset{Version: "test", Namespaces: []Namespace{{ID: "a"}, {ID: "b"}}, Codes: []Code{{Namespace: "a", Value: "x", EntityID: "one", ValidFrom: "2020-01-01", ValidTo: &end}, {Namespace: "b", Value: "z", EntityID: "one", ValidFrom: "2020-01-01"}, {Namespace: "b", Value: "y", EntityID: "one", ValidFrom: "2020-01-01"}}}
	resolver, _ := New(data)
	result := resolver.Crosswalk("a", "x", "b", "2024-01-01")
	if result.Status != "ambiguous" || len(result.Candidates) != 2 {
		t.Fatalf("ambiguity was guessed away: %#v", result)
	}
}
