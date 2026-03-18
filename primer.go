package main

import (
	"regexp"
)

type Primer struct {
	ID       string    `yaml:"id"`
	AIReview string    `yaml:"ai_review"`
	Type     string    `yaml:"type"`
	Filters  FilterSet `yaml:",inline"`
	Content  string
}

type PrimerMatch struct {
	Primer Primer
	Files  []string
}

func LoadPrimers(searchPaths []string, repo string, headSHA string, oh *OutputHandler) ([]Primer, error) {
	scanner := NewScanner(searchPaths, repo, headSHA, oh)
	results, err := scanner.Load("primer", func() interface{} { return &Primer{} })
	if err != nil && len(results) == 0 {
		return nil, err
	}
	if err != nil {
		oh.Printf("Warning: issues encountered while loading primers: %v\n", err)
	}

	var primers []Primer
	for _, res := range results {
		p := res.Raw.(*Primer)
		p.Content = res.Instructions
		primers = append(primers, *p)
	}

	return primers, nil
}

func (rc *RunConfig) FindMatches(personaContext *PRContext) []PrimerMatch {
	var matches []PrimerMatch

	// Pre-compile regexes for all primers
	type compiledPrimer struct {
		primer Primer
		fs     *FilterSet
	}
	var compiledPrimers []compiledPrimer
	for _, p := range rc.Primers {
		fs := p.Filters
		for _, r := range fs.RawRegexFilters {
			re, err := regexp.Compile(r)
			if err != nil {
				rc.OutputHandler.Printf("    Warning: invalid regex %s in primer %s: %v\n", r, p.ID, err)
				continue
			}
			fs.RegexFilters = append(fs.RegexFilters, re)
		}
		compiledPrimers = append(compiledPrimers, compiledPrimer{primer: p, fs: &fs})
	}

	for _, cp := range compiledPrimers {
		var matchedFiles []string
		for _, fileCtx := range personaContext.Files {
			if fileCtx.Matches(FileMatchOptions{
				FilterSet:  cp.fs,
				Branch:     personaContext.Branch,
				CommitDate: personaContext.CommitDate,
			}) {
				matchedFiles = append(matchedFiles, fileCtx.Filename)
			}
		}

		if len(matchedFiles) > 0 {
			matches = append(matches, PrimerMatch{
				Primer: cp.primer,
				Files:  matchedFiles,
			})
		}
	}

	return matches
}
